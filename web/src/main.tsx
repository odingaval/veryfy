import { createRoot } from "react-dom/client";
import App from "./app/App.tsx";
import { SolanaProvider } from "./app/components/SolanaProvider";
import "./styles/index.css";

createRoot(document.getElementById("root")!).render(
  <SolanaProvider>
    <App />
  </SolanaProvider>
);