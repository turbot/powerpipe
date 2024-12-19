import copy from "copy-to-clipboard";
import { classNames } from "@powerpipe/utils/styles";
import {
  createContext,
  useCallback,
  useContext,
  useEffect,
  useState,
} from "react";
import Icon from "@powerpipe/components/Icon";

type ICopyToClipboardContext = {
  doCopy: boolean;
  setDoCopy: (value: boolean) => void;
};

const CopyToClipboardContext = createContext<ICopyToClipboardContext | null>(
  null,
);

const CopyToClipboardProvider = ({ children }) => {
  const [doCopy, setDoCopy] = useState(false);
  return (
    <CopyToClipboardContext.Provider value={{ doCopy, setDoCopy }}>
      {children({ setDoCopy })}
    </CopyToClipboardContext.Provider>
  );
};

const CopyToClipboard = ({
  data,
  className = "text-foreground-light",
  stopPropagation = true,
}) => {
  const context = useContext(CopyToClipboardContext);
  const { doCopy, setDoCopy } = context
    ? context
    : ({} as ICopyToClipboardContext);
  const [copySuccess, setCopySuccess] = useState(false);

  const handleCopy = useCallback(
    async (e) => {
      if (e && stopPropagation) {
        e.stopPropagation();
      }
      const copyOutput = copy(data);
      if (copyOutput) {
        setCopySuccess(true);
      }
    },
    [data, setCopySuccess, stopPropagation],
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

  useEffect(() => {
    const triggerCopy = async () => {
      // @ts-ignore
      await handleCopy();
      setDoCopy(false);
    };
    if (doCopy) {
      triggerCopy();
    }
  }, [handleCopy, doCopy, setDoCopy]);

  return (
    <>
      {!copySuccess && (
        <Icon
          icon="content_copy"
          className={classNames("h-5 w-5 cursor-pointer", className)}
          onClick={(e) => handleCopy(e)}
        />
      )}
      {copySuccess && (
        <Icon
          icon="materialsymbols-solid:content_copy"
          className={classNames("h-5 w-5 text-ok", className)}
        />
      )}
    </>
  );
};

export default CopyToClipboard;

export { CopyToClipboardProvider };
