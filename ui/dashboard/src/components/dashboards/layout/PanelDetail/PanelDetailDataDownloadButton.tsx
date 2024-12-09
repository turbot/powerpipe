import NeutralButton from "@powerpipe/components/forms/NeutralButton";
import useDownloadPanelData from "@powerpipe/hooks/useDownloadPanelData";
import { noop } from "@powerpipe/utils/func";

const PanelDetailDataDownloadButton = ({ panelDefinition, size }) => {
  const { download, processing } = useDownloadPanelData(
    panelDefinition,
    panelDefinition.data,
  );

  return (
    <NeutralButton
      disabled={processing}
      onClick={processing ? noop : () => download()}
      size={size}
    >
      <>Download</>
    </NeutralButton>
  );
};

export default PanelDetailDataDownloadButton;
