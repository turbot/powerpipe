import usePrefixedSearchParams from "@powerpipe/hooks/usePrefixedSearchParams";
import {
  createContext,
  ReactNode,
  useCallback,
  useContext,
  useState,
} from "react";
import { KeyValueStringPairs } from "@powerpipe/components/dashboards/common/types";
import { useSearchParams } from "react-router-dom";

interface IDashboardInputsContext {
  inputs: Record<string, string>;
  lastChangedInput: string | null;
  setLastChangedInput: (name: string | null) => void;
  updateInput: (name: string, value: string, recordHistory: boolean) => void;
  deleteInput: (name: string, recordHistory: boolean) => void;
  setInputs: (values: KeyValueStringPairs, recordHistory: boolean) => void;
  clearInputs: (recordHistory: boolean) => void;
}

interface DashboardInputsProviderProps {
  children: ReactNode;
}

const DashboardInputsContext = createContext<IDashboardInputsContext | null>(
  null,
);

export const DashboardInputsProvider = ({
  children,
}: DashboardInputsProviderProps) => {
  const [lastChangedInput, setLastChangedInput] = useState<string | null>(null);
  const [searchParams, setSearchParams] = useSearchParams();
  const inputs = usePrefixedSearchParams("input.");

  const updateInput = useCallback(
    (name: string, value: string, recordHistory: boolean) => {
      const currentValue = searchParams.get(name);
      searchParams.set(name, value);
      setSearchParams(searchParams, { replace: !recordHistory });
      if (value === currentValue) {
        return;
      }
      setLastChangedInput(name);
    },
    [searchParams, setLastChangedInput, setSearchParams],
  );

  const deleteInput = useCallback(
    (name: string, recordHistory: boolean) => {
      searchParams.delete(name);
      setSearchParams(searchParams, { replace: !recordHistory });
      setLastChangedInput(name);
    },
    [searchParams, setLastChangedInput, setSearchParams],
  );

  const setInputs = useCallback(
    (values: KeyValueStringPairs, recordHistory: boolean) => {
      for (const key of Object.keys(inputs)) {
        searchParams.delete(key);
      }
      for (const [name, value] of Object.entries(values) || {}) {
        searchParams.set(name, value);
      }
      setSearchParams(searchParams, { replace: !recordHistory });
      setLastChangedInput(null);
    },
    [inputs, searchParams, setLastChangedInput, setSearchParams],
  );

  const clearInputs = useCallback(
    (recordHistory: boolean) => {
      for (const key of Object.keys(inputs)) {
        searchParams.delete(key);
      }
      setSearchParams(searchParams, { replace: !recordHistory });
      setLastChangedInput(null);
    },
    [inputs, searchParams, setLastChangedInput, setSearchParams],
  );

  return (
    <DashboardInputsContext.Provider
      value={{
        inputs,
        lastChangedInput,
        updateInput,
        deleteInput,
        setLastChangedInput,
        setInputs,
        clearInputs,
      }}
    >
      {children}
    </DashboardInputsContext.Provider>
  );
};

export const useDashboardInputs = () => {
  const context = useContext(DashboardInputsContext);
  if (!context) {
    throw new Error(
      "useDashboardInputs must be used within a DashboardInputsContext",
    );
  }
  return context;
};
