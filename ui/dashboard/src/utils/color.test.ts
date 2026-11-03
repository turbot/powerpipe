import { colorToRgb, hexToRgb, ColorGenerator } from "./color";

const themeColors = { alert: "#ff0000", info: "#0000ff", ok: "#00ff00" };

describe("color", () => {
  describe("hexToRgb", () => {
    test("converts a hex colour", () => {
      expect(hexToRgb("#ff8800")).toEqual([255, 136, 0]);
    });
  });

  describe("colorToRgb", () => {
    test("converts a hex colour", () => {
      expect(colorToRgb("#ff8800", themeColors)).toEqual([255, 136, 0]);
    });

    test("converts a CSS colour name", () => {
      expect(colorToRgb("red", themeColors)).toEqual([255, 0, 0]);
    });

    test("resolves a named theme colour", () => {
      expect(colorToRgb("ok", themeColors)).toEqual([0, 255, 0]);
    });

    test("returns undefined for a name it does not know", () => {
      expect(colorToRgb("not-a-colour", themeColors)).toBeUndefined();
    });
  });

  describe("ColorGenerator", () => {
    test("produces a hex string for an ansi256 code", () => {
      const generator = new ColorGenerator(16, 0);
      expect(generator.toColorObject(16)).toEqual({
        ansi256: 16,
        hex: "#000000",
      });
      expect(generator.toColorObject(51)).toEqual({
        ansi256: 51,
        hex: "#00FFFF",
      });
    });
  });
});
