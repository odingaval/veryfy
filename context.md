Here's a complete project brief your coding agent can use as context.

---

# Veritas — Project Brief

## What We're Building

Veritas is a decentralized professional license verification system built on Solana. It allows licensing bodies (like KMPDC for doctors, LSK for lawyers, NTSA for drivers) to issue tamper-proof licenses on-chain, and allows anyone — a hospital, a court, a patient — to instantly verify whether a license is genuine, active, and unexpired.

The core problem it solves: professional licenses in Kenya are easy to forge, slow to verify manually, and have no single trusted source of truth. A hospital hiring a doctor today must call KMPDC directly, a process that is slow, gameable, and doesn't scale. Veritas replaces that with cryptographic proof.

---

## The Problem in One Paragraph

Fake professional credentials are a documented problem in Kenya and across East Africa. A forged medical license, a fake advocate certificate, or an expired driving license presented as valid all carry real-world consequences — patient harm, legal fraud, road safety failures. The current verification process relies on phone calls, physical letters, and manual cross-checks with licensing bodies. It is slow, inconsistent, and easy to circumvent. There is no public, tamper-proof registry that any verifier can query in real time.

---

## How It Works — The Three Actors

**Issuer** — A registered licensing body (KMPDC, LSK, NTSA). They connect their wallet, fill in a professional's license details, and issue the credential on-chain. Each license is stored as a hashed record on Solana, signed by the issuer's wallet.

**Holder** — The licensed professional (doctor, lawyer, driver). They receive a license record they can share as a license number or QR code. They do not need to interact with the app directly for MVP.

**Verifier** — Anyone who needs to confirm a license is real: a hospital HR manager, a law firm, a traffic officer, a patient. They enter the license details into the verifier dashboard and get an instant result — no account needed, no wallet needed.

---

## The Verification Mechanism

Verification works through deterministic hashing. When a license is issued, the system generates a SHA-256 hash of the license's core fields:

```
SHA256("licenseNumber|holderName|issuerWallet|licenseType|expiryDate")
```

This hash is stored on-chain. When a verifier submits the same details, the system recomputes the hash and checks whether it exists on-chain. If someone tampers with any field — changes a name, alters an expiry date, swaps the license number — the hash breaks and verification fails. The blockchain never stores personal data, only the hash.

---

## License Status Values

Every verification returns one of four statuses:

- **VALID** — hash found on-chain, not revoked, not expired
- **INVALID** — hash not found (fake or tampered credential)
- **REVOKED** — hash found but issuer has explicitly revoked it
- **EXPIRED** — hash found, not revoked, but expiry date has passed

---

## License Types

The system is designed to handle multiple professional license types from day one:

- `MEDICAL` — issued by KMPDC (doctors, dentists, clinical officers)
- `LEGAL` — issued by LSK (advocates)
- `DRIVING` — issued by NTSA

Adding a new license type requires no contract changes — it is just a string field on the license record.

---

## What We Are Not Building

To stay focused for the hackathon, the following are explicitly out of scope:

- A holder-facing mobile app or wallet
- Payment or subscription features
- Integration with any real government API
- National ID or biometric verification
- Multi-signature or DAO-based governance
- Mainnet deployment (devnet only)

---

## Tech Stack

| Layer | Technology | Notes |
|---|---|---|
| Blockchain | Solana (devnet) | Fast, cheap, good for demo |
| Smart contract | Anchor (Rust) | Standard Solana framework |
| Backend | Go | REST API, hashing, Solana RPC |
| Frontend | React + TypeScript | Vite, two dashboards |
| Wallet | Phantom | Issuer auth only |
| Hashing | SHA-256 | stdlib, deterministic |

---

## On-Chain Data Model

Two account types live on-chain:

**Issuer account** — PDA derived from `[b"issuer", wallet.key()]`
```
name: String           // "KMPDC", "LSK"
license_type: String   // "MEDICAL", "LEGAL"
wallet: Pubkey
registered_at: i64
```

**License account** — PDA derived from `[b"license", license_hash]`
```
holder_wallet: Pubkey
license_hash: [u8; 32]
issuer: Pubkey
license_type: String
license_number: String
expiry_date: i64
is_revoked: bool
issued_at: i64
```

---

## API Surface (Go backend)

```
POST   /issuers/register     Register a new licensing body
POST   /licenses/issue       Issue a new license (issuer only)
POST   /licenses/verify      Verify a license by details (public)
POST   /licenses/revoke      Revoke an existing license (issuer only)
GET    /health               Health check
```

---

## Frontend Pages

**`/verify`** — Public verifier dashboard. No wallet required. A verifier enters license details and gets back a status card. This is the primary demo-facing page.

**`/issue`** — Issuer dashboard. Phantom wallet required. Issuer selects license type, fills in holder details, issues on-chain. Also shows a table of previously issued licenses with revoke buttons.

**`/`** — Landing page or redirect to `/verify`.

---

## Critical Integration Contract

The single most important thing both backend and frontend must agree on before writing code:

**The hash input string format:**
```
{licenseNumber}|{holderName}|{issuerWallet}|{licenseType}|{expiryDate}
```

**Example:**
```
SHA256("KE/MED/12345|John Kamau|7xKp...wallet|MEDICAL|2026-12-31")
```

`expiryDate` is always `YYYY-MM-DD`. Fields joined with pipe `|`. No trailing pipe. This must be identical in the Go backend and the Anchor contract. Any divergence here breaks verification entirely.

---

## Demo Script (for judges)

The demo tells a story in under 2 minutes:

1. KMPDC wallet connects → issues Dr. John Kamau's medical license on-chain
2. Hospital verifier searches license number → **VALID** ✅
3. Same details but name changed to "Jon Kamau" → **INVALID** ❌ (hash breaks)
4. KMPDC revokes a license → same query returns **REVOKED** ❌
5. Show a license with a past expiry date → **EXPIRED** ❌

Five outcomes, one clear story: you cannot fake, tamper, or reuse a Veritas license.

---

## Suggested Build Order

1. Deploy hello-world Anchor program to devnet → confirms toolchain works
2. Define and deploy full contract with all three instructions
3. Go backend: hash service first (unit test it), then Solana client, then handlers
4. Frontend: types and API layer first, then verifier dashboard, then issuer dashboard
5. End-to-end test: issue via frontend → verify via frontend
6. Seed devnet with demo data (KMPDC, LSK issuers + sample licenses)
7. Polish demo flow, prepare pitch

---

That's everything the agent needs to have full context. Hand it the folder structure from earlier alongside this and it can work on any layer independently without asking for background.