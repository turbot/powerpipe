import { classNames } from "./styles";

const getIconForType = (type, icon) => {
  if (!type && !icon) {
    return null;
  }

  if (icon) {
    return icon;
  }

  switch (type) {
    case "alert":
      return "materialsymbols-solid:error";
    // return "materialsymbols-solid:priority_high";
    case "ok":
      // return "materialsymbols-solid:check_circle";
      return "materialsymbols-solid:check";
    case "info":
      // return "materialsymbols-solid:info";
      return "materialsymbols-solid:info_i";
    case "severity":
      return "materialsymbols-solid:warning";
    case "skip":
      return "materialsymbols-solid:arrow_right_alt";
    default:
      return null;
  }
};

const getIconStyles = (type) => {
  if (!type) {
    return {};
  }

  switch (type) {
    case "alert":
      return { fontWeight: "bold" };
    case "ok":
      return { fontVariationSettings: "'wght' 700" };
    case "info":
      return {};
    case "severity":
      return {};
    case "skip":
      return {};
    default:
      return {};
  }
};

const getIconClasses = (type) => {
  const baseClasses = "text-3xl opacity-100";
  switch (type) {
    case "alert":
      return classNames(baseClasses, "text-alert");
    case "info":
      return classNames(baseClasses, "text-info");
    case "ok":
      return classNames(baseClasses, "text-ok");
    case "severity":
      return classNames(baseClasses, "text-severity");
    default:
      return classNames(baseClasses, "text-skip");
  }
};

const getIconWrapperClasses = (type) => {
  const baseClasses = "flex items-center justify-center h-9 w-9 rounded-full";
  switch (type) {
    case "alert":
      return classNames(baseClasses, "bg-alert/15");
    case "info":
      return classNames(baseClasses, "bg-info/15");
    case "ok":
      return classNames(baseClasses, "bg-ok/15");
    case "severity":
      return classNames(baseClasses, "bg-severity/15");
    default:
      return classNames(baseClasses, "bg-skip/15");
  }
};

const getWrapperClasses = (type) => {
  const baseClasses = "border-l-4 rounded-r-md"; // rounded-l-md rounded-r-md rounded-b-md";
  // const baseClasses = "rounded-md"; // rounded-l-md rounded-r-md rounded-b-md";
  switch (type) {
    case "alert":
      // return classNames(baseClasses, "bg-alert/5");
      return classNames(baseClasses, "border-alert");
    case "info":
      // return classNames(baseClasses, "bg-info/5");
      return classNames(baseClasses, "border-info");
    case "ok":
      // return classNames(baseClasses, "bg-ok/5");
      return classNames(baseClasses, "border-ok");
    case "severity":
      // return classNames(baseClasses, "bg-severity/5");
      return classNames(baseClasses, "border-severity");
    case "skip":
      // return classNames(baseClasses, "bg-skip/5");
      return classNames(baseClasses, "border-skip");
    default:
      return "rounded-md";
  }
};

export {
  getIconClasses,
  getIconForType,
  getIconStyles,
  getIconWrapperClasses,
  getWrapperClasses,
};
