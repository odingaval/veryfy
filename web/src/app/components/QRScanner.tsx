import { useState } from "react";
import { Button } from "./ui/button";
import { Card, CardContent } from "./ui/card";
import { QrCode, Camera, X } from "lucide-react";
import type { LicenseType } from "./LicenseCard";

interface QRScannerProps {
  onScan: (data: {
    licenseNumber: string;
    holderName: string;
    issuerWallet: string;
    licenseType: LicenseType;
    expiryDate: string;
  }) => void;
  onClose: () => void;
}

export function QRScanner({ onScan, onClose }: QRScannerProps) {
  const [scanning, setScanning] = useState(false);

  const handleMockScan = () => {
    setScanning(true);

    // Simulate scanning delay
    setTimeout(() => {
      // Mock QR code data
      const mockData = {
        licenseNumber: "KE/MED/98765",
        holderName: "Sarah Ochieng",
        issuerWallet: "9vKXtg2CW87d97TXJSDpbD5jBkheTqA83TZRuJosgDfG",
        licenseType: "MEDICAL" as LicenseType,
        expiryDate: "2027-12-31",
      };

      onScan(mockData);
      setScanning(false);
    }, 2000);
  };

  return (
    <Card className="w-full">
      <CardContent className="pt-6">
        <div className="space-y-4">
          <div className="flex items-center justify-between">
            <div className="flex items-center gap-2">
              <QrCode className="w-5 h-5" />
              <h3 className="font-semibold">Scan QR Code</h3>
            </div>
            <Button variant="ghost" size="sm" onClick={onClose}>
              <X className="w-4 h-4" />
            </Button>
          </div>

          {/* Mock Camera View */}
          <div className="relative aspect-square bg-muted rounded-lg overflow-hidden flex items-center justify-center">
            {scanning ? (
              <div className="text-center space-y-4">
                <div className="w-16 h-16 border-4 border-primary border-t-transparent rounded-full animate-spin mx-auto" />
                <p className="text-sm text-muted-foreground">
                  Scanning QR code...
                </p>
              </div>
            ) : (
              <div className="text-center space-y-4 p-8">
                <Camera className="w-16 h-16 mx-auto text-muted-foreground" />
                <div className="space-y-2">
                  <p className="text-sm font-medium">
                    Position QR code within the frame
                  </p>
                  <p className="text-xs text-muted-foreground">
                    The camera will automatically scan the license QR code
                  </p>
                </div>
              </div>
            )}

            {/* Scan Frame Overlay */}
            {!scanning && (
              <div className="absolute inset-0 flex items-center justify-center pointer-events-none">
                <div className="w-48 h-48 border-2 border-primary rounded-lg relative">
                  {/* Corner markers */}
                  <div className="absolute top-0 left-0 w-4 h-4 border-t-4 border-l-4 border-primary" />
                  <div className="absolute top-0 right-0 w-4 h-4 border-t-4 border-r-4 border-primary" />
                  <div className="absolute bottom-0 left-0 w-4 h-4 border-b-4 border-l-4 border-primary" />
                  <div className="absolute bottom-0 right-0 w-4 h-4 border-b-4 border-r-4 border-primary" />
                </div>
              </div>
            )}
          </div>

          {/* Demo Button */}
          <div className="space-y-2">
            <Button
              onClick={handleMockScan}
              disabled={scanning}
              className="w-full"
            >
              <Camera className="w-4 h-4 mr-2" />
              {scanning ? "Scanning..." : "Simulate QR Scan (Demo)"}
            </Button>
            <p className="text-xs text-center text-muted-foreground">
              In production, this would access your device camera
            </p>
          </div>
        </div>
      </CardContent>
    </Card>
  );
}
