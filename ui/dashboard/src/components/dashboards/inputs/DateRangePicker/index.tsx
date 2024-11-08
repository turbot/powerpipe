import dayjs from "dayjs";
import useDeepCompareEffect from "use-deep-compare-effect";
import utc from "dayjs/plugin/utc";
import { DashboardActions } from "@powerpipe/types";
import { DayPicker } from "react-day-picker";
import {
  IInput,
  InputProps,
} from "@powerpipe/components/dashboards/inputs/types";
import { registerInputComponent } from "@powerpipe/components/dashboards/inputs";
import { useDashboard } from "@powerpipe/hooks/useDashboard";
import { useState, useRef, useEffect } from "react";
import "react-day-picker/dist/style.css";
import "react-time-picker/dist/TimePicker.css";
import "./DateRangePicker.css";
dayjs.extend(utc);

const presets = [
  { label: "1h", value: "1h" },
  { label: "3h", value: "3h" },
  { label: "6h", value: "6h" },
  { label: "12h", value: "12h" },
  { label: "1d", value: "1d" },
  { label: "7d", value: "7d" },
  { label: "Custom", value: "custom" },
];

const timeOptions = {
  minutes: [5, 10, 15, 30, 45],
  hours: [1, 2, 3, 6, 8, 12],
  days: [1, 2, 3, 4, 5, 6],
  weeks: [1, 2, 3, 4],
};

