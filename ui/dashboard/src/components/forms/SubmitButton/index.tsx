import Button, { ButtonProps } from "../Button";
import { classNames } from "@powerpipe/utils/styles";
import { forwardRef } from "react";

const SubmitButton = forwardRef(
  (
    {
      children,
      className = "",
      disabled = false,
      onClick,
      size = "md",
      title,
      type,
    }: ButtonProps,
    ref,
  ) => {
    return (
      <Button
        ref={ref}
        className={classNames(
          "bg-info border border-info hover:bg-steampipe-blue-dark hover:border-steampipe-blue-dark disabled:bg-info",
          className,
        )}
        disabled={disabled}
        onClick={onClick}
        size={size}
        title={title}
        type={type}
      >
        {children}
      </Button>
    );
  },
);

export default SubmitButton;
