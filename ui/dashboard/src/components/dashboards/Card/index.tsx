// Ensure Table is loaded & registered first
import "@powerpipe/components/dashboards/Table";
import DashboardIcon from "../common/DashboardIcon";
import has from "lodash/has";
import IntegerDisplay from "@powerpipe/components/IntegerDisplay";
import isNumber from "lodash/isNumber";
import isObject from "lodash/isObject";
import LoadingIndicator from "../LoadingIndicator";
import useDeepCompareEffect from "use-deep-compare-effect";
import useTemplateRender from "@powerpipe/hooks/useTemplateRender";
import { BasePrimitiveProps, ExecutablePrimitiveProps } from "../common";
import {
  CardDataProcessor,
  CardDiffState,
  CardType,
} from "../data/CardDataProcessor";
import { classNames } from "@powerpipe/utils/styles";
import { injectSearchPathPrefix } from "@powerpipe/utils/url";
import { getComponent, registerComponent } from "../index";
import {
  getIconClasses,
  getIconStyles,
  getWrapperClasses,
} from "@powerpipe/utils/card";
import { IDiffProperties } from "../data/types";
import { PanelProperties } from "@powerpipe/types";
import { useDashboard } from "@powerpipe/hooks/useDashboard";
import { useEffect, useState } from "react";

const Table = getComponent("table");

export interface CardProperties extends IDiffProperties {
  label?: string;
  value?: any;
  icon?: string;
  href?: string;
  loading?: boolean;
}

export type CardProps = PanelProperties &
  Omit<BasePrimitiveProps, "display_type"> &
  ExecutablePrimitiveProps & {
    display_type?: CardType;
    properties: CardProperties;
  };

type CardState = {
  loading: boolean;
  label: string | null;
  value: any | null;
  type: CardType;
  icon: string | null;
  href: string | null;
  diff?: CardDiffState;
};

// TODO diffing
// Need to know we're in diff mode
// Need data to diff against
// Need to be able to diff said data against current data
// Need to try to infer state of the change as best as possible
// e.g. a card going from alarm 10 to alarm 5 is good, so it's down 100% / green
//      a card going from alarm 10 to alarm 20 is bad, so it's up 100% / red
//      a card going from alarm 10 to alarm 10 is neutral, so it's no change
//      a card going from alarm 10 to ok 10 is good, so it's no change in value, but change in state

const useCardState = ({
  data,
  display_type,
  properties,
  status,
}: CardProps) => {
  const [calculatedProperties, setCalculatedProperties] = useState<CardState>(
    new CardDataProcessor().getDefaultState(status, properties, display_type),
  );

  useDeepCompareEffect(() => {
    const diff = new CardDataProcessor();
    const newState = diff.buildCardState(
      data,
      display_type,
      properties,
      status,
    );
    setCalculatedProperties(newState);
    setCalculatedProperties(
      diff.buildCardState(data, display_type, properties, status),
    );
  }, [data, display_type, properties, setCalculatedProperties, status]);

  return calculatedProperties;
};

const Label = ({ value }) => {
  if (!value) {
    return null;
  }

  if (isObject(value)) {
    return JSON.stringify(value);
  }

  return value;
};

const Value = ({ loading, value }) => {
  if (loading || value === null || value === undefined) {
    return (
      <DashboardIcon
        className="h-8 w-8"
        icon="materialsymbols-outline:remove"
      />
    );
  }

  if (isNumber(value)) {
    return <IntegerDisplay num={value} startAt="100k" />;
  }

  return <Label value={value} />;
};

const CardDiffDisplay = ({ diff }: { diff: CardDiffState | undefined }) => {
  if (!diff || diff.direction === "none") {
    return null;
  }

  return (
    <span
      className={classNames(
        "inline-flex rounded-lg px-2 font-medium md:mt-2 space-x-1 text-lg",
        diff.status === "ok" ? "text-ok" : null,
        diff.status === "alert" ? "text-alert" : null,
        // diff.status === "alert" ? "text-severity" : null
      )}
    >
      <span aria-hidden="true" className={classNames("self-end")}>
        {diff.direction === "up" ? "↑" : diff.direction === "down" ? "↓" : null}
      </span>
      <span className="sr-only">
        {" "}
        {diff.direction === "up" ? "Increased" : "Decreased"} by{" "}
      </span>
      {(diff.direction === "up" || diff.direction === "down") && (
        <>
          {/*@ts-ignore*/}
          <div className="flex items-baseline space-x-1">
            <span className="text-lg font-semibold">
              {/*@ts-ignore*/}
              <IntegerDisplay num={diff.value || null} />
            </span>
            {diff.value_percent !== undefined && (
              <span className="text-sm ">({diff.value_percent}%)</span>
            )}
          </div>
        </>
      )}
    </span>
  );
};

