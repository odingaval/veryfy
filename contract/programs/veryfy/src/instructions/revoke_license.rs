use anchor_lang::prelude::*;
use crate::{
    state::{license::{License, LicenseStatus}, issuer::Issuer},
    errors::VeryfyError,
};

/// Context for revoking a license
#[derive(Accounts)]
#[instruction(asset_hash: [u8; 32])]
pub struct RevokeLicense<'info> {
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
        mut,
        seeds = [b"issuer", issuer.authority.as_ref()],
        bump = issuer.bump,
        has_one = authority @ VeryfyError::UnauthorizedIssuer,
    )]
    pub issuer: Account<'info, Issuer>,

    pub system_program: Program<'info, System>,
}

pub fn revoke_license(
    ctx: Context<RevokeLicense>,
    _asset_hash: [u8; 32],
) -> Result<()> {
    let license = &mut ctx.accounts.license;
    // Ensure the license was issued by this issuer (has_one constraint already checks)
    license.status = LicenseStatus::Revoked;
    Ok(())
}
