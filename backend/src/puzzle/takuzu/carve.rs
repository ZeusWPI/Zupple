use rand::prelude::SliceRandom;
use rand::thread_rng;
use crate::puzzle::takuzu::{Takuzu, TakuzuCell, TakuzuDifficulty, TakuzuPuzzle};
use crate::puzzle::takuzu::error::TakuzuError;
use crate::puzzle::takuzu::error::TakuzuError::CantSolve;
use crate::puzzle::takuzu::solve::{solve_difficulty};
use crate::puzzle::takuzu::TakuzuCell::Empty;
use crate::puzzle::takuzu::TakuzuDifficulty::Easy;

impl Takuzu {
    pub fn generate_puzzle(&mut self, diff: TakuzuDifficulty, filled_precentage: Option<f64>) -> Result<(), TakuzuError> {
        let cell_count = (self.size as u16).pow(2);
        let minimum_filled_cells = match filled_precentage {
            None => (diff.default_fill_percentage() * f64::from(cell_count)).ceil() as u16,
            Some(percentage) => (percentage * f64::from(cell_count)).ceil() as u16
        };

        let mut best_mask: Vec<Vec<bool>> = vec![];
        let mut best_filled = cell_count + 1;
        for _ in 0..128 {
            let mut mask = full_mask(self.size);
            let mut filled = cell_count;


            // Build list of all coordinates
            let mut positions = Vec::new();
            for r in 0..self.size {
                for c in 0..self.size {
                    positions.push((r, c));
                }
            }

            // Shuffle positions randomly
            positions.shuffle(&mut thread_rng());

            for (r, c) in positions {
                if filled - 1 < minimum_filled_cells {
                    break;
                }

                mask[r as usize][c as usize] = false;

                if !self.puzzle_solvable(&mask, diff) {
                    mask[r as usize][c as usize] = true;
                    continue;
                }

                filled -= 1;
            }
            if !self.mask_matches_difficulty(&mask, diff) {
                continue
            }

            if filled < best_filled {
                best_mask = mask.to_vec();
                best_filled = filled;
            }

            if filled == minimum_filled_cells {
                self.puzzles.push(TakuzuPuzzle {
                    difficulty: diff,
                    mask,
                });
                return Ok(());
            }

            if best_mask.len() != self.size as usize {
                return Err(CantSolve)
            }
        }
        self.puzzles.push(TakuzuPuzzle {
            difficulty: diff,
            mask: best_mask
        });
        Ok(())
    }

    pub fn puzzle_solvable(&self, puzzle: &Vec<Vec<bool>>, diff: TakuzuDifficulty) -> bool {
        let board = self.masked_board(puzzle);
        solve_difficulty(self.size, &board, diff).is_ok()
    }

    fn masked_board(&self, mask: &Vec<Vec<bool>>) -> Vec<Vec<TakuzuCell>> {
        let mut masked = self.grid.to_vec();
        for r in 0..self.size {
            for c in 0..self.size {
                if !mask[r as usize][c as usize] {
                    masked[r as usize][c as usize] = Empty
                }
            }
        }
        masked
    }

    fn mask_matches_difficulty(&self, mask: &Vec<Vec<bool>>, diff: TakuzuDifficulty) -> bool {
        match diff.lower() {
            Ok(lower) => !self.puzzle_solvable(mask, lower),
            Err(_) => self.puzzle_solvable(mask, diff)
        }
    }
}

fn full_mask(size: u8) -> Vec<Vec<bool>> {
    vec![vec![true; size as usize]; size as usize]
}