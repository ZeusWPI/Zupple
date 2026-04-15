#[derive(Debug)]
pub enum TakuzuError {
    NoLowerDifficulty,
    InvalidSize,
    CantSolve,
    NoUniqueSolution
}