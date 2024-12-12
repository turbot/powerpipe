import NeutralButton from "@powerpipe/components/forms/NeutralButton";
import useDownloadPanelData from "@powerpipe/hooks/useDownloadPanelData";
import { noop } from "@powerpipe/utils/func";
import { useDashboardPanelDetail } from "@powerpipe/hooks/useDashboardPanelDetail";

const PanelDetailDataDownloadButton = ({ panelDefinition, size }) => {
  const { panelOverrideData } = useDashboardPanelDetail();
  const { download, processing } = useDownloadPanelData(panelDefinition);

  return (
    <NeutralButton
      disabled={processing}
      onClick={
        processing
          ? noop
          : () => download(panelOverrideData || panelDefinition.data)
      }
      size={size}
    >
      <>Download</>
    </NeutralButton>
  );
};

export default PanelDetailDataDownloadButton;
