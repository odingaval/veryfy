import React, { useState } from "react";
import { useVeryfyApi } from "../../hooks/useVeryfyApi";
import { QRScanner } from "../components/QRScanner";
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from "../components/ui/card";
import { Input } from "../components/ui/input";
import { Button } from "../components/ui/button";
import { CheckCircle, XCircle, Clock, Loader2 } from "lucide-react";
import { Badge } from "../components/ui/badge";
import type { VerifyLicenseParams } from "../../types";

export function PublicVerify() {
  const [licenseInput, setLicenseInput] = useState("");
  const [loading, setLoading] = useState(false);
  const [result, setResult] = useState<any | null>(null);
  
  const veryfyApi = useVeryfyApi();

  const performVerification = async (params: VerifyLicenseParams) => {
    setLoading(true);
    try {
      const { status, details } = await veryfyApi.verifyLicense(params);
      const normalizedStatus = status === "VALID" ? "verified" : status.toLowerCase();

      if (!details) {
        setResult({ status: "invalid", notFound: true });
        return;
      }

      setResult({
        status: normalizedStatus,
        organization: "On-Chain Veryfy Network",
        issuer: details.issuerId,
        issueDate: details.issuedDate.split("T")[0],
        expiryDate: details.expiryDate ? details.expiryDate.split("T")[0] : "Never",
        txHash: details.verificationHash,
        category: details.licenseType,
        holder: details.holderName,
      });
    } catch (e) {
      setResult({ status: "invalid", notFound: true });
    } finally {
      setLoading(false);
    }
  };

  const handleSearch = async () => {
    if (!licenseInput) return;
    await performVerification({ licenseHash: licenseInput });
  };

  const handleQRScan = async (data: VerifyLicenseParams) => {
    setLicenseInput("");
    await performVerification(data);
  };

  return (
    <div className="min-h-screen flex flex-col items-center pt-24 pb-12 bg-[#05050A]">
      {/* Header */}
      <h1 className="text-4xl md:text-5xl font-extrabold text-white mb-8 text-center">
        Verify a License
      </h1>

      {/* QR Scanner */}
      <div className="w-full max-w-xl mb-8">
        <QRScanner onScan={handleQRScan} onClose={() => {}} />
      </div>

      {/* Manual Search */}
      <div className="w-full max-w-xl flex gap-2 mb-8">
        <Input
          placeholder="Enter license hash (Base58 or hex)"
          value={licenseInput}
          onChange={e => setLicenseInput(e.target.value)}
        />
        <Button onClick={handleSearch} disabled={loading} className="whitespace-nowrap">
          {loading ? <Loader2 className="w-4 h-4 mr-2 animate-spin" /> : "Search"}
        </Button>
      </div>

      {/* Result Card */}
      {result && (
        <Card className="glass-panel w-full max-w-2xl">
          <CardHeader className="flex flex-col space-y-2">
            <div className="flex items-center justify-between">
              <CardTitle className="text-2xl text-white">
                {result.status === "verified" && <CheckCircle className="inline w-6 h-6 text-emerald-400 mr-2" />}
                {result.status === "expired" && <Clock className="inline w-6 h-6 text-amber-400 mr-2" />}
                {result.status === "revoked" && <XCircle className="inline w-6 h-6 text-red-500 mr-2" />}
                License {result.status.toUpperCase()}
              </CardTitle>
              <Badge
                variant={result.status === "verified" ? "default" : "destructive"}
                className="text-sm font-medium"
              >
                {result.status.toUpperCase()}
              </Badge>
            </div>
            <CardDescription className="text-white/70">
              Verification details for <span className="font-medium text-white">{licenseInput || "scanned license data"}</span>
            </CardDescription>
          </CardHeader>
          <CardContent className="grid grid-cols-2 gap-4 text-white/80 pt-4">
            <div><strong>Organization:</strong> {result.organization}</div>
            <div><strong>Issuer:</strong> {result.issuer}</div>
            <div><strong>Category:</strong> {result.category}</div>
            <div><strong>Issued:</strong> {result.issueDate}</div>
            <div><strong>Expires:</strong> {result.expiryDate}</div>
            <div className="flex items-center">
              <strong>Tx Hash:</strong>
              <span className="ml-1 font-mono text-xs break-all">{result.txHash}</span>
            </div>
          </CardContent>
        </Card>
      )}
    </div>
  );
}
