import { DatetimeRange } from "@powerpipe/hooks/useDashboardDatetimeRange";

const isRelativeUrl = (url) => {
  return (
    new URL(document.baseURI).origin === new URL(url, document.baseURI).origin
  );
};

const injectTimeRange = (url: string, timeRange: DatetimeRange | null | undefined) => {
  if (!timeRange) {
    return url;
  }
  return injectSearchParam(url, "time_range", JSON.stringify(timeRange));
}

const injectSearchPathPrefix = (url: string, searchPathPrefix: string[]) => {
  if (!searchPathPrefix.length) {
    return url;
  }
  return injectSearchParam(
    url,
    "search_path_prefix",
    searchPathPrefix.join(","),
  );
};

const injectSearchParam = (url: string, key: string, value: string) => {
  let parsedUrl: URL;
  if (isRelativeUrl(url)) {
    parsedUrl = new URL(url, document.baseURI);
  } else {
    parsedUrl = new URL(url);
  }
  parsedUrl.searchParams.set(key, value);
  return `${parsedUrl.pathname}${parsedUrl.search}`;
};

export { isRelativeUrl, injectSearchPathPrefix, injectTimeRange };
