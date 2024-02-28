import { useMemo } from "react";
import { useSearchParams } from "react-router-dom";

const useDashboardSearchPathPrefix = () => {
  const [searchParams] = useSearchParams();

  return useMemo(() => {
    const rawSearchPathPrefix = searchParams.get("search_path_prefix");

    if (!!rawSearchPathPrefix) {
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

export default useDashboardSearchPathPrefix;
