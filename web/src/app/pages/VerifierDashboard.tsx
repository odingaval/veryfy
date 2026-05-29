import { useState } from "react";
import { Button } from "../components/ui/button";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "../components/ui/card";
import { Input } from "../components/ui/input";
import { Label } from "../components/ui/label";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "../components/ui/select";
import { LicenseCard, License, LicenseType } from "../components/LicenseCard";
import { LicenseStatus } from "../components/StatusBadge";
import { QRScanner } from "../components/QRScanner";
import { Search, ArrowLeft, Loader2, QrCode, Keyboard } from "lucide-react";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "../components/ui/tabs";

interface VerifierDashboardProps {
  onNavigate: (page: "landing" | "issue") => void;
}

export function VerifierDashboard({ onNavigate }: VerifierDashboardProps) {
  const [loading, setLoading] = useState(false);
  const [verifiedLicense, setVerifiedLicense] = useState<License | null>(null);
  const [inputMode, setInputMode] = useState<"manual" | "qr">("manual");
  const [formData, setFormData] = useState({
    licenseNumber: "",
    holderName: "",
    issuerWallet: "",
    licenseType: "" as LicenseType | "",
    expiryDate: "",
  });

  const handleVerify = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);

    // Simulate API call
    await new Promise((resolve) => setTimeout(resolve, 1500));

    // Mock verification result
    const mockLicense: License = {
      licenseNumber: formData.licenseNumber,
      holderName: formData.holderName,
      holderWallet: "7xKXtg2CW87d97TXJSDpbD5jBkheTqA83TZRuJosgAsU",
      licenseType: formData.licenseType as LicenseType,
      expiryDate: formData.expiryDate,
      issuerName: "Kenya Medical Practitioners and Dentists Council",
      status: new Date(formData.expiryDate) > new Date() ? "VALID" : "EXPIRED",
      issuedAt: "2024-01-15",
      licenseHash: "0x1a2b3c4d5e6f7g8h9i0j1k2l3m4n5o6p7q8r9s0t1u2v3w4x5y6z",
    };

    setVerifiedLicense(mockLicense);
    setLoading(false);
  };

  const handleInputChange = (field: string, value: string) => {
    setFormData((prev) => ({ ...prev, [field]: value }));
  };

  const handleQRScan = (data: {
    licenseNumber: string;
    holderName: string;
    issuerWallet: string;
    licenseType: LicenseType;
    expiryDate: string;
  }) => {
    setFormData(data);
    setInputMode("manual");
    // Auto-verify after QR scan
    setTimeout(() => {
      handleVerify(new Event("submit") as any);
    }, 100);
  };

  return (
    <div className="min-h-screen bg-background">
      {/* Header */}
      <div className="border-b border-border">
        <div className="container mx-auto px-4 py-4">
          <div className="flex items-center justify-between">
            <div className="flex items-center gap-4">
              <Button
                variant="ghost"
                size="sm"
                onClick={() => onNavigate("landing")}
              >
                <ArrowLeft className="w-4 h-4 mr-2" />
                Back
              </Button>
              <h1 className="text-2xl font-bold">Verify License</h1>
            </div>
            <Button
              variant="outline"
              onClick={() => onNavigate("issue")}
            >
              Issue License
            </Button>
          </div>
        </div>
      </div>

      {/* Main Content */}
      <div className="container mx-auto px-4 py-8">
        <div className="max-w-4xl mx-auto grid md:grid-cols-2 gap-8">
          {/* Verification Form */}
          <div>
            <Card>
              <CardHeader>
                <div className="flex items-start justify-between">
                  <div>
                    <CardTitle>Verify License</CardTitle>
                    <CardDescription>
                      Scan a QR code or manually enter license details to verify
                    </CardDescription>
                  </div>
                  {inputMode === "manual" && (
                    <Button
                      variant="outline"
                      size="sm"
                      onClick={() => setInputMode("qr")}
                    >
                      <QrCode className="w-4 h-4 mr-2" />
                      Scan QR
                    </Button>
                  )}
                </div>
              </CardHeader>
              <CardContent>
                <Tabs value={inputMode} onValueChange={(v) => setInputMode(v as "manual" | "qr")}>
                  <TabsList className="grid w-full grid-cols-2 mb-4">
                    <TabsTrigger value="manual">
                      <Keyboard className="w-4 h-4 mr-2" />
                      Manual Entry
                    </TabsTrigger>
                    <TabsTrigger value="qr">
                      <QrCode className="w-4 h-4 mr-2" />
                      Scan QR Code
                    </TabsTrigger>
                  </TabsList>

                  <TabsContent value="manual">
                    <form onSubmit={handleVerify} className="space-y-4">
                  <div className="space-y-2">
                    <Label htmlFor="licenseType">License Type</Label>
                    <Select
                      value={formData.licenseType}
                      onValueChange={(value) =>
                        handleInputChange("licenseType", value)
                      }
                    >
                      <SelectTrigger id="licenseType">
                        <SelectValue placeholder="Select license type" />
                      </SelectTrigger>
                      <SelectContent>
                        <SelectItem value="MEDICAL">Medical License</SelectItem>
                        <SelectItem value="LEGAL">Legal License</SelectItem>
                        <SelectItem value="DRIVING">Driving License</SelectItem>
                      </SelectContent>
                    </Select>
                  </div>

                  <div className="space-y-2">
                    <Label htmlFor="licenseNumber">License Number</Label>
                    <Input
                      id="licenseNumber"
                      placeholder="e.g., KE/MED/12345"
                      value={formData.licenseNumber}
                      onChange={(e) =>
                        handleInputChange("licenseNumber", e.target.value)
                      }
                      required
                    />
                  </div>

                  <div className="space-y-2">
                    <Label htmlFor="holderName">License Holder Name</Label>
                    <Input
                      id="holderName"
                      placeholder="e.g., John Kamau"
                      value={formData.holderName}
                      onChange={(e) =>
                        handleInputChange("holderName", e.target.value)
                      }
                      required
                    />
                  </div>

                  <div className="space-y-2">
                    <Label htmlFor="issuerWallet">Issuer Wallet</Label>
                    <Input
                      id="issuerWallet"
                      placeholder="Solana wallet address"
                      value={formData.issuerWallet}
                      onChange={(e) =>
                        handleInputChange("issuerWallet", e.target.value)
                      }
                      required
                    />
                  </div>

                  <div className="space-y-2">
                    <Label htmlFor="expiryDate">Expiry Date</Label>
                    <Input
                      id="expiryDate"
                      type="date"
                      value={formData.expiryDate}
                      onChange={(e) =>
                        handleInputChange("expiryDate", e.target.value)
                      }
                      required
                    />
                  </div>

                      <Button
                        type="submit"
                        className="w-full"
                        disabled={loading}
                      >
                        {loading ? (
                          <>
                            <Loader2 className="w-4 h-4 mr-2 animate-spin" />
                            Verifying...
                          </>
                        ) : (
                          <>
                            <Search className="w-4 h-4 mr-2" />
                            Verify License
                          </>
                        )}
                      </Button>
                    </form>
                  </TabsContent>

                  <TabsContent value="qr">
                    <QRScanner
                      onScan={handleQRScan}
                      onClose={() => setInputMode("manual")}
                    />
                  </TabsContent>
                </Tabs>
              </CardContent>
            </Card>
          </div>

          {/* Verification Result */}
          <div>
            {verifiedLicense ? (
              <div className="space-y-4">
                <h2 className="text-xl font-semibold">Verification Result</h2>
                <LicenseCard license={verifiedLicense} showHash />
              </div>
            ) : (
              <Card className="h-full flex items-center justify-center min-h-[400px]">
                <CardContent className="text-center space-y-2">
                  <Search className="w-12 h-12 mx-auto text-muted-foreground" />
                  <p className="text-muted-foreground">
                    Enter license details and click verify to see results
                  </p>
                </CardContent>
              </Card>
            )}
          </div>
        </div>
      </div>
    </div>
  );
}
