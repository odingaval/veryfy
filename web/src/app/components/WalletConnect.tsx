import { Button } from "./ui/button";
import { Wallet, LogOut } from "lucide-react";
import { useState } from "react";

interface WalletConnectProps {
  onConnect?: (publicKey: string) => void;
  onDisconnect?: () => void;
}

export function WalletConnect({ onConnect, onDisconnect }: WalletConnectProps) {
  const [connected, setConnected] = useState(false);
  const [publicKey, setPublicKey] = useState<string | null>(null);

  const handleConnect = () => {
    // Mock connection - in real app this would use Phantom wallet adapter
    const mockPublicKey = "7xKXtg2CW87d97TXJSDpbD5jBkheTqA83TZRuJosgAsU";
    setPublicKey(mockPublicKey);
    setConnected(true);
    onConnect?.(mockPublicKey);
  };

  const handleDisconnect = () => {
    setPublicKey(null);
    setConnected(false);
    onDisconnect?.();
  };

  const truncateAddress = (address: string) => {
    return `${address.slice(0, 4)}...${address.slice(-4)}`;
  };

  if (connected && publicKey) {
    return (
      <div className="flex items-center gap-2">
        <div className="px-3 py-1.5 bg-muted rounded-md text-sm font-mono">
          {truncateAddress(publicKey)}
        </div>
        <Button
          variant="outline"
          size="sm"
          onClick={handleDisconnect}
        >
          <LogOut className="w-4 h-4 mr-2" />
          Disconnect
        </Button>
      </div>
    );
  }

  return (
    <Button onClick={handleConnect}>
      <Wallet className="w-4 h-4 mr-2" />
      Connect Wallet
    </Button>
  );
}
