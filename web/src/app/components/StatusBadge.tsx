import { Badge } from "./ui/badge";
import { CheckCircle2, XCircle, AlertCircle, Clock } from "lucide-react";

export type LicenseStatus = "VALID" | "INVALID" | "REVOKED" | "EXPIRED";

interface StatusBadgeProps {
  status: LicenseStatus;
}

export function StatusBadge({ status }: StatusBadgeProps) {
  const config = {
    VALID: {
      label: "Valid",
      className: "bg-success text-success-foreground hover:bg-success/90",
      icon: CheckCircle2,
    },
    INVALID: {
      label: "Invalid",
      className: "bg-muted text-muted-foreground hover:bg-muted/90",
      icon: XCircle,
    },
    REVOKED: {
      label: "Revoked",
      className: "bg-destructive text-destructive-foreground hover:bg-destructive/90",
      icon: AlertCircle,
    },
    EXPIRED: {
      label: "Expired",
      className: "bg-warning text-warning-foreground hover:bg-warning/90",
      icon: Clock,
    },
  };

  const { label, className, icon: Icon } = config[status];

  return (
    <Badge className={className}>
      <Icon className="w-3 h-3 mr-1" />
      {label}
    </Badge>
  );
}
