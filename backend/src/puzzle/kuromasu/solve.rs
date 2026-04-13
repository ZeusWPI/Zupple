use crate::puzzle::kuromasu::grid::{Axis, Grid};
use crate::puzzle::kuromasu::{Cell, KuromasuError, BLOCKED, FILLED, UNKNOWN};

#[derive(Debug, Clone, Copy, PartialEq, Eq)]
pub enum Difficulty {
    Easy,
    Medium,
    Hard,
    Expert,
}

impl Difficulty {
    pub fn lower(self) -> Result<Self, KuromasuError> {
        match self {
            Self::Medium => Ok(Self::Easy),
            Self::Hard => Ok(Self::Medium),
            Self::Expert => Ok(Self::Hard),
            Self::Easy => Err(KuromasuError::NoLowerDifficulty),
        }
    }

    pub fn default_filled_percentage(self) -> f64 {
        match self {
            Self::Easy => 0.60,
            Self::Medium => 0.55,
            Self::Hard => 0.45,
            Self::Expert => 0.35,
        }
    }

    fn mode(self) -> SolveMode {
        match self {
            Self::Easy => SolveMode {
                allow_saturation: true,
                allow_exact_fill: true,
                allow_single_direction: true,
                allow_sees_nothing: true,
                ..SolveMode::default()
            },
            Self::Medium => SolveMode {
                allow_saturation: true,
                allow_exact_fill: true,
                allow_single_direction: true,
                allow_sees_nothing: true,
                allow_candidate_elimination: true,
                allow_forced_prefix: true,
                ..SolveMode::default()
            },
            Self::Hard => SolveMode {
                allow_saturation: true,
                allow_exact_fill: true,
                allow_single_direction: true,
                allow_sees_nothing: true,
                allow_candidate_elimination: true,
                allow_forced_prefix: true,
                allow_forced_prefix_global: true,
                ..SolveMode::default()
            },
            Self::Expert => SolveMode {
                allow_saturation: true,
                allow_exact_fill: true,
                allow_single_direction: true,
                allow_sees_nothing: true,
                allow_candidate_elimination: true,
                allow_forced_prefix: true,
                allow_forced_prefix_global: true,
                allow_backtracking: true,
            },
        }
    }
}

#[derive(Debug, Default, Clone, Copy, PartialEq, Eq)]
struct SolveMode {
    allow_saturation: bool,
    allow_exact_fill: bool,
    allow_single_direction: bool,
    allow_sees_nothing: bool,
    allow_candidate_elimination: bool,
    allow_forced_prefix: bool,
    allow_forced_prefix_global: bool,
    allow_backtracking: bool,
}

#[derive(Debug, Clone, Copy, PartialEq, Eq)]
struct Direction {
    axis: Axis,
    left: bool,
    open: usize,
}

const DIRECTIONS: [Direction; 4] = [
    Direction {
        axis: Axis::Row,
        left: true,
        open: 0,
    },
    Direction {
        axis: Axis::Row,
        left: false,
        open: 0,
    },
    Direction {
        axis: Axis::Col,
        left: true,
        open: 0,
    },
    Direction {
        axis: Axis::Col,
        left: false,
        open: 0,
    },
];

pub(super) fn solve_difficulty(
    size: usize,
    board: &[Cell],
    difficulty: Difficulty,
) -> (Vec<Cell>, bool) {
    let mut grid = Grid::new(size, board.to_vec());
    let solvable = solve_with_rules(&mut grid, difficulty.mode());

    (grid.into_cells(), solvable)
}

pub(super) fn solve_unique(size: usize, board: &[Cell]) -> bool {
    let mut grid = Grid::new(size, board.to_vec());
    let mut solutions = 0;

    count_solutions(&mut grid, 0, &mut solutions);

    solutions == 1
}

