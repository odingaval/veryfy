import { useMemo } from 'react';
import { ConnectionProvider, WalletProvider } from '@solana/wallet-adapter-react';
import { WalletModalProvider } from '@solana/wallet-adapter-react-ui';
import '@solana/wallet-adapter-react-ui/styles.css';

export function SolanaProvider({ children }: { children: React.ReactNode }) {
  // Use local test validator for development
  const endpoint = useMemo(() => import.meta.env.VITE_SOLANA_RPC_URL || "http://127.0.0.1:8899", []);
  
  const wallets = useMemo(() => [
    // Add wallets here if you want to support specific ones,
    // though the standard wallet adapter handles modern standard wallets automatically
  ], []);

  return (
    <ConnectionProvider endpoint={endpoint}>
      <WalletProvider wallets={wallets} autoConnect>
        <WalletModalProvider>{children}</WalletModalProvider>
      </WalletProvider>
    </ConnectionProvider>
  );
}
