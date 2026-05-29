// types/index.ts
export type LicenseType = "MEDICAL" | "LEGAL" | "DRIVING";
export type LicenseStatus = "VALID" | "INVALID" | "REVOKED" | "EXPIRED";

export interface License {
  licenseNumber: string;
  holderName: string;
  holderWallet: string;
  licenseType: LicenseType;
  expiryDate: string;
  issuerName: string;
  status: LicenseStatus;
}

export interface IssueLicenseParams {
  licenseNumber: string;
  holderName: string;
  holderWallet: string;
  licenseType: LicenseType;
  expiryDate: string;
  issuerWallet: string;
}

export interface VerifyLicenseParams {
  licenseNumber?: string;
  holderName?: string;
  issuerWallet?: string;
  licenseType?: LicenseType;
  expiryDate?: string;
  licenseHash?: string;
  qrCodeData?: string;
}
