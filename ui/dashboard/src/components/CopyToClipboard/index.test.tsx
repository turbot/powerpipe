import CopyToClipboard from "./index";
import { act, fireEvent, render, screen } from "@testing-library/react";

// The icon set loads lazily from a context this test does not need; render
// each icon as a plain element named after the icon it would show.
jest.mock("@powerpipe/components/Icon", () => ({
  __esModule: true,
  default: ({ icon, onClick }) => (
    <span data-testid={icon} onClick={onClick} />
  ),
}));

// The real copy-to-clipboard package runs; only the browser APIs it calls are
// stubbed: the async Clipboard API, execCommand and window.prompt.

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

const COPY_ICON = "content_copy";
const SUCCESS_ICON = "materialsymbols-solid:content_copy";

describe("CopyToClipboard", () => {
  test("shows the success icon once the clipboard write resolves", async () => {
    const writeText = jest.fn().mockResolvedValue(undefined);
    const { prompt } = stubBrowser({ writeText, execCommand: jest.fn() });
    render(<CopyToClipboard data="some text" />);

    await act(async () => {
      fireEvent.click(screen.getByTestId(COPY_ICON));
    });

    expect(writeText).toHaveBeenCalledWith("some text");
    expect(prompt).not.toHaveBeenCalled();
    expect(screen.getByTestId(SUCCESS_ICON)).toBeTruthy();
    expect(screen.queryByTestId(COPY_ICON)).toBeNull();
  });

  test("keeps the copy icon and offers the prompt when every copy path fails", async () => {
    const writeText = jest.fn().mockRejectedValue(new Error("denied"));
    const execCommand = jest.fn().mockReturnValue(false);
    const { prompt } = stubBrowser({ writeText, execCommand });
    render(<CopyToClipboard data="some text" />);

    await act(async () => {
      fireEvent.click(screen.getByTestId(COPY_ICON));
    });

    expect(writeText).toHaveBeenCalledWith("some text");
    expect(execCommand).toHaveBeenCalledWith("copy");
    expect(prompt).toHaveBeenCalledWith(expect.any(String), "some text");
    expect(screen.queryByTestId(SUCCESS_ICON)).toBeNull();
    expect(screen.getByTestId(COPY_ICON)).toBeTruthy();
  });
});
