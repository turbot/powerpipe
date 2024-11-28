import Icon from "@powerpipe/components/Icon";
import { createPortal } from "react-dom";
import { Markdown } from "@powerpipe/components/dashboards/Text";
import { Popover } from "@headlessui/react";
import { ThemeProvider, ThemeWrapper } from "@powerpipe/hooks/useTheme";
import { useState } from "react";
import { usePopper } from "react-popper";

const DocumentationView = ({
  documentation,
}: {
  documentation: string | undefined;
}) => {
  const [popperElement, setPopperElement] = useState(null);
  const [referenceElement, setReferenceElement] = useState(null);
  const { styles, attributes } = usePopper(referenceElement, popperElement, {
    placement: "bottom-start",
  });

  if (!documentation) {
    return null;
  }

  return (
    <Popover className="relative">
      {/*@ts-ignore*/}
      <Popover.Button ref={setReferenceElement} as="div">
        <Icon icon="info" className="h-4 w-4 cursor-pointer" />
      </Popover.Button>
      <Popover.Panel className="absolute z-10 pt-px">
        {createPortal(
          <ThemeProvider>
            <ThemeWrapper>
              <div
                // @ts-ignore
                ref={setPopperElement}
                style={{ ...styles.popper }}
                {...attributes.popper}
              >
                <div className="border border-dashboard rounded-md bg-dashboard-panel mt-1 p-3 space-y-3 min-w-60 max-w-96">
                  <Markdown value={documentation} />
                </div>
              </div>
            </ThemeWrapper>
          </ThemeProvider>,
          // @ts-ignore as this element definitely exists
          document.getElementById("portals"),
        )}
      </Popover.Panel>
    </Popover>
  );
};

export default DocumentationView;
