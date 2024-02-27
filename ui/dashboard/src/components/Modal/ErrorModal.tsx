import ErrorMessage from "@powerpipe/components/ErrorMessage";
import Modal from "@powerpipe/components/Modal";
import NeutralButton from "@powerpipe/components/forms/NeutralButton";
import { ErrorIcon } from "@powerpipe/constants/icons";
import { useState } from "react";

const ErrorModal = ({ error, title }) => {
  const [show, setShow] = useState(true);
  return show ? (
    <Modal
      actions={[
        <NeutralButton key="close" onClick={() => setShow(false)}>
          <>Close</>
        </NeutralButton>,
      ]}
      icon={<ErrorIcon className="h-8 w-8 text-red-600" aria-hidden="true" />}
      children={
        <p className="w-full sm:w-11/12 text-sm text-foreground-light break-words whitespace-pre-wrap">
          <div className="break-all">
            <ErrorMessage error={error} />
          </div>
        </p>
      }
      onClose={async () => {
        setShow(false);
      }}
      title={title}
    />
  ) : null;
};

export default ErrorModal;
