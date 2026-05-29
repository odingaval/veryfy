# Veryfy Current State

Last updated after backend Stage 3.

Current branch:

```txt
feat/backend
```

Recent backend commits:

```txt
6edb5fc feat(api): add license models and hashing
7e6ef1f feat(api): add backend foundation
```

Note: `implementation_plan.md` exists locally as an untracked planning document and has not been committed.

## Project Shape

The repo currently has scaffolded folders for:

- `contract/` - Solana / Anchor program
- `api/` - Go backend
- `web/` - React frontend

Only the backend has real implementation work so far. The contract and frontend are still folder/file scaffolds.

## Backend Module

The Go backend lives in `api/`.

Module path:

```txt
github.com/odingaval/veryfy/api
```

Run backend tests with:

```bash
cd api
GOCACHE=/tmp/veryfy-go-cache go test ./...
```

The explicit `GOCACHE` is useful in the Codex sandbox because the default Go build cache path may be read-only.

## Implemented Backend Pieces

### Server Foundation

Implemented:

- `api/cmd/server/main.go`
- Standard library HTTP server using `http.NewServeMux`.
- `GET /health`.
- Graceful shutdown on interrupt/SIGTERM.
- Config loading.
- Middleware chain:
  - CORS
  - request logging
  - panic recovery
  - request timeout
- JSON response envelope helpers.

The live server could not be smoke-tested in the sandbox because binding to a local port was blocked, but handler and package tests pass.

### Config

Implemented in:

```txt
api/internal/config/config.go
```

Config currently reads:

```txt
PORT
SOLANA_RPC_URL
PROGRAM_ID
ADMIN_KEYPAIR_PATH
SOLANA_CLUSTER
ALLOWED_ORIGINS
REQUEST_TIMEOUT_SECONDS
```

Defaults:

- `PORT`: `8080`
- `ALLOWED_ORIGINS`: `http://localhost:5173`
- `REQUEST_TIMEOUT_SECONDS`: 15 seconds

### Response Envelope

Implemented in:

```txt
api/internal/httpjson/response.go
```

Success responses should use:

```json
{
  "data": {},
  "error": null
}
```

Error responses should use:

```json
{
  "data": null,
  "error": {
    "code": "INVALID_REQUEST",
    "message": "expiryDate is required"
  }
}
```

Future handlers should use the helper functions here instead of writing JSON manually.

### Models And Validation

Implemented in:

```txt
api/internal/models/license.go
api/internal/models/requests.go
api/internal/models/validation.go
```

Current request DTOs:

- `IssueLicenseRequest`
- `VerifyLicenseRequest`
- `RevokeLicenseRequest`
- `RegisterIssuerRequest`

Current response DTOs:

- `IssueLicenseResponse`
- `VerifyLicenseResponse`
- `RevokeLicenseResponse`
- `RegisterIssuerResponse`

Current verification statuses:

```txt
VALID
INVALID
REVOKED
EXPIRED
```

Important license type decision:

`licenseType` is intentionally not a hardcoded enum. It is an issuer-defined string slug so the system can support new license categories without code changes.

Current `licenseType` validation:

- 2-64 characters
- uppercase letters allowed
- numbers allowed
- underscores allowed
- hyphens allowed
- spaces and lowercase letters rejected

Examples accepted:

```txt
MEDICAL
LEGAL
DRIVING
PHARMACY
ENGINEERING
NURSING_COUNCIL
PILOT-LICENSE
```

Wallet validation currently checks that wallet fields look like 32-byte Solana public keys encoded as base58. This is lightweight validation and does not yet call Solana RPC.

License hashes are represented in the API as 64-character hex strings.

### Hashing

Implemented in:

```txt
api/internal/services/hash.go
```

The hash input contract is critical and must not drift:

```txt
licenseNumber|holderName|issuerWallet|licenseType|expiryDate
```

Example:

```txt
KE/MED/12345|John Kamau|11111111111111111111111111111111|MEDICAL|2026-12-31
```

Hash algorithm:

```txt
SHA-256
```

The hash service returns both:

- raw `[32]byte`
- lowercase hex string

No trimming or normalization is currently applied inside the hash service. Inputs should be validated before hashing. If normalization is added later, frontend, backend, and contract must all agree on it.

### Repository Boundary

Implemented in:

```txt
api/internal/repositories/interfaces.go
api/internal/repositories/memory.go
```

Current repository interfaces:

