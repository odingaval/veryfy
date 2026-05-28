Let me read the frontend skill first before laying this out.Here's the full folder structure and task breakdown, written so a coding agent can work on each area independently.

---

## Folder Structure

```
veritas/
├── contract/                          # Solana / Anchor program
│   ├── programs/
│   │   └── veritas/
│   │       └── src/
│   │           ├── lib.rs             # Program entry point
│   │           ├── instructions/
│   │           │   ├── mod.rs
│   │           │   ├── register_issuer.rs
│   │           │   ├── issue_license.rs
│   │           │   └── revoke_license.rs
│   │           ├── state/
│   │           │   ├── mod.rs
│   │           │   ├── issuer.rs      # Issuer account struct
│   │           │   └── license.rs     # License account struct
│   │           └── errors.rs
│   ├── tests/
│   │   └── veritas.ts                 # Anchor integration tests
│   ├── Anchor.toml
│   └── Cargo.toml
│
├── api/                               # Go backend
│   ├── cmd/
│   │   └── server/
│   │       └── main.go                # Entry point, wires everything
│   ├── internal/
│   │   ├── handlers/
│   │   │   ├── issuer.go              # POST /issuers/register
│   │   │   ├── license.go             # POST /licenses/issue, /verify, /revoke
│   │   │   └── health.go              # GET /health
│   │   ├── services/
│   │   │   ├── hash.go                # SHA-256 hashing logic
│   │   │   └── solana.go              # Solana RPC client wrapper
│   │   ├── models/
│   │   │   └── license.go             # Shared request/response structs
│   │   └── middleware/
│   │       └── cors.go
│   ├── go.mod
│   └── go.sum
│
└── web/                               # React frontend
    ├── src/
    │   ├── main.tsx
    │   ├── App.tsx                    # Routes
    │   ├── components/
    │   │   ├── ui/                    # Reusable: Button, Card, Badge, Input
    │   │   ├── LicenseCard.tsx        # Displays a single license result
    │   │   ├── StatusBadge.tsx        # VALID / INVALID / REVOKED / EXPIRED
    │   │   └── WalletConnect.tsx      # Phantom wallet button
    │   ├── pages/
    │   │   ├── IssuerDashboard.tsx    # Issue + revoke licenses
    │   │   ├── VerifierDashboard.tsx  # Look up licenses
    │   │   └── Landing.tsx            # Optional: explain what the app does
    │   ├── hooks/
    │   │   ├── useWallet.ts           # Phantom wallet connection
    │   │   └── useLicense.ts          # API calls: issue, verify, revoke
    │   ├── lib/
    │   │   └── api.ts                 # All fetch calls to Go API, typed
    │   ├── types/
    │   │   └── index.ts               # Shared TypeScript types
    │   └── styles/
    │       └── globals.css
    ├── index.html
    ├── package.json
    └── vite.config.ts
```

---

## Contract Tasks

**1. Define state structs** (`state/issuer.rs`, `state/license.rs`)
- `Issuer`: `wallet: Pubkey`, `name: String`, `license_type: String`, `registered_at: i64`
- `License`: `holder_wallet: Pubkey`, `license_hash: [u8; 32]`, `issuer: Pubkey`, `license_type: String`, `license_number: String`, `expiry_date: i64`, `is_revoked: bool`, `issued_at: i64`

**2. Write `register_issuer` instruction**
- Only callable by a known admin keypair (hardcode for hackathon)
- Creates an `Issuer` PDA derived from `[b"issuer", wallet.key()]`
- Stores name and license type

**3. Write `issue_license` instruction**
- Caller must be a registered issuer (check PDA exists)
- Creates a `License` PDA derived from `[b"license", license_hash]`
- Stores all license fields

**4. Write `revoke_license` instruction**
- Caller must be the original issuer
- Sets `is_revoked = true` on the License PDA

**5. Read instruction (no transaction needed)**
- Fetch License PDA by hash via RPC — this is a read, not a contract instruction
- Expiry is computed client-side: compare `expiry_date` vs current timestamp

**6. Write Anchor tests** (`tests/veritas.ts`)
- Test: register issuer → issue license → verify → revoke → verify again
- Deploy to devnet, save program ID to `Anchor.toml`

---

## Backend Tasks (Go)

**1. Project setup**
- Init Go module: `github.com/yourteam/veritas-api`
- Install `github.com/gagliardetto/solana-go` for Solana RPC
- Install `github.com/go-chi/chi` for routing
- Set up `.env`: `SOLANA_RPC_URL`, `PROGRAM_ID`, `ADMIN_KEYPAIR_PATH`

