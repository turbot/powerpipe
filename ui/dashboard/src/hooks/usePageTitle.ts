import useDeepCompareEffect from "use-deep-compare-effect";

const getPageTitle = (titleParts: (string | undefined | null)[] = []) => {
  return [...titleParts, "Powerpipe"].filter((v) => !!v).join(" | ");
};

const usePageTitle = (
  titleParts: (string | undefined | null)[] = [],
  skip = false,
) => {
  useDeepCompareEffect(() => {
    if (skip) {
      return;
    }
    document.title = getPageTitle(titleParts);
  }, [skip, titleParts]);
};

export default usePageTitle;

export { getPageTitle };
