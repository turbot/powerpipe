import { useEffect, useRef } from "react";

function usePrevious<T>(value: T): T {
  const ref: any = useRef<T | undefined>(undefined);
  useEffect(() => {
    ref.current = value;
  }, [value]);
  return ref.current;
}

export default usePrevious;
