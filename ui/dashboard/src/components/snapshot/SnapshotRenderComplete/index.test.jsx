import React from "react";
import SnapshotRenderComplete from "./index.tsx";
import "@testing-library/jest-dom";
import {
  DashboardDataModeCLISnapshot,
  DashboardDataModeCloudSnapshot,
  DashboardDataModeLive,
} from "@powerpipe/types";
import { DashboardContext } from "@powerpipe/hooks/useDashboardState";
import { render } from "@testing-library/react";

test("return null when data mode is CLI snapshot", async () => {
  // ARRANGE
  const { container } = render(
    <DashboardContext.Provider
      value={{ dataMode: DashboardDataModeCLISnapshot }}
    >
      <SnapshotRenderComplete />
    </DashboardContext.Provider>,
  );

  // ASSERT
  expect(container).toBeEmptyDOMElement();
});

test("return null when data mode is Cloud snapshot", async () => {
  // ARRANGE
  const { container } = render(
    <DashboardContext.Provider
      value={{ dataMode: DashboardDataModeCloudSnapshot }}
    >
      <SnapshotRenderComplete />
    </DashboardContext.Provider>,
  );

  // ASSERT
  expect(container).toBeEmptyDOMElement();
});

test("return null when data mode is live, but state not complete", async () => {
  // ARRANGE
  const { container } = render(
    <DashboardContext.Provider
      value={{ dataMode: DashboardDataModeLive, state: "running" }}
    >
      <SnapshotRenderComplete />
    </DashboardContext.Provider>,
  );

  // ASSERT
  expect(container).toBeEmptyDOMElement();
});

test("return div when data mode is live and state is complete", async () => {
  // ARRANGE
  render(
    <DashboardContext.Provider
      value={{ dataMode: DashboardDataModeLive, state: "complete" }}
    >
      <SnapshotRenderComplete />
    </DashboardContext.Provider>,
  );

  // ASSERT
  expect(document.querySelector("#snapshot-complete")).toBeTruthy();
});
