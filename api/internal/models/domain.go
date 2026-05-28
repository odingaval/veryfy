package models

import "time"

// LicenseAccount is the backend's decoded view of a stored license.
type LicenseAccount struct {
	LicenseHash   LicenseHash
	LicenseNumber string
	HolderName    string
	HolderWallet  string
	IssuerWallet  string
	LicenseType   LicenseType
	ExpiryDate    string
	IsRevoked     bool
	IssuedAt      time.Time
}

// IssuerAccount is the backend's decoded view of a registered issuer.
type IssuerAccount struct {
	Name         string
	Wallet       string
	LicenseType  LicenseType
	RegisteredAt time.Time
}

// IssueLicenseParams contains service-to-repository issue inputs.
type IssueLicenseParams struct {
	LicenseHash   LicenseHash
	LicenseNumber string
	HolderName    string
	HolderWallet  string
	IssuerWallet  string
	LicenseType   LicenseType
	ExpiryDate    string
}

// RevokeLicenseParams contains service-to-repository revoke inputs.
type RevokeLicenseParams struct {
	LicenseHash  LicenseHash
	IssuerWallet string
}

// RegisterIssuerParams contains service-to-repository issuer registration inputs.
type RegisterIssuerParams struct {
	Name        string
	Wallet      string
	LicenseType LicenseType
}
