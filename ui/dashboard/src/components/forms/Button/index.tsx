import { classNames } from "@powerpipe/utils/styles";
import { forwardRef, ReactNode } from "react";

export type ButtonProps = {
  children: ReactNode;
  className?: string;
  disabled?: boolean;
  onClick(): void | Promise<void>;
  size?: "sm" | "md" | "lg";
  title?: string;
  type?: "button" | "submit";
};

const Button = forwardRef(
  (
    {
      children,
      className,
      disabled,
      onClick = async () => {},
      size = "md",
      title,
      type = "button",
    }: ButtonProps,
    ref,
  ) => {
    let sizeClass;
    switch (size) {
      case "sm":
        sizeClass = "py-1 px-2 font-sm";
        break;
      case "lg":
        sizeClass = "py-3 px-6 font-sm";
        break;
      default:
        sizeClass = "py-2 px-3";
    }
    return (
      <button
        // @ts-ignore
        ref={ref}
        className={classNames(
          sizeClass,
          className,
          "rounded-md shadow-sm whitespace-nowrap focus:outline-none disabled:opacity-50 disabled:cursor-default",
        )}
        disabled={disabled}
        onClick={onClick}
        title={title}
        type={type}
      >
        {children}
      </button>
    );
  },
);

export default Button;
