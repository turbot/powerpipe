import ErrorModal from "@powerpipe/components/Modal/ErrorModal";
import { useDashboardState } from "@powerpipe/hooks/useDashboardState";

type WorkspaceErrorModalProps = {
  error: any;
};

const WorkspaceErrorModal = ({ error }: WorkspaceErrorModalProps) => (
  <ErrorModal error={error} title="Workspace Error" />
);

const WorkspaceErrorModalWrapper = () => {
  const { error } = useDashboardState();
  if (!error) {
    return null;
  }
  return <WorkspaceErrorModal error={error} />;
};

export default WorkspaceErrorModalWrapper;
