import { useCallback, useEffect, useState } from "react";
import copy from "copy-to-clipboard";

const useCopyToClipboard = () => {
  const [copySuccess, setCopySuccess] = useState(false);

  const handleCopy = useCallback(
    async (data: string) => {
      // copy-to-clipboard 4 resolves to false when the copy failed. When the
      // browser refuses clipboard access, fall back to a prompt holding the
      // text, as version 3 did.
      const copied = await copy(data, { fallbackToPrompt: true });
      if (copied) {
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
