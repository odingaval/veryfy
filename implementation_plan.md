# Veryfy Integration Implementation Plan

This plan replaces the earlier backend-only plan. The contract and frontend have now landed, and the safest MVP path is to **avoid changing contract code unless it is clearly broken**. We will treat the current Anchor program as the integration source of truth and adapt the backend/frontend around it.

## Current Integration Decision

- The frontend wallet signs issuer transactions directly.
- The contract stores generic license proof data on-chain.
- The backend keeps REST helpers, validation, hashing, and optional Solana read/verification support.
- The backend should not hold issuer private keys for the MVP.

## Fixed Contract Interface

Use the existing contract instructions:

```txt
register_issuer(name)
issue_license(asset_hash, expiry)
revoke_license(asset_hash)
```

Use the existing PDA seeds:

```txt
issuer PDA  = ["issuer", issuerWallet]
license PDA = ["license", assetHashBytes]
```

Use the existing on-chain `License` fields:

```txt
holder: Pubkey
issuer: Pubkey
status: Active | Revoked | Expired
expiry: i64
asset_hash: [u8; 32]
bump: u8
```

Important implication: the contract does not store `licenseNumber`, `holderName`, `licenseType`, or `issuerWallet` as separate license fields. Veryfy gives meaning to `asset_hash` by hashing those off-chain fields deterministically.

## Shared Hash Contract

Every layer must derive the same 32-byte `asset_hash` from:

```txt
licenseNumber|holderName|issuerWallet|licenseType|expiryDate
```

Rules:

- Algorithm: SHA-256.
- `expiryDate`: `YYYY-MM-DD`.
- No trailing pipe.
- No trimming or normalization unless all layers are updated together.
- Raw 32 bytes are passed to the contract as `asset_hash`.
- Hex is preferred for backend/API display.
- Base58 may be used by the frontend for QR/license IDs, but it must decode to the same 32 bytes.

## Response Contract

Backend REST responses use envelopes:

```json
{
  "data": {},
  "error": null
}
```

Frontend REST clients must unwrap `data`. Direct Solana hooks may return plain objects internally, but UI code should not mix envelope and non-envelope assumptions in the same path.

## Implementation Stages

## Stage 1: Freeze Integration Contract

Status: implemented.

Goal: document and lock the cross-layer contract without changing Anchor code.

Tasks:

- Record the current contract instructions, PDA seeds, on-chain account fields, and hash rule.
- Decide frontend wallet signing remains the MVP transaction path.
- Decide backend does not sign issuer transactions for MVP.
- Keep `issue`, `revoke`, and `register` backend endpoints available for local/memory development, but do not treat them as the primary on-chain path yet.

Deliverable:

- Updated `implementation_plan.md` reflects the current cross-layer integration strategy.

## Stage 2: Shared Frontend Hashing And API Shape

Status: implemented.

Goal: remove the frontend/backend hash mismatch and align REST response handling.

Tasks:

- Replace the frontend mock hash with real SHA-256 over the shared pipe string.
- Use the connected wallet public key as `issuerWallet` when issuing.
- Allow verification from either a shared-field payload or an existing license hash.
- Keep raw on-chain `asset_hash` bytes identical regardless of display encoding.
- Update the frontend REST API helper to unwrap backend `{ data, error }` envelopes.

Deliverable:

- Frontend direct Solana issuing/verifying uses the same hash rule as the Go backend.
- Frontend REST helper is compatible with backend response envelopes.

## Stage 3: Backend Solana Read Verification

Status: implemented.

Goal: implement useful Stage 6 behavior without requiring backend transaction signing.

Tasks:

- Implement PDA derivation for license accounts from `asset_hash`.
- Implement Solana account fetch for `License`.
- Decode the Anchor account discriminator and fields.
- Map contract status and expiry to API statuses:
  - missing account -> `INVALID`
  - revoked -> `REVOKED`
  - expired timestamp -> `EXPIRED`
  - active and unexpired -> `VALID`
- Add tests for PDA derivation, account decoding, and status mapping.

Deliverable:

- Backend can verify licenses against Solana by hash/details.

## Stage 4: End-To-End Demo Integration

Goal: make the product path coherent for demo.

Tasks:

- Frontend issue flow creates a license and displays/share the resulting license ID.
- Frontend verify flow can verify by QR/license ID.
- Optional manual verify flow recomputes hash from full details.
- Revoke flow updates on-chain status and the verifier shows `REVOKED`.
- Document localnet/devnet environment variables and demo steps.

Deliverable:

- One clean demo path: issue -> verify -> tamper -> revoke -> verify again.

## Stage 5: Optional Backend Write Integration

Goal: only if needed, add server-side transaction submission.

Tasks:

- Decide which server keypair signs transactions and why.
- Implement `register_issuer`, `issue_license`, and `revoke_license` transaction builders.
- Keep issuer-wallet security explicit; avoid silently using admin keys for issuer actions.

Deliverable:

- Backend write endpoints can submit real Solana transactions, if the project chooses that architecture.
