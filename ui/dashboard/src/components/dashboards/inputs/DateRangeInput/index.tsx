import DatetimeRangePicker from "@powerpipe/components/dashboards/DatetimeRange/DatetimeRangePicker";
import dayjs from "dayjs";
import { DashboardDataModeLive } from "@powerpipe/types";
import {
  IInput,
  InputProps,
} from "@powerpipe/components/dashboards/inputs/types";
import { registerInputComponent } from "@powerpipe/components/dashboards/inputs";
import { useDashboardInputs } from "@powerpipe/hooks/useDashboardInputs";
import { useDashboardState } from "@powerpipe/hooks/useDashboardState";
import { useEffect, useMemo } from "react";

const defaultValue = {
  from: dayjs().subtract(7, "day"),
  to: null,
  relative: "7d",
};

const DateRangeInput = (props: InputProps) => {
  const { dataMode } = useDashboardState();
  const { inputs, updateInput } = useDashboardInputs();
  const stateValue = inputs[props.name];

  const value = useMemo(() => {
    if (stateValue) {
      try {
        const parsed = JSON.parse(stateValue);
        return {
          from: parsed.from
            ? dayjs(parsed.from)
            : dayjs().subtract(7, "day").utc(),
          to: parsed.to ? dayjs(parsed.to) : null,
          relative: parsed.relative || "7d",
        };
      } catch (err) {
        console.error("Parse error", err);
        return defaultValue;
      }
    } else {
      return defaultValue;
    }
  }, [stateValue]);

  useEffect(() => {
    if (stateValue) {
      return;
    }
    updateInput(
      props.name,
      // [defaultValue.from.format("YYYY-MM-DD"), null, "7d"].join(","),
      JSON.stringify(defaultValue),
      !!stateValue,
    );
  }, [stateValue]);

  const onInputChange = (
    from: string,
    to?: string | null,
    relative?: string | null,
  ) => {
    updateInput(
      props.name,
      // [from, to, relative].join(","),
      JSON.stringify({
        from,
        to,
        relative,
      }),
      !!stateValue,
    );
  };

  return (
    <DatetimeRangePicker
      from={value.from}
      to={value.to}
      relative={value.relative}
      disabled={dataMode !== DashboardDataModeLive}
      withTime={false}
      onChange={onInputChange}
    />
  );
};

const definition: IInput = {
  type: "date_range",
  component: DateRangeInput,
};

registerInputComponent(definition.type, definition);

export default definition;
