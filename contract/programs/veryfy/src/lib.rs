use anchor_lang::prelude::*;

mod errors;
mod state;
mod instructions;

use crate::errors::VeryfyError;
use crate::instructions::issue_license::{IssueLicense, issue_license as issue_license_handler};

declare_id!("YourProgramIdHere");

#[program]
pub mod veryfy {
    use super::*;

    /// Issue a new license NFT on‑chain.
    ///
    /// * `asset_hash` – 32‑byte hash identifying the off‑chain asset.
    /// * `expiry` – Unix timestamp; 0 means never expires.
    pub fn issue_license(
        ctx: Context<IssueLicense>,
        asset_hash: [u8; 32],
        expiry: i64,
    ) -> Result<()> {
        issue_license_handler(ctx, asset_hash, expiry)
    }

    // Future instructions (revoke_license, register_issuer) will be added here.
}
