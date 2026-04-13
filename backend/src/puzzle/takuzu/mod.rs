mod error;

use crate::puzzle::takuzu::error::TakuzuError;
use crate::puzzle::takuzu::TakuzuDifficulty::{Easy, Hard, Medium, Extreme};

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