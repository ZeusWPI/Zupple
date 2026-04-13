mod error;
mod generate;
mod grid;
mod rng;
mod solve;

pub use error::KuromasuError;
pub use solve::Difficulty;
use std::fmt;

use crate::puzzle::kuromasu::grid::Grid;
use crate::puzzle::kuromasu::rng::Rng;

pub type Cell = u8;
pub type Mask = Vec<Cell>;

pub const MAX_SIZE: usize = 32;
pub const BLOCKED: Cell = MAX_SIZE as Cell + 1;
pub const UNKNOWN: Cell = BLOCKED + 1;
pub const FILLED: Cell = UNKNOWN + 1;

#[derive(Debug, Clone, PartialEq, Eq)]
pub struct Kuromasu {
    size: usize,
    field: Vec<Cell>,
}

impl Kuromasu {
    pub fn new(size: usize) -> Result<Self, KuromasuError> {
        Grid::validate_size(size)?;
        Ok(Self::generate(size))
    }

    pub fn from_field(size: usize, field: Vec<Cell>) -> Result<Self, KuromasuError> {
        Grid::validate_size(size)?;
        if field.len() != size * size {
            return Err(KuromasuError::InvalidFieldLength);
        }

        Ok(Self { size, field })
    }

    pub fn size(&self) -> usize {
        self.size
    }

    pub fn field(&self) -> &[Cell] {
        &self.field
    }

    pub fn puzzle(&self, difficulty: Difficulty) -> Result<Mask, KuromasuError> {
        self.grid().generate_mask(difficulty, None)
    }

    pub fn puzzle_with_filled_percentage(
        &self,
        difficulty: Difficulty,
        filled_percentage: f64,
    ) -> Result<Mask, KuromasuError> {
        self.grid()
            .generate_mask(difficulty, Some(filled_percentage))
    }

    pub fn puzzle_string(&self, mask: &[Cell]) -> Result<String, KuromasuError> {
        self.grid().render(Some(mask))
    }

    fn grid(&self) -> Grid {
        Grid::new(self.size, self.field.clone())
    }
}

impl fmt::Display for Kuromasu {
    fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result {
        let s = self.grid().render(None).map_err(|_| fmt::Error)?;
        f.write_str(&s)
    }
}

#[cfg(test)]
mod tests {
    use crate::puzzle::kuromasu::{BLOCKED, Difficulty, Kuromasu};

    #[test]
    fn generates_a_sized_field() {
        let puzzle = Kuromasu::with_seed(4, 42).expect("valid puzzle");

        assert_eq!(puzzle.size(), 4);
        assert_eq!(puzzle.field().len(), 16);
        assert!(puzzle.field().iter().all(|cell| *cell <= BLOCKED));
    }

    #[test]
    fn rejects_invalid_size() {
        assert!(Kuromasu::new(0).is_err());
        assert!(Kuromasu::new(33).is_err());
    }

    #[test]
    fn renders_masked_puzzle() {
        let puzzle = Kuromasu::from_field(2, vec![BLOCKED, 1, 1, 2]).expect("valid field");

        let rendered = puzzle.puzzle_string(&[1, 0, 1, 1]).expect("valid mask");

        assert!(rendered.contains("XX"));
        assert!(rendered.contains("  "));
    }

    #[test]
    fn exposes_typed_difficulty_defaults() {
        assert_eq!(Difficulty::Easy.default_filled_percentage(), 0.60);
        assert_eq!(Difficulty::Expert.lower(), Ok(Difficulty::Hard));
    }
}
