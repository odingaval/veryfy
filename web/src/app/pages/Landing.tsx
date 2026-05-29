import { Button } from "../components/ui/button";
import { Card, CardContent } from "../components/ui/card";
import { Shield, Search, FileCheck, Lock } from "lucide-react";

interface LandingProps {
  onNavigate: (page: "verify" | "issue") => void;
}

export function Landing({ onNavigate }: LandingProps) {
  return (
    <div className="min-h-screen bg-gradient-to-b from-background to-muted/20">
      {/* Hero Section */}
      <div className="container mx-auto px-4 py-16">
        <div className="max-w-4xl mx-auto text-center space-y-6">
          <div className="inline-flex items-center justify-center w-16 h-16 bg-primary text-primary-foreground rounded-2xl mb-4">
            <Shield className="w-8 h-8" />
          </div>
          <h1 className="text-4xl md:text-5xl font-bold tracking-tight">
            Veritas
          </h1>
          <p className="text-xl text-muted-foreground max-w-2xl mx-auto">
            Blockchain-powered document verification system for licenses and
            credentials. Secure, transparent, and tamper-proof verification on
            Solana.
          </p>

          <div className="flex flex-col sm:flex-row gap-4 justify-center pt-4">
            <Button
              size="lg"
              onClick={() => onNavigate("verify")}
              className="text-base"
            >
              <Search className="w-5 h-5 mr-2" />
              Verify License
            </Button>
            <Button
              size="lg"
              variant="outline"
              onClick={() => onNavigate("issue")}
              className="text-base"
            >
              <FileCheck className="w-5 h-5 mr-2" />
              Issue License
            </Button>
          </div>
        </div>

        {/* Features */}
        <div className="grid md:grid-cols-3 gap-6 mt-20 max-w-5xl mx-auto">
          <Card>
            <CardContent className="pt-6 space-y-3">
              <div className="w-12 h-12 bg-primary/10 text-primary rounded-lg flex items-center justify-center">
                <Lock className="w-6 h-6" />
              </div>
              <h3 className="font-semibold">Blockchain Security</h3>
              <p className="text-sm text-muted-foreground">
                All licenses are stored on Solana blockchain, ensuring
                immutability and transparency of verification records.
              </p>
            </CardContent>
          </Card>

          <Card>
            <CardContent className="pt-6 space-y-3">
              <div className="w-12 h-12 bg-primary/10 text-primary rounded-lg flex items-center justify-center">
                <Search className="w-6 h-6" />
              </div>
              <h3 className="font-semibold">Instant Verification</h3>
              <p className="text-sm text-muted-foreground">
                Verify any license in seconds by entering the license details.
                Real-time status checks against on-chain data.
              </p>
            </CardContent>
          </Card>

          <Card>
            <CardContent className="pt-6 space-y-3">
              <div className="w-12 h-12 bg-primary/10 text-primary rounded-lg flex items-center justify-center">
                <FileCheck className="w-6 h-6" />
              </div>
              <h3 className="font-semibold">Trusted Issuers</h3>
              <p className="text-sm text-muted-foreground">
                Only authorized institutions can issue licenses. Each issuer is
                verified and registered on-chain.
              </p>
            </CardContent>
          </Card>
        </div>

        {/* Supported License Types */}
        <div className="mt-16 max-w-3xl mx-auto">
          <h2 className="text-center text-2xl font-semibold mb-8">
            Supported License Types
          </h2>
          <div className="grid sm:grid-cols-3 gap-4">
            {["Medical Licenses", "Legal Licenses", "Driving Licenses"].map(
              (type) => (
                <div
                  key={type}
                  className="p-4 border border-border rounded-lg text-center bg-card"
                >
                  <p className="font-medium">{type}</p>
                </div>
              )
            )}
          </div>
        </div>
      </div>
    </div>
  );
}
