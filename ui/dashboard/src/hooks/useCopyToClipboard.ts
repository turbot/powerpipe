import { useCallback, useEffect, useState } from "react";
import * as copy from "copy-to-clipboard";

const useCopyToClipboard = () => {
  const [copySuccess, setCopySuccess] = useState(false);

  const handleCopy = useCallback(
    (data) => {
      // @ts-ignore
      const copyOutput = copy(data);
      if (copyOutput) {
        setCopySuccess(true);
      }
    },
    [setCopySuccess],
  );

  useEffect(() => {
    let timeoutId;
    if (copySuccess) {
      timeoutId = setTimeout(() => {
        setCopySuccess(false);
      }, 1000);
    }
    return () => clearTimeout(timeoutId);
  }, [copySuccess]);

  return { copy: handleCopy, copySuccess };
};

export default useCopyToClipboard;