const ValueWithDiff = ({ loading, value, diff }) => (
  <div className="flex items-center space-x-2">
    {" "}
    {/* Adjusts flex to align elements inline */}
    <Value loading={loading} value={value} />
    {diff && <CardDiffDisplay diff={diff} />}
  </div>
);

const Card = (props: CardProps) => {
  const ExternalLink = getComponent("external_link");
  const state = useCardState(props);
  const { searchPathPrefix } = useDashboard();
  const [renderError, setRenderError] = useState<string | null>(null);
  const [renderedHref, setRenderedHref] = useState<string | null>(
    state.href || null,
  );
  const { ready: templateRenderReady, renderTemplates } = useTemplateRender();

  useEffect(() => {
    if ((state.loading || !state.href) && (renderError || renderedHref)) {
      setRenderError(null);
      setRenderedHref(null);
    }
  }, [state.loading, state.href, renderError, renderedHref]);

  useDeepCompareEffect(() => {
    if (!templateRenderReady || state.loading || !state.href) {
      return;
    }

    const renderData = { ...state };
    if (props.data && props.data.columns && props.data.rows) {
      const row = props.data.rows[0];
      props.data.columns.forEach((col) => {
        if (!has(renderData, col.name)) {
          renderData[col.name] = row[col.name];
        }
      });
    }

    const doRender = async () => {
      const renderedResults = await renderTemplates(
        { card: state.href as string },
        [renderData],
      );
      if (
        !renderedResults ||
        renderedResults.length === 0 ||
        !renderedResults[0].card
      ) {
        setRenderedHref(null);
        setRenderError(null);
      } else if (renderedResults[0].card.result) {
        const withSearchPathPrefix = injectSearchPathPrefix(
          renderedResults[0].card.result,
          searchPathPrefix,
        );
        setRenderedHref(withSearchPathPrefix);
        setRenderError(null);
      } else if (renderedResults[0].card.error) {
        setRenderError(renderedResults[0].card.error as string);
        setRenderedHref(null);
      }
    };
    doRender();
  }, [renderTemplates, templateRenderReady, state, props.data]);

  const card = (
    <div
      className={classNames(
        "overflow-hidden bg-dashboard-panel text-foreground print:bg-white print:text-black shadow-sm p-3 pr-5",
        getWrapperClasses(state.type),
        !state.icon ? "pl-4" : undefined,
      )}
    >
      <div className="flex space-x-3">
        {(state.loading || state.icon) && (
          <div className={classNames("shrink-0 grow-0")}>
            {state.loading ? (
              <LoadingIndicator className="h-8 w-8" />
            ) : (
              <DashboardIcon
                className={classNames("h-8 w-8", getIconClasses(state.type))}
                aria-hidden="true"
                icon={state.icon}
                style={getIconStyles(state.type)}
              />
            )}
          </div>
        )}
        <div className="grow mt-0.5 min-w-0">
          <dt>
            <p className="text-lg truncate" title={state.label || undefined}>
              {state.loading ? "Loading..." : state.label}
            </p>
          </dt>
          <dd className="font-semibold text-3xl mt-1 mb-1">
            <ValueWithDiff
              loading={state.loading}
              value={state.value}
              diff={state.diff}
            />
          </dd>
        </div>
      </div>
    </div>
  );

  if (renderedHref) {
    return <ExternalLink to={renderedHref}>{card}</ExternalLink>;
  }

  return card;
};

const CardWrapper = (props: CardProps) => {
  if (props.display_type === "table") {
    // @ts-ignore
    return <Table {...props} />;
  }
  return <Card {...props} />;
};

registerComponent("card", CardWrapper);

export default CardWrapper;
