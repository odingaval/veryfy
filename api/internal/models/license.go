package models

// LicenseType identifies an issuer-defined licensing category.
type LicenseType string

// LicenseStatus is the public verification result returned by the API.
type LicenseStatus string

const (
	LicenseStatusValid   LicenseStatus = "VALID"
	LicenseStatusInvalid LicenseStatus = "INVALID"
	LicenseStatusRevoked LicenseStatus = "REVOKED"
	LicenseStatusExpired LicenseStatus = "EXPIRED"
)

// LicenseHash is the canonical API representation of a SHA-256 license hash.
type LicenseHash struct {
	Bytes [32]byte
	Hex   string
}
