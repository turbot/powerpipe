// import DatePickerProvider, {
//   Title,
//   Header,
//   WeekDays,
//   DaySlots,
// } from "headless-react-datepicker";
import dayjs from "dayjs";
import NeutralButton from "@powerpipe/components/forms/NeutralButton";
import SubmitButton from "@powerpipe/components/forms/SubmitButton";
import useDeepCompareEffect from "use-deep-compare-effect";
import utc from "dayjs/plugin/utc";
import { classNames } from "@powerpipe/utils/styles";
import { DashboardActions } from "@powerpipe/types";
import { DayPicker, getDefaultClassNames } from "react-day-picker";
import {
  IInput,
  InputProps,
} from "@powerpipe/components/dashboards/inputs/types";
import { Popover, Tab } from "@headlessui/react";
import { registerInputComponent } from "@powerpipe/components/dashboards/inputs";
import { useDashboard } from "@powerpipe/hooks/useDashboard";
import { useEffect, useState } from "react";
// import "headless-react-datepicker/dist/styles.css";
import "react-day-picker/dist/style.css";
import "react-time-picker/dist/TimePicker.css";
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

const CustomDatePicker = ({
  duration,
  tempState,
  unitOfTime,
  setDuration,
  setTempState,
  setUnitOfTime,
  onApply,
  onCancel,
  onTimeOptionClick,
}) => {
  const defaultDayPickerClassNames = getDefaultClassNames();
  const [tab, setTab] = useState("relative");
  const tabClasses = (selected) =>
    classNames(
      "flex-1 py-2 cursor-pointer bg-dashboard-panel border text-center rounded-md focus:outline-none",
      selected
        ? "font-bold border border-divide"
        : "text-foreground-light font-light border-dashboard-panel",
    );
  const presetClasses = (selected) =>
    classNames(
      "py-1.5 px-2.5 rounded-md cursor-pointer border bg-dashboard-panel",
      selected
        ? "font-bold border-divide"
        : "text-foreground-light font-light border-dashboard-panel",
    );

  return (
    <div className="border border-dashboard-panel rounded-md bg-dashboard p-3 space-y-3">
      <Tab.Group
        selectedIndex={tab === "relative" ? 0 : 1}
        onChange={(index) => setTab(index === 0 ? "relative" : "absolute")}
      >
        <Tab.List className="flex gap-[10px] mb-[10px]">
          <Tab className={({ selected }) => tabClasses(selected)}>Relative</Tab>
          <Tab className={({ selected }) => tabClasses(selected)}>Absolute</Tab>
        </Tab.List>
        <Tab.Panels>
          <Tab.Panel>
            {/* Content for Relative Tab */}
            <div className="space-y-3">
              <div className="space-y-1">
                <label>Minutes</label>
                <div className="flex space-x-2">
                  {timeOptions.minutes.map((min) => (
                    <div
                      key={min}
                      onClick={() => onTimeOptionClick(min, "minute")}
                      className={presetClasses(
                        duration === min && unitOfTime === "minute",
                      )}
                    >
                      {min}
                    </div>
                  ))}
                </div>
              </div>
              <div className="space-y-1">
                <label>Hours</label>
                <div className="flex space-x-2">
                  {timeOptions.hours.map((hour) => (
                    <div
                      key={hour}
                      onClick={() => onTimeOptionClick(hour, "hour")}
                      className={presetClasses(
                        duration === hour && unitOfTime === "hour",
                      )}
                    >
                      {hour}
                    </div>
                  ))}
                </div>
              </div>
              <div className="space-y-1">
                <label>Days</label>
                <div className="flex space-x-2">
                  {timeOptions.days.map((day) => (
                    <div
                      key={day}
                      onClick={() => onTimeOptionClick(day, "day")}
                      className={presetClasses(
                        duration === day && unitOfTime === "day",
                      )}
                    >
                      {day}
                    </div>
                  ))}
                </div>
              </div>
              <div className="space-y-1">
                <label>Weeks</label>
                <div className="flex space-x-2">
                  {timeOptions.weeks.map((week) => (
                    <div
                      key={week}
                      onClick={() => onTimeOptionClick(week, "week")}
                      className={presetClasses(
                        duration === week && unitOfTime === "week",
                      )}
                    >
                      {week}
                    </div>
                  ))}
                </div>
              </div>
              <div className="flex items-center space-x-3">
                <label>Duration</label>
                <input
                  type="number"
                  min={1}
                  max={999999999999}
                  value={duration}
                  onChange={(e) => setDuration(Number(e.target.value))}
                  className="flex-grow border border-divide rounded-md p-2 bg-dashboard-panel"
                />
                <select
                  value={unitOfTime}
                  onChange={(e) => setUnitOfTime(e.target.value)}
                  className="block p-2 border border-divide rounded-md bg-dashboard"
                >
                  <option value="minute">Minutes</option>
                  <option value="hour">Hours</option>
                  <option value="day">Days</option>
                  <option value="week">Weeks</option>
                </select>
              </div>
            </div>
          </Tab.Panel>
          <Tab.Panel>
            {/* Content for Absolute Tab */}
            <div className="space-y-3">
              <div className="">
                {/*<DatePickerProvider isRange>*/}
                {/*  <Title />*/}
                {/*  <Header />*/}
                {/*  <WeekDays />*/}
                {/*  <DaySlots />*/}
                {/*</DatePickerProvider>*/}
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
                  className="mx-auto bg-dashboard-panel dark:bg-dashboard text-foreground dark:text-foreground-light rounded-md p-2"
                  classNames={{
                    months_dropdown: classNames(
                      defaultDayPickerClassNames.months_dropdown,
                      "text-sm",
                    ),
                    months_caption: classNames(
                      defaultDayPickerClassNames.months_caption,
                      "text-sm",
                    ),
                    years_dropdown: classNames(
                      defaultDayPickerClassNames.years_dropdown,
                      "text-sm",
                    ),
                    range_start: classNames(
                      defaultDayPickerClassNames.range_start,
                      "bg-dashboard",
                    ),
                    range_middle: classNames(
                      defaultDayPickerClassNames.range_middle,
                      "bg-dashboard",
                    ),
                    range_end: classNames(
                      defaultDayPickerClassNames.range_end,
                      "bg-dashboard",
                    ),
                  }}
                  captionLayout="dropdown"
                  pagedNavigation
                />
              </div>
              <div className="flex space-x-2">
                <div className="flex-grow space-y-3">
                  <div className="">
                    <label>Start date</label>
                    <input
                      type="date"
                      value={tempState.from.format("YYYY-MM-DD")}
                      onChange={(e) =>
                        setTempState((previous) => ({
                          ...previous,
                          from: dayjs(e.target.value).utc(),
                        }))
                      }
                      className="bg-dashboard-panel text-foreground dark:bg-dashboard dark:text-foreground-light border border-table-border rounded p-2 w-full"
                    />
                  </div>
                  <div>
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
                      className="bg-dashboard-panel text-foreground dark:bg-dashboard dark:text-foreground-light border border-table-border rounded p-2 w-full"
                    />
                  </div>
                </div>
                <div className="flex-grow space-y-3">
                  <div>
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
                      className="bg-dashboard-panel text-foreground dark:bg-dashboard dark:text-foreground-light border border-table-border rounded p-2 w-full"
                    />
                  </div>
                  <div className="w-full">
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
                      className="bg-dashboard-panel text-foreground dark:bg-dashboard dark:text-foreground-light border border-table-border rounded p-2 w-full"
                    />
                  </div>
                </div>
              </div>
            </div>
          </Tab.Panel>
        </Tab.Panels>
      </Tab.Group>
      <div className="flex space-x-2 justify-end">
        <NeutralButton onClick={onCancel}>Cancel</NeutralButton>
        <SubmitButton onClick={onApply}>Apply</SubmitButton>
      </div>
    </div>
  );
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

  const [duration, setDuration] = useState(1);
  const [unitOfTime, setUnitOfTime] = useState("hours");

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
    <div className="flex flex-col">
      <div className="inline-flex space-x-2">
        {presets.map((preset) => {
          const presetClassName = classNames(
            "py-1.5 px-2.5 rounded-md cursor-pointer border bg-dashboard-panel",
            state.relative === preset.value ||
              (!presets.find((p) => p.value === state.relative) &&
                preset.value === "custom")
              ? "font-bold border-divide"
              : "text-foreground-light font-light border-dashboard-panel",
          );
          if (preset.value === "custom") {
            return (
              <Popover key={preset.value} className="relative">
                <Popover.Button as="div" className={presetClassName}>
                  {preset.label}
                </Popover.Button>
                <Popover.Panel className="absolute z-10 pt-px">
                  {({ close }) => (
                    <CustomDatePicker
                      duration={duration}
                      tempState={tempState}
                      unitOfTime={unitOfTime}
                      setDuration={setDuration}
                      setTempState={setTempState}
                      setUnitOfTime={setUnitOfTime}
                      onApply={() => {
                        handleApply();
                        close();
                      }}
                      onCancel={() => {
                        handleCancel();
                        close();
                      }}
                      onTimeOptionClick={handleTimeOptionClick}
                    />
                  )}
                </Popover.Panel>
              </Popover>
            );
          }
          return (
            <div
              key={preset.value}
              onClick={() => handlePresetChange(preset.value)}
              className={presetClassName}
            >
              {preset.label}
            </div>
          );
        })}
      </div>
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
