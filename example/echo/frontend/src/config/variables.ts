type w = Window & typeof globalThis & {import: {meta: {env: DynVars}}};

interface DynVars {
  BE_ORIGIN: string;
  APP_ROOT: string;
  MODE: string;
}

const dynVars: DynVars = {
  BE_ORIGIN:
    import.meta.env.VITE_BE_ORIGIN ||
    (window as w)['import'].meta.env.BE_ORIGIN,
  APP_ROOT:
    import.meta.env.VITE_APP_ROOT || (window as w)['import'].meta.env.APP_ROOT,
  MODE: import.meta.env.MODE || (window as w)['import'].meta.env.MODE,
};

// change them as you see fit!
const staticVars = {
  TITLE: 'AstronLab - Code your dream!',
  DESCRIPTION: 'We help you to build your tech dream. Sky is the limit!',
  API_VER_PATH: '/api/v1',
  PROD_SITE: 'https://astronlab.com/', // change this to actual production site
  SHORT_NAME: 'AstronLab', // maximum of 12 characters
  OG_IMG_URL: 'https://static.astronlab.com/projects/astronlab/equal-ratio.svg',
  THEME_COLOR: '#06c',
  OUT_DIR: './dist',
};

export const vars = {
  env: dynVars,
  ...staticVars,
  BE_BASE: dynVars.BE_ORIGIN + dynVars.APP_ROOT + 'dashboard', // this could be a different url depending deployment. i.e. https://dash.beurl.com
  NAME: staticVars.TITLE, // maximum of 45 characters
  OG_TITLE: staticVars.TITLE + ' | ' + staticVars.DESCRIPTION, // note: this should be different from the default title; ref: https://twitter.com/jon_neal/status/1428721238071988237
  OG_IMG_ALT: staticVars.DESCRIPTION,
  IMG_URL_512: dynVars.APP_ROOT + '_gen/icon-512x512.png', // this will be generated automatically in the Layout.astro
};
