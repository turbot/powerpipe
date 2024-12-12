// Ensure Table is loaded & registered first
import "@powerpipe/components/dashboards/Table";
import Card, { CardProps } from "@powerpipe/components/dashboards/Card";
import CheckGrouping from "../CheckGrouping";
import CustomizeViewSummary from "../CustomizeViewSummary";
import DashboardTitle from "@powerpipe/components/dashboards/titles/DashboardTitle";
import Error from "@powerpipe/components/dashboards/Error";
import FilterCardWrapper from "@powerpipe/components/dashboards/grouping/FilterCardWrapper";
import Grid from "@powerpipe/components/dashboards/layout/Grid";
import Panel from "@powerpipe/components/dashboards/layout/Panel";
import PanelControls from "@powerpipe/components/dashboards/layout/Panel/PanelControls";
import useFilterConfig from "@powerpipe/hooks/useFilterConfig";
import {
  BenchmarkTreeProps,
  CheckDisplayGroup,
  CheckNode,
  CheckSummary,
} from "../common";
import { CardType } from "@powerpipe/components/dashboards/data/CardDataProcessor";
import { default as BenchmarkType } from "../common/Benchmark";
import {
  getComponent,
  registerComponent,
} from "@powerpipe/components/dashboards";
import {
  GroupingProvider,
  useBenchmarkGrouping,
} from "@powerpipe/hooks/useBenchmarkGrouping";
import { noop } from "@powerpipe/utils/func";
import {
  PanelControlsProvider,
  usePanelControls,
} from "@powerpipe/hooks/usePanelControls";
import { PanelDefinition } from "@powerpipe/types";
import { useDashboardPanelDetail } from "@powerpipe/hooks/useDashboardPanelDetail";
import { useEffect, useMemo, useState } from "react";
import { Width } from "@powerpipe/components/dashboards/common";

const Table = getComponent("table");

type BenchmarkTableViewProps = {
  benchmark: BenchmarkType;
  definition: PanelDefinition;
};

type InnerCheckProps = {
  benchmark: BenchmarkType;
  definition: PanelDefinition;
  grouping: CheckNode;
  groupingConfig: CheckDisplayGroup[];
  firstChildSummaries: CheckSummary[];
  withTitle: boolean;
};

