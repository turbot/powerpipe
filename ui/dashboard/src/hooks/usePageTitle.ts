import useDeepCompareEffect from "use-deep-compare-effect";

const getPageTitle = (titleParts: (string | undefined | null)[] = []) => {
  return [...titleParts, "Powerpipe"].filter((v) => !!v).join(" | ");
};

const usePageTitle = (titleParts: (string | undefined | null)[] = []) => {
  useDeepCompareEffect(() => {
    document.title = getPageTitle(titleParts);
  }, [titleParts]);
};

export default usePageTitle;

export { getPageTitle };
