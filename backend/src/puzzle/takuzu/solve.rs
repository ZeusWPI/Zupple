use crate::puzzle::takuzu::error::TakuzuError;
use crate::puzzle::takuzu::{Takuzu, TakuzuCell, TakuzuDifficulty};
use crate::puzzle::takuzu::error::TakuzuError::{CantSolve, NoUniqueSolution};
use crate::puzzle::takuzu::generate::CellCounts;
use crate::puzzle::takuzu::TakuzuCell::Empty;

struct SolveMode {
    allow_triplets:     bool, // Avoid triplets
    allow_balance:      bool, // Same amount of X and O in a row / column
    allow_uniqueness:   bool, // Each row / column is unique
    allow_backtracking: bool, // Full backtracking
}

impl TakuzuDifficulty {
    fn get_solve_rules(&self) -> SolveMode {
        match self {
            TakuzuDifficulty::Easy => SolveMode {
                allow_triplets: true,
                allow_balance: false,
                allow_uniqueness: false,
                allow_backtracking: false,
            },
            TakuzuDifficulty::Medium => SolveMode {
                allow_triplets: true,
                allow_balance: true,
                allow_uniqueness: false,
                allow_backtracking: false,
            },
            TakuzuDifficulty::Hard => SolveMode {
                allow_triplets: true,
                allow_balance: true,
                allow_uniqueness: true,
                allow_backtracking: false,
            },
            TakuzuDifficulty::Extreme => SolveMode {
                allow_triplets: true,
                allow_balance: true,
                allow_uniqueness: true,
                allow_backtracking: true,
            }
        }
    }
}


pub fn solve_difficulty(size: u8, board: &Vec<Vec<TakuzuCell>>, diff: TakuzuDifficulty) -> Result<Vec<Vec<TakuzuCell>>, TakuzuError> {

    let rules = diff.get_solve_rules();
    let mut grid = board.to_vec();

    loop {
        let mut changed = false;

        if rules.allow_triplets && apply_triplets(size, &mut grid) {
            changed = true;
        }
        if rules.allow_balance && apply_balance(size, &mut grid) {
            changed = true;
        }
        if rules.allow_uniqueness && apply_uniqueness(size, &mut grid) {
            changed = true;
        }

        if is_valid_solution(size, &grid) {
            return Ok(grid);
        }

        if !changed {
            break;
        }
    }

    if !rules.allow_backtracking {
        return Err(CantSolve);
    }

    solve_by_backtracking(size, grid)
}

fn solve_by_backtracking(size: u8, mut board: Vec<Vec<TakuzuCell>>) -> Result<Vec<Vec<TakuzuCell>>, TakuzuError> {
    let mut solutions: u8 = 0;
    fn solve(board: &mut Vec<Vec<TakuzuCell>>, size: u8, row: usize, col: usize, solutions: &mut u8) {
        if *solutions > 1 {
            return
        }
        // Finished entire board
        if row == size as usize {
            if is_valid_solution(size, board) {
                return
            }
        }

        // Move to next row
        let (next_row, next_col) = if col + 1 == size as usize {
            (row + 1, 0)
        } else {
            (row, col + 1)
        };

        if board[row][col] != Empty {
            solve(board, size, next_row, next_col, solutions);
            return
        }

        for candidate in [TakuzuCell::O, TakuzuCell::X] {
            if can_place_cell(board, size, col, row, candidate) {
                board[row][col] = candidate;
                solve(board, size, next_row, next_col, solutions);
            }
        }

        board[row][col] = Empty;
    }

    solve(&mut board, size, 0, 0, &mut solutions);
    match solutions {
        2.. => Err(NoUniqueSolution),
        1 => Ok(board),
        0 => Err(CantSolve)
    }
}

