// Ensure Table is loaded & registered first
import "@powerpipe/components/dashboards/Table";
import Card, { CardProps } from "@powerpipe/components/dashboards/Card";
import CustomizeViewSummary from "../CustomizeViewSummary";
import DashboardTitle from "@powerpipe/components/dashboards/titles/DashboardTitle";
import DetectionGrouping from "../DetectionGrouping";
import Error from "@powerpipe/components/dashboards/Error";
import Grid from "@powerpipe/components/dashboards/layout/Grid";
import Panel from "@powerpipe/components/dashboards/layout/Panel";
import PanelControls from "@powerpipe/components/dashboards/layout/Panel/PanelControls";
import useGroupingFilterConfig from "@powerpipe/hooks/useGroupingFilterConfig";
import usePanelControls from "@powerpipe/hooks/usePanelControls";
import {
  DetectionBenchmarkTreeProps,
  DetectionDisplayGroup,
  DetectionFilter,
  DetectionNode,
  DetectionSummary,
} from "@powerpipe/components/dashboards/grouping/common";
import { CardType } from "@powerpipe/components/dashboards/data/CardDataProcessor";
import { default as DetectionBenchmarkType } from "../common/DetectionBenchmark";
import {
  getComponent,
  registerComponent,
} from "@powerpipe/components/dashboards";
import {
  GroupingProvider,
  useDetectionGrouping,
} from "@powerpipe/hooks/useDetectionGrouping";
import { noop } from "@powerpipe/utils/func";
import { DashboardActions, PanelDefinition, PanelsMap } from "@powerpipe/types";
import { useDashboard } from "@powerpipe/hooks/useDashboard";
import { useEffect, useMemo, useState } from "react";
import { useSearchParams } from "react-router-dom";
import { validateFilter } from "../CheckFilterEditor";
import { Width } from "@powerpipe/components/dashboards/common";
import DateRangePicker from "./DateRangePicker/DateRangePicker";

const Table = getComponent("table");

type BenchmarkTableViewProps = {
  benchmark: DetectionBenchmarkType;
  definition: PanelDefinition;
};

type InnerCheckProps = {
  benchmark: DetectionBenchmarkType;
  definition: PanelDefinition;
  diff_panels?: PanelsMap;
  grouping: DetectionNode;
  groupingConfig: DetectionDisplayGroup[];
  firstChildSummaries: DetectionSummary[];
  diffFirstChildSummaries: DetectionSummary[] | undefined;
  diffGrouping: DetectionNode | null;
  showControls: boolean;
  withTitle: boolean;
};

const DetectionBenchmark = (props: InnerCheckProps) => {
  const { expressions } = useGroupingFilterConfig();
  const { cliMode, dispatch, selectedPanel } = useDashboard();
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
          dispatch({ type: DashboardActions.SHOW_CUSTOMIZE_BENCHMARK_PANEL }),
        component: <CustomizeViewSummary />,
        title: "Filter & Group",
      },
    ]);
  }, [cliMode, dispatch, setCustomControls]);

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

    let diffTotalSummary: DetectionSummary | null = null;
    if (
      props.diffFirstChildSummaries &&
      props.diffFirstChildSummaries?.length > 0
    ) {
      diffTotalSummary = props.diffFirstChildSummaries.reduce(
        (cumulative, current) => {
          cumulative.total += current.total;
          return cumulative;
        },
        { total: 0 },
      );
    }

    return [
      {
        name: `${props.definition.name}.container.summary.total`,
        width: 2,
        display_type: totalSummary.total > 0 ? "alert" : "ok",
        properties: {
          label: "Total Detections",
          value: totalSummary.total,
          icon:
            totalSummary.total > 0
              ? "materialsymbols-solid:circle_notifications"
              : "materialsymbols-solid:check_circle",
        },
        diff_panel: !!diffTotalSummary
          ? {
              name: `${props.definition.name}.container.summary.total.diff`,
              width: 2,
              display_type: diffTotalSummary.total > 0 ? "ok" : "skip",
              properties: {
                label: "Total Detections",
                value: totalSummary.total,
                icon:
                  totalSummary.total > 0
                    ? "materialsymbols-solid:circle_notifications"
                    : "materialsymbols-solid:check_circle",
              },
            }
          : null,
      },
    ];
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
    let newFilter: DetectionFilter;
    if (expressionHasFilter) {
      newFilter = {
        operator: "and",
        expressions: expressions?.filter((expr) => expr.type !== "status"),
      } as DetectionFilter;
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
      } as DetectionFilter;
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
      <Grid name="temp">
        <DateRangePicker />
      </Grid>
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
                parentType="detection_benchmark"
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

type DetectionBenchmarkWrapperProps = PanelDefinition & {
  diff_panels?: PanelsMap;
  showControls: boolean;
  withTitle: boolean;
};

const DetectionBenchmarkWrapper = (props: DetectionBenchmarkWrapperProps) => {
  return (
    <GroupingProvider definition={props} diff_panels={props.diff_panels}>
      <Inner showControls={props.showControls} withTitle={props.withTitle} />
    </GroupingProvider>
  );
};

registerComponent("detection_benchmark", DetectionBenchmarkWrapper);

export default DetectionBenchmarkWrapper;
