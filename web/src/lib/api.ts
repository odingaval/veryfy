// lib/api.ts
import type { License, IssueLicenseParams, VerifyLicenseParams } from '../types';

const API_URL = import.meta.env.VITE_API_URL ?? 'http://localhost:8080';

interface ApiError {
  code: string;
  message: string;
}

interface ApiEnvelope<T> {
  data: T | null;
  error: ApiError | null;
}

async function readApiResponse<T>(res: Response): Promise<T> {
  const payload = await res.json();

  if (!res.ok) {
    const message = payload?.error?.message ?? payload?.message ?? 'Request failed';
    throw new Error(message);
  }

  if (payload && typeof payload === 'object' && 'data' in payload && 'error' in payload) {
    const envelope = payload as ApiEnvelope<T>;
    if (envelope.error) {
      throw new Error(envelope.error.message);
    }
    return envelope.data as T;
  }

  return payload as T;
}

export const api = {
  async issueLicense(params: IssueLicenseParams): Promise<{ licenseHash: string; txSignature: string }> {
    const res = await fetch(`${API_URL}/licenses/issue`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(params),
    });
    return readApiResponse(res);
  },

  async verifyLicense(params: VerifyLicenseParams): Promise<{ status: string; licenseHash: string; details?: License }> {
    const res = await fetch(`${API_URL}/licenses/verify`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(params),
    });
    return readApiResponse(res);
  },

  async revokeLicense(licenseHash: string, issuerWallet: string): Promise<{ txSignature: string }> {
    const res = await fetch(`${API_URL}/licenses/revoke`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ licenseHash, issuerWallet }),
    });
    return readApiResponse(res);
  }
};
