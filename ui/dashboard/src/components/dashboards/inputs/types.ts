import {
  BasePrimitiveProps,
  ExecutablePrimitiveProps,
} from "@powerpipe/components/dashboards/common";
import { ComponentType } from "react";
import { PanelDefinition } from "@powerpipe/types";

export type BaseInputProps = PanelDefinition &
  BasePrimitiveProps &
  ExecutablePrimitiveProps;

export type SelectOption = {
  label: React.ReactNode;
  value: string | null;
  tags?: object;
};

export type SelectInputOption = {
  name: string;
  label?: string;
};

export type InputProperties = {
  label?: string;
  options?: SelectInputOption[];
  placeholder?: string;
  unqualified_name: string;
};

export type InputProps = BaseInputProps & {
  display_type?: InputType;
  properties: InputProperties;
};

export type InputType =
  | "combo"
  | "date"
  | "datetime"
  | "date_range"
  | "datetime_range"
  | "hidden"
  | "multicombo"
  | "multiselect"
  | "select"
  | "table"
  | "text"
  | "time";

export type IInput = {
  type: InputType;
  component: ComponentType<any>;
};
