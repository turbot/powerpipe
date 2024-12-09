import Icon from "@powerpipe/components/Icon";
import { classNames } from "@powerpipe/utils/styles";
import { createPortal } from "react-dom";
import { ReactNode, useMemo, useState } from "react";
import { ThemeProvider, ThemeWrapper } from "@powerpipe/hooks/useTheme";
import { usePopper } from "react-popper";

export interface PanelControlProps {
  action: (e: any) => Promise<void>;
  component?: ReactNode;
  disabled?: boolean;
  icon: string;
  title: string;
}

const PanelControl = ({
  action,
  component,
  disabled,
  icon,
  title,
}: PanelControlProps) => {
  return (
    <div
      className={classNames(
        "flex items-center space-x-2 px-2 py-1.5 bg-dashboard-panel first:rounded-tl-[4px] first:rounded-bl-[4px] last:rounded-tr-[4px] last:rounded-br-[4px] hover:bg-dashboard",
        disabled
          ? "cursor-not-allowed text-foreground-light"
          : "cursor-pointer text-foreground",
      )}
      onClick={async (e) => {
        e.stopPropagation();
        if (disabled) {
          return;
        }
        await action(e);
      }}
      title={title}
    >
      {component}
      {icon && <Icon className="w-4.5 h-4.5" icon={icon} />}
    </div>
  );
};

const PanelControls = ({ controls, referenceElement, withOffset = false }) => {
  const [popperElement, setPopperElement] = useState(null);
  // Need to define memoized / stable modifiers else the usePopper hook will infinitely re-render
  const noFlip = useMemo(() => ({ name: "flip", enabled: false }), []);
  const offset = useMemo(() => {
    return {
      name: "offset",
      options: {
        // For some reason the height of the popper is not correct unless scrollbars are visible.
        // I've sunk too much time trying to find the root cause, but luckily I only
        // need to modify this along a fixed offset, so can hard-code this for now.
        offset: withOffset ? [-14.125, -14.125] : [0, -28.25],
        // offset: ({ popper }) => {
        // const offset = -popper.height / 2;
        // return [offset, offset];
        // },
      },
    };
  }, [withOffset]);
  const { styles, attributes } = usePopper(referenceElement, popperElement, {
    modifiers: [noFlip, offset],
    placement: "top-end",
  });

  return (
    <>
      {createPortal(
        <ThemeProvider>
          <ThemeWrapper>
            <div
              // @ts-ignore
              ref={setPopperElement}
              style={{ ...styles.popper }}
              {...attributes.popper}
            >
              <div className="flex border border-black-scale-3 rounded-md divide-x divide-divide">
                {controls.map((control, idx) => (
                  <PanelControl
                    key={idx}
                    disabled={control.disabled}
                    action={control.action}
                    component={control.component}
                    icon={control.icon}
                    title={control.title}
                  />
                ))}
              </div>
            </div>
          </ThemeWrapper>
        </ThemeProvider>,
        // @ts-ignore as this element definitely exists
        document.getElementById("portals"),
      )}
    </>
  );
};

export { PanelControl };

export default PanelControls;
