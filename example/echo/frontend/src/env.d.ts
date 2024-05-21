/// <reference types="astro/client" />
/// <reference types="vite-plugin-pwa/info" />
/// <reference types="vite-plugin-pwa/client" />

interface ImportMetaEnv {
  VITE_APP_ROOT: string;
  VITE_BE_ORIGIN: string;
}

interface ImportMeta {
  readonly env: ImportMetaEnv;
}

declare module '*.svg' {
  const content: string;
}

declare module '*.json' {
  const content: unknown;
  export default content;
}
