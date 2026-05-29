use anchor_lang::prelude::*;

#[error_code]
pub enum VeryfyError {
    #[msg("Unauthorized issuer")] // Thrown when the signer is not the registered issuer
    UnauthorizedIssuer,

    #[msg("Math overflow")] // Used for safe arithmetic on counters
    MathOverflow,

    #[msg("Feature not yet implemented")] // Placeholder for stub instructions
    NotImplemented,
}
