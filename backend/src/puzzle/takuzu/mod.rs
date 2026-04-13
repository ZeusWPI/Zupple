mod error;
mod generate;

use std::fmt::{Display, Formatter, Write};
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

#[derive(Clone, Copy, PartialEq, Eq, Debug)]
pub enum TakuzuCell {
    O, X, Empty
}
impl Display for TakuzuCell {
    fn fmt(&self, f: &mut Formatter<'_>) -> std::fmt::Result {
        match self {
            TakuzuCell::O => f.write_str("O"),
            TakuzuCell::X => f.write_str("X"),
            TakuzuCell::Empty => f.write_str("_")
        }
    }
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

impl Display for Takuzu {
    fn fmt(&self, f: &mut Formatter<'_>) -> std::fmt::Result {
        f.write_str("Takuzu: ")?;
        f.write_fmt(format_args!("size {}", self.size))?;


        let mut divider = "\n ".to_owned();

        for _ in 0..self.size {
            divider.push_str("----");
        }
        divider.push_str("-\n");

        for row in 0..self.size {
            f.write_str(&divider)?;
            for col in 0..self.size {
                let value = self.grid[row as usize][col as usize];
                f.write_fmt(format_args!(" | {}", value))?;
            }
            f.write_str(" |")?;
        }
        f.write_str(&divider)?;
        Ok(())
    }
}