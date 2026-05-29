// lib/api.ts
import { License, IssueLicenseParams, VerifyLicenseParams } from '../types';

const API_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080';

export const api = {
  async issueLicense(params: IssueLicenseParams): Promise<{ licenseHash: string; txSignature: string }> {
    const res = await fetch(`${API_URL}/licenses/issue`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(params),
    });
    if (!res.ok) throw new Error('Failed to issue license');
    return res.json();
  },

  async verifyLicense(params: VerifyLicenseParams): Promise<{ status: string; details: License }> {
    const res = await fetch(`${API_URL}/licenses/verify`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(params),
    });
    if (!res.ok) throw new Error('Verification failed');
    return res.json();
  },

  async revokeLicense(licenseHash: string, issuerWallet: string): Promise<{ txSignature: string }> {
    const res = await fetch(`${API_URL}/licenses/revoke`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ licenseHash, issuerWallet }),
    });
    if (!res.ok) throw new Error('Revocation failed');
    return res.json();
  }
};
