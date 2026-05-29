import { Card, CardContent, CardHeader, CardTitle } from "./ui/card";
import { StatusBadge, LicenseStatus } from "./StatusBadge";
import { Separator } from "./ui/separator";
import { Shield, User, Calendar, Hash, Building2 } from "lucide-react";

export type LicenseType = "MEDICAL" | "LEGAL" | "DRIVING";

export interface License {
  licenseNumber: string;
  holderName: string;
  holderWallet: string;
  licenseType: LicenseType;
  expiryDate: string;
  issuerName: string;
  status: LicenseStatus;
  issuedAt?: string;
  licenseHash?: string;
}

interface LicenseCardProps {
  license: License;
  showHash?: boolean;
}

const LICENSE_TYPE_LABELS = {
  MEDICAL: "Medical License",
  LEGAL: "Legal License",
  DRIVING: "Driving License",
};

export function LicenseCard({ license, showHash = false }: LicenseCardProps) {
  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleDateString("en-US", {
      year: "numeric",
      month: "long",
      day: "numeric",
    });
  };

  return (
    <Card className="w-full">
      <CardHeader className="flex flex-row items-start justify-between space-y-0">
        <div className="space-y-1">
          <CardTitle className="flex items-center gap-2">
            <Shield className="w-5 h-5" />
            {LICENSE_TYPE_LABELS[license.licenseType]}
          </CardTitle>
          <p className="text-sm text-muted-foreground">
            License #{license.licenseNumber}
          </p>
        </div>
        <StatusBadge status={license.status} />
      </CardHeader>
      <CardContent className="space-y-4">
        <div className="grid gap-3">
          <div className="flex items-start gap-3">
            <User className="w-4 h-4 mt-0.5 text-muted-foreground" />
            <div className="space-y-0.5">
              <p className="text-sm font-medium">License Holder</p>
              <p className="text-sm text-muted-foreground">{license.holderName}</p>
              <p className="text-xs text-muted-foreground font-mono break-all">
                {license.holderWallet}
              </p>
            </div>
          </div>

          <Separator />

          <div className="flex items-start gap-3">
            <Building2 className="w-4 h-4 mt-0.5 text-muted-foreground" />
            <div className="space-y-0.5">
              <p className="text-sm font-medium">Issuing Authority</p>
              <p className="text-sm text-muted-foreground">{license.issuerName}</p>
            </div>
          </div>

          <Separator />

          <div className="grid grid-cols-2 gap-4">
            <div className="flex items-start gap-2">
              <Calendar className="w-4 h-4 mt-0.5 text-muted-foreground" />
              <div className="space-y-0.5">
                <p className="text-sm font-medium">Expires</p>
                <p className="text-sm text-muted-foreground">
                  {formatDate(license.expiryDate)}
                </p>
              </div>
            </div>

            {license.issuedAt && (
              <div className="flex items-start gap-2">
                <Calendar className="w-4 h-4 mt-0.5 text-muted-foreground" />
                <div className="space-y-0.5">
                  <p className="text-sm font-medium">Issued</p>
                  <p className="text-sm text-muted-foreground">
                    {formatDate(license.issuedAt)}
                  </p>
                </div>
              </div>
            )}
          </div>

          {showHash && license.licenseHash && (
            <>
              <Separator />
              <div className="flex items-start gap-3">
                <Hash className="w-4 h-4 mt-0.5 text-muted-foreground" />
                <div className="space-y-0.5">
                  <p className="text-sm font-medium">License Hash</p>
                  <p className="text-xs text-muted-foreground font-mono break-all">
                    {license.licenseHash}
                  </p>
                </div>
              </div>
            </>
          )}
        </div>
      </CardContent>
    </Card>
  );
}
