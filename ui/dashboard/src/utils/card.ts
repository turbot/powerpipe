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
    case "severity":
      return "materialsymbols-solid:warning";
    case "skip":
      return "materialsymbols-solid:arrow_circle_right";
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

const getWrapperClasses = (type) => {
  const baseClasses = "border-l-4 rounded-r-md";
  switch (type) {
    case "alert":
      return classNames(baseClasses, "border-alert");
    case "info":
      return classNames(baseClasses, "border-info");
    case "ok":
      return classNames(baseClasses, "border-ok");
    case "severity":
      return classNames(baseClasses, "border-severity");
    case "skip":
      return classNames(baseClasses, "border-skip");
    default:
      return "rounded-md";
  }
};

export { getIconClasses, getIconForType, getIconStyles, getWrapperClasses };
