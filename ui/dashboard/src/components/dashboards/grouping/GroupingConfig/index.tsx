import GroupingEditor from "../GroupingEditor";
import useGroupingConfig from "@powerpipe/hooks/useGroupingConfig";

type GroupingConfigProps = {
  panelName: string;
};

const GroupingConfig = ({ panelName }: GroupingConfigProps) => {
  const { grouping, update } = useGroupingConfig(panelName);

  return <GroupingEditor config={grouping} onApply={update} />;
};

export default GroupingConfig;
