/// <reference types="vite/client" />

interface ImportMetaEnv {
  readonly VITE_POLYGON_IO_API_KEY: string;

  readonly VITE_GITHUB_NAME: string;
  readonly VITE_GITHUB_REPO: string;
}

interface ImportMeta {
  readonly env: ImportMetaEnv;
}
