// Ensure Table is loaded & registered first
import "@powerpipe/components/dashboards/Table";
import Card, { CardProps } from "@powerpipe/components/dashboards/Card";
import CustomizeViewSummary from "../CustomizeViewSummary";
import DashboardTitle from "@powerpipe/components/dashboards/titles/DashboardTitle";
import DetectionGrouping from "../DetectionGrouping";
import Error from "@powerpipe/components/dashboards/Error";
import FilterCardWrapper from "@powerpipe/components/dashboards/grouping/FilterCardWrapper";
import Grid from "@powerpipe/components/dashboards/layout/Grid";
import Panel from "@powerpipe/components/dashboards/layout/Panel";
import PanelControls from "@powerpipe/components/dashboards/layout/Panel/PanelControls";
import useDownloadDetectionBenchmarkData from "@powerpipe/hooks/useDownloadDetectionBenchmarkData";
import useFilterConfig from "@powerpipe/hooks/useFilterConfig";
import { CardType } from "@powerpipe/components/dashboards/data/CardDataProcessor";
import { DateRangePicker } from "@powerpipe/components/dashboards/inputs/DateRangePickerInput";
import { default as DetectionBenchmarkType } from "../common/DetectionBenchmark";
import {
  DetectionBenchmarkTreeProps,
  DetectionDisplayGroup,
  DetectionNode,
  DetectionSummary,
} from "@powerpipe/components/dashboards/grouping/common";
import {
  GroupingProvider,
  useDetectionGrouping,
} from "@powerpipe/hooks/useDetectionGrouping";
import { noop } from "@powerpipe/utils/func";
import {
  PanelControlsProvider,
  usePanelControls,
} from "@powerpipe/hooks/usePanelControls";
import { PanelDefinition } from "@powerpipe/types";
import { registerComponent } from "@powerpipe/components/dashboards";
import { TableViewWrapper as Table } from "@powerpipe/components/dashboards/Table";
import { useDashboardPanelDetail } from "@powerpipe/hooks/useDashboardPanelDetail";
import { useEffect, useMemo, useState } from "react";
import { Width } from "@powerpipe/components/dashboards/common";

type DetectionBenchmarkTableViewProps = {
  benchmark: DetectionBenchmarkType;
  definition: PanelDefinition;
};

type InnerCheckProps = {
  benchmark: DetectionBenchmarkType;
  definition: PanelDefinition;
  grouping: DetectionNode;
  groupingConfig: DetectionDisplayGroup[];
  firstChildSummaries: DetectionSummary[];
  showControls: boolean;
  withTitle: boolean;
};

