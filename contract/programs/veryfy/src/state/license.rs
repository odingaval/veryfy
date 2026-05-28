use anchor_lang::prelude::*;

#[derive(AnchorSerialize, AnchorDeserialize, Clone, PartialEq, Eq)]
pub enum LicenseStatus {
    Active,
    Revoked,
    Expired,
}

/// On‑chain License account
#[account]
pub struct License {
    /// The wallet that holds the license NFT
    pub holder: Pubkey,
    /// PDA of the issuer that created this license
    pub issuer: Pubkey,
    /// Current lifecycle status
    pub status: LicenseStatus,
    /// Unix timestamp of expiry (0 = never expires)
    pub expiry: i64,
    /// Hash of the off‑chain asset (e.g., file, video)
    pub asset_hash: [u8; 32],
    /// PDA bump seed
    pub bump: u8,
}

impl License {
    /// Approximate size of the account for rent exemption (excluding 8‑byte discriminator)
    pub const MAX_SIZE: usize = 32   // holder Pubkey
        + 32                          // issuer Pubkey
        + 1                           // enum discriminant
        + 8                           // expiry i64
        + 32                          // asset_hash
        + 1;                          // bump
}
