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
      return "materialsymbols-solid:circle_notifications";
    case "ok":
      return "materialsymbols-solid:check_circle";
    case "info":
      return "materialsymbols-solid:info";
    case "skip":
      return "materialsymbols-solid:arrow_circle_right";
    case "severity-critical":
      return "materialsymbols-solid:pulse_alert";
    case "severity-high":
      return "materialsymbols-solid:warning";
    case "severity-medium":
      return "materialsymbols-solid:campaign";
    case "severity-low":
      return "materialsymbols-solid:info";
    case "severity":
      return "materialsymbols-solid:warning";
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
    case "severity-critical":
      return classNames(baseClasses, "text-alert");
    case "severity-high":
      return classNames(baseClasses, "text-orange");
    case "severity-medium":
      return classNames(baseClasses, "text-severity");
    case "severity-low":
      return classNames(baseClasses, "text-info");
    default:
      return classNames(baseClasses, "text-skip");
  }
};

const getWrapperClasses = (type) => {
  const baseClasses = "border-l-4 rounded-r-md";
  switch (type) {
    case "alert":
      return classNames(baseClasses, "border-alert");
    case "info":
      return classNames(baseClasses, "border-info");
    case "ok":
      return classNames(baseClasses, "border-ok");
    case "skip":
      return classNames(baseClasses, "border-skip");
    case "severity-critical":
      return classNames(baseClasses, "border-alert");
    case "severity-high":
      return classNames(baseClasses, "border-orange");
    case "severity-medium":
      return classNames(baseClasses, "border-severity");
    case "severity-low":
      return classNames(baseClasses, "border-info");
    case "severity":
      return classNames(baseClasses, "border-severity");
    default:
      return "rounded-md";
  }
};

export { getIconClasses, getIconForType, getIconStyles, getWrapperClasses };