const DetectionBenchmark = (props: InnerCheckProps) => {
  const {
    filter: { expressions },
  } = useFilterConfig(props.definition?.name);
  const [referenceElement, setReferenceElement] = useState(null);
  const [showBenchmarkControls, setShowBenchmarkControls] = useState(false);
  const { panelControls: benchmarkControls, setCustomControls } =
    usePanelControls();
  const { download, processing } = useDownloadDetectionBenchmarkData(
    props.benchmark,
  );
  const { selectSidePanel, selectedPanel } = useDashboardPanelDetail();

  useEffect(() => {
    setCustomControls([
      {
        key: "filter-and-group",
        title: "Filter & Group",
        component: <CustomizeViewSummary panelName={props.definition.name} />,
        action: async () =>
          selectSidePanel({
            panel: props.definition,
          }),
      },
      {
        key: "download-data",
        disabled: processing || props.grouping.status !== "complete",
        title: "Download data",
        icon: "arrow-down-tray",
        action: download,
      },
    ]);
  }, [
    processing,
    props.grouping.status,
    props.definition.name,
    setCustomControls,
  ]);

  const summaryCards = useMemo(() => {
    if (!props.grouping) {
      return [];
    }

    const totalSummary = props.firstChildSummaries.reduce(
      (cumulative, current) => {
        cumulative.total += current.total;
        return cumulative;
      },
      { total: 0 },
    );

    const summary_cards = [
      {
        name: `${props.definition.name}.container.summary.total`,
        width: 2,
        display_type:
          totalSummary.total > 0
            ? "alert"
            : props.grouping.status === "complete"
              ? "ok"
              : null,
        properties: {
          loading:
            totalSummary.total === 0 && props.grouping.status === "running",
          label: "Total",
          value: totalSummary.total,
          icon:
            totalSummary.total > 0
              ? "materialsymbols-solid:circle_notifications"
              : props.grouping.status === "complete"
                ? "materialsymbols-solid:check_circle"
                : null,
        },
      },
    ];

    const severity_summary = props.grouping.severity_summary;
    const criticalRaw = severity_summary["critical"];
    const highRaw = severity_summary["high"];
    const mediumRaw = severity_summary["medium"];
    const lowRaw = severity_summary["low"];
    const critical = criticalRaw || 0;
    const high = highRaw || 0;
    const medium = mediumRaw || 0;
    const low = lowRaw || 0;

    // If we have at least 1 critical result
    const criticalTotal = critical;
    summary_cards.push({
      name: `${props.definition.name}.container.summary.severity.critical`,
      width: 2,
      display_type: criticalTotal > 0 ? "severity-critical" : "",
      properties: {
        loading: criticalTotal === 0 && props.grouping.status === "running",
        label: "Critical",
        value: criticalTotal,
        icon: "materialsymbols-solid:pulse_alert",
      },
    });

    // If we have at least 1 high result
    const highTotal = high;
    summary_cards.push({
      name: `${props.definition.name}.container.summary.severity.high`,
      width: 2,
      display_type: highTotal > 0 ? "severity-high" : "",
      properties: {
        loading: highTotal === 0 && props.grouping.status === "running",
        label: "High",
        value: highTotal,
        icon: "materialsymbols-solid:warning",
      },
    });

    // If we have at least 1 medium result
    const mediumTotal = medium;
    summary_cards.push({
      name: `${props.definition.name}.container.summary.severity.medium`,
      width: 2,
      display_type: mediumTotal > 0 ? "severity-medium" : "",
      properties: {
        loading: mediumTotal === 0 && props.grouping.status === "running",
        label: "Medium",
        value: mediumTotal,
        icon: "materialsymbols-solid:campaign",
      },
    });

    // If we have at least 1 low result
    const lowTotal = low;
    summary_cards.push({
      name: `${props.definition.name}.container.summary.severity.low`,
      width: 2,
      display_type: lowTotal > 0 ? "severity-low" : "",
      properties: {
        loading: lowTotal === 0 && props.grouping.status === "running",
        label: "Low",
        value: lowTotal,
        icon: "materialsymbols-solid:info",
      },
    });

    return summary_cards;
  }, [props.firstChildSummaries, props.grouping, props.definition.name]);

  if (!props.grouping) {
    return null;
  }

  return (
    <Grid
      name={props.definition.name}
      width={props.definition.width}
      events={{
        onMouseEnter: props.showControls
          ? () => setShowBenchmarkControls(true)
          : noop,
        onMouseLeave: () => setShowBenchmarkControls(false),
      }}
      setRef={setReferenceElement}
    >
      {/*Don't show when in panel detail view*/}
      {!selectedPanel && (
        <DashboardTitle
          title={props.definition.title}
          controls={
            showBenchmarkControls ? (
              <PanelControls
                referenceElement={referenceElement}
                controls={benchmarkControls}
              />
            ) : null
          }
        />
      )}
      <Grid name="temp">
        <DateRangePicker
          name="input.detection_range"
          panel_type="input"
          properties={{
            name: "detection_range",
            unqualified_name: "input.detection_range",
            type: "text",
          }}
        />
      </Grid>
      <Grid name={`${props.definition.name}.container.summary`}>
        {summaryCards
          .filter(({ name }) => {
            // Always include the total card
            if (name === `${props.definition.name}.container.summary.total`) {
              return true;
            }
            const severityFilter = expressions?.find(
              (expr) => expr.type === "severity",
            );
            const severityType = name.split(".")[name.split(".").length - 1];
            if (severityFilter && severityFilter.operator === "equal") {
              return severityType === severityFilter.value;
            } else if (
              severityFilter &&
              severityFilter.operator === "not_equal"
            ) {
              return severityType !== severityFilter.value;
            } else if (severityFilter && severityFilter.operator === "in") {
              return severityFilter.value?.includes(severityType);
            } else if (severityFilter && severityFilter.operator === "not_in") {
              return !severityFilter.value?.includes(severityType);
            }
            return true;
          })
          .map((summaryCard) => {
            const cardProps: CardProps = {
              name: summaryCard.name,
              dashboard: props.definition.dashboard,
              display_type: summaryCard.display_type as CardType,
              panel_type: "card",
              properties: summaryCard.properties,
              status: "complete",
              width: summaryCard.width as Width,
            };
            return (
              <Panel
                key={summaryCard.name}
                definition={cardProps}
                parentType="benchmark"
                showControls={false}
              >
                <FilterCardWrapper
                  cardName={summaryCard.name}
                  panelName={props.definition.name}
                  dimension="severity"
                  expressions={expressions}
                >
                  <Card {...cardProps} />
                </FilterCardWrapper>
              </Panel>
            );
          })}
      </Grid>
      <Grid name={`${props.definition.name}.container.tree`}>
        <DetectionTree
          name={`${props.definition.name}.container.tree.results`}
          dashboard={props.definition.dashboard}
          panel_type="benchmark_tree"
          properties={{
            grouping: props.grouping,
            first_child_summaries: props.firstChildSummaries,
          }}
          status="complete"
        />
      </Grid>
    </Grid>
  );
};

