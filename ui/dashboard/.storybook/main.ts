const path = require("path");
const set = require("lodash/set");

module.exports = {
  stories: ["../src/**/*.mdx", "../src/**/*.stories.@(js|jsx|ts|tsx)"],

  addons: [
    "@storybook/addon-links",
    "@storybook/preset-create-react-app",
    // Temporarily disabled for Storybook 9 - incompatible with @storybook/types import
    // "storybook-dark-mode",
    // Temporarily disabled for Storybook 9 - incompatible with @storybook/preview-api import
    // "storybook-addon-react-router-v6",
    "@chromatic-com/storybook",
    "@storybook/addon-docs"
  ],

  typescript: {
    check: false,
    checkOptions: {},
    reactDocgen: "react-docgen-typescript",
    reactDocgenTypescriptOptions: {
      shouldExtractLiteralValuesFromEnum: true,
      propFilter: (prop) =>
        prop.parent ? !/node_modules/.test(prop.parent.fileName) : true,
    },
  },

  webpackFinal: async (config) => {
    if (!config.resolve) {
      config.resolve = {};
    }
    config.resolve.alias = {
      ...config.resolve.alias,
      "@powerpipe": path.resolve(__dirname, "../src"),
    };
    // Webpack 5 no longer includes Node.js polyfills by default
    config.resolve.fallback = {
      ...config.resolve.fallback,
      fs: false,
      path: false,
    };
    return config;
  },

  framework: {
    name: "@storybook/react-webpack5",
    options: {},
  },

  docs: {},
};
