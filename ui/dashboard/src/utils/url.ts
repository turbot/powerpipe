const isRelativeUrl = (url) => {
  return (
    new URL(document.baseURI).origin === new URL(url, document.baseURI).origin
  );
};

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

export { isRelativeUrl, injectSearchPathPrefix };
