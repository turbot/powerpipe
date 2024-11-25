// Ensure Table is loaded & registered first
import "@powerpipe/components/dashboards/Table";
import Card, { CardProps } from "@powerpipe/components/dashboards/Card";
import CheckGrouping from "../CheckGrouping";
import CustomizeViewSummary from "../CustomizeViewSummary";
import DashboardTitle from "@powerpipe/components/dashboards/titles/DashboardTitle";
import Error from "@powerpipe/components/dashboards/Error";
import Grid from "@powerpipe/components/dashboards/layout/Grid";
import Panel from "@powerpipe/components/dashboards/layout/Panel";
import PanelControls from "@powerpipe/components/dashboards/layout/Panel/PanelControls";
import useFilterConfig from "@powerpipe/hooks/useFilterConfig";
import usePanelControls from "@powerpipe/hooks/usePanelControls";
import {
  BenchmarkTreeProps,
  CheckDisplayGroup,
  Filter,
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
import { DashboardActions, PanelDefinition, PanelsMap } from "@powerpipe/types";
import { useDashboard } from "@powerpipe/hooks/useDashboard";
import { useEffect, useMemo, useState } from "react";
import { useSearchParams } from "react-router-dom";
import { validateFilter } from "@powerpipe/components/dashboards/grouping/FilterEditor";
import { Width } from "@powerpipe/components/dashboards/common";

const Table = getComponent("table");

type BenchmarkTableViewProps = {
  benchmark: BenchmarkType;
  definition: PanelDefinition;
};

type InnerCheckProps = {
  benchmark: BenchmarkType;
  definition: PanelDefinition;
  diff_panels?: PanelsMap;
  grouping: CheckNode;
  groupingConfig: CheckDisplayGroup[];
  firstChildSummaries: CheckSummary[];
  diffFirstChildSummaries: CheckSummary[] | undefined;
  diffGrouping: CheckNode | null;
  showControls: boolean;
  withTitle: boolean;
};

const Benchmark = (props: InnerCheckProps) => {
  const {
    filter: { expressions },
  } = useFilterConfig(props.definition?.name);
  const { dispatch, selectedPanel } = useDashboard();
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
  const [showBenchmarkControls, setShowBenchmarkControls] = useState(false);
  const definitionWithData = useMemo(() => {
    return {
      ...props.definition,
      data: benchmarkDataTable,
    };
  }, [benchmarkDataTable, props.definition]);
  const { panelControls: benchmarkControls, setCustomControls } =
    usePanelControls(definitionWithData, props.showControls);

  useEffect(() => {
    setCustomControls([
      {
        action: async () =>
          dispatch({
            type: DashboardActions.SHOW_CUSTOMIZE_BENCHMARK_PANEL,
            panel_name: props.definition.name,
          }),
        component: <CustomizeViewSummary panelName={props.definition.name} />,
        title: "Filter & Group",
      },
    ]);
  }, [dispatch, props.definition.name, setCustomControls]);

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

    let diffTotalSummary: CheckSummary | null = null;
    if (
      props.diffFirstChildSummaries &&
      props.diffFirstChildSummaries?.length > 0
    ) {
      diffTotalSummary = props.diffFirstChildSummaries.reduce(
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
    }

    const summary_cards = [
      {
        name: `${props.definition.name}.container.summary.ok`,
        width: 2,
        display_type: totalSummary.ok > 0 ? "ok" : "skip",
        properties: {
          label: "OK",
          value: totalSummary.ok,
          // icon: "materialsymbols-solid:check",
          icon: "materialsymbols-solid:check_circle",
          // icon: "materialsymbols-outline:check_circle",
        },
        diff_panel: !!diffTotalSummary
          ? {
              name: `${props.definition.name}.container.summary.ok.diff`,
              width: 2,
              display_type: diffTotalSummary.ok > 0 ? "ok" : "skip",
              properties: {
                label: "OK",
                value: diffTotalSummary.ok,
                // icon: "materialsymbols-solid:check",
                icon: "materialsymbols-solid:check_circle",
                // icon: "materialsymbols-outline:check_circle",
              },
            }
          : null,
      },
      {
        name: `${props.definition.name}.container.summary.alarm`,
        width: 2,
        display_type: totalSummary.alarm > 0 ? "alert" : "skip",
        properties: {
          label: "Alarm",
          value: totalSummary.alarm,
          icon: "materialsymbols-solid:circle_notifications",
          // icon: "materialsymbols-solid:circle_notifications",
          // icon: "materialsymbols-outline:circle_notifications",
        },
        diff_panel: !!diffTotalSummary
          ? {
              name: `${props.definition.name}.container.summary.alarm.diff`,
              width: 2,
              display_type: diffTotalSummary.alarm > 0 ? "alert" : "skip",
              properties: {
                label: "Alarm",
                value: diffTotalSummary.alarm,
                icon: "materialsymbols-solid:circle_notifications",
                // icon: "materialsymbols-solid:circle_notifications",
                // icon: "materialsymbols-outline:circle_notifications",
              },
            }
          : null,
      },
      {
        name: `${props.definition.name}.container.summary.error`,
        width: 2,
        display_type: totalSummary.error > 0 ? "alert" : "skip",
        properties: {
          label: "Error",
          value: totalSummary.error,
          // icon: "materialsymbols-solid:priority_high",
          icon: "materialsymbols-solid:error",
          // icon: "materialsymbols-outline:error",
        },
        diff_panel: !!diffTotalSummary
          ? {
              name: `${props.definition.name}.container.summary.error.diff`,
              width: 2,
              display_type: diffTotalSummary.error > 0 ? "alert" : "skip",
              properties: {
                label: "Error",
                value: diffTotalSummary.error,
                // icon: "materialsymbols-solid:priority_high",
                icon: "materialsymbols-solid:error",
                // icon: "materialsymbols-outline:error",
              },
            }
          : null,
      },
      {
        name: `${props.definition.name}.container.summary.info`,
        width: 2,
        display_type: totalSummary.info > 0 ? "info" : "skip",
        properties: {
          label: "Info",
          value: totalSummary.info,
          // icon: "materialsymbols-solid:info_i",
          icon: "materialsymbols-solid:info",
          // icon: "materialsymbols-outline:info",
        },
        diff_panel: !!diffTotalSummary
          ? {
              name: `${props.definition.name}.container.summary.info.diff`,
              width: 2,
              display_type: diffTotalSummary.info > 0 ? "info" : "skip",
              properties: {
                label: "Info",
                value: diffTotalSummary.info,
                // icon: "materialsymbols-solid:info_i",
                icon: "materialsymbols-solid:info",
                // icon: "materialsymbols-outline:info",
              },
            }
          : null,
      },
      {
        name: `${props.definition.name}.container.summary.skip`,
        width: 2,
        display_type: "skip",
        properties: {
          label: "Skipped",
          value: totalSummary.skip,
          // icon: "materialsymbols-solid:arrow_right_alt",
          icon: "materialsymbols-solid:arrow_circle_right",
          // icon: "materialsymbols-outline:arrow_circle_right",
        },
        diff_panel: !!diffTotalSummary
          ? {
              name: `${props.definition.name}.container.summary.skip.diff`,
              width: 2,
              display_type: "skip",
              properties: {
                label: "Skipped",
                value: diffTotalSummary.skip,
                // icon: "materialsymbols-solid:arrow_right_alt",
                icon: "materialsymbols-solid:arrow_circle_right",
                // icon: "materialsymbols-solid:arrow_circle_right",
                // icon: "materialsymbols-outline:arrow_circle_right",
              },
            }
          : null,
      },
    ];

    const severity_summary = props.grouping.severity_summary;
    const criticalRaw = severity_summary["critical"];
    const highRaw = severity_summary["high"];
    const critical = criticalRaw || 0;
    const high = highRaw || 0;

    // Calc diff vs previous
    const diff_severity_summary = props.diffGrouping?.severity_summary;
    let diffCriticalRaw, diffHighRaw, diffCritical, diffHigh;
    if (diff_severity_summary) {
      diffCriticalRaw = diff_severity_summary["critical"];
      diffHighRaw = diff_severity_summary["high"];
      diffCritical = diffCriticalRaw || 0;
      diffHigh = diffHighRaw || 0;
    }

    // If we have at least 1 critical or undefined control defined in this run
    if (criticalRaw !== undefined || highRaw !== undefined) {
      const total = critical + high;
      const diffTotal = diffCritical + diffHigh;
      summary_cards.push({
        name: `${props.definition.name}.container.summary.severity`,
        width: 2,
        display_type: total > 0 ? "severity" : "",
        properties: {
          label: "Critical / High",
          value: total,
          icon: "materialsymbols-solid:warning",
        },
        diff_panel: diff_severity_summary
          ? {
              name: `${props.definition.name}.container.summary.severity.diff`,
              width: 2,
              display_type: diffTotal > 0 ? "severity" : "",
              properties: {
                label: "Critical / High",
                value: diffTotal,
                icon: "materialsymbols-solid:warning",
              },
            }
          : null,
      });
    }
    return summary_cards;
  }, [
    props.firstChildSummaries,
    props.diffFirstChildSummaries,
    props.grouping,
    props.diffGrouping,
    props.definition.name,
  ]);

  const [, setSearchParams] = useSearchParams();

  if (!props.grouping) {
    return null;
  }

  const toggleFilter = (filterName: string) => () => {
    const split = filterName.split(".");
    filterName = split[split.length - 1];
    const expressionHasFilter = !!expressions?.find(
      (expr) => expr.type === "status",
    );
    let newFilter: Filter;
    if (expressionHasFilter) {
      newFilter = {
        operator: "and",
        expressions: expressions?.filter((expr) => expr.type !== "status"),
      } as Filter;
      if (validateFilter(newFilter)) {
        setSearchParams((prev) => {
          const newParams = new URLSearchParams(prev);
          const asJson = JSON.stringify(newFilter);
          newParams.set("where", asJson);
          return newParams;
        });
      } else {
        setSearchParams((prev) => {
          const newParams = new URLSearchParams(prev);
          newParams.delete("where");
          return newParams;
        });
      }
    } else {
      newFilter = {
        operator: "and",
        expressions: expressions
          ?.filter((expr) => !!expr.type)
          .concat({
            type: "status",
            value: filterName,
            operator: "equal",
            title: filterName,
          }),
      } as Filter;
      if (validateFilter(newFilter)) {
        setSearchParams((prev) => {
          const newParams = new URLSearchParams(prev);
          const asJson = JSON.stringify(newFilter);
          newParams.set("where", asJson);
          return newParams;
        });
      }
    }
  };

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
      <Grid name={`${props.definition.name}.container.summary`}>
        {summaryCards
          .filter(({ name }) => {
            const statusFromExpressions = expressions?.find(
              (expr) => expr.type === "status",
            )?.value;
            if (statusFromExpressions) {
              return name.includes(statusFromExpressions);
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
                <span
                  className="cursor-pointer"
                  onClick={toggleFilter(summaryCard.name)}
                >
                  {/*@ts-ignore*/}
                  <Card {...cardProps} diff_panel={summaryCard.diff_panel} />
                </span>
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

const Inner = ({ showControls, withTitle }) => {
  const {
    benchmark,
    definition,
    grouping,
    groupingsConfig,
    firstChildSummaries,
    diffFirstChildSummaries,
    diffGrouping,
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
        groupingConfig={groupingsConfig}
        firstChildSummaries={firstChildSummaries}
        showControls={showControls}
        withTitle={withTitle}
        diffFirstChildSummaries={diffFirstChildSummaries}
        diffGrouping={diffGrouping}
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
  diff_panels?: PanelsMap;
  showControls: boolean;
  withTitle: boolean;
};

const BenchmarkWrapper = (props: BenchmarkProps) => {
  return (
    <GroupingProvider definition={props} diff_panels={props.diff_panels}>
      <Inner showControls={props.showControls} withTitle={props.withTitle} />
    </GroupingProvider>
  );
};

registerComponent("benchmark", BenchmarkWrapper);

export default BenchmarkWrapper;
