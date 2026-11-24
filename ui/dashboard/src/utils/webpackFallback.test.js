// elliptic reached the dashboard bundle through crypto-browserify, which was
// webpack's fallback for Node's `crypto`. jq-wasm is the only importer, and
// only in its Node branch, so the fallback is `false` and the polyfill is
// not a dependency. These tests keep it that way.
const cracoConfig = require("../../craco.config.js");
const packageJson = require("../../package.json");

const fallbackFor = () =>
  cracoConfig.webpack.configure({
    module: { rules: [] },
    resolve: { plugins: [{ constructor: { name: "ModuleScopePlugin" } }] },
  }).resolve.fallback;

describe("webpack fallbacks for Node modules", () => {
  test("crypto resolves to nothing, so no crypto polyfill is bundled", () => {
    expect(fallbackFor().crypto).toBe(false);
  });

  test("the other fallbacks stay polyfilled", () => {
    const fallback = fallbackFor();
    ["buffer", "path", "stream", "vm"].forEach((name) =>
      expect(typeof fallback[name]).toBe("string"),
    );
  });

  test("crypto-browserify and the pins only it needed are not declared", () => {
    const declared = {
      ...packageJson.dependencies,
      ...packageJson.devDependencies,
      ...packageJson.resolutions,
    };
    ["crypto-browserify", "elliptic", "pbkdf2"].forEach((name) =>
      expect(declared).not.toHaveProperty(name),
    );
  });
});
