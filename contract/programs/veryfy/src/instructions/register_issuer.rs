use anchor_lang::prelude::*;
use crate::{
    state::issuer::Issuer,

};

/// Context for registering a new issuer
#[derive(Accounts)]
#[instruction(name: String)]
pub struct RegisterIssuer<'info> {
    #[account(mut, signer)]
    pub payer: Signer<'info>, // pays rent for the new issuer PDA

    #[account(
        init,
        payer = payer,
        space = 8 + Issuer::MAX_SIZE,
        seeds = [b"issuer", payer.key().as_ref()],
        bump,
    )]
    pub issuer: Account<'info, Issuer>,

    pub system_program: Program<'info, System>,
}

pub fn register_issuer(
    ctx: Context<RegisterIssuer>,
    name: String,
) -> Result<()> {
    let issuer = &mut ctx.accounts.issuer;
    issuer.authority = ctx.accounts.payer.key();
    // Truncate name if too long (max 64 bytes) – Anchor will enforce size later
    issuer.name = name;
    issuer.issued_count = 0;
    issuer.bump = ctx.bumps.issuer;
    Ok(())
}
