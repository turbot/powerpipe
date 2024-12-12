import { useMemo } from "react";
import { useSearchParams } from "react-router-dom";

const usePrefixedSearchParams = (prefix: string) => {
  const [searchParams] = useSearchParams();

  return useMemo(() => {
    const result: Record<string, string> = {};
    for (const [key, value] of searchParams.entries()) {
      if (key.startsWith(prefix)) {
        result[key] = value;
      }
    }
    return result;
  }, [searchParams, prefix]);
};

export default usePrefixedSearchParams;
