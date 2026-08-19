//! Minimal CRE keystone receiver: stores the latest {key, value} pair written by a
//! CRE workflow. Mirrors the receiver role played by EVM `MyWriterTrigger.sol` in other
//! building blocks, but for Solana: it accepts a keystone-forwarder CPI, decodes a Borsh
//! `KvEntry` payload, and persists it in a `KvStore` account.
//!
//! This program does not implement DON/report verification itself — that is the
//! keystone forwarder's job (Chainlink already deploys and operates the forwarder).
//! `on_report` only checks that the CPI came from the forwarder program recorded at
//! `initialize`, via the `forwarder_authority` PDA the forwarder signs with.

#![allow(deprecated)] // anchor-lang 0.31 #[program] uses AccountInfo::realloc

use anchor_lang::prelude::*;

declare_id!("aPkfGwVg9Lj3yF2CSgxeRmMn7uU9tKWT6PSHTDf2ZX6");

#[program]
pub mod kv_store_receiver {
    use super::*;

    /// One-time setup: creates the `kv_store` account and records the keystone forwarder
    /// program id that is allowed to CPI into `on_report`.
    ///
    /// `forwarder_program` must be the on-chain keystone forwarder program id you were
    /// given for your target CRE environment (Chainlink deploys and operates this
    /// program — you do not deploy it yourself).
    pub fn initialize(ctx: Context<Initialize>, forwarder_program: Pubkey) -> Result<()> {
        require!(
            forwarder_program != Pubkey::default(),
            KvStoreError::InvalidForwarderProgram
        );
        let kv_store = &mut ctx.accounts.kv_store;
        kv_store.forwarder_program = forwarder_program;
        kv_store.key = String::new();
        kv_store.value = String::new();
        kv_store.updated_at = 0;
        kv_store.update_count = 0;
        Ok(())
    }

    /// Called by the keystone forwarder's CPI. Decodes the Borsh-encoded `KvEntry`
    /// report and stores it.
    pub fn on_report(ctx: Context<OnReport>, _metadata: Vec<u8>, report: Vec<u8>) -> Result<()> {
        verify_forwarder_cpi(
            &ctx.accounts.state,
            &ctx.accounts.forwarder_authority,
            &ctx.accounts.kv_store,
        )?;

        let mut reader: &[u8] = report.as_slice();
        let entry = KvEntry::deserialize(&mut reader)
            .map_err(|_| error!(KvStoreError::InvalidReportPayload))?;

        let kv_store = &mut ctx.accounts.kv_store;
        kv_store.key = entry.key;
        kv_store.value = entry.value;
        kv_store.updated_at = Clock::get()?.unix_timestamp;
        kv_store.update_count = kv_store.update_count.saturating_add(1);

        emit!(KvUpdated {
            key: kv_store.key.clone(),
            value: kv_store.value.clone(),
            updated_at: kv_store.updated_at,
            update_count: kv_store.update_count,
        });

        msg!(
            "kv_store_receiver on_report key={} value={} update_count={}",
            kv_store.key,
            kv_store.value,
            kv_store.update_count
        );

        Ok(())
    }
}

/// Borsh layout of the report payload written by the CRE workflow (see
/// `my-workflow/main.ts`).
#[derive(AnchorSerialize, AnchorDeserialize, Clone)]
pub struct KvEntry {
    pub key: String,
    pub value: String,
}

fn verify_forwarder_cpi(
    state: &UncheckedAccount,
    forwarder_authority: &Signer,
    kv_store: &Account<KvStore>,
) -> Result<()> {
    let forwarder_program = kv_store.forwarder_program;
    require!(
        forwarder_program != Pubkey::default(),
        KvStoreError::ForwarderProgramNotConfigured
    );

    require_keys_eq!(
        *state.to_account_info().owner,
        forwarder_program,
        KvStoreError::MismatchedForwarderProgram
    );

    let state_key = state.key();
    let seeds: &[&[u8]] = &[b"forwarder", state_key.as_ref(), crate::ID.as_ref()];
    let (expected_authority, _bump) = Pubkey::find_program_address(seeds, &forwarder_program);

    require_keys_eq!(
        expected_authority,
        forwarder_authority.key(),
        KvStoreError::InvalidForwarderAuthority
    );

    Ok(())
}

#[error_code]
pub enum KvStoreError {
    #[msg("forwarder_program must be a non-default pubkey")]
    InvalidForwarderProgram,
    #[msg("Call initialize with the keystone forwarder program id before on_report")]
    ForwarderProgramNotConfigured,
    #[msg("Forwarder state account owner does not match the forwarder_program stored at initialize")]
    MismatchedForwarderProgram,
    #[msg("forwarder_authority is not the PDA for this state, receiver, and forwarder_program")]
    InvalidForwarderAuthority,
    #[msg("Report payload must be Borsh KvEntry (key: string, value: string)")]
    InvalidReportPayload,
}

#[event]
pub struct KvUpdated {
    pub key: String,
    pub value: String,
    pub updated_at: i64,
    pub update_count: u64,
}

#[account]
#[derive(InitSpace)]
pub struct KvStore {
    /// Keystone forwarder program id allowed to call `on_report` (set at `initialize`).
    pub forwarder_program: Pubkey,
    #[max_len(64)]
    pub key: String,
    #[max_len(200)]
    pub value: String,
    pub updated_at: i64,
    pub update_count: u64,
}

#[derive(Accounts)]
pub struct Initialize<'info> {
    #[account(init, payer = signer, space = 8 + KvStore::INIT_SPACE)]
    pub kv_store: Account<'info, KvStore>,

    #[account(mut)]
    pub signer: Signer<'info>,

    pub system_program: Program<'info, System>,
}

#[derive(Accounts)]
pub struct OnReport<'info> {
    /// CHECK: forwarder state account; owner must match `kv_store.forwarder_program`
    /// (checked in `verify_forwarder_cpi`).
    pub state: UncheckedAccount<'info>,

    /// PDA signer supplied by the forwarder CPI; verified in `verify_forwarder_cpi`.
    pub forwarder_authority: Signer<'info>,

    #[account(mut)]
    pub kv_store: Account<'info, KvStore>,
}