- `LicenseRepository`
- `IssuerRepository`

Current repository errors:

- `ErrLicenseNotFound`
- `ErrIssuerNotFound`
- `ErrDuplicateLicense`
- `ErrDuplicateIssuer`
- `ErrUnauthorizedIssuer`

`MemoryRepository` is available for tests and local development before the real Solana repository exists. It currently enforces:

- issuer must be registered before issuing
- issuer can only issue its registered `licenseType`
- duplicate license hashes are rejected
- only the original issuer can revoke a license

### Service Layer

Implemented in:

```txt
api/internal/services/license.go
api/internal/services/issuer.go
```

`LicenseService` currently supports:

- issue license
- verify license
- revoke license

`IssuerService` currently supports:

- register issuer

Verification behavior:

- license hash not found -> `INVALID`
- license found and revoked -> `REVOKED`
- license found and expired -> `EXPIRED`
- license found, active, and unexpired -> `VALID`

## Implemented Tests

Current tests cover:

- Config defaults and environment overrides.
- Health handler JSON response.
- License type validation.
- Expiry date validation.
- License hash hex parsing.
- Request validation.
- Hash input string order.
- Hash determinism.
- Hash changes when any hash field changes.
- Issuer registration service flow.
- License issue service flow.
- License revoke service flow.
- License verification service statuses:
  - `VALID`
  - `INVALID`
  - `REVOKED`
  - `EXPIRED`

Latest verification command:

```bash
cd api
GOCACHE=/tmp/veryfy-go-cache go test ./...
```

Expected result: all packages pass.

## Not Implemented Yet

The following are still pending:

- `POST /issuers/register`
- `POST /licenses/issue`
- `POST /licenses/verify`
- `POST /licenses/revoke`
- Real Solana repository.
- Anchor instruction calls.
- PDA derivation.
- Account decoding.
- Devnet integration tests.
- `.env.example`.
- Local run documentation.

Current placeholder files:

```txt
api/internal/handlers/issuer.go
api/internal/handlers/license.go
api/internal/services/solana.go
```

## Next Backend Stage

Next stage should be Stage 4 from `implementation_plan.md`:

- Implement HTTP handlers for:
  - `POST /issuers/register`
  - `POST /licenses/issue`
  - `POST /licenses/verify`
  - `POST /licenses/revoke`
- Decode request JSON.
- Validate request bodies.
- Call services.
- Return response envelopes.
- Map service/repository errors to HTTP status codes.

Do not start real Solana integration before the HTTP handlers are wired and tested against the memory repository.

## Code Practices For Future Agents

Follow these rules when continuing backend work:

- Keep `cmd/server/main.go` thin. It should wire dependencies and start the server only.
- Keep handlers thin. Handlers decode, validate, call services, and encode responses.
- Keep business rules in services.
- Hide Solana details behind repository interfaces.
- Use manual constructor injection. Do not add a DI framework.
- Define interfaces where they are consumed.
- Pass `context.Context` through handlers, services, and repositories.
- Do not use package-level mutable globals for dependencies.
- Keep normal tests independent of devnet.
- Add integration tests for Solana later behind explicit environment flags.
- Use the existing JSON response envelope.
- Use the existing hash service for license hashes. Do not recompute hashes ad hoc in handlers.
- Preserve the exact hash field order.
- Keep `licenseType` flexible. Do not reintroduce a hardcoded list of allowed license types.
- Prefer small focused files over large catch-all files.
- Run `gofmt` before finishing backend changes.
- Run `GOCACHE=/tmp/veryfy-go-cache go test ./...` from `api/` before reporting completion.

## Cross-Team Contracts To Preserve

Frontend/backend/contract must agree on:

```txt
licenseNumber|holderName|issuerWallet|licenseType|expiryDate
```

`expiryDate` is an API string in:

```txt
YYYY-MM-DD
```

`licenseHash` is exposed by the backend API as lowercase hex.

Verification status strings must remain:

```txt
VALID
INVALID
REVOKED
EXPIRED
```

`licenseType` comes from frontend requests at runtime. For the MVP, the frontend may render it as a user-entered or dropdown-selected value. Later, it can be fetched from registered issuers.

## Known Sandbox Note

The backend server may fail to bind to a port inside the Codex sandbox with:

```txt
socket: operation not permitted
```

That does not necessarily indicate an application bug. Use handler tests and package tests for verification in the sandbox.
