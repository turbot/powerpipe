import { useDashboardTheme } from "@powerpipe/hooks/useDashboardTheme";
import { useEffect, useState } from "react";

const useSelectInputStyles = () => {
  const [, setRandomVal] = useState(0);
  const { theme, wrapperRef } = useDashboardTheme();

  // This is annoying, but unless I force a refresh the theme doesn't stay in sync when you switch
  useEffect(() => setRandomVal(Math.random()), [theme.name]);

  if (!wrapperRef) {
    return null;
  }

  // @ts-ignore
  const style = window.getComputedStyle(wrapperRef);
  const background = style.getPropertyValue("--color-dashboard");
  const backgroundPanel = style.getPropertyValue("--color-dashboard-panel");
  const foreground = style.getPropertyValue("--color-foreground");
  const blackScale3 = style.getPropertyValue("--color-black-scale-3");

  return {
    clearIndicator: (provided) => ({
      ...provided,
      cursor: "pointer",
    }),
    control: (provided, state) => {
      return {
        ...provided,
        backgroundColor: backgroundPanel,
        borderColor: state.isFocused
          ? "#2684FF !important"
          : `${blackScale3} !important`,
        boxShadow: "none",
      };
    },
    dropdownIndicator: (provided) => ({
      ...provided,
      cursor: "pointer",
    }),
    input: (provided) => {
      return {
        ...provided,
        color: foreground,
      };
    },
    singleValue: (provided) => {
      return {
        ...provided,
        color: foreground,
      };
    },
    menu: (provided) => {
      return {
        ...provided,
        backgroundColor: backgroundPanel,
        border: `1px solid ${blackScale3}`,
        boxShadow: "none",
        marginTop: 0,
        marginBottom: 0,
        maxWidth: "400px",
      };
    },
    menuList: (provided) => {
      return {
        ...provided,
        paddingTop: 0,
        paddingBottom: 0,
      };
    },
    menuPortal: (base) => ({ ...base, zIndex: 9999 }),
    option: (provided, state) => {
      return {
        ...provided,
        backgroundColor: state.isFocused ? background : "none",
        color: foreground,
      };
    },
    placeholder: (provided) => {
      return {
        ...provided,
        whiteSpace: "nowrap",
      };
    },
    group: (provided) => ({
      ...provided,
      borderTop: `1px solid ${blackScale3}`,
    }),
  };
};

export default useSelectInputStyles;
