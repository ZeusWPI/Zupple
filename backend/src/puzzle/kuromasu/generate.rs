use crate::puzzle::kuromasu::grid::Grid;
use crate::puzzle::kuromasu::rng::Rng;
use crate::puzzle::kuromasu::solve::{Difficulty, solve_difficulty, solve_unique};
use crate::puzzle::kuromasu::{BLOCKED, Cell, Kuromasu, KuromasuError, Mask};

const BLOCKED_CHANCE: f64 = 0.30;
const MASK_ATTEMPTS: usize = 128;

impl Kuromasu {
    pub(super) fn generate(size: usize) -> Self {
        let mut grid = Grid::empty(size);

        for row in 0..grid.size() {
            for col in 0..grid.size() {
                if rng.gen_bool(BLOCKED_CHANCE) {
                    grid.set(row, col, BLOCKED);
                }
            }
        }

        for row in 0..grid.size() {
            for col in 0..grid.size() {
                if grid.at(row, col) == BLOCKED {
                    continue;
                }

                let sees = grid.sees(row, col) as Cell;
                grid.set(row, col, if sees == 0 { BLOCKED } else { sees });
            }
        }

        Self {
            size,
            field: grid.into_cells(),
        }
    }
}

impl Grid {
    pub(super) fn generate_mask(
        &self,
        difficulty: Difficulty,
        filled_percentage: Option<f64>,
    ) -> Result<Mask, KuromasuError> {
        let minimum_filled_percentage = resolve_filled_percentage(difficulty, filled_percentage)?;
        let cell_count = self.size() * self.size();
        let minimum_filled_cells = min_filled_cells(cell_count, minimum_filled_percentage);
        let mut rng = Rng::from_entropy();

        let mut best_mask = None;
        let mut best_filled = cell_count + 1;

        for _ in 0..MASK_ATTEMPTS {
            let mut mask = full_mask(cell_count);
            let mut filled = cell_count;

            for index in rng.permutation(cell_count) {
                if filled - 1 < minimum_filled_cells {
                    break;
                }

                mask[index] = 0;
                if !self.mask_solvable(&mask, difficulty) {
                    mask[index] = 1;
                    continue;
                }

                filled -= 1;
            }

            if !self.mask_matches_difficulty(&mask, difficulty) {
                continue;
            }

            if filled < best_filled {
                best_filled = filled;
                best_mask = Some(mask.clone());
            }

            if filled == minimum_filled_cells {
                return Ok(mask);
            }
        }

        best_mask.ok_or(KuromasuError::UnreachableMask)
    }

    fn mask_matches_difficulty(&self, mask: &[Cell], difficulty: Difficulty) -> bool {
        match difficulty.lower() {
            Ok(lower) => !self.mask_solvable(mask, lower),
            Err(KuromasuError::NoLowerDifficulty) => self.mask_solvable(mask, difficulty),
            Err(_) => false,
        }
    }

    fn mask_solvable(&self, mask: &[Cell], difficulty: Difficulty) -> bool {
        let Ok(board) = self.masked_board(mask) else {
            return false;
        };

        let (solution, solvable) = solve_difficulty(self.size(), &board, difficulty);
        solvable && solve_unique(self.size(), &board) && self.same_shape(&solution)
    }
}

fn full_mask(cell_count: usize) -> Mask {
    vec![1; cell_count]
}

fn resolve_filled_percentage(
    difficulty: Difficulty,
    filled_percentage: Option<f64>,
) -> Result<f64, KuromasuError> {
    let Some(value) = filled_percentage else {
        return Ok(difficulty.default_filled_percentage());
    };

    if !(0.0..=1.0).contains(&value) {
        return Err(KuromasuError::InvalidFilledPercentage);
    }

    Ok(value)
}

fn min_filled_cells(cell_count: usize, filled_percentage: f64) -> usize {
    (filled_percentage * cell_count as f64).ceil() as usize
}

#[cfg(test)]
mod tests {
    use crate::puzzle::kuromasu::grid::Grid;
    use crate::puzzle::kuromasu::{BLOCKED, Difficulty};

    #[test]
    fn mask_solvable_accepts_hidden_numbered_cells() {
        let grid = Grid::new(
            4,
            vec![
                BLOCKED, 1, BLOCKED, 3, 5, 4, 5, 6, 2, BLOCKED, 3, 4, 5, 3, 5, 6,
            ],
        );
        let mask = vec![1, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1];

        assert!(grid.mask_solvable(&mask, Difficulty::Expert));
    }

    #[test]
    fn filled_percentage_must_be_in_range() {
        let grid = Grid::new(2, vec![BLOCKED, 1, 1, 2]);

        assert!(grid.generate_mask(Difficulty::Easy, Some(-0.1)).is_err());
        assert!(grid.generate_mask(Difficulty::Easy, Some(1.1)).is_err());
    }
}