const DetectionTree = (props: DetectionBenchmarkTreeProps) => {
  if (!props.properties || !props.properties.first_child_summaries) {
    return null;
  }

  return <DetectionGrouping node={props.properties.grouping} />;
};

const DetectionBenchmarkTableView = ({
  benchmark,
  definition,
}: DetectionBenchmarkTableViewProps) => {
  const benchmarkDataTable = useMemo(
    () => benchmark.get_data_table(),
    [benchmark],
  );

  return (
    <Panel
      definition={{
        name: definition.name,
        dashboard: definition.dashboard,
        panel_type: "table",
        width: definition.width,
        children: definition.children,
        data: benchmarkDataTable,
        status: benchmarkDataTable ? "complete" : "running",
      }}
      parentType="benchmark"
    >
      <Table
        name={`${definition.name}.table`}
        panel_type="table"
        data={benchmarkDataTable}
        filterEnabled
      />
    </Panel>
  );
};

const Inner = ({ showControls, withTitle }) => {
  const {
    benchmark,
    definition,
    grouping,
    groupingConfig,
    firstChildSummaries,
  } = useDetectionGrouping();

  if (!definition || !benchmark || !grouping) {
    return null;
  }

  if (!definition.display_type || definition.display_type === "benchmark") {
    return (
      <DetectionBenchmark
        benchmark={benchmark}
        definition={definition}
        grouping={grouping}
        groupingConfig={groupingConfig}
        firstChildSummaries={firstChildSummaries}
        showControls={showControls}
        withTitle={withTitle}
      />
    );
    // @ts-ignore
  } else if (definition.display_type === "table") {
    return (
      <DetectionBenchmarkTableView
        benchmark={benchmark}
        definition={definition}
      />
    );
  } else {
    return (
      <Panel
        definition={{
          name: definition.name,
          dashboard: definition.dashboard,
          panel_type: "benchmark",
          width: definition.width,
          status: "error",
        }}
        parentType="benchmark"
      >
        <Error
          error={`Unsupported benchmark type ${definition.display_type}`}
        />
      </Panel>
    );
  }
};

type DetectionBenchmarkWrapperProps = {
  definition: PanelDefinition;
  benchmarkChildren?: PanelDefinition[] | undefined;
  showControls: boolean;
  withTitle: boolean;
};

const DetectionBenchmarkWrapper = (props: DetectionBenchmarkWrapperProps) => {
  return (
    <GroupingProvider
      definition={props.definition}
      benchmarkChildren={props.benchmarkChildren}
    >
      <PanelControlsProvider
        definition={props.definition}
        enabled={props.showControls}
      >
        <Inner showControls={props.showControls} withTitle={props.withTitle} />
      </PanelControlsProvider>
    </GroupingProvider>
  );
};

registerComponent("detection_benchmark", DetectionBenchmarkWrapper);

export default DetectionBenchmarkWrapper;
