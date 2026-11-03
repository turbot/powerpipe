import GroupingEditor from "../GroupingEditor";
import useGroupingConfig from "@powerpipe/hooks/useGroupingConfig";

type GroupingConfigProps = {
  panelName: string;
};

const GroupingConfig = ({ panelName }: GroupingConfigProps) => {
  const { grouping, defaultGrouping, update } = useGroupingConfig(panelName);

  return (
    <GroupingEditor
      config={grouping}
      defaultConfig={defaultGrouping}
      onApply={update}
    />
  );
};

export default GroupingConfig;
