/// <reference types="vite/client" />

interface ImportMetaEnv {
  /** URL of the backend API (fallback is http://localhost:8080) */
  readonly VITE_API_URL?: string;
  // add other VITE_* variables here as needed
}

interface ImportMeta {
  readonly env: ImportMetaEnv;
}
