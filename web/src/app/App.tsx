import React, { useState } from "react";
import { Toaster } from "sonner";
import { Landing } from "./pages/Landing";
import { IssuerDashboard } from "./pages/IssuerDashboard";
import { PublicVerify } from "./pages/PublicVerify";
import { ShieldCheck } from "lucide-react";

export const App: React.FC = () => {
  const [page, setPage] = useState<"landing" | "issue" | "verify">("landing");

  const renderPage = () => {
    switch (page) {
      case "issue":
        return <IssuerDashboard onNavigate={setPage} />;
              case "verify":
          return <PublicVerify />;
      default:
        return <Landing onNavigate={setPage} />;
    }
  };

  return (
    <div className="min-h-screen flex flex-col font-sans transition-colors duration-300">
      <Toaster 
        position="top-right" 
        theme="dark" 
        toastOptions={{
          style: {
            background: 'rgba(15, 15, 25, 0.9)',
            border: '1px solid rgba(255,255,255,0.1)',
            backdropFilter: 'blur(10px)',
            color: 'white'
          }
        }} 
      />

      {/* Floating Glass Header */}
      <header className="fixed top-4 left-0 right-0 z-50 px-4">
        <div className="container mx-auto max-w-5xl">
          <div className="glass-panel h-16 flex items-center justify-between px-6 shadow-2xl">
            <div className="flex items-center space-x-3 cursor-pointer group" onClick={() => setPage("landing")}>
              <div className="w-10 h-10 bg-gradient-to-tr from-indigo-500 to-purple-500 rounded-xl flex items-center justify-center shadow-lg group-hover:scale-105 transition-transform">
                <ShieldCheck className="w-6 h-6 text-white" />
              </div>
              <span className="text-2xl font-bold tracking-tight text-white drop-shadow-md">Veryfy</span>
            </div>

            <nav className="hidden md:flex items-center space-x-2">
              <button
                onClick={() => setPage("landing")}
                className={`px-4 py-2 rounded-lg text-sm font-semibold transition-all ${page === "landing" ? "bg-white/10 text-white" : "text-white/60 hover:text-white hover:bg-white/5"}`}
              >
                Overview
              </button>
              <button
                onClick={() => setPage("verify")}
                className={`px-4 py-2 rounded-lg text-sm font-semibold transition-all ${page === "verify" ? "bg-white/10 text-white" : "text-white/60 hover:text-white hover:bg-white/5"}`}
              >
                Verifier
              </button>
              <button
                onClick={() => setPage("issue")}
                className={`px-4 py-2 rounded-lg text-sm font-semibold transition-all ${page === "issue" ? "bg-white/10 text-white" : "text-white/60 hover:text-white hover:bg-white/5"}`}
              >
                Issuer Dashboard
              </button>
            </nav>
          </div>
        </div>
      </header>

      {/* Main Content Area */}
      <main className="flex-1 flex flex-col pt-24">
        {renderPage()}
      </main>

      {/* Footer */}
      <footer className="py-8 text-center text-sm text-white/40 border-t border-white/5 bg-[#05050A]">
        © {new Date().getFullYear()} Veryfy Protocol. Powered by <span className="text-purple-400 font-semibold">Solana</span>.
      </footer>
    </div>
  );
};

export default App;