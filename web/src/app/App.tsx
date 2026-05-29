import { useState } from "react";
import { Landing } from "./pages/Landing";
import { VerifierDashboard } from "./pages/VerifierDashboard";
import { IssuerDashboard } from "./pages/IssuerDashboard";

type Page = "landing" | "verify" | "issue";

export default function App() {
  const [currentPage, setCurrentPage] = useState<Page>("landing");

  const handleNavigate = (page: Page) => {
    setCurrentPage(page);
  };

  return (
    <div className="size-full">
      {currentPage === "landing" && <Landing onNavigate={handleNavigate} />}
      {currentPage === "verify" && <VerifierDashboard onNavigate={handleNavigate} />}
      {currentPage === "issue" && <IssuerDashboard onNavigate={handleNavigate} />}
    </div>
  );
}