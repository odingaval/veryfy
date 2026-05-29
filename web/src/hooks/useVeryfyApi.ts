import { useConnection, useWallet } from '@solana/wallet-adapter-react';
import { Program, AnchorProvider, Idl, BN } from '@coral-xyz/anchor';
import { PublicKey, SYSVAR_RENT_PUBKEY, SystemProgram } from '@solana/web3.js';
import bs58 from 'bs58';
import idl from '../idl/veryfy.json';
import type { IssueLicenseParams, VerifyLicenseParams } from '../types';

const PROGRAM_ID = new PublicKey(idl.address);

function requireHashFields(params: VerifyLicenseParams): Required<Pick<
  VerifyLicenseParams,
  "licenseNumber" | "holderName" | "issuerWallet" | "licenseType" | "expiryDate"
>> {
  const requiredFields = ["licenseNumber", "holderName", "issuerWallet", "licenseType", "expiryDate"] as const;

  for (const field of requiredFields) {
    if (!params[field]) {
      throw new Error(`${field} is required to derive a license hash`);
    }
  }

  return params as Required<Pick<
    VerifyLicenseParams,
    "licenseNumber" | "holderName" | "issuerWallet" | "licenseType" | "expiryDate"
  >>;
}

function licenseHashInput(params: Required<Pick<
  VerifyLicenseParams,
  "licenseNumber" | "holderName" | "issuerWallet" | "licenseType" | "expiryDate"
>>): string {
  return `${params.licenseNumber}|${params.holderName}|${params.issuerWallet}|${params.licenseType}|${params.expiryDate}`;
}

async function sha256Bytes(input: string): Promise<Uint8Array> {
  const digest = await crypto.subtle.digest("SHA-256", new TextEncoder().encode(input));
  return new Uint8Array(digest);
}

async function hashLicenseData(params: VerifyLicenseParams): Promise<Uint8Array> {
  return sha256Bytes(licenseHashInput(requireHashFields(params)));
}

function decodeLicenseHash(value: string): Uint8Array {
  if (/^[0-9a-fA-F]{64}$/.test(value)) {
    const bytes = new Uint8Array(32);
    for (let i = 0; i < bytes.length; i += 1) {
      bytes[i] = Number.parseInt(value.slice(i * 2, i * 2 + 2), 16);
    }
    return bytes;
  }

  return bs58.decode(value);
}

export function useVeryfyApi() {
  const { connection } = useConnection();
  const wallet = useWallet();

  const getProgram = () => {
    if (!wallet.publicKey || !wallet.signTransaction) throw new Error("Wallet not connected");
    const provider = new AnchorProvider(connection, wallet as any, { commitment: "confirmed" });
    return new Program(idl as Idl, provider);
  };

  const issueLicense = async (params: IssueLicenseParams): Promise<{ licenseHash: string; txSignature: string }> => {
    if (!wallet.publicKey) throw new Error("Wallet not connected");
    const program = getProgram();
    const issuerWallet = wallet.publicKey.toString();

    const assetHash = await hashLicenseData({
      licenseNumber: params.licenseNumber,
      holderName: params.holderName,
      issuerWallet,
      licenseType: params.licenseType,
      expiryDate: params.expiryDate,
    });

    const expiry = params.expiryDate ? new Date(params.expiryDate).getTime() / 1000 : 0;

    // Derive PDAs using TextEncoder and Uint8Array instead of Buffer for browser compatibility
    const [issuerPda] = PublicKey.findProgramAddressSync(
      [new TextEncoder().encode("issuer"), wallet.publicKey.toBytes()],
      PROGRAM_ID
    );

    const [licensePda] = PublicKey.findProgramAddressSync(
      [new TextEncoder().encode("license"), assetHash],
      PROGRAM_ID
    );

    try {
      // First, try to register the issuer if it doesn't exist
      try {
        await (program.account as any).issuer.fetch(issuerPda);
      } catch (e) {
        // Doesn't exist, create it
        await program.methods.registerIssuer("Authority Name")
          .accounts({
            payer: wallet.publicKey,
            issuer: issuerPda,
            systemProgram: SystemProgram.programId,
          })
          .rpc();
      }

      // Now issue the license
      const tx = await program.methods.issueLicense(assetHash, new BN(expiry))
        .accounts({
          payer: wallet.publicKey,
          license: licensePda,
          issuer: issuerPda,
          authority: wallet.publicKey,
          systemProgram: SystemProgram.programId,
          rent: SYSVAR_RENT_PUBKEY,
        } as any)
        .rpc();

      return {
        licenseHash: bs58.encode(assetHash),
        txSignature: tx
      };
    } catch (e) {
      console.error(e);
      throw new Error("Failed to execute on-chain transaction");
    }
  };

  const verifyLicense = async (params: VerifyLicenseParams): Promise<{ status: string; details: any }> => {
    // This can be called without a wallet connected, just using connection
    const provider = new AnchorProvider(connection, {} as any, { commitment: "confirmed" });
    const program = new Program(idl as Idl, provider);

    let assetHash: Uint8Array;
    if (params.licenseHash || params.qrCodeData) {
      const encoded = params.licenseHash || params.qrCodeData || "";
      assetHash = decodeLicenseHash(encoded);
    } else {
      assetHash = await hashLicenseData(params);
    }

    const [licensePda] = PublicKey.findProgramAddressSync(
      [new TextEncoder().encode("license"), assetHash],
      PROGRAM_ID
    );

    try {
      const licenseAccount = await (program.account as any).license.fetch(licensePda);
      const now = Date.now() / 1000;
      let status = "VALID";
      if (Object.keys(licenseAccount.status)[0] === "revoked") status = "REVOKED";
      if (licenseAccount.expiry.toNumber() > 0 && licenseAccount.expiry.toNumber() < now) status = "EXPIRED";

      return {
        status,
        details: {
          licenseNumber: "Verified On-Chain",
          holderName: licenseAccount.holder.toString(),
          licenseType: "BLOCKCHAIN",
          status,
          issuedDate: new Date().toISOString(),
          expiryDate: licenseAccount.expiry.toNumber() === 0 ? null : new Date(licenseAccount.expiry.toNumber() * 1000).toISOString(),
          issuerId: licenseAccount.issuer.toString(),
          verificationHash: licensePda.toString()
        } as any
      };
    } catch (e) {
      console.error(e);
      return { status: "INVALID", details: null as any };
    }
  };

  const revokeLicense = async (licenseHashBase58: string): Promise<{ txSignature: string }> => {
    if (!wallet.publicKey) throw new Error("Wallet not connected");
    const program = getProgram();

    const assetHash = decodeLicenseHash(licenseHashBase58);

    const [issuerPda] = PublicKey.findProgramAddressSync(
      [new TextEncoder().encode("issuer"), wallet.publicKey.toBytes()],
      PROGRAM_ID
    );

    const [licensePda] = PublicKey.findProgramAddressSync(
      [new TextEncoder().encode("license"), assetHash],
      PROGRAM_ID
    );

    const tx = await program.methods.revokeLicense(Array.from(assetHash))
      .accounts({
        authority: wallet.publicKey,
        license: licensePda,
        issuer: issuerPda,
        authorityAccount: issuerPda, // as per the IDL accounts
        systemProgram: SystemProgram.programId,
      } as any)
      .rpc();

    return { txSignature: tx };
  };

  return { issueLicense, verifyLicense, revokeLicense };
}
