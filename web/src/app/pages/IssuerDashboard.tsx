import { useState } from "react";
import { Button } from "../components/ui/button";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "../components/ui/card";
import { Input } from "../components/ui/input";
import { Label } from "../components/ui/label";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "../components/ui/select";
import { WalletConnect } from "../components/WalletConnect";
import { LicenseType } from "../components/LicenseCard";
import { StatusBadge, LicenseStatus } from "../components/StatusBadge";
import { QRScanner } from "../components/QRScanner";
import { ArrowLeft, FileCheck, Loader2, Trash2, QrCode } from "lucide-react";
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "../components/ui/table";
import { Dialog, DialogContent, DialogDescription, DialogHeader, DialogTitle, DialogTrigger } from "../components/ui/dialog";

interface IssuerDashboardProps {
  onNavigate: (page: "landing" | "verify") => void;
}

interface IssuedLicense {
  id: string;
  licenseNumber: string;
  holderName: string;
  licenseType: LicenseType;
  status: LicenseStatus;
  issuedDate: string;
}

export function IssuerDashboard({ onNavigate }: IssuerDashboardProps) {
  const [connected, setConnected] = useState(false);
  const [loading, setLoading] = useState(false);
  const [issuedLicenses, setIssuedLicenses] = useState<IssuedLicense[]>([
    {
      id: "1",
      licenseNumber: "KE/MED/12345",
      holderName: "John Kamau",
      licenseType: "MEDICAL",
      status: "VALID",
      issuedDate: "2024-01-15",
    },
    {
      id: "2",
      licenseNumber: "KE/LEG/67890",
      holderName: "Jane Wanjiku",
      licenseType: "LEGAL",
      status: "VALID",
      issuedDate: "2024-02-20",
    },
  ]);

  const [formData, setFormData] = useState({
    licenseType: "" as LicenseType | "",
    holderName: "",
    licenseNumber: "",
    holderWallet: "",
    expiryDate: "",
  });

  const handleIssueLicense = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);

    // Simulate API call
    await new Promise((resolve) => setTimeout(resolve, 1500));

    // Add to issued licenses
    const newLicense: IssuedLicense = {
      id: Date.now().toString(),
      licenseNumber: formData.licenseNumber,
      holderName: formData.holderName,
      licenseType: formData.licenseType as LicenseType,
      status: "VALID",
      issuedDate: new Date().toISOString().split("T")[0],
    };

    setIssuedLicenses((prev) => [newLicense, ...prev]);

    // Reset form
    setFormData({
      licenseType: "",
      holderName: "",
      licenseNumber: "",
      holderWallet: "",
      expiryDate: "",
    });

    setLoading(false);
  };

  const handleRevoke = async (id: string) => {
    setIssuedLicenses((prev) =>
      prev.map((license) =>
        license.id === id ? { ...license, status: "REVOKED" as LicenseStatus } : license
      )
    );
  };

  const handleInputChange = (field: string, value: string) => {
    setFormData((prev) => ({ ...prev, [field]: value }));
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
              <h1 className="text-2xl font-bold">Issue License</h1>
            </div>
            <div className="flex items-center gap-4">
              <Button
                variant="outline"
                onClick={() => onNavigate("verify")}
              >
                Verify License
              </Button>
              <WalletConnect
                onConnect={() => setConnected(true)}
                onDisconnect={() => setConnected(false)}
              />
            </div>
          </div>
        </div>
      </div>

      {/* Main Content */}
      <div className="container mx-auto px-4 py-8 space-y-8">
        {!connected ? (
          <Card className="max-w-2xl mx-auto">
            <CardContent className="pt-6 text-center space-y-4">
              <div className="w-16 h-16 bg-muted rounded-full flex items-center justify-center mx-auto">
                <FileCheck className="w-8 h-8 text-muted-foreground" />
              </div>
              <div className="space-y-2">
                <h2 className="text-xl font-semibold">Connect Your Wallet</h2>
                <p className="text-muted-foreground">
                  Please connect your authorized issuer wallet to issue and manage licenses.
                </p>
              </div>
            </CardContent>
          </Card>
        ) : (
          <>
            {/* Issue License Form */}
            <div className="max-w-2xl mx-auto">
              <Card>
                <CardHeader>
                  <CardTitle>Issue New License</CardTitle>
                  <CardDescription>
                    Create a new blockchain-verified license for a license holder.
                  </CardDescription>
                </CardHeader>
                <CardContent>
                  <form onSubmit={handleIssueLicense} className="space-y-4">
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

                    <div className="grid md:grid-cols-2 gap-4">
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
                        <Label htmlFor="holderName">Holder Name</Label>
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
                    </div>

                    <div className="space-y-2">
                      <Label htmlFor="holderWallet">Holder Wallet Address</Label>
                      <Input
                        id="holderWallet"
                        placeholder="Solana wallet address"
                        value={formData.holderWallet}
                        onChange={(e) =>
                          handleInputChange("holderWallet", e.target.value)
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
                          Issuing License...
                        </>
                      ) : (
                        <>
                          <FileCheck className="w-4 h-4 mr-2" />
                          Issue License
                        </>
                      )}
                    </Button>
                  </form>
                </CardContent>
              </Card>
            </div>

            {/* Issued Licenses Table */}
            <div className="max-w-6xl mx-auto">
              <Card>
                <CardHeader>
                  <CardTitle>Issued Licenses</CardTitle>
                  <CardDescription>
                    Manage and revoke licenses you have issued.
                  </CardDescription>
                </CardHeader>
                <CardContent>
                  <Table>
                    <TableHeader>
                      <TableRow>
                        <TableHead>License Number</TableHead>
                        <TableHead>Holder Name</TableHead>
                        <TableHead>Type</TableHead>
                        <TableHead>Issued Date</TableHead>
                        <TableHead>Status</TableHead>
                        <TableHead className="text-right">Actions</TableHead>
                      </TableRow>
                    </TableHeader>
                    <TableBody>
                      {issuedLicenses.map((license) => (
                        <TableRow key={license.id}>
                          <TableCell className="font-mono text-sm">
                            {license.licenseNumber}
                          </TableCell>
                          <TableCell>{license.holderName}</TableCell>
                          <TableCell className="capitalize">
                            {license.licenseType.toLowerCase()}
                          </TableCell>
                          <TableCell>
                            {new Date(license.issuedDate).toLocaleDateString()}
                          </TableCell>
                          <TableCell>
                            <StatusBadge status={license.status} />
                          </TableCell>
                          <TableCell className="text-right">
                            <Button
                              variant="ghost"
                              size="sm"
                              onClick={() => handleRevoke(license.id)}
                              disabled={license.status === "REVOKED"}
                            >
                              <Trash2 className="w-4 h-4 mr-2" />
                              Revoke
                            </Button>
                          </TableCell>
                        </TableRow>
                      ))}
                    </TableBody>
                  </Table>
                </CardContent>
              </Card>
            </div>
          </>
        )}
      </div>
    </div>
  );
}
