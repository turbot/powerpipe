import Icon from "@powerpipe/components/Icon";
import { createPortal } from "react-dom";
import { Markdown } from "@powerpipe/components/dashboards/Text";
import { Popover } from "@headlessui/react";
import { ThemeProvider, ThemeWrapper } from "@powerpipe/hooks/useTheme";
import { useEffect, useState } from "react";
import { usePopper } from "react-popper";

const Documentation = ({
  documentation,
  onOpen,
  onClose,
}: {
  documentation: string;
  onOpen: () => void;
  onClose: () => void;
}) => {
  useEffect(() => {
    onOpen();
    return onClose;
  }, []);

  return (
    <div className="border border-dashboard rounded-md bg-dashboard-panel text-foreground mt-1 p-3 space-y-3 min-w-60 max-w-xl max-h-96 overflow-y-auto">
      <Markdown value={documentation} />
    </div>
  );
};

const DocumentationView = ({
  documentation,
  onOpen,
  onClose,
}: {
  documentation: string | undefined;
  onOpen: () => void;
  onClose: () => void;
}) => {
  const [popperElement, setPopperElement] = useState(null);
  const [referenceElement, setReferenceElement] = useState(null);
  const { styles, attributes } = usePopper(referenceElement, popperElement, {
    placement: "bottom-start",
    modifiers: [
      {
        name: "flip",
        options: {
          fallbackPlacements: ["top-start", "right-start"],
        },
      },
    ],
  });

  if (!documentation) {
    return null;
  }

  return (
    <Popover className="relative">
      {/*@ts-ignore*/}
      <Popover.Button ref={setReferenceElement} as="div">
        <Icon
          title="View documentation"
          icon="help"
          className="h-5 w-5 cursor-pointer"
        />
      </Popover.Button>
      <Popover.Overlay className="fixed inset-0 bg-black opacity-40" />
      <Popover.Panel className="absolute z-10 pt-px">
        {createPortal(
          <ThemeProvider>
            <ThemeWrapper>
              <div
                // @ts-ignore
                ref={setPopperElement}
                style={{ ...styles.popper }}
                {...attributes.popper}
                onClick={(e) => e.stopPropagation()}
              >
                <Documentation
                  documentation={documentation}
                  onOpen={onOpen}
                  onClose={onClose}
                />
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
