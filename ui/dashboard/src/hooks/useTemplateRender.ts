import {
  KeyValuePairs,
  TemplatesMap,
} from "@powerpipe/components/dashboards/common/types";
import { renderInterpolatedTemplates } from "@powerpipe/utils/template";
import { useCallback, useEffect, useState } from "react";

const useTemplateRender = () => {
  const [jq, setJq] = useState<any | null>(null);

  // Dynamically import jq-web from its own bundle
  useEffect(() => {
    const loadJq = async () => {
      try {
        const module = await import("jq-wasm");
        setJq(module);
      } catch (error) {
        console.error("Error loading jq-web:", error);
      }
    };

    loadJq();
  }, []);

  const renderTemplates = useCallback(
    async (templates: TemplatesMap, data: KeyValuePairs[]) => {
      if (!jq) {
        return [];
      }
      return renderInterpolatedTemplates(templates, data, jq);
    },
    [jq],
  );

  return {
    renderTemplates,
    ready: !!jq,
  };
};

export default useTemplateRender;