fn solve_with_rules(grid: &mut Grid, mode: SolveMode) -> bool {
    loop {
        let mut changed = false;

        if mode.allow_saturation && apply_saturation(grid) {
            changed = true;
        }
        if mode.allow_exact_fill && apply_exact_fill(grid) {
            changed = true;
        }
        if mode.allow_single_direction && apply_single_direction(grid) {
            changed = true;
        }
        if mode.allow_sees_nothing && apply_sees_nothing(grid) {
            changed = true;
        }
        if mode.allow_candidate_elimination && apply_candidate_elimination(grid) {
            changed = true;
        }
        if mode.allow_forced_prefix && apply_forced_prefix(grid) {
            changed = true;
        }
        if mode.allow_forced_prefix_global && apply_forced_prefix_global(grid) {
            changed = true;
        }

        if grid.solved() {
            return true;
        }

        if !changed {
            break;
        }
    }

    mode.allow_backtracking && solve_by_backtracking(grid, 0)
}

fn get_dirs(grid: &Grid, row: usize, col: usize) -> [Direction; 4] {
    let mut dirs = DIRECTIONS;

    for direction in &mut dirs {
        let (axis, index, pos) = direction.line(row, col);
        direction.open = if direction.left {
            grid.line_open_left(axis, index, pos)
        } else {
            grid.line_open_right(axis, index, pos)
        };
    }

    dirs.sort_by_key(|direction| direction.open);
    dirs
}

impl Direction {
    fn line(self, row: usize, col: usize) -> (Axis, usize, usize) {
        match self.axis {
            Axis::Row => (Axis::Row, row, col),
            Axis::Col => (Axis::Col, col, row),
        }
    }

    fn sees(self, grid: &Grid, row: usize, col: usize) -> usize {
        let (axis, index, pos) = self.line(row, col);
        if self.left {
            grid.line_sees_left(axis, index, pos)
        } else {
            grid.line_sees_right(axis, index, pos)
        }
    }

    fn fill(self, grid: &mut Grid, row: usize, col: usize, max: usize) -> usize {
        let (axis, index, pos) = self.line(row, col);
        if self.left {
            grid.line_fill_left(axis, index, pos, max)
        } else {
            grid.line_fill_right(axis, index, pos, max)
        }
    }
}

fn apply_saturation(grid: &mut Grid) -> bool {
    let mut changed = false;

    for row in 0..grid.size() {
        for col in 0..grid.size() {
            let cell = grid.at(row, col);
            if cell >= BLOCKED || grid.sees(row, col) != usize::from(cell) {
                continue;
            }

            if grid.add_boundary(row, col) {
                changed = true;
            }
        }
    }

    changed
}

fn apply_exact_fill(grid: &mut Grid) -> bool {
    let mut changed = false;

    for row in 0..grid.size() {
        for col in 0..grid.size() {
            let cell = grid.at(row, col);
            if cell >= BLOCKED {
                continue;
            }

            let sees = grid.sees(row, col);
            if sees == usize::from(cell) {
                continue;
            }

            let open = grid.open(row, col);
            if sees + open != usize::from(cell) {
                continue;
            }

            if grid.fill(row, col) > 0 {
                changed = true;
            }
        }
    }

    changed
}

fn apply_single_direction(grid: &mut Grid) -> bool {
    let mut changed = false;

    for row in 0..grid.size() {
        for col in 0..grid.size() {
            let cell = grid.at(row, col);
            if cell >= BLOCKED {
                continue;
            }

            let sees = grid.sees(row, col);
            if sees == usize::from(cell) {
                continue;
            }

            let dirs = get_dirs(grid, row, col);
            if dirs[3].open == 0 || dirs[2].open > 0 {
                continue;
            }

            let max = usize::from(cell) - sees + dirs[3].sees(grid, row, col);
            if dirs[3].fill(grid, row, col, max) > 0 {
                changed = true;
            }
        }
    }

    changed
}

fn apply_sees_nothing(grid: &mut Grid) -> bool {
    let mut changed = false;

    for row in 0..grid.size() {
        for col in 0..grid.size() {
            if grid.at(row, col) != UNKNOWN {
                continue;
            }

            if grid.line_in_bounds(col as isize - 1) && grid.at(row, col - 1) != BLOCKED {
                continue;
            }
            if grid.line_in_bounds(col as isize + 1) && grid.at(row, col + 1) != BLOCKED {
                continue;
            }
            if grid.line_in_bounds(row as isize - 1) && grid.at(row - 1, col) != BLOCKED {
                continue;
            }
            if grid.line_in_bounds(row as isize + 1) && grid.at(row + 1, col) != BLOCKED {
                continue;
            }

            grid.set(row, col, BLOCKED);
            changed = true;
        }
    }

    changed
}

