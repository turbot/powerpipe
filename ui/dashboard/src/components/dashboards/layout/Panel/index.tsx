import PanelStatus from "./PanelStatus";
import PanelControls from "./PanelControls";
import PanelInformation from "./PanelInformation";
import PanelProgress from "./PanelProgress";
import PanelTitle from "@powerpipe/components/dashboards/titles/PanelTitle";
import Placeholder from "@powerpipe/components/dashboards/Placeholder";
import { BaseChartProps } from "@powerpipe/components/dashboards/charts/types";
import { CardProps } from "@powerpipe/components/dashboards/Card";
import { classNames } from "@powerpipe/utils/styles";
import { DashboardPanelType, PanelDefinition } from "@powerpipe/types";
import { FlowProps } from "@powerpipe/components/dashboards/flows/types";
import { getResponsivePanelWidthClass } from "@powerpipe/utils/layout";
import { GraphProps } from "@powerpipe/components/dashboards/graphs/types";
import { HierarchyProps } from "@powerpipe/components/dashboards/hierarchies/types";
import { ImageProps } from "@powerpipe/components/dashboards/Image";
import { InputProps } from "@powerpipe/components/dashboards/inputs/types";
import { memo, useState } from "react";
import { PanelProvider, usePanel } from "@powerpipe/hooks/usePanel";
import { ReactNode } from "react";
import { registerComponent } from "@powerpipe/components/dashboards";
import { TableProps } from "@powerpipe/components/dashboards/Table";
import { TextProps } from "@powerpipe/components/dashboards/Text";
import { useDashboardPanelDetail } from "@powerpipe/hooks/useDashboardPanelDetail";
import { usePanelControls } from "@powerpipe/hooks/usePanelControls";

type PanelProps = {
  children: ReactNode;
  className?: string;
  definition:
    | BaseChartProps
    | CardProps
    | FlowProps
    | GraphProps
    | HierarchyProps
    | ImageProps
    | InputProps
    | PanelDefinition
    | TableProps
    | TextProps;
  parentType: DashboardPanelType;
  showControls?: boolean;
  showPanelError?: boolean;
  showPanelStatus?: boolean;
  forceBackground?: boolean;
};

const Panel = ({
  children,
  className,
  definition,
  showControls = true,
  showPanelError = true,
  showPanelStatus = true,
  forceBackground = false,
}: PanelProps) => {
  const { selectedPanel } = useDashboardPanelDetail();
  const { inputPanelsAwaitingValue } = usePanel();
  const [referenceElement, setReferenceElement] = useState(null);
  const baseStyles = classNames(
    "relative col-span-12",
    getResponsivePanelWidthClass(definition.width),
    "overflow-auto",
  );
  const { panelControls, showPanelControls, setShowPanelControls } =
    usePanelControls();

  if (inputPanelsAwaitingValue.length > 0) {
    return null;
  }

  const PlaceholderComponent = Placeholder.component;

  const shouldShowLoader =
    showPanelStatus &&
    definition.status !== "cancelled" &&
    definition.status !== "error" &&
    definition.status !== "complete";

  const showPanelContents =
    !definition.error || (definition.error && !showPanelError);

  return (
    <div
      // @ts-ignore
      ref={setReferenceElement}
      id={definition.name}
      className={baseStyles}
      onMouseEnter={
        !definition.properties?.embedded && showControls
          ? () => setShowPanelControls(true)
          : undefined
      }
      onMouseLeave={() => setShowPanelControls(false)}
    >
      <section
        aria-labelledby={
          definition.title ? `${definition.name}-title` : undefined
        }
        className={classNames(
          "col-span-12",
          definition?.properties?.embedded ? null : "m-0.5",
          forceBackground ||
            (definition.panel_type !== "image" &&
              definition.panel_type !== "card" &&
              definition.panel_type !== "input") ||
            ((definition.panel_type === "image" ||
              definition.panel_type === "card" ||
              definition.panel_type === "input") &&
              definition.display_type === "table")
            ? "bg-dashboard-panel print:bg-white shadow-sm rounded-md"
            : null,
          definition?.properties?.embedded ? "shadow-none rounded-none" : null,
        )}
      >
        {showPanelControls && (
          <PanelControls
            referenceElement={referenceElement}
            controls={panelControls}
            withOffset
          />
        )}
        {definition.title && (
          <div
            className={classNames(
              definition.panel_type === "input" &&
                definition.display_type !== "table" &&
                !forceBackground
                ? "pl-0 pr-2 sm:pr-4 py-2"
                : "px-4 py-4",
            )}
          >
            <PanelTitle name={definition.name} title={definition.title} />
          </div>
        )}

        <div
          className={classNames(
            "relative",
            definition.title &&
              ((definition.panel_type !== "input" &&
                // @ts-ignore
                definition.status !== "complete") ||
                (definition.panel_type !== "input" &&
                  definition.panel_type !== "table") ||
                (definition.panel_type === "table" &&
                  definition.display_type === "line"))
              ? "border-t border-divide"
              : null,
            selectedPanel ||
              (definition.panel_type === "table" &&
                definition.display_type !== "line") ||
              definition.display_type === "table"
              ? "overflow-x-auto"
              : "overflow-x-hidden",
            className,
          )}
        >
          <PanelProgress className={definition.title ? null : "rounded-t-md"} />
          <PanelInformation />
          <PlaceholderComponent
            animate={definition.status === "running"}
            ready={!shouldShowLoader}
          >
            <>
              {((showPanelError && definition.status === "error") ||
                showPanelStatus) && (
                <PanelStatus
                  definition={definition as PanelDefinition}
                  showPanelError={showPanelError}
                />
              )}
              {showPanelContents && children}
            </>
          </PlaceholderComponent>
        </div>
      </section>
    </div>
  );
};

const PanelWrapper = memo((props: PanelProps) => {
  const { children, ...rest } = props;
  return (
    <PanelProvider
      definition={props.definition}
      parentType={props.parentType}
      showControls={props.showControls}
    >
      <Panel {...rest}>{children}</Panel>
    </PanelProvider>
  );
});

registerComponent("panel", PanelWrapper);

export default PanelWrapper;
