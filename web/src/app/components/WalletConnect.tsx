import { WalletMultiButton } from "@solana/wallet-adapter-react-ui";
import { useWallet } from "@solana/wallet-adapter-react";
import { useEffect } from "react";

interface WalletConnectProps {
  onConnect?: (publicKey: string) => void;
  onDisconnect?: () => void;
}

export function WalletConnect({ onConnect, onDisconnect }: WalletConnectProps) {
  const { publicKey, connected } = useWallet();

  useEffect(() => {
    if (connected && publicKey) {
      onConnect?.(publicKey.toString());
    } else {
      onDisconnect?.();
    }
  }, [connected, publicKey, onConnect, onDisconnect]);

  // The WalletMultiButton handles its own UI (Connect/Disconnect dropdown)
  // We can wrap it slightly if needed, but it's typically fine out of the box.
  return (
    <div className="wallet-adapter-wrapper">
      <WalletMultiButton className="!bg-white !text-black hover:!bg-gray-100 !rounded-xl !font-semibold !transition-all !h-10" />
    </div>
  );
}
