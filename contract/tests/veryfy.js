const anchor = require("@coral-xyz/anchor");
const { assert } = require("chai");

describe("veryfy", () => {
  // Configure the client to use the local cluster.
  anchor.setProvider(anchor.AnchorProvider.env());

  const program = anchor.workspace.Veryfy;
  const provider = anchor.getProvider();

  // Test accounts
  const issuerKeypair = anchor.web3.Keypair.generate();
  
  // Asset hash mock (32 bytes)
  const assetHash = Array.from({ length: 32 }, () => Math.floor(Math.random() * 256));

  it("Is initialized!", async () => {
    // Airdrop SOL to the issuer for fees
    const signature = await provider.connection.requestAirdrop(
      issuerKeypair.publicKey,
      2 * anchor.web3.LAMPORTS_PER_SOL
    );
    const latestBlockHash = await provider.connection.getLatestBlockhash();
    await provider.connection.confirmTransaction({
      blockhash: latestBlockHash.blockhash,
      lastValidBlockHeight: latestBlockHash.lastValidBlockHeight,
      signature: signature,
    });
  });

  it("Registers an issuer", async () => {
    // Derive the PDA for the issuer
    const [issuerPda, _issuerBump] = anchor.web3.PublicKey.findProgramAddressSync(
      [Buffer.from("issuer"), issuerKeypair.publicKey.toBuffer()],
      program.programId
    );

    const issuerName = "Veryfy Official";

    // Call register_issuer instruction
    await program.methods
      .registerIssuer(issuerName)
      .accounts({
        issuer: issuerPda,
        payer: issuerKeypair.publicKey,
        systemProgram: anchor.web3.SystemProgram.programId,
      })
      .signers([issuerKeypair])
      .rpc();

    // Fetch the account and assert its values
    const issuerAccount = await program.account.issuer.fetch(issuerPda);
    assert.strictEqual(issuerAccount.name, issuerName);
    assert.ok(issuerAccount.authority.equals(issuerKeypair.publicKey));
    assert.strictEqual(issuerAccount.issuedCount.toNumber(), 0);
  });

  it("Issues a license", async () => {
    // Derive PDAs
    const [issuerPda] = anchor.web3.PublicKey.findProgramAddressSync(
      [Buffer.from("issuer"), issuerKeypair.publicKey.toBuffer()],
      program.programId
    );
    const [licensePda] = anchor.web3.PublicKey.findProgramAddressSync(
      [Buffer.from("license"), Buffer.from(assetHash)],
      program.programId
    );

    const expiryTime = new anchor.BN(Date.now() / 1000 + 3600); // Expires in 1 hour

    await program.methods
      .issueLicense(assetHash, expiryTime)
      .accounts({
        license: licensePda,
        issuer: issuerPda,
        authority: issuerKeypair.publicKey,
        payer: issuerKeypair.publicKey,
        systemProgram: anchor.web3.SystemProgram.programId,
      })
      .signers([issuerKeypair])
      .rpc();

    // Fetch the license account and verify
    const licenseAccount = await program.account.license.fetch(licensePda);
    assert.ok(licenseAccount.issuer.equals(issuerPda));
    assert.ok(licenseAccount.holder.equals(issuerKeypair.publicKey));
    assert.strictEqual(licenseAccount.assetHash.toString(), assetHash.toString());
    assert.deepEqual(licenseAccount.status, { active: {} }); // Assuming the enum is { active: {} }
    
    // Check if the issuer's issued_count was incremented
    const issuerAccount = await program.account.issuer.fetch(issuerPda);
    assert.strictEqual(issuerAccount.issuedCount.toNumber(), 1);
  });

  it("Revokes a license", async () => {
    const [issuerPda] = anchor.web3.PublicKey.findProgramAddressSync(
      [Buffer.from("issuer"), issuerKeypair.publicKey.toBuffer()],
      program.programId
    );
    const [licensePda] = anchor.web3.PublicKey.findProgramAddressSync(
      [Buffer.from("license"), Buffer.from(assetHash)],
      program.programId
    );

    await program.methods
      .revokeLicense(assetHash)
      .accounts({
        license: licensePda,
        issuer: issuerPda,
        authority: issuerKeypair.publicKey,
      })
      .signers([issuerKeypair])
      .rpc();

    // Fetch the license account and verify its status is revoked
    const licenseAccount = await program.account.license.fetch(licensePda);
    assert.deepEqual(licenseAccount.status, { revoked: {} });
  });

  it("Fails to revoke a license with unauthorized authority", async () => {
    const maliciousKeypair = anchor.web3.Keypair.generate();
    
    // Airdrop SOL to malicious user
    const signature = await provider.connection.requestAirdrop(
      maliciousKeypair.publicKey,
      1 * anchor.web3.LAMPORTS_PER_SOL
    );
    const latestBlockHash = await provider.connection.getLatestBlockhash();
    await provider.connection.confirmTransaction({
      blockhash: latestBlockHash.blockhash,
      lastValidBlockHeight: latestBlockHash.lastValidBlockHeight,
      signature: signature,
    });

    const [issuerPda] = anchor.web3.PublicKey.findProgramAddressSync(
      [Buffer.from("issuer"), issuerKeypair.publicKey.toBuffer()],
      program.programId
    );
    const [licensePda] = anchor.web3.PublicKey.findProgramAddressSync(
      [Buffer.from("license"), Buffer.from(assetHash)],
      program.programId
    );

    try {
      await program.methods
        .revokeLicense(assetHash)
        .accounts({
          license: licensePda,
          issuer: issuerPda,
          authority: maliciousKeypair.publicKey,
        })
        .signers([maliciousKeypair])
        .rpc();
      
      assert.fail("Transaction should have failed");
    } catch (err) {
      assert.include(err.message, "UnauthorizedIssuer", "Expected an UnauthorizedIssuer error");
    }
  });
});
