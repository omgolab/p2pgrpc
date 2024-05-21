/* eslint-disable node/no-missing-import */
/* eslint-disable node/no-unpublished-import */
import sitemap from '@astrojs/sitemap';
import svelte from '@astrojs/svelte';
import tailwind from '@astrojs/tailwind';
import astroPWA from '@vite-pwa/astro';
import compress from 'astro-compress';
import critters from 'astro-critters';
import { defineConfig } from 'astro/config';

// Helper imports
import { manifest } from './src/config/build';
import { vars } from './src/config/variables';

// https://astro.build/config
const isDev = process.env.NODE_ENV !== 'production';

export default defineConfig({
  trailingSlash: 'always',
  site: vars.PROD_SITE,
  integrations: [
    svelte(),
    tailwind(),
    !isDev && sitemap(),
    !isDev && critters({
      Critters: {
        publicPath: vars.env.APP_ROOT,
        pruneSource: true,
        minimumExternalSize: 4096,
        fonts: !isDev,
        reduceInlineStyles: !isDev,
        inlineFonts: !isDev,
      }
    }),
    astroPWA({
      registerType: 'autoUpdate',
      manifest: manifest,
      workbox: {
        globPatterns: ['**/*.{js,css,html,ico,png,svg,json,json,xml,txt}'],
      },
    }),
    !isDev && compress()
  ],
  prefetch: {
    prefetchAll: false
  },
  output: 'static',
  base: vars.env.APP_ROOT,
  outDir: vars.OUT_DIR,
  vite: {
    build: {
      minify: 'terser',
      terserOptions: {
        compress: {
          drop_console: !isDev,
          drop_debugger: !isDev,
        },
      },
    },
    resolve: {
      alias: {
        '@': '/src',
        // needed for runtime
        demo_entities: '/node_modules/protos-pkg/gen/ts/entities/demo/',
        demo_services: '/node_modules/protos-pkg/gen/ts/services/demo/',
        commons_entities: '/node_modules/protos-pkg/gen/ts/entities/commons/',
      },
    },
  },
  build: {
    format: 'directory',
  },
  compressHTML: !isDev,
  server: {
    open: !isDev,
    port: isDev ? 3000 : 35000,
  },
  image: {
    service: {
      entrypoint: 'astro/assets/services/sharp',
    },
    remotePatterns: [{ hostname: '*.astronlab.com' }],
  }
});