const Benchmark = (props: InnerCheckProps) => {
  const {
    filter: { expressions },
  } = useFilterConfig(props.definition?.name);
  const { selectedPanel } = useDashboardPanelDetail();
  const benchmarkDataTable = useMemo(() => {
    if (
      !props.benchmark ||
      !props.grouping ||
      props.grouping.status !== "complete"
    ) {
      return undefined;
    }
    return props.benchmark.get_data_table();
  }, [props.benchmark, props.grouping]);
  const [referenceElement, setReferenceElement] = useState(null);
  const {
    enabled: panelControlsEnabled,
    panelControls: benchmarkControls,
    showPanelControls,
    setCustomControls,
    setPanelData,
    setShowPanelControls,
  } = usePanelControls();
  const { selectFilterAndGroupPanel } = useDashboardPanelDetail();

  useEffect(() => {
    setCustomControls([
      {
        key: "filter-and-group",
        title: "Filter & Group",
        component: <CustomizeViewSummary panelName={props.definition.name} />,
        action: async () => selectFilterAndGroupPanel(props.definition.name),
      },
    ]);
  }, [props.definition.name, setCustomControls]);

  useEffect(() => {
    if (!benchmarkDataTable) {
      return;
    }
    setPanelData(benchmarkDataTable);
  }, [benchmarkDataTable, setPanelData]);

  const summaryCards = useMemo(() => {
    if (!props.grouping) {
      return [];
    }

    const totalSummary = props.firstChildSummaries.reduce(
      (cumulative, current) => {
        cumulative.error += current.error;
        cumulative.alarm += current.alarm;
        cumulative.ok += current.ok;
        cumulative.info += current.info;
        cumulative.skip += current.skip;
        return cumulative;
      },
      { error: 0, alarm: 0, ok: 0, info: 0, skip: 0 },
    );

    const summary_cards = [
      {
        name: `${props.definition.name}.container.summary.status.ok`,
        width: 2,
        display_type: totalSummary.ok > 0 ? "ok" : "skip",
        properties: {
          label: "OK",
          value: totalSummary.ok,
          icon: "materialsymbols-solid:check_circle",
        },
      },
      {
        name: `${props.definition.name}.container.summary.status.alarm`,
        width: 2,
        display_type: totalSummary.alarm > 0 ? "alert" : "skip",
        properties: {
          label: "Alarm",
          value: totalSummary.alarm,
          icon: "materialsymbols-solid:circle_notifications",
        },
      },
      {
        name: `${props.definition.name}.container.summary.status.error`,
        width: 2,
        display_type: totalSummary.error > 0 ? "alert" : "skip",
        properties: {
          label: "Error",
          value: totalSummary.error,
          icon: "materialsymbols-solid:error",
        },
      },
      {
        name: `${props.definition.name}.container.summary.status.info`,
        width: 2,
        display_type: totalSummary.info > 0 ? "info" : "skip",
        properties: {
          label: "Info",
          value: totalSummary.info,
          icon: "materialsymbols-solid:info",
        },
      },
      {
        name: `${props.definition.name}.container.summary.status.skip`,
        width: 2,
        display_type: "skip",
        properties: {
          label: "Skipped",
          value: totalSummary.skip,
          icon: "materialsymbols-solid:arrow_circle_right",
        },
      },
    ];

    const severity_summary = props.grouping.severity_summary;
    const criticalRaw = severity_summary["critical"];
    const highRaw = severity_summary["high"];
    const critical = criticalRaw || 0;
    const high = highRaw || 0;

    // If we have at least 1 critical or undefined control defined in this run
    if (criticalRaw !== undefined || highRaw !== undefined) {
      const total = critical + high;
      summary_cards.push({
        name: `${props.definition.name}.container.summary.severity`,
        width: 2,
        display_type: total > 0 ? "severity" : "",
        properties: {
          label: "Critical / High",
          value: total,
          icon: "materialsymbols-solid:warning",
        },
      });
    }
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
        onMouseEnter: panelControlsEnabled
          ? () => setShowPanelControls(true)
          : noop,
        onMouseLeave: () => setShowPanelControls(false),
      }}
      setRef={setReferenceElement}
    >
      {/*Don't show when in panel detail view*/}
      {!selectedPanel && (
        <DashboardTitle
          title={props.definition.title}
          controls={
            showPanelControls ? (
              <PanelControls
                referenceElement={referenceElement}
                controls={benchmarkControls}
              />
            ) : null
          }
        />
      )}
      <Grid name={`${props.definition.name}.container.summary`}>
        {summaryCards
          .filter(({ name }) => {
            const statusFilter = expressions?.find(
              (expr) => expr.type === "status",
            );
            const severityType = name.split(".")[name.split(".").length - 1];
            if (statusFilter && statusFilter.operator === "equal") {
              return severityType === statusFilter.value;
            } else if (statusFilter && statusFilter.operator === "not_equal") {
              return severityType !== statusFilter.value;
            } else if (statusFilter && statusFilter.operator === "in") {
              return statusFilter.value?.includes(severityType);
            } else if (statusFilter && statusFilter.operator === "not_in") {
              return !statusFilter.value?.includes(severityType);
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
              >
                <FilterCardWrapper
                  cardName={summaryCard.name}
                  panelName={props.definition.name}
                  dimension="status"
                  expressions={expressions}
                >
                  <Card {...cardProps} />
                </FilterCardWrapper>
              </Panel>
            );
          })}
      </Grid>
      <Grid name={`${props.definition.name}.container.tree`}>
        <BenchmarkTree
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

const BenchmarkTree = (props: BenchmarkTreeProps) => {
  if (!props.properties || !props.properties.first_child_summaries) {
    return null;
  }

  return <CheckGrouping node={props.properties.grouping} />;
};

const BenchmarkTableView = ({
  benchmark,
  definition,
}: BenchmarkTableViewProps) => {
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
      />
    </Panel>
  );
};

const Inner = ({ withTitle }) => {
  const {
    benchmark,
    definition,
    grouping,
    groupingConfig,
    firstChildSummaries,
  } = useBenchmarkGrouping();

  if (!definition || !benchmark || !grouping) {
    return null;
  }

  if (!definition.display_type || definition.display_type === "benchmark") {
    return (
      <Benchmark
        benchmark={benchmark}
        definition={definition}
        grouping={grouping}
        groupingConfig={groupingConfig}
        firstChildSummaries={firstChildSummaries}
        withTitle={withTitle}
      />
    );
    // @ts-ignore
  } else if (definition.display_type === "table") {
    return <BenchmarkTableView benchmark={benchmark} definition={definition} />;
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

type BenchmarkProps = PanelDefinition & {
  showControls: boolean;
  withTitle: boolean;
};

const BenchmarkWrapper = (props: BenchmarkProps) => {
  return (
    <GroupingProvider definition={props}>
      <PanelControlsProvider definition={props} enabled={props.showControls}>
        <Inner withTitle={props.withTitle} />
      </PanelControlsProvider>
    </GroupingProvider>
  );
};

registerComponent("benchmark", BenchmarkWrapper);

export default BenchmarkWrapper;
