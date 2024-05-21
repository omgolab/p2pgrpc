/* eslint-disable node/no-unpublished-require */
/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ['./src/**/*.{astro,html,js,jsx,md,mdx,svelte,ts,tsx,vue}'],
  theme: {
    screens: {
      '2xs': '0px',
      xs: '320px',
      sm: '640px',
      md: '768px',
      lg: '1024px',
      xl: '1280px',
      '2xl': '1536px',
    },
    extend: {
      screens: {
        ls: {
          raw: '(min-width: 769px) and (max-width: 1024px) and (max-height: 540px)',
        },
        tablet: {
          raw: '(min-width: 769px) and (max-width: 1024px) and (min-height: 541px)',
        },
        smls: {
          raw: '(min-width: 640px) and (max-width: 768px) and (max-height: 540px)',
        },
        'sm-tablet': {
          raw: '(min-width: 640px) and (max-width: 768px) and (min-height: 541px)',
        },
      },

      colors: {
        'solid-purple': '#9900FF',
        'solid-orange': '#FF6633',
        'grey-background': '#E5E5E5',
        'ping-color-light': '#F2BF72',
        'ping-color-default': '#FF9900',
        'jitter-color-light': '#F29872',
        'jitter-color-default': '#FF4C00',
        'upload-color-light': '#EA72F2',
        'upload-color-default': '#F000FF',
        'download-color-light': '#BF72F2',
        'download-color-default': '#9900FF',
      },
    },
  },
  plugins: [require('daisyui'), require('@tailwindcss/typography')],

};
