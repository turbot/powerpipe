const CircularDependencyPlugin = require("circular-dependency-plugin");
const path = require("path");
const set = require("lodash/set");
const webpack = require("webpack");

module.exports = {
  jest: {
    configure: {
      moduleNameMapper: {
        "^@powerpipe/(.*)$": "<rootDir>/src/$1",
      },
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
