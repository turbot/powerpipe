import { DashboardActions } from "../../types";
import { useDashboard } from "../../hooks/useDashboard";
import { useNavigate } from "react-router-dom";
import { useRef } from "react";

const SnapshotDiffButton = () => {
  const { dispatch, snapshot } = useDashboard();
  const fileInputRef = useRef<HTMLInputElement | null>(null);
  const navigate = useNavigate();

  if (!snapshot) {
    return null;
  }

  return (
    <>
      <button
        type="button"
        className="rounded-md bg-info px-2.5 py-1.5 text-sm font-semibold text-white"
        onClick={() => {
          fileInputRef.current?.click();
        }}
      >
        Diff
      </button>
      <input
        ref={fileInputRef}
        accept="application/json, .pps, .sps"
        className="hidden"
        id="snapshot-diff"
        name="snapshot-diff"
        type="file"
        onChange={(e) => {
          const files = e.target.files;
          if (!files || files.length === 0) {
            return;
          }
          const fr = new FileReader();
          fr.onload = () => {
            if (!fr.result) {
              return;
            }

            e.target.value = "";
            try {
              navigate(`/snapshot/diff`);
              const data = JSON.parse(fr.result.toString());
              dispatch({
                type: DashboardActions.GET_SNAPSHOT_DIFF,
                snapshot: data,
              });
            } catch (err: any) {
              console.log(err);
              dispatch({
                type: DashboardActions.WORKSPACE_ERROR,
                error: "Unable to load snapshot:" + err.message,
              });
            }
          };
          fr.readAsText(files[0]);
        }}
      />
    </>
  );
};

export default SnapshotDiffButton;
