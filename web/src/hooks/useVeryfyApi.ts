import { useConnection, useWallet } from '@solana/wallet-adapter-react';
import { Program, AnchorProvider, Idl, BN } from '@coral-xyz/anchor';
import { PublicKey, SystemProgram } from '@solana/web3.js';
import bs58 from 'bs58';
import idl from '../idl/veryfy.json';
import type { IssueLicenseParams, VerifyLicenseParams } from '../types';

const PROGRAM_ID = new PublicKey(idl.address);

// Helper to hash off-chain metadata (e.g. license info) into a 32-byte array
function hashLicenseData(params: any): number[] {
  // Simple deterministic hash mock. In production, use SHA256 of JSON string
  const str = JSON.stringify(params);
  const arr = new Array(32).fill(0);
  for (let i = 0; i < str.length; i++) {
    arr[i % 32] = (arr[i % 32] + str.charCodeAt(i)) % 256;
  }
  return arr;
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

    const assetHash = hashLicenseData({
      licenseType: params.licenseType,
      holderName: params.holderName,
      licenseNumber: params.licenseNumber,
    });

    const expiry = params.expiryDate ? new Date(params.expiryDate).getTime() / 1000 : 0;

    // Derive PDAs using TextEncoder and Uint8Array instead of Buffer for browser compatibility
    const [issuerPda] = PublicKey.findProgramAddressSync(
      [new TextEncoder().encode("issuer"), wallet.publicKey.toBytes()],
      PROGRAM_ID
    );

    const [licensePda] = PublicKey.findProgramAddressSync(
      [new TextEncoder().encode("license"), new Uint8Array(assetHash)],
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
        } as any)
        .rpc();

      return {
        licenseHash: bs58.encode(new Uint8Array(assetHash)),
        txSignature: tx
      };
    } catch (e: any) {
      console.error("Smart Contract Error:", e);
      throw new Error(e.message || "Failed to execute on-chain transaction");
    }
  };

  const verifyLicense = async (params: VerifyLicenseParams): Promise<{ status: string; details: any }> => {
    // This can be called without a wallet connected, just using connection
    const provider = new AnchorProvider(connection, {} as any, { commitment: "confirmed" });
    const program = new Program(idl as Idl, provider);

    // Recompute hash
    const assetHash = bs58.decode(params.licenseHash || params.qrCodeData || "");
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
    } catch (e: any) {
      console.error("Verification Error:", e);
      return { status: "INVALID", details: null as any };
    }
  };

  const revokeLicense = async (licenseHashBase58: string): Promise<{ txSignature: string }> => {
    if (!wallet.publicKey) throw new Error("Wallet not connected");
    const program = getProgram();

    const assetHash = bs58.decode(licenseHashBase58);

    const [issuerPda] = PublicKey.findProgramAddressSync(
      [new TextEncoder().encode("issuer"), wallet.publicKey.toBytes()],
      PROGRAM_ID
    );

    const [licensePda] = PublicKey.findProgramAddressSync(
      [new TextEncoder().encode("license"), assetHash],
      PROGRAM_ID
    );

    try {
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
    } catch (e: any) {
      console.error("Revoke Error:", e);
      throw new Error(e.message || "Failed to execute on-chain transaction");
    }
  };

  return { issueLicense, verifyLicense, revokeLicense };
}
