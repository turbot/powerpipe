import { classNames } from "@powerpipe/utils/styles";
import { ClearIcon, SearchIcon } from "@powerpipe/constants/icons";
import { forwardRef } from "react";
import { ThemeNames } from "@powerpipe/hooks/useTheme";
import { useDashboardTheme } from "@powerpipe/hooks/useDashboardTheme";

const SearchInput = forwardRef(
  (
    {
      className,
      disabled = false,
      placeholder,
      readOnly = false,
      setValue,
      value,
    }: {
      className?: string;
      disabled?: boolean;
      placeholder?: string;
      readOnly?: boolean;
      setValue: (value: string) => void;
      value: string;
    },
    ref,
  ) => {
    const { theme } = useDashboardTheme();
    return (
      <div className="relative">
        <div className="pointer-events-none absolute inset-y-0 left-0 pl-3 flex items-center text-foreground-light text-sm">
          <SearchIcon className="h-4 w-4" />
        </div>
        <input
          ref={ref}
          className={classNames(
            className,
            "flex-1 block w-full bg-dashboard-panel rounded-md border px-8 overflow-x-auto text-sm md:text-base disabled:bg-black-scale-1 focus:ring-0",
            theme.name === ThemeNames.STEAMPIPE_DARK
              ? "border-gray-700"
              : "border-[#e7e9ed]",
          )}
          disabled={disabled}
          onChange={(e) => setValue(e.target.value)}
          placeholder={placeholder}
          readOnly={readOnly}
          type="text"
          value={value}
        />
        {value && (
          <div
            className="absolute inset-y-0 right-0 pr-3 flex items-center cursor-pointer text-foreground"
            onClick={() => setValue("")}
          >
            <ClearIcon className="h-4 w-4" />
          </div>
        )}
      </div>
    );
  },
);

export default SearchInput;