**2. Hash service** (`internal/services/hash.go`)
- Function: `GenerateLicenseHash(licenseNumber, holderName, issuer, expiryDate string) [32]byte`
- Uses `crypto/sha256` from stdlib
- Must be deterministic — same inputs always same output
- Write unit tests for this, it's critical

**3. Solana service** (`internal/services/solana.go`)
- `IssueLicense(ctx, params)` — builds + sends transaction to Anchor program
- `RevokeLicense(ctx, licenseHash)` — sends revoke transaction
- `FetchLicense(ctx, licenseHash) (*License, error)` — reads PDA account data
- `RegisterIssuer(ctx, params)` — registers a new issuing body

**4. Handlers** (`internal/handlers/`)

`POST /licenses/issue`
```
body: { licenseNumber, holderName, holderWallet, licenseType, expiryDate, issuerWallet }
→ hash the details
→ call solana.IssueLicense
→ return { licenseHash, txSignature }
```

`POST /licenses/verify`
```
body: { licenseNumber, holderName, issuerWallet, licenseType, expiryDate }
→ recompute hash from same inputs
→ call solana.FetchLicense
→ check is_revoked, check expiry
→ return { status: "VALID" | "INVALID" | "REVOKED" | "EXPIRED", details }
```

`POST /licenses/revoke`
```
body: { licenseHash, issuerWallet }
→ call solana.RevokeLicense
→ return { txSignature }
```

`POST /issuers/register`
```
body: { name, licenseType, wallet }
→ call solana.RegisterIssuer
→ return { txSignature }
```

**5. CORS middleware** (`internal/middleware/cors.go`)
- Allow `localhost:5173` (Vite dev server) during development

**6. Wire everything in `main.go`**
- Load env, init router, mount handlers, start server on `:8080`

---

## Frontend Tasks (React)

**1. Project setup**
- `npm create vite@latest web -- --template react-ts`
- Install: `@solana/web3.js`, `@solana/wallet-adapter-react`, `@solana/wallet-adapter-phantom`
- Install: `react-router-dom`, `react-hot-toast` for notifications

**2. Types** (`src/types/index.ts`)
```ts
type LicenseType = "MEDICAL" | "LEGAL" | "DRIVING"

type LicenseStatus = "VALID" | "INVALID" | "REVOKED" | "EXPIRED"

interface License {
  licenseNumber: string
  holderName: string
  holderWallet: string
  licenseType: LicenseType
  expiryDate: string
  issuerName: string
  status: LicenseStatus
}
```

**3. API layer** (`src/lib/api.ts`)
- All fetch calls live here, nothing else imports from `fetch` directly
- `issueLicense(params)`, `verifyLicense(params)`, `revokeLicense(hash)`
- Base URL from `import.meta.env.VITE_API_URL`

**4. Hooks** (`src/hooks/`)
- `useWallet.ts` — wraps Phantom adapter, exposes `connect()`, `publicKey`, `connected`
- `useLicense.ts` — wraps API calls with loading/error state

**5. Reusable components** (`src/components/`)
- `StatusBadge.tsx` — colour-coded pill: green/red/orange/grey per status
- `LicenseCard.tsx` — displays full license details + status after verification
- `WalletConnect.tsx` — Phantom connect button, shows truncated address when connected

**6. Issuer Dashboard** (`src/pages/IssuerDashboard.tsx`)
- Wallet must be connected to use this page
- Form: license type dropdown (MEDICAL / LEGAL / DRIVING), holder name, license number, holder wallet, expiry date
- On submit: calls `issueLicense`, shows success toast + transaction link to Solana explorer
- Table below: lists issued licenses with a Revoke button per row

**7. Verifier Dashboard** (`src/pages/VerifierDashboard.tsx`)
- No wallet needed — this is public
- Form: license number + holder name + issuer + license type + expiry date
- On submit: calls `verifyLicense`, renders `LicenseCard` with result
- This is your demo-facing page — keep it clean and fast

**8. Routing** (`src/App.tsx`)
```
/           → Landing or redirect to /verify
/verify     → VerifierDashboard
/issue      → IssuerDashboard (wallet-gated)
```

---

## Shared contract between frontend and backend

Make sure both sides agree on this before anyone writes a line:

```
Verify inputs: licenseNumber + holderName + issuerWallet + licenseType + expiryDate
Hash function: SHA-256 of above fields, joined with pipes
e.g. SHA256("KE/MED/12345|John Kamau|KMPDC_WALLET|MEDICAL|2026-12-31")
```

This is the most likely integration bug. Nail it early.

---

What do you want to go deeper on first — the Anchor contract, the Go backend, or the React frontend?