fn apply_candidate_elimination(grid: &mut Grid) -> bool {
    let mut changed = false;

    for row in 0..grid.size() {
        for col in 0..grid.size() {
            if grid.at(row, col) != UNKNOWN {
                continue;
            }

            grid.set(row, col, FILLED);
            if !grid.valid() {
                grid.set(row, col, BLOCKED);
                changed = true;
                continue;
            }

            grid.set(row, col, BLOCKED);
            if !grid.valid() {
                grid.set(row, col, FILLED);
                changed = true;
                continue;
            }

            grid.set(row, col, UNKNOWN);
        }
    }

    changed
}

fn apply_forced_prefix(grid: &mut Grid) -> bool {
    let mut changed = false;

    for row in 0..grid.size() {
        for col in 0..grid.size() {
            let cell = grid.at(row, col);
            if cell >= BLOCKED {
                continue;
            }

            let sees = grid.sees(row, col);
            if sees == usize::from(cell) {
                continue;
            }

            let dirs = get_dirs(grid, row, col);
            let other_max = dirs[0].open + dirs[1].open + dirs[2].open;
            let cells_left = usize::from(cell).saturating_sub(sees + other_max);
            if cells_left == 0 {
                continue;
            }

            let max = cells_left + dirs[3].sees(grid, row, col);
            if dirs[3].fill(grid, row, col, max) > 0 {
                changed = true;
            }
        }
    }

    changed
}

fn apply_forced_prefix_global(grid: &mut Grid) -> bool {
    let mut changed = false;

    for row in 0..grid.size() {
        for col in 0..grid.size() {
            let cell = grid.at(row, col);
            if cell >= BLOCKED {
                continue;
            }

            let sees = grid.sees(row, col);
            if sees == usize::from(cell) {
                continue;
            }

            let dirs = get_dirs(grid, row, col);
            let other_max = dirs[0].open + dirs[1].open + dirs[2].open;
            let mut cells_left = usize::from(cell) as isize - sees as isize - other_max as isize;
            let mut simulated_cells = Vec::with_capacity(other_max);

            for direction in dirs.iter().take(3) {
                let (axis, index, pos) = direction.line(row, col);
                let positions = directional_positions(grid.size(), pos, direction.left);
                let mut considered_unknowns = 0;

                for line_pos in positions {
                    let candidate = grid.line_at(axis, index, line_pos);
                    if candidate == BLOCKED {
                        break;
                    }

                    if candidate != UNKNOWN {
                        continue;
                    }

                    if considered_unknowns == direction.open {
                        break;
                    }
                    considered_unknowns += 1;

                    grid.line_set(axis, index, line_pos, FILLED);
                    if !grid.valid() {
                        grid.line_set(axis, index, line_pos, UNKNOWN);
                        cells_left += 1;
                        break;
                    }

                    let (sim_row, sim_col) = match axis {
                        Axis::Row => (index, line_pos),
                        Axis::Col => (line_pos, index),
                    };
                    simulated_cells.push((sim_row, sim_col));
                }
            }

            if cells_left > 0 {
                let max = cells_left as usize + dirs[3].sees(grid, row, col);
                if dirs[3].fill(grid, row, col, max) > 0 {
                    changed = true;
                }
            }

            for (sim_row, sim_col) in simulated_cells {
                if grid.at(sim_row, sim_col) == FILLED {
                    grid.set(sim_row, sim_col, UNKNOWN);
                }
            }
        }
    }

    changed
}

fn directional_positions(size: usize, pos: usize, left: bool) -> Box<dyn Iterator<Item = usize>> {
    if left {
        Box::new((0..pos).rev())
    } else {
        Box::new(pos + 1..size)
    }
}

