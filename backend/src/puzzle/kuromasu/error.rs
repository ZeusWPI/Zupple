use thiserror::Error;

#[derive(Debug, Clone, PartialEq, Error)]
pub enum KuromasuError {
    #[error("size must be positive")]
    InvalidSize,

    #[error("size can be at most {0}")]
    SizeTooLarge(usize),

    #[error("field length does not match the board size")]
    InvalidFieldLength,

    #[error("filled percentage must be between 0 and 1")]
    InvalidFilledPercentage,

    #[error("invalid mask length for board size")]
    InvalidMaskLength,

    #[error("mask values must be 0 or 1")]
    InvalidMaskValue,

    #[error(
        "unable to generate a mask that matches the requested difficulty and filled percentage"
    )]
    UnreachableMask,

    #[error("difficulty has no lower difficulty")]
    NoLowerDifficulty,
}
