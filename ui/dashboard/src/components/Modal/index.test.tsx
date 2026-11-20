import Modal from "./index";
import { act, fireEvent, render, screen } from "@testing-library/react";

// jsdom has no matchMedia and no ResizeObserver; the theme provider and the
// dialog's focus handling read them. The real Dialog and Transition run.
beforeAll(() => {
  window.matchMedia = ((query: string) => ({
    matches: false,
    media: query,
    onchange: null,
    addEventListener: () => {},
    removeEventListener: () => {},
    addListener: () => {},
    removeListener: () => {},
    dispatchEvent: () => false,
  })) as unknown as typeof window.matchMedia;
  (window as any).ResizeObserver = class {
    observe() {}
    unobserve() {}
    disconnect() {}
  };
});

describe("Modal", () => {
  test("renders the dialog with its title, body, actions and backdrop", async () => {
    await act(async () => {
      render(
        <Modal
          title="Something went wrong"
          actions={[<button key="a">Dismiss</button>]}
          onClose={async () => {}}
        >
          <p>The body text</p>
        </Modal>,
      );
    });

    const dialog = screen.getByRole("dialog");
    expect(dialog).toBeTruthy();
    expect(screen.getByText("Something went wrong").tagName).toBe("H2");
    expect(screen.getByText("The body text")).toBeTruthy();
    expect(screen.getByText("Dismiss")).toBeTruthy();
    // The backdrop is the translucent full-screen layer behind the panel.
    expect(dialog.querySelector(".fixed.inset-0.bg-gray-500")).not.toBeNull();
  });

  test("the close button calls onClose", async () => {
    const onClose = jest.fn(async () => {});
    await act(async () => {
      render(
        <Modal title="Title" onClose={onClose}>
          <p>Body</p>
        </Modal>,
      );
    });

    await act(async () => {
      fireEvent.click(screen.getByText("Close").closest("button")!);
    });

    expect(onClose).toHaveBeenCalledTimes(1);
  });

  test("a click outside the panel does not close it unless allowed", async () => {
    const onClose = jest.fn(async () => {});
    await act(async () => {
      render(
        <Modal title="Title" onClose={onClose}>
          <p>Body</p>
        </Modal>,
      );
    });

    await act(async () => {
      fireEvent.mouseDown(document.body);
      fireEvent.click(document.body);
    });

    expect(onClose).not.toHaveBeenCalled();
    expect(screen.getByRole("dialog")).toBeTruthy();
  });

  test("a click outside the panel calls onClose when allowed", async () => {
    const onClose = jest.fn(async () => {});
    await act(async () => {
      render(
        <Modal title="Title" allowClickAway onClose={onClose}>
          <p>Body</p>
        </Modal>,
      );
    });

    await act(async () => {
      fireEvent.pointerDown(document.body);
      fireEvent.mouseDown(document.body);
      fireEvent.pointerUp(document.body);
      fireEvent.mouseUp(document.body);
      fireEvent.click(document.body);
    });

    expect(onClose).toHaveBeenCalled();
  });
});
