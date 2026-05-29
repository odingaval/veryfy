// src/app/pages/Landing.tsx
import { Shield, Search, FileCheck, ChevronRight } from "lucide-react";
import { Button } from "../components/ui/button";

export function Landing({ onNavigate }: { onNavigate: (page: "verify" | "issue" | "landing") => void }) {
  return (
    <div className="min-h-screen relative flex flex-col items-center justify-start pt-24 pb-12 overflow-hidden">
      {/* Hero */}
      <section className="container mx-auto px-4 text-center max-w-4xl animate-float">
        <div className="inline-flex items-center justify-center p-4 bg-white/5 rounded-2xl mb-6 shadow-lg border border-white/10">
          <Shield className="w-12 h-12 text-indigo-400" />
        </div>
        <h1 className="text-6xl md:text-8xl font-extrabold text-white drop-shadow-lg">
          <span className="text-gradient">Trust Every Organization</span>
        </h1>
        <p className="mt-6 text-xl text-white/70 max-w-2xl mx-auto">
          Blockchain‑powered credential verification that is secure, transparent and instantly verifiable.
        </p>
        <div className="mt-8 flex flex-col sm:flex-row gap-4 justify-center">
          <Button className="button-premium flex items-center group" onClick={() => onNavigate("verify")}> 
            Verify License <ChevronRight className="ml-2 w-5 h-5 group-hover:translate-x-1 transition-transform"/>
          </Button>
          <Button className="button-premium-outline flex items-center" onClick={() => onNavigate("issue")}> 
            Issue Credential <FileCheck className="ml-2 w-5 h-5"/>
          </Button>
        </div>
      </section>

      {/* How It Works */}
      <section className="container mx-auto px-4 mt-24 max-w-5xl">
        <h2 className="text-3xl font-bold text-white mb-8 text-center">How It Works</h2>
        <div className="grid md:grid-cols-3 gap-8">
          {["Issue", "Scan", "Validate"].map((step, i) => (
            <div key={i} className="glass-panel p-6 text-center hover:scale-105 transition-transform">
              <div className="mb-4 flex items-center justify-center w-14 h-14 mx-auto bg-indigo-500/20 rounded-full">
                {step === "Issue" && <Shield className="w-6 h-6 text-white"/>}
                {step === "Scan" && <Search className="w-6 h-6 text-white"/>}
                {step === "Validate" && <FileCheck className="w-6 h-6 text-white"/>}
              </div>
              <h3 className="text-xl font-semibold text-white mb-2">{step}</h3>
              <p className="text-white/60 text-sm">{step === "Issue" ? "Authorities issue tamper‑proof licenses on‑chain" : step === "Scan" ? "Anyone scans the QR code with a mobile device" : "The blockchain instantly verifies authenticity"}</p>
            </div>
          ))}
        </div>
      </section>

      {/* Trust Metrics */}
      <section className="container mx-auto px-4 mt-24 max-w-5xl">
        <h2 className="text-3xl font-bold text-white mb-8 text-center">Trust Metrics</h2>
        <div className="grid md:grid-cols-4 gap-6">
          {[
            {label: "Licenses Issued", value: "12,340", icon: <Shield className="w-5 h-5 text-indigo-400"/>},
            {label: "Daily Verifications", value: "8,721", icon: <Search className="w-5 h-5 text-emerald-400"/>},
            {label: "Avg. Verification Time", value: "0.8 s", icon: <FileCheck className="w-5 h-5 text-emerald-400"/>},
            {label: "Revocation Rate", value: "0.4%", icon: <ChevronRight className="w-5 h-5 text-amber-400"/>},
          ].map((item, i) => (
            <div key={i} className="glass-panel p-6 text-center">
              <div className="mb-3 flex items-center justify-center text-2xl">{item.icon}</div>
              <p className="text-3xl font-bold text-white">{item.value}</p>
              <p className="text-white/60 mt-1">{item.label}</p>
            </div>
          ))}
        </div>
      </section>

      {/* Footer */}
      <footer className="w-full mt-24 py-8 bg-[#030213]/70 text-white/50 text-center">
        © {new Date().getFullYear()} Veryfy Protocol. Powered by <span className="text-purple-400 font-semibold">Solana</span>.
      </footer>
    </div>
  );
}