const DateRangePicker = (props: InputProps) => {
  const { dispatch, selectedDashboardInputs } = useDashboard();
  const stateValue = selectedDashboardInputs[props.name];

  const [state, setState] = useState<{
    from: dayjs.Dayjs;
    to?: dayjs.Dayjs | null;
    relative?: string | null;
    showCustom?: boolean;
  }>(() => {
    const stateValue = selectedDashboardInputs[props.name];
    if (stateValue) {
      try {
        const parsed = JSON.parse(stateValue);
        const fallback = dayjs();
        return {
          from: parsed.from
            ? dayjs(parsed.from)
            : fallback.subtract(1, "day").utc(),
          to: parsed.to ? dayjs(parsed.to) : null,
          relative: parsed.relative || "1d",
        };
      } catch (err) {
        console.log("Parse error", err);
        const now = dayjs();
        return {
          from: now.subtract(1, "day").utc(),
          to: null,
          relative: "1d",
        };
      }
    } else {
      const now = dayjs();
      return {
        from: now.subtract(1, "day").utc(),
        to: null,
        relative: "1d",
      };
    }
  });

  useEffect(() => {
    if (stateValue) {
      return;
    }
    dispatch({
      type: DashboardActions.SET_DASHBOARD_INPUT,
      name: props.name,
      value: JSON.stringify({
        from: dayjs().subtract(1, "day").utc(),
        to: null,
        relative: "1d",
      }),
      recordInputsHistory: !!stateValue,
    });
  }, []);

  useDeepCompareEffect(() => {
    if (state.showCustom) {
      return;
    }
    dispatch({
      type: DashboardActions.SET_DASHBOARD_INPUT,
      name: props.name,
      value: JSON.stringify({
        from: state.from,
        to: state.to,
        relative: state.relative,
      }),
      recordInputsHistory: !!stateValue,
    });
  }, [state]);

  const [tempState, setTempState] = useState<{
    from: dayjs.Dayjs;
    to?: dayjs.Dayjs | null;
    relative?: string | null;
    showCustom?: boolean;
  }>(state);

  useEffect(() => {
    setTempState(() => state);
  }, [state]);

  const [tab, setTab] = useState("relative");
  const [duration, setDuration] = useState(1);
  const [unitOfTime, setUnitOfTime] = useState("hours");
  const customButtonRef = useRef(null);

  const handlePresetChange = (preset) => {
    switch (preset) {
      case "1h":
        setDuration(1);
        setUnitOfTime("hour");
        setState((previous) => ({
          ...previous,
          from: dayjs().subtract(1, "hour").utc(),
          to: null,
          relative: "1h",
          showCustom: false,
        }));
        break;
      case "3h":
        setDuration(3);
        setUnitOfTime("hour");
        setState((previous) => ({
          ...previous,
          from: dayjs().subtract(3, "hour").utc(),
          to: null,
          relative: "3h",
          showCustom: false,
        }));
        break;
      case "6h":
        setDuration(6);
        setUnitOfTime("hour");
        setState((previous) => ({
          ...previous,
          from: dayjs().subtract(6, "hour").utc(),
          to: null,
          relative: "6h",
          showCustom: false,
        }));
        break;
      case "12h":
        setDuration(12);
        setUnitOfTime("hour");
        setState((previous) => ({
          ...previous,
          from: dayjs().subtract(12, "hour").utc(),
          to: null,
          relative: "12h",
          showCustom: false,
        }));
        break;
      case "1d":
        setDuration(1);
        setUnitOfTime("day");
        setState((previous) => ({
          ...previous,
          from: dayjs().subtract(1, "day").utc(),
          to: null,
          relative: "1d",
          showCustom: false,
        }));
        break;
      case "7d":
        setDuration(7);
        setUnitOfTime("day");
        setState((previous) => ({
          ...previous,
          from: dayjs().subtract(7, "day").utc(),
          to: null,
          relative: "7d",
          showCustom: false,
        }));
        break;
      case "custom":
        setState((previous) => ({
          ...previous,
          showCustom: true,
        }));
    }
  };

  const handleApply = () => {
    setState({ ...tempState, showCustom: false });
  };

  const handleCancel = () => {
    setState({ ...state, showCustom: false });
  };

  const handleTimeOptionClick = (value, unit) => {
    setDuration(value);
    setUnitOfTime(unit);
    setTempState((previous) => ({
      ...previous,
      from: dayjs().subtract(value, unit).utc(),
      to: null,
      relative: `${value}${unit === "minute" ? "m" : unit === "hour" ? "h" : unit === "day" ? "d" : unit === "week" ? "w" : ""}`,
      showCustom: false,
    }));
  };

  return (
    <div className="flex flex-col max-w-23vw">
      <div className="presets">
        {presets.map((preset) => (
          <button
            key={preset.value}
            onClick={() => handlePresetChange(preset.value)}
            className={`preset-button ${state.relative === preset.value ? "active" : ""}`}
            ref={preset.value === "custom" ? customButtonRef : null}
          >
            {preset.label}
          </button>
        ))}
      </div>

      {state.showCustom && (
        <div
          className="custom-popover"
          style={{
            position: "absolute",
            top:
              customButtonRef.current?.getBoundingClientRect().bottom +
              window.scrollY,
            left:
              customButtonRef.current?.getBoundingClientRect().left +
              window.scrollX,
            zIndex: 1000,
            border: "1px solid #ddd",
            borderRadius: "5px",
            backgroundColor: "#fff",
            padding: "20px",
            boxShadow: "0 2px 8px rgba(0, 0, 0, 0.15)",
          }}
        >
          <div className="tabs">
            <button
              className={`tab-button ${tab === "relative" ? "active" : ""}`}
              onClick={() => setTab("relative")}
            >
              Relative
            </button>
            <button
              className={`tab-button ${tab === "absolute" ? "active" : ""}`}
              onClick={() => setTab("absolute")}
            >
              Absolute
            </button>
          </div>

          {tab === "absolute" ? (
            <div className="absolute-panel">
              <div className="calendar-container">
                <DayPicker
                  mode="range"
                  selected={{
                    from: tempState.from.utc().toDate(),
                    to: tempState.to?.utc().toDate(),
                  }}
                  onSelect={({ from, to }) => {
                    const newFrom = new Date(
                      from.getFullYear(),
                      from.getMonth(),
                      from.getDate(),
                      tempState.from.hour(),
                      tempState.from.minute(),
                      tempState.from.second(),
                    );
                    const newTo = new Date(
                      to.getFullYear(),
                      to.getMonth(),
                      to.getDate(),
                      tempState.to?.hour() || 0,
                      tempState.to?.minute() || 0,
                      tempState.to?.second() || 0,
                    );
                    const parsedFrom = dayjs(newFrom).utc();
                    const parsedTo = dayjs(newTo).utc();
                    setTempState((previous) => ({
                      ...previous,
                      from: parsedFrom,
                      to: parsedTo,
                    }));
                  }}
                  className="single-day-picker"
                  captionLayout="dropdown"
                  pagedNavigation
                />
              </div>
              <div className="time-inputs">
                <div className="time-input-container">
                  <label>Start date</label>
                  <input
                    type="date"
                    value={state.from.format("YYYY-MM-DD")}
                    onChange={(e) =>
                      setTempState((previous) => ({
                        ...previous,
                        from: dayjs(e.target.value).utc(),
                      }))
                    }
                  />
                  <label>Start time</label>
                  <input
                    type="time"
                    value={`${tempState.from.hour()}:${tempState.from.minute()}:${tempState.from.second()}`}
                    step="1"
                    onChange={(e) => {
                      setTempState((previous) => ({
                        ...previous,
                        from: dayjs(
                          `${tempState.from.format("YYYY")}-${tempState.from.format("MM")}-${tempState.from.format("DD")} ${e.target.value}`,
                        ),
                      }));
                    }}
                  />
                </div>
                <div className="time-input-container">
                  <label>End date</label>
                  <input
                    type="date"
                    value={tempState.to?.format("YYYY-MM-DD") || undefined}
                    onChange={(e) =>
                      setTempState((previous) => ({
                        ...previous,
                        to: dayjs(e.target.value).utc(),
                      }))
                    }
                  />
                  <label>End time</label>
                  <input
                    type="time"
                    value={
                      tempState.to
                        ? `${tempState.to.hour()}:${tempState.to.minute()}:${tempState.to.second()}`
                        : `00:00:00`
                    }
                    step="1"
                    onChange={(e) => {
                      const toTime = tempState.to || dayjs();
                      setTempState((previous) => ({
                        ...previous,
                        from: dayjs(
                          `${toTime.format("YYYY")}-${toTime.format("MM")}-${toTime.format("DD")} ${e.target.value}`,
                        ),
                      }));
                    }}
                  />
                </div>
              </div>
            </div>
          ) : (
            <div className="relative-panel">
              <div className="time-option-groups">
                <label>Minutes</label>
                <div className="option-group">
                  {timeOptions.minutes.map((min) => (
                    <button
                      key={min}
                      onClick={() => handleTimeOptionClick(min, "minute")}
                      className={`time-option ${duration === min && unitOfTime === "minute" ? "active" : ""}`}
                    >
                      {min}
                    </button>
                  ))}
                </div>

                <label>Hours</label>
                <div className="option-group">
                  {timeOptions.hours.map((hour) => (
                    <button
                      key={hour}
                      onClick={() => handleTimeOptionClick(hour, "hour")}
                      className={`time-option ${duration === hour && unitOfTime === "hour" ? "active" : ""}`}
                    >
                      {hour}
                    </button>
                  ))}
                </div>

                <label>Days</label>
                <div className="option-group">
                  {timeOptions.days.map((day) => (
                    <button
                      key={day}
                      onClick={() => handleTimeOptionClick(day, "day")}
                      className={`time-option ${duration === day && unitOfTime === "day" ? "active" : ""}`}
                    >
                      {day}
                    </button>
                  ))}
                </div>

                <label>Weeks</label>
                <div className="option-group">
                  {timeOptions.weeks.map((week) => (
                    <button
                      key={week}
                      onClick={() => handleTimeOptionClick(week, "week")}
                      className={`time-option ${duration === week && unitOfTime === "week" ? "active" : ""}`}
                    >
                      {week}
                    </button>
                  ))}
                </div>
              </div>

              <div className="duration-input">
                <label>Duration</label>
                <input
                  type="number"
                  min={1}
                  max={999999999999}
                  value={duration}
                  onChange={(e) => setDuration(Number(e.target.value))}
                />
                <select
                  value={unitOfTime}
                  onChange={(e) => setUnitOfTime(e.target.value)}
                >
                  <option value="minute">Minutes</option>
                  <option value="hour">Hours</option>
                  <option value="day">Days</option>
                  <option value="week">Weeks</option>
                </select>
              </div>
            </div>
          )}

          <div className="actions">
            <button className="apply-button" onClick={handleApply}>
              Apply
            </button>
            <button className="cancel-button" onClick={handleCancel}>
              Cancel
            </button>
          </div>
        </div>
      )}
    </div>
  );
};

const definition: IInput = {
  type: "date_range",
  component: DateRangePicker,
};

registerInputComponent(definition.type, definition);

export default definition;

export { DateRangePicker };
