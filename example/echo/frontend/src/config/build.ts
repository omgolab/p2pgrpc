/* eslint-disable node/no-unpublished-import */
import {vars} from './variables';
import type {ManifestOptions} from 'vite-plugin-pwa';

/**
 * Defines the configuration for PWA webmanifest.
 */
export const manifest: Partial<ManifestOptions> = {
  name: vars.NAME,
  short_name: vars.SHORT_NAME,
  start_url: vars.env.APP_ROOT,
  icons: [
    {
      src: vars.OG_IMG_URL,
      sizes: 'any',
      type: 'image/svg+xml',
      purpose: 'any',
    },
    {
      src: vars.IMG_URL_512,
      sizes: '512x512',
      type: 'image/png',
      purpose: 'maskable',
    },
  ],
  display: 'standalone',
  scope: vars.env.APP_ROOT,
  theme_color: vars.THEME_COLOR,
  shortcuts: [
    {
      name: vars.NAME,
      short_name: vars.SHORT_NAME,
      description: vars.DESCRIPTION,
      url: vars.env.APP_ROOT,
      icons: [
        {
          src: vars.IMG_URL_512,
          sizes: '512x512',
          type: 'image/png',
          purpose: 'any',
        },
      ],
    },
  ],
  description: vars.DESCRIPTION,
};
