#![allow(ambiguous_glob_reexports)]
use anchor_lang::prelude::*;

pub mod errors;
pub mod state;
pub mod events;
pub mod instructions;

pub use instructions::*;


use crate::instructions::{issue_license as issue_license_handler, revoke_license as revoke_license_handler, register_issuer as register_issuer_handler};

declare_id!("2ePXnx39sHSU5NKpTcRwxhGnND64kHnLZQsy3iavMyPf");

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

    /// Revoke an existing license. Only the issuing authority can call this.
    pub fn revoke_license(
        ctx: Context<RevokeLicense>,
        asset_hash: [u8; 32],
    ) -> Result<()> {
        revoke_license_handler(ctx, asset_hash)
    }

    /// Register a new issuer on‑chain.
    pub fn register_issuer(
        ctx: Context<RegisterIssuer>,
        name: String,
    ) -> Result<()> {
        register_issuer_handler(ctx, name)
    }
}
