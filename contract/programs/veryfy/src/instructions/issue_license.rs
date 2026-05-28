use anchor_lang::prelude::*;
use crate::{
    state::{license::{License, LicenseStatus}, issuer::Issuer},
    errors::VeryfyError,
};

/// Context for the `issue_license` instruction
#[derive(Accounts)]
#[instruction(asset_hash: [u8; 32], expiry: i64)]
pub struct IssueLicense<'info> {
    #[account(mut, signer)]
    pub payer: Signer<'info>, // pays rent for the new accounts

    #[account(
        init,
        payer = payer,
        space = 8 + License::MAX_SIZE,
        seeds = [b"license", &asset_hash],
        bump,
    )]
    pub license: Account<'info, License>,

    #[account(
        mut,
        seeds = [b"issuer", issuer.authority.as_ref()],
        bump = issuer.bump,
        has_one = authority @ VeryfyError::UnauthorizedIssuer,
    )]
    pub issuer: Account<'info, Issuer>,

    /// CHECK: PDA authority for the issuer; not read directly
    #[account(seeds = [b"issuer", issuer.authority.as_ref()], bump = issuer.bump)]
    pub authority: UncheckedAccount<'info>,

    pub system_program: Program<'info, System>,
    pub rent: Sysvar<'info, Rent>,
}

pub fn issue_license(
    ctx: Context<IssueLicense>,
    asset_hash: [u8; 32],
    expiry: i64,
) -> Result<()> {
    // Populate License account
    let license = &mut ctx.accounts.license;
    license.holder = ctx.accounts.payer.key();
    license.issuer = ctx.accounts.issuer.key();
    license.status = LicenseStatus::Active;
    license.expiry = expiry;
    license.asset_hash = asset_hash;
    license.bump = *ctx.bumps.get("license").unwrap();

    // Update issuer analytics safely
    ctx.accounts.issuer.issued_count = ctx
        .accounts
        .issuer
        .issued_count
        .checked_add(1)
        .ok_or(VeryfyError::MathOverflow)?;

    Ok(())
}
