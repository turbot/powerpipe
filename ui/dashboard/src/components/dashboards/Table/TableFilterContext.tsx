import { createContext, ReactNode, useContext } from "react";

type ITableFilterContext = {
  additiveFilters: [];
  exclusionFilter: [];

  setPanelInformation: (information: ReactNode) => void;
  setShowPanelControls: (show: boolean) => void;
  setShowPanelInformation: (show: boolean) => void;
};

const TableFilterContext = createContext<ITableFilterContext | null>(null);

const TableFilterProvider = ({
  children,
  definition,
  parentType,
  showControls,
}: TableFilterProviderProps) => {
  return (
    <TableFilterContext.Provider
      value={{
        definition,
      }}
    >
      {children}
    </TableFilterContext.Provider>
  );
};

const useTableFilter = () => {
  const context = useContext(TableFilterContext);
  if (context === undefined) {
    throw new Error("usePanel must be used within a TableFilterContext");
  }
  return context as ITableFilterContext;
};

export { TableFilterContext, TableFilterProvider, useTableFilter };
