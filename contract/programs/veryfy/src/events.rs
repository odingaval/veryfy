use anchor_lang::prelude::*;

#[event]
pub struct LicenseIssued {
    pub license: Pubkey,
    pub issuer: Pubkey,
    pub holder: Pubkey,
    pub asset_hash: [u8; 32],
    pub expiry: i64,
}
