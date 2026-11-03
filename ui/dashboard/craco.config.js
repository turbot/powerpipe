const CircularDependencyPlugin = require("circular-dependency-plugin");
const path = require("path");
const set = require("lodash/set");
const webpack = require("webpack");

module.exports = {
  jest: {
    // A function, not an object: craco appends an object's arrays to CRA's
    // defaults, and we need to replace transformIgnorePatterns.
    configure: (jestConfig) => {
      jestConfig.moduleNameMapper = {
        ...jestConfig.moduleNameMapper,
        "^@powerpipe/(.*)$": "<rootDir>/src/$1",
        // Jest 27 does not read package.json "exports"
        "^react-router-dom$":
          "<rootDir>/node_modules/react-router-dom/dist/index.js",
        "^react-router/dom$":
          "<rootDir>/node_modules/react-router/dist/development/dom-export.js",
      };
      // color-convert v3 and its color-name dependency ship as ES modules
      // only, so Jest has to run them through Babel like our own source.
      // Every other package in node_modules is still left untransformed.
      jestConfig.transformIgnorePatterns = [
        "[/\\\\]node_modules[/\\\\](?!(color-convert|color-name)[/\\\\]).+\\.(js|jsx|mjs|cjs|ts|tsx)$",
        "^.+\\.module\\.(css|sass|scss)$",
      ];
      return jestConfig;
    },
  },
  webpack: {
    alias: {
      "@powerpipe": path.resolve(__dirname, "src"),
    },
    configure: (webpackConfig) => {
      // Enable WebAssembly support
      webpackConfig.experiments = {
        asyncWebAssembly: true,
        ...webpackConfig.experiments,
      };

      // Exclude .wasm files from the file-loader
      const fileLoader = webpackConfig.module.rules.find((rule) =>
        rule.oneOf?.some((r) => r.loader?.includes("file-loader")),
      );
      if (fileLoader) {
        fileLoader.oneOf.forEach((r) => {
          if (r.loader && r.loader.includes("file-loader")) {
            r.exclude.push(/\.wasm$/);
          }
        });
      }

      const scopePluginIndex = webpackConfig.resolve.plugins.findIndex(
        ({ constructor }) =>
          constructor && constructor.name === "ModuleScopePlugin",
      );

      webpackConfig.resolve.plugins.splice(scopePluginIndex, 1);

      webpackConfig = set(
        webpackConfig,
        "resolve.fallback.buffer",
        require.resolve("buffer"),
      );
      webpackConfig = set(
        webpackConfig,
        "resolve.fallback.crypto",
        require.resolve("crypto-browserify"),
      );
      webpackConfig = set(
        webpackConfig,
        "resolve.fallback.path",
        require.resolve("path-browserify"),
      );
      webpackConfig = set(
        webpackConfig,
        "resolve.fallback.stream",
        require.resolve("stream-browserify"),
      );
      webpackConfig = set(
        webpackConfig,
        "resolve.fallback.vm",
        require.resolve("vm-browserify"),
      );
      webpackConfig = set(webpackConfig, "resolve.fallback.fs", false);
      webpackConfig = set(webpackConfig, "experiments.asyncWebAssembly", true);
      // webpackConfig = set(webpackConfig, "experiments.syncWebAssembly", true);
      return webpackConfig;
    },
    plugins: [
      new webpack.ProvidePlugin({
        Buffer: ["buffer", "Buffer"],
        process: "process/browser.js",
      }),
      new CircularDependencyPlugin({
        // exclude detection of files based on a RegExp
        exclude: /node_modules/,
        // include specific files based on a RegExp
        include: /src/,
        // add errors to webpack instead of warnings
        failOnError: true,
        // allow import cycles that include an asyncronous import,
        // e.g. via import(/* webpackMode: "weak" */ './file.js')
        allowAsyncCycles: false,
        // set the current working directory for displaying module paths
        cwd: process.cwd(),
      }),
    ],
  },
};