fn solve_by_backtracking(grid: &mut Grid, index: usize) -> bool {
    if index == grid.size() * grid.size() {
        return grid.solved() && grid.valid();
    }

    if grid.cells()[index] != UNKNOWN {
        return solve_by_backtracking(grid, index + 1);
    }

    grid.cells_mut()[index] = FILLED;
    if grid.valid() && solve_by_backtracking(grid, index + 1) {
        return true;
    }

    grid.cells_mut()[index] = BLOCKED;
    if grid.valid() && solve_by_backtracking(grid, index + 1) {
        return true;
    }

    grid.cells_mut()[index] = UNKNOWN;
    false
}

fn count_solutions(grid: &mut Grid, index: usize, solutions: &mut usize) {
    if *solutions > 1 {
        return;
    }

    if index == grid.size() * grid.size() {
        if grid.solved() && grid.valid() {
            *solutions += 1;
        }
        return;
    }

    if grid.cells()[index] != UNKNOWN {
        count_solutions(grid, index + 1, solutions);
        return;
    }

    grid.cells_mut()[index] = FILLED;
    if grid.valid() {
        count_solutions(grid, index + 1, solutions);
    }

    grid.cells_mut()[index] = BLOCKED;
    if grid.valid() {
        count_solutions(grid, index + 1, solutions);
    }

    grid.cells_mut()[index] = UNKNOWN;
}

#[cfg(test)]
mod tests {
    use crate::puzzle::kuromasu::grid::Grid;
    use crate::puzzle::kuromasu::solve::{solve_unique, solve_with_rules, SolveMode};
    use crate::puzzle::kuromasu::{BLOCKED, UNKNOWN};

    #[test]
    fn saturation_solves_expected_board() {
        let mut grid = Grid::new(
            4,
            vec![
                UNKNOWN, 1, UNKNOWN, 3, 5, 4, 5, 6, 2, UNKNOWN, 3, 4, 5, 3, 5, 6,
            ],
        );

        assert!(solve_with_rules(
            &mut grid,
            SolveMode {
                allow_saturation: true,
                ..SolveMode::default()
            },
        ));
    }

    #[test]
    fn exact_fill_solves_expected_board() {
        let mut grid = Grid::new(
            4,
            vec![
                BLOCKED, 1, BLOCKED, 3, 5, UNKNOWN, UNKNOWN, UNKNOWN, UNKNOWN, BLOCKED, 3, UNKNOWN,
                UNKNOWN, 3, 5, 6,
            ],
        );

        assert!(solve_with_rules(
            &mut grid,
            SolveMode {
                allow_exact_fill: true,
                ..SolveMode::default()
            },
        ));
    }

    #[test]
    fn single_direction_solves_expected_board() {
        let mut grid = Grid::new(
            4,
            vec![
                BLOCKED, 1, BLOCKED, 3, UNKNOWN, UNKNOWN, UNKNOWN, 6, UNKNOWN, BLOCKED, 3, 4, 5, 3,
                5, 6,
            ],
        );

        assert!(solve_with_rules(
            &mut grid,
            SolveMode {
                allow_single_direction: true,
                ..SolveMode::default()
            },
        ));
    }

    #[test]
    fn sees_nothing_solves_expected_board() {
        let mut grid = Grid::new(
            4,
            vec![
                BLOCKED, BLOCKED, BLOCKED, 3, BLOCKED, UNKNOWN, BLOCKED, 3, 1, BLOCKED, 2, 4, 4, 3,
                4, 6,
            ],
        );

        assert!(solve_with_rules(
            &mut grid,
            SolveMode {
                allow_sees_nothing: true,
                ..SolveMode::default()
            },
        ));
    }

    #[test]
    fn candidate_elimination_solves_expected_board() {
        let mut grid = Grid::new(
            4,
            vec![
                BLOCKED, 1, BLOCKED, 3, 5, 4, 5, 6, 2, UNKNOWN, 3, 4, 5, 3, 5, 6,
            ],
        );

        assert!(solve_with_rules(
            &mut grid,
            SolveMode {
                allow_candidate_elimination: true,
                ..SolveMode::default()
            },
        ));
    }

    #[test]
    fn detects_non_unique_board() {
        let board = vec![UNKNOWN; 4];

        assert!(!solve_unique(2, &board));
    }
}
