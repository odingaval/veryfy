package models

// IssueLicenseRequest is the request body for POST /licenses/issue.
type IssueLicenseRequest struct {
	LicenseNumber string      `json:"licenseNumber"`
	HolderName    string      `json:"holderName"`
	HolderWallet  string      `json:"holderWallet"`
	LicenseType   LicenseType `json:"licenseType"`
	ExpiryDate    string      `json:"expiryDate"`
	IssuerWallet  string      `json:"issuerWallet"`
}

// VerifyLicenseRequest is the request body for POST /licenses/verify.
type VerifyLicenseRequest struct {
	LicenseNumber string      `json:"licenseNumber"`
	HolderName    string      `json:"holderName"`
	LicenseType   LicenseType `json:"licenseType"`
	ExpiryDate    string      `json:"expiryDate"`
	IssuerWallet  string      `json:"issuerWallet"`
}

// RevokeLicenseRequest is the request body for POST /licenses/revoke.
type RevokeLicenseRequest struct {
	LicenseHash  string `json:"licenseHash"`
	IssuerWallet string `json:"issuerWallet"`
}

// RegisterIssuerRequest is the request body for POST /issuers/register.
type RegisterIssuerRequest struct {
	Name        string      `json:"name"`
	LicenseType LicenseType `json:"licenseType"`
	Wallet      string      `json:"wallet"`
}

// IssueLicenseResponse is returned after a successful license issuance.
type IssueLicenseResponse struct {
	LicenseHash string `json:"licenseHash"`
	TxSignature string `json:"txSignature"`
}

// VerifyLicenseResponse is returned after a verification attempt.
type VerifyLicenseResponse struct {
	Status      LicenseStatus `json:"status"`
	LicenseHash string        `json:"licenseHash"`
}

// RevokeLicenseResponse is returned after a successful revocation.
type RevokeLicenseResponse struct {
	TxSignature string `json:"txSignature"`
}

// RegisterIssuerResponse is returned after successful issuer registration.
type RegisterIssuerResponse struct {
	TxSignature string `json:"txSignature"`
}
