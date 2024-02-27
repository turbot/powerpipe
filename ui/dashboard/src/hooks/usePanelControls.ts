import useDownloadPanelData from "./useDownloadPanelData";
import useSelectPanel from "./useSelectPanel";
import { IPanelControl } from "@powerpipe/components/dashboards/layout/Panel/PanelControls";
import { PanelDefinition } from "@powerpipe/types";
import { useCallback, useEffect, useState } from "react";

const usePanelControls = (definition: PanelDefinition, show = false) => {
  const { download } = useDownloadPanelData(definition);
  const { select } = useSelectPanel(definition);

  const downloadPanelData = useCallback(
    async (e) => {
      e.stopPropagation();
      await download();
    },
    [download],
  );

  const getBasePanelControls = useCallback(() => {
    const controls: IPanelControl[] = [];
    if (!show || !definition) {
      return controls;
    }
    if (definition.data) {
      controls.push({
        action: downloadPanelData,
        icon: "arrow-down-tray",
        title: "Download data",
      });
    }
    controls.push({
      action: select,
      icon: "arrows-pointing-out",
      title: "View detail",
    });
    return controls;
  }, [definition, downloadPanelData, select, show]);

  const [panelControls, setPanelControls] = useState(getBasePanelControls());
  const [customControls, setCustomControls] = useState<IPanelControl[]>([]);

  useEffect(() => {
    // console.log({
    //   customControls,
    //   definition,
    //   getBasePanelControls,
    //   setPanelControls,
    //   show,
    // });
    setPanelControls([...customControls, ...getBasePanelControls()]);
  }, [
    customControls,
    definition,
    getBasePanelControls,
    setPanelControls,
    show,
  ]);

  return { panelControls, setCustomControls };
};

export default usePanelControls;
