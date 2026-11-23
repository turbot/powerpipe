import Modal from "./index";
import { fireEvent, render, screen } from "@testing-library/react";

// The theme hook reads the browser's colour-scheme preference; jsdom has no
// matchMedia, so answer "no preference" and leave everything else real.
beforeAll(() => {
  Object.defineProperty(window, "matchMedia", {
    writable: true,
    configurable: true,
    value: (query: string) => ({
      matches: false,
      media: query,
      onchange: null,
      addEventListener: () => {},
      removeEventListener: () => {},
      addListener: () => {},
      removeListener: () => {},
      dispatchEvent: () => false,
    }),
  });
});

const clickOutside = (element: Element) => {
  fireEvent.pointerDown(element);
  fireEvent.mouseDown(element);
  fireEvent.pointerUp(element);
  fireEvent.mouseUp(element);
  fireEvent.click(element);
};

describe("Modal", () => {
  test("shows its title and children inside a dialog", () => {
    render(
      <Modal onClose={jest.fn()} title="Delete snapshot">
        <p>Are you sure?</p>
      </Modal>,
    );

    expect(screen.getByRole("dialog")).toBeTruthy();
    expect(screen.getByText("Delete snapshot")).toBeTruthy();
    expect(screen.getByText("Are you sure?")).toBeTruthy();
  });

  test("calls onClose from the close button", () => {
    const onClose = jest.fn();
    render(
      <Modal onClose={onClose} title="Delete snapshot">
        <p>Are you sure?</p>
      </Modal>,
    );

    fireEvent.click(screen.getByText("Close").closest("button")!);

    expect(onClose).toHaveBeenCalledTimes(1);
  });

  test("a click outside the panel closes it when click-away is allowed", () => {
    const onClose = jest.fn();
    render(
      <Modal allowClickAway onClose={onClose} title="Delete snapshot">
        <p>Are you sure?</p>
      </Modal>,
    );

    clickOutside(screen.getByRole("dialog"));

    expect(onClose).toHaveBeenCalled();
  });

  test("a click outside the panel does nothing by default", () => {
    const onClose = jest.fn();
    render(
      <Modal onClose={onClose} title="Delete snapshot">
        <p>Are you sure?</p>
      </Modal>,
    );

    clickOutside(screen.getByRole("dialog"));

    expect(onClose).not.toHaveBeenCalled();
    expect(screen.getByText("Delete snapshot")).toBeTruthy();
  });

  test("a click inside the panel does not close it", () => {
    const onClose = jest.fn();
    render(
      <Modal allowClickAway onClose={onClose} title="Delete snapshot">
        <p>Are you sure?</p>
      </Modal>,
    );

    clickOutside(screen.getByText("Are you sure?"));

    expect(onClose).not.toHaveBeenCalled();
  });
});
