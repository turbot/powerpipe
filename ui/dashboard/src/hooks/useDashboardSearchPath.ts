import { useMemo } from "react";
import { useSearchParams } from "react-router-dom";

const useDashboardSearchPath = () => {
  const [searchParams] = useSearchParams();
  return useMemo(() => {
    const rawSearchPath = searchParams.get("search_path");
    const rawSearchPathPrefix = searchParams.get("search_path_prefix");

    // For now, search path wins over prefix
    if (
      (!!rawSearchPath && !!rawSearchPathPrefix) ||
      (!rawSearchPathPrefix && !!rawSearchPath)
    ) {
      try {
        return rawSearchPath.split(",");
      } catch (error) {
        console.error("Error parsing search path", error);
        return [];
      }
    } else if (!!rawSearchPathPrefix) {
      try {
        return rawSearchPathPrefix.split(",");
      } catch (error) {
        console.error("Error parsing search path prefix", error);
        return [];
      }
    } else {
      return [];
    }
  }, [searchParams]);
};

export default useDashboardSearchPath;
