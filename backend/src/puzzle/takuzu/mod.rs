mod error;
mod generate;

use crate::puzzle::takuzu::error::TakuzuError;
use crate::puzzle::takuzu::error::TakuzuError::InvalidSize;
use crate::puzzle::takuzu::generate::generate_grid;
use crate::puzzle::takuzu::TakuzuDifficulty::{Easy, Hard, Medium, Extreme};

#[derive(Debug)]
pub enum TakuzuDifficulty {
    Easy, Medium, Hard, Extreme
}
impl TakuzuDifficulty {
    pub fn lower(&self) -> Result<TakuzuDifficulty, TakuzuError> {
        match self {
            Medium => Ok(Easy),
            Hard => Ok(Medium),
            Extreme => Ok(Hard),
            _ => Err(TakuzuError::NoLowerDifficulty)
        }
    }
}

#[derive(Clone, Copy, PartialEq, Eq)]
#[derive(Debug)]
pub enum TakuzuCell {
    O, X, Empty
}

#[derive(Debug)]
struct TakuzuPuzzle {
    difficulty: TakuzuDifficulty,
    mask: Vec<Vec<TakuzuCell>>
}
#[derive(Debug)]
pub struct Takuzu {
    puzzles: Vec<TakuzuPuzzle>,
    grid: Vec<Vec<TakuzuCell>>,
    size: u8
}
impl Takuzu {
    pub fn new(size: u8) -> Result<Takuzu, TakuzuError> {
        if size % 2 == 1 || size > 16 {
            return Err(InvalidSize);
        }
        Ok(Takuzu {
            puzzles: vec![],
            size,
            grid: generate_grid(size)
        })
    }

    pub fn generate_puzzle(&mut self, difficulty: TakuzuDifficulty) {
        todo!()
    }
}

