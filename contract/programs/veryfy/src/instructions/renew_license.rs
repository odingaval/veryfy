// src/instructions/renew_license.rs
use anchor_lang::prelude::*;
use crate::{
    state::license::{License, LicenseStatus},
    state::issuer::Issuer,
    errors::VeryfyError,
    events::LicenseRenewed,
};

/// Context for renewing a license's expiry
#[derive(Accounts)]
#[instruction(asset_hash: [u8; 32], new_expiry: i64)]
pub struct RenewLicense<'info> {
    /// The authority that originally issued the license (signer)
    #[account(mut, signer)]
    pub authority: Signer<'info>,

    #[account(
        mut,
        seeds = [b"license", &asset_hash],
        bump = license.bump,
        has_one = issuer @ VeryfyError::UnauthorizedIssuer,
    )]
    pub license: Account<'info, License>,

    #[account(
        seeds = [b"issuer", authority.key().as_ref()],
        bump = issuer.bump,
    )]
    pub issuer: Account<'info, Issuer>,

    /// CHECK: PDA authority for the issuer – not read directly
    #[account(seeds = [b"issuer", authority.key().as_ref()], bump = issuer.bump)]
    pub authority_account: UncheckedAccount<'info>,
}

pub fn renew_license(
    ctx: Context<RenewLicense>,
    asset_hash: [u8; 32],
    new_expiry: i64,
) -> Result<()> {
    // Ensure license is active
    if ctx.accounts.license.status != LicenseStatus::Active {
        return Err(VeryfyError::LicenseAlreadyRevoked);
    }
    // Update expiry safely (prevent overflow)
    ctx.accounts.license.expiry = new_expiry;
    // Emit event
    emit!(LicenseRenewed {
        license: ctx.accounts.license.key(),
        issuer: ctx.accounts.issuer.key(),
        new_expiry,
    });
    Ok(())
}