fn apply_triplets(size: u8, board: &mut Vec<Vec<TakuzuCell>>) -> bool {
    if size == 2 {
        return false;
    }

    let mut changed = false;
    for r in 0..size as usize {
        for c in 0..size as usize {
            if board[r][c] == Empty {
                if r>=2 && board[r-1][c] == board[r-2][c] && board[r-2][c] != Empty {
                    board[r][c] = board[r-1][c].other();
                    changed = true;
                }
                if c>=2 && board[r][c-1] == board[r][c-2] && board[r][c-2] != Empty {
                    board[r][c] = board[r][c-1].other();
                    changed = true;
                }
                if r>=1 && r < (size - 1) as usize && board[r-1][c] == board[r+1][c] && board[r-1][c] != Empty {
                    board[r][c] = board[r-1][c].other();
                    changed = true;
                }
                if c>=1 && c < (size - 1) as usize && board[r][c-1] == board[r][c+1] && board[r][c-1] != Empty {
                    board[r][c] = board[r][c-1].other();
                    changed = true;
                }
                if r < (size - 2) as usize && board[r+1][c] == board[r+2][c] && board[r+2][c] != Empty {
                    board[r][c] = board[r+1][c].other();
                    changed = true;
                }
                if c < (size - 2) as usize && board[r][c+1] == board[r][c+2] && board[r][c+2] != Empty {
                    board[r][c] = board[r][c+1].other();
                    changed = true;
                }
            }
        }
    }
    changed
}
fn apply_balance(size: u8, board: &mut Vec<Vec<TakuzuCell>>) -> bool {
    let mut changed = false;
    for i in 0..size as usize {
        let mut row_counts = CellCounts { x_count: 0, o_count: 0 };
        let mut col_counts = CellCounts { x_count: 0, o_count: 0 };
        for j in 0..size as usize {
            row_counts.add(board[i][j]);
            col_counts.add(board[j][i]);
        }
        let row_completion_value = row_counts.completion_value(size);
        let col_completion_value = col_counts.completion_value(size);
        if row_completion_value.is_some() || col_completion_value.is_some() {
            for j in 0..size as usize {
                if board[i][j] == Empty && let Some(fill_value) = row_completion_value {
                    board[i][j] = fill_value;
                    changed = true;
                }
                if board[j][i] == Empty && let Some(fill_value) = col_completion_value {
                    board[j][i] = fill_value;
                    changed = true;
                }
            }
        }
    }
    changed
}
fn apply_uniqueness(size: u8, board: &mut Vec<Vec<TakuzuCell>>) -> bool {
    let mut changed = false;

    for is_column in [false, true] {
        for full_idx in 0..size as usize {
            let is_full = (0..size as usize).all(|j| {
                let cell = if is_column {
                    board[j][full_idx]
                } else {
                    board[full_idx][j]
                };
                cell != Empty
            });

            if !is_full {
                continue;
            }

            for test_idx in 0..size as usize {
                if test_idx == full_idx {
                    continue;
                }

                let mut empty_pos = None;
                let mut matches = true;

                for j in 0..size as usize {
                    let test_cell = if is_column {
                        board[j][test_idx]
                    } else {
                        board[test_idx][j]
                    };

                    let full_cell = if is_column {
                        board[j][full_idx]
                    } else {
                        board[full_idx][j]
                    };

                    match test_cell {
                        Empty => {
                            if empty_pos.is_some() {
                                matches = false;
                                break;
                            }
                            empty_pos = Some(j);
                        }
                        value if value == full_cell => {}
                        _ => {
                            matches = false;
                            break;
                        }
                    }
                }

                if matches {
                    if let Some(pos) = empty_pos {
                        let fill_value = if is_column {
                            board[pos][full_idx].other()
                        } else {
                            board[full_idx][pos].other()
                        };

                        if is_column {
                            board[pos][test_idx] = fill_value;
                        } else {
                            board[test_idx][pos] = fill_value;
                        }

                        changed = true;
                    }
                }
            }
        }
    }

    changed
}

fn is_valid_solution(size: u8, board: &Vec<Vec<TakuzuCell>>) -> bool {

    // Check rows
    for i in 0..size as usize {
        for j in (i + 1)..size as usize {
            if board[i] == board[j] {
                return false;
            }
        }
    }

    // Check columns
    for i in 0..size as usize {
        for j in (i + 1)..size as usize {
            let same = (0..size).all(|row| board[row as usize][i] == board[row as usize][j]);

            if same {
                return false;
            }
        }
    }

    true
}

fn can_place_cell(
    board: &Vec<Vec<TakuzuCell>>,
    size: u8,
    r: usize,
    c: usize,
    value: TakuzuCell,
) -> bool {
    // Row check (triplet rule)
    for offset in TRIPLET_OFFSETS {
        let a = r as isize + offset[0];
        let b = r as isize + offset[1];

        if a >= 0 && b >= 0 && a < size as isize && b < size as isize {
            let a = a as usize;
            let b = b as usize;

            if board[a][c] == value && board[b][c] == value {
                return false;
            }
        }
    }

    // Column check (triplet rule)
    for offset in TRIPLET_OFFSETS {
        let a = c as isize + offset[0];
        let b = c as isize + offset[1];

        if a >= 0 && b >= 0 && a < size as isize && b < size as isize {
            let a = a as usize;
            let b = b as usize;

            if board[r][a] == value && board[r][b] == value {
                return false;
            }
        }
    }

    // Count check (balance rule)
    let mut row_counts = CellCounts {x_count: 0, o_count: 0};
    let mut col_counts = CellCounts {x_count: 0, o_count: 0};

    for i in 0..size as usize {
        row_counts.add(board[r][i]);
        col_counts.add(board[i][c]);
    }

    row_counts.add(value);
    col_counts.add(value);
    col_counts.can_balance_with_remaining(size/2) && row_counts.can_balance_with_remaining(size / 2)
}

const TRIPLET_OFFSETS: [[isize; 2]; 3] = [
    [-1, -2],
    [-1, 1],
    [1, 2],
];