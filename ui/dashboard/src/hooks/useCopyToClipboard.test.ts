import useCopyToClipboard from "./useCopyToClipboard";
import { act, renderHook } from "@testing-library/react";

// These tests run the real copy-to-clipboard package and stub only the browser
// APIs it calls: the async Clipboard API, execCommand and window.prompt.

const originalClipboard = Object.getOwnPropertyDescriptor(
  navigator,
  "clipboard",
);
const originalSecureContext = Object.getOwnPropertyDescriptor(
  window,
  "isSecureContext",
);
const originalExecCommand = document.execCommand;
const originalPrompt = window.prompt;

const stubBrowser = ({
  writeText,
  execCommand,
}: {
  writeText: jest.Mock;
  execCommand: jest.Mock;
}) => {
  Object.defineProperty(window, "isSecureContext", {
    configurable: true,
    value: true,
  });
  Object.defineProperty(navigator, "clipboard", {
    configurable: true,
    value: { writeText },
  });
  document.execCommand = execCommand;
  const prompt = jest.fn();
  window.prompt = prompt;
  return { prompt };
};

afterEach(() => {
  if (originalClipboard) {
    Object.defineProperty(navigator, "clipboard", originalClipboard);
  } else {
    // @ts-ignore
    delete navigator.clipboard;
  }
  if (originalSecureContext) {
    Object.defineProperty(window, "isSecureContext", originalSecureContext);
  } else {
    // @ts-ignore
    delete window.isSecureContext;
  }
  document.execCommand = originalExecCommand;
  window.prompt = originalPrompt;
});

describe("useCopyToClipboard", () => {
  test("shows success once the clipboard write resolves", async () => {
    const writeText = jest.fn().mockResolvedValue(undefined);
    const { prompt } = stubBrowser({ writeText, execCommand: jest.fn() });
    const { result } = renderHook(() => useCopyToClipboard());

    await act(async () => {
      await result.current.copy("some text");
    });

    expect(writeText).toHaveBeenCalledWith("some text");
    expect(prompt).not.toHaveBeenCalled();
    expect(result.current.copySuccess).toBe(true);
  });

  test("shows no success and offers the prompt when every copy path fails", async () => {
    const writeText = jest.fn().mockRejectedValue(new Error("denied"));
    const execCommand = jest.fn().mockReturnValue(false);
    const { prompt } = stubBrowser({ writeText, execCommand });
    const { result } = renderHook(() => useCopyToClipboard());

    await act(async () => {
      await result.current.copy("some text");
    });

    expect(writeText).toHaveBeenCalledWith("some text");
    expect(execCommand).toHaveBeenCalledWith("copy");
    expect(prompt).toHaveBeenCalledWith(expect.any(String), "some text");
    expect(result.current.copySuccess).toBe(false);
  });
});
