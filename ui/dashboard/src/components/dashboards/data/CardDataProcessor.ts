import get from "lodash/get";
import isNumber from "lodash/isNumber";
import { CardProperties } from "@powerpipe/components/dashboards/Card";
import { DashboardPanelType, DashboardRunState } from "@powerpipe/types";
import { getColumn, hasData } from "@powerpipe/utils/data";
import { getIconForType } from "@powerpipe/utils/card";
import { IPanelDiff } from "@powerpipe/components/dashboards/data/types";
import {
  isNumericCol,
  LeafNodeData,
} from "@powerpipe/components/dashboards/common";

export interface CardDiffState extends IPanelDiff {
  value?: number;
  value_percent?: "infinity" | number;
  direction: "none" | "up" | "down";
  status?: "ok" | "alert" | "severity" | null;
}

export interface CardState {
  loading: boolean;
  label: string | null;
  value: any | null;
  value_number: number | null;
  type: CardType;
  icon: string | null;
  href: string | null;
  diff?: CardDiffState;
}

export type CardDataFormat = "simple" | "formal";

export type CardType = "alert" | "info" | "ok" | "severity" | "table" | null;

export class CardDataProcessor {
  constructor() {}

  getDefaultState = (
    status: DashboardRunState,
    properties: CardProperties,
    display_type: CardType | undefined,
  ): CardState => {
    return {
      loading: status === "running" || !!properties?.loading,
      label: properties.label || null,
      value: isNumber(properties.value)
        ? properties.value.toLocaleString()
        : properties.value || null,
      value_number: isNumber(properties.value) ? properties.value : null,
      type: display_type || null,
      icon: getIconForType(display_type, properties.icon),
      href: properties.href || null,
    };
  };

  buildCardState(
    data: LeafNodeData | undefined,
    display_type: CardType | undefined,
    properties: CardProperties,
    status: DashboardRunState,
  ): CardState {
    if (!data || !hasData(data)) {
      const state = this.getDefaultState(status, properties, display_type);
      state.diff = { direction: "none" };
      return state;
    }

    return this.parseData(data, display_type, properties);
  }

  parseData(
    data: LeafNodeData,
    display_type: CardType | undefined,
    properties: CardProperties,
  ): CardState {
    const dataFormat = this.getDataFormat(data);
    if (dataFormat === "simple") {
      const firstCol = data.columns[0];
      const isNumericValue = isNumericCol(firstCol.data_type);
      const row = data.rows[0];
      const value = row[firstCol.name];
      return {
        loading: false,
        label: firstCol.name,
        value:
          value !== null && value !== undefined && isNumericValue
            ? value.toLocaleString()
            : value,
        value_number: isNumericValue && isNumber(value) ? value : null,
        type: display_type || null,
        icon: getIconForType(display_type, properties.icon),
        href: properties.href || null,
      };
    } else {
      const diffColumn = data.columns.find((c) => c.name === "__diff");
      const diffColumnValue = diffColumn
        ? get(data, `rows[0].${diffColumn.name}`, null)
        : null;
      const hasDiff = diffColumnValue !== null && diffColumnValue !== "none";
      const formalLabel = get(data, "rows[0].label", null);
      const formalValue = get(data, `rows[0].value`, null);
      const formalDiffValue = get(data, `rows[0].value_diff`, null);
      const formalType = get(data, `rows[0].type`, null);
      const formalIcon = get(data, `rows[0].icon`, null);
      const formalHref = get(data, `rows[0].href`, null);
      const valueCol = getColumn(data.columns, "value");
      const valueDiffCol = hasDiff
        ? getColumn(data.columns, "value_diff")
        : null;
      const isNumericValue = !!valueCol && isNumericCol(valueCol.data_type);
      const isNumericDiffValue =
        !!valueDiffCol && isNumericCol(valueDiffCol.data_type);
      const value =
        formalValue !== null && formalValue !== undefined && isNumericValue
          ? formalValue.toLocaleString()
          : formalValue;
      const value_number =
        formalValue && isNumericValue && isNumber(formalValue)
          ? formalValue
          : null;
      let value_number_diff;
      if (valueDiffCol) {
        value_number_diff =
          formalDiffValue && isNumericDiffValue && isNumber(formalDiffValue)
            ? formalDiffValue
            : null;
      }

      return {
        loading: false,
        label: formalLabel,
        value,
        value_number,
        diff: this.diff(
          hasDiff,
          value_number,
          value_number_diff,
          formalType || display_type,
        ),
        type: formalType || display_type || null,
        icon: getIconForType(
          formalType || display_type,
          formalIcon || properties.icon,
        ),
        href: formalHref || properties.href || null,
      };
    }
  }

  diff(
    hasDiff: boolean,
    currentValue: number | null,
    previousValue: number | null,
    displayType: CardType | undefined,
  ): CardDiffState {
    // If the columns aren't numeric then we can't diff...
    if (!hasDiff || currentValue === null || previousValue === null) {
      return {
        direction: "none",
      };
    }

    const direction =
      currentValue > previousValue
        ? "up"
        : currentValue === previousValue
          ? "none"
          : "down";

    let value: number;
    let value_percent: "infinity" | number;
    let status: "ok" | "alert" | "severity" | null = null;
    if (direction === "up") {
      value = currentValue - previousValue;
      value_percent = Math.ceil((value / previousValue) * 100);
      status =
        displayType === "alert"
          ? "alert"
          : displayType === "ok"
            ? "ok"
            : displayType === "severity"
              ? "severity"
              : null;
    } else if (direction === "down") {
      value = previousValue - currentValue;
      value_percent = Math.ceil((value / previousValue) * 100);
      console.log({ previousValue, currentValue, value, value_percent });
      status =
        displayType === "alert"
          ? "ok"
          : displayType === "ok"
            ? "alert"
            : displayType === "severity"
              ? "ok"
              : null;
    } else {
      value = 0;
      value_percent = 0;
    }

    return {
      value,
      value_percent,
      direction,
      status,
    };
  }

  getDataFormat = (data: LeafNodeData | undefined): CardDataFormat => {
    if (!!data && data.columns.length > 1) {
      return "formal";
    }
    return "simple";
  };

  get panel_type(): DashboardPanelType {
    return "card";
  }
}
