use anchor_lang::prelude::*;

#[account]
pub struct Issuer {
    /// Authority that controls this issuer (owner wallet)
    pub authority: Pubkey,
    /// Optional human‑readable name (max ~64 bytes)
    pub name: String,
    /// Counter of how many licenses this issuer has created
    pub issued_count: u64,
    /// PDA bump seed
    pub bump: u8,
}

impl Issuer {
    /// Approximate size of the account for rent exemption (excluding 8‑byte discriminator)
    pub const MAX_SIZE: usize = 32   // authority Pubkey
        + 4 + 64                     // name string (length prefix + up to 64 bytes)
        + 8                          // issued_count u64
        + 1; // bump
}
