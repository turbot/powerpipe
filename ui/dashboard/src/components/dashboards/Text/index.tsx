import Error from "../Error";
import gfm from "remark-gfm"; // Support for strikethrough, tables, tasklists and URLs
import ReactMarkdown from "react-markdown";
import rehypeExternalLinks from "rehype-external-links";
import {
  BasePrimitiveProps,
  ExecutablePrimitiveProps,
} from "@powerpipe/components/dashboards/common";
import { classNames } from "@powerpipe/utils/styles";
import { PanelDefinition } from "@powerpipe/types";
import { registerComponent } from "../index";

const getLongPanelClasses = () => {
  // switch (type) {
  // case "alert":
  //   return "p-2 border border-alert bg-alert-light border overflow-hidden sm:rounded-md";
  // default:
  return "overflow-hidden sm:rounded-md";
  // }
};

const getShortPanelClasses = () => {
  // switch (type) {
  //   case "alert":
  //     return "p-2 border border-alert bg-alert-light prose prose-sm sm:rounded-md max-w-none";
  //   default:
  return "prose prose-sm sm:rounded-md max-w-none";
  // }
};

export type TextProps = PanelDefinition &
  BasePrimitiveProps &
  ExecutablePrimitiveProps & {
    display_type?: "raw" | "markdown" | "html";
    properties: {
      value: string;
    };
  };

const Markdown = ({ value }) => {
  if (!value) {
    return null;
  }

  // Use standard prose styles from Tailwind
  // Do not restrict width, that's the job of the wrapping panel
  const isLong = value.split("\n").length > 3;
  const panelClasses = isLong ? getLongPanelClasses() : getShortPanelClasses();
  const proseHeadings =
    "prose-h1:text-3xl prose-h2:text-2xl prose-h3:text-xl prose-h3:mt-1 p-4 text-foreground prose-h1:text-foreground prose-h2:text-foreground prose-h4:text-foreground prose-h4:text-foreground prose-h5:text-foreground prose-h6:text-foreground";
  const markdown = (
    <ReactMarkdown
      rehypePlugins={[[rehypeExternalLinks, { target: "_blank" }]]}
      remarkPlugins={[gfm]}
    >
      {value}
    </ReactMarkdown>
  );

  return (
    <>
      {isLong ? (
        <div className={panelClasses}>
          <div
            className={classNames(
              "p-2 sm:p-1 prose prose-sm max-w-none break-keep",
              proseHeadings,
            )}
          >
            {markdown}
          </div>
        </div>
      ) : (
        <article
          className={classNames(panelClasses, "break-keep", proseHeadings)}
        >
          {markdown}
        </article>
      )}
    </>
  );
};

const Raw = ({ value }) => {
  if (!value) {
    return null;
  }
  return (
    <pre className="whitespace-pre-wrap break-keep text-foreground">
      {value}
    </pre>
  );
};

const renderText = (type, value) => {
  switch (type) {
    case "markdown":
      return <Markdown value={value} />;
    case "raw":
      return <Raw value={value} />;
    default:
      return <Error error={`Unsupported text type ${type}`} />;
  }
};

const Text = (props: TextProps) =>
  renderText(
    props.display_type || "markdown",
    props.properties ? props.properties.value : null,
  );

registerComponent("text", Text);

export { Markdown };

export default Text;
