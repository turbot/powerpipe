import dayjs from "dayjs";
import NeutralButton from "@powerpipe/components/forms/NeutralButton";
import SubmitButton from "@powerpipe/components/forms/SubmitButton";
import utc from "dayjs/plugin/utc";
import { classNames } from "@powerpipe/utils/styles";
import { createPortal } from "react-dom";
import { DayPicker, getDefaultClassNames } from "react-day-picker";
import { parseDate } from "@powerpipe/utils/date";
import { Popover, Tab } from "@headlessui/react";
import { ThemeProvider, ThemeWrapper } from "@powerpipe/hooks/useTheme";
import { useEffect, useState } from "react";
import { usePopper } from "react-popper";
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
}: {
  duration: number;
  tempState: {
    from: dayjs.Dayjs;
    to?: dayjs.Dayjs | null;
    relative?: string | null;
    showCustom?: boolean;
  };
  unitOfTime: string;
  setDuration: (duration: number) => void;
  setTempState: (state: {
    from: dayjs.Dayjs;
    to?: dayjs.Dayjs | null;
    relative?: string | null;
    showCustom?: boolean;
  }) => void;
  setUnitOfTime: (unit: string) => void;
  onApply: () => void;
  onCancel: () => void;
  onTimeOptionClick: (value: number, unit: string) => void;
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
                    from: parseDate(tempState.from)?.toDate(),
                    to: tempState.to
                      ? parseDate(tempState.to)?.toDate()
                      : undefined,
                  }}
                  onSelect={(selectEvent) => {
                    if (!selectEvent) {
                      return;
                    }
                    const { from, to } = selectEvent;
                    const newFrom = parseDate(
                      new Date(
                        from.getFullYear(),
                        from.getMonth(),
                        from.getDate(),
                      ),
                    );
                    const newTo = parseDate(
                      new Date(
                        to.getFullYear(),
                        to.getMonth(),
                        to.getDate(),
                        23,
                        59,
                        59,
                      ),
                    );
                    setTempState({
                      ...tempState,
                      relative: "custom",
                      from: newFrom,
                      to: newTo,
                    });
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
                      defaultValue={parseDate(tempState.from)?.format(
                        "YYYY-MM-DD",
                      )}
                      onChange={(e) =>
                        setTempState({
                          ...tempState,
                          relative: "custom",
                          from: dayjs(e.target.value).utc(),
                        })
                      }
                      className="bg-dashboard-panel text-foreground dark:bg-dashboard dark:text-foreground-light border border-table-border rounded p-2 w-full"
                    />
                  </div>
                  <div>
                    <label>Start time</label>
                    <input
                      type="time"
                      defaultValue={parseDate(tempState.from)?.format(
                        "HH:mm:ss",
                      )}
                      step="1"
                      onChange={(e) => {
                        setTempState({
                          ...tempState,
                          relative: "custom",
                          from: dayjs(
                            `${parseDate(tempState.from)?.format("YYYY-MM-DD")} ${e.target.value}`,
                          ),
                        });
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
                      defaultValue={
                        tempState.to
                          ? parseDate(tempState.to)?.format("YYYY-MM-DD")
                          : undefined
                      }
                      onChange={(e) =>
                        setTempState({
                          ...tempState,
                          relative: "custom",
                          to: dayjs(e.target.value).utc(),
                        })
                      }
                      className="bg-dashboard-panel text-foreground dark:bg-dashboard dark:text-foreground-light border border-table-border rounded p-2 w-full"
                    />
                  </div>
                  <div className="w-full">
                    <label>End time</label>
                    <input
                      type="time"
                      defaultValue={
                        tempState.to
                          ? parseDate(tempState.to)?.format("HH:mm:ss")
                          : `00:00:00`
                      }
                      step="1"
                      onChange={(e) => {
                        const toTime = tempState.to
                          ? parseDate(tempState.to)
                          : dayjs();
                        setTempState({
                          ...tempState,
                          relative: "custom",
                          from: dayjs(
                            `${toTime?.format("YYYY-MM-DD")} ${e.target.value}`,
                          ),
                        });
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

const DatetimeRangePicker = ({
  from,
  to,
  relative,
  disabled,
  onChange,
}: {
  from: dayjs.Dayjs;
  to?: dayjs.Dayjs | null;
  relative?: string | null;
  disabled: boolean;
  onChange: (
    from: dayjs.Dayjs,
    to?: dayjs.Dayjs | null,
    relative?: string | null,
  ) => void;
}) => {
  const [popperElement, setPopperElement] = useState(null);
  const [referenceElement, setReferenceElement] = useState(null);
  const { styles, attributes } = usePopper(referenceElement, popperElement, {
    placement: "bottom-start",
    modifiers: [
      {
        name: "flip",
        options: {
          fallbackPlacements: ["top-start", "right-start"],
        },
      },
    ],
  });

  const [state, setState] = useState<{
    from: dayjs.Dayjs;
    to?: dayjs.Dayjs | null;
    relative?: string | null;
    showCustom?: boolean;
  }>({ from, to, relative, showCustom: false });

  //   return {
  //     from:
  //   }
  //
  //   if (stateValue) {
  //     try {
  //       const parsed = JSON.parse(stateValue);
  //       return {
  //         from: parsed.from
  //           ? dayjs(parsed.from)
  //           : dayjs.subtract(7, "day").utc(),
  //         to: parsed.to ? dayjs(parsed.to) : null,
  //         relative: parsed.relative || "7d",
  //       };
  //     } catch (err) {
  //       console.error("Parse error", err);
  //       const now = dayjs();
  //       return {
  //         from: now.subtract(7, "day").utc(),
  //         to: null,
  //         relative: "7d",
  //       };
  //     }
  //   } else {
  //     const now = dayjs();
  //     return {
  //       from: now.subtract(7, "day").utc(),
  //       to: null,
  //       relative: "7d",
  //     };
  //   }
  // });

  // useEffect(() => {
  //   if (stateValue) {
  //     return;
  //   }
  //   updateInput(
  //     props.name,
  //     JSON.stringify({
  //       from: dayjs().subtract(7, "day").utc(),
  //       to: null,
  //       relative: "7d",
  //     }),
  //     !!stateValue,
  //   );
  // }, [stateValue]);

  useEffect(() => {
    if (state.showCustom) {
      return;
    }
    onChange(
      state.from.toISOString(),
      state.to?.toISOString() || null,
      state.relative,
    );
    // updateInput(
    //   props.name,
    //   JSON.stringify({
    //     from: state.from,
    //     to: state.to,
    //     relative: state.relative,
    //   }),
    //   !!stateValue,
    // );
  }, [state.from, state.to, state.relative, state.showCustom]);

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
    <div className="flex rounded-md border border-black-scale-3">
      {presets.map((preset, index) => {
        const presetClassName = classNames(
          "py-1.5 px-2.5 rounded-md",
          state.showCustom ? null : "border border-t-0 border-b-0",
          index === 0 ? "border-l-0" : null,
          index === presets.length - 1 ? "border-r-0" : null,
          disabled ? null : "cursor-pointer",
          state.relative === preset.value ||
            (!presets.find((p) => p.value === state.relative) &&
              preset.value === "custom")
            ? "bg-dashboard border-divide"
            : "text-foreground-light border-dashboard-panel",
        );
        if (preset.value === "custom") {
          return (
            <Popover key={preset.value} className="relative">
              <Popover.Button
                ref={setReferenceElement}
                as="div"
                className={presetClassName}
                disabled={disabled}
              >
                {preset.label}
              </Popover.Button>
              <Popover.Panel className="absolute z-10 pt-px">
                {({ close }) =>
                  createPortal(
                    <ThemeProvider>
                      <ThemeWrapper>
                        <div
                          // @ts-ignore
                          ref={setPopperElement}
                          style={{ ...styles.popper }}
                          {...attributes.popper}
                          onClick={(e) => e.stopPropagation()}
                        >
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
                        </div>
                      </ThemeWrapper>
                    </ThemeProvider>,
                    // @ts-ignore as this element definitely exists
                    document.getElementById("portals"),
                  )
                }
              </Popover.Panel>
            </Popover>
          );
        }
        return (
          <div
            key={preset.value}
            onClick={
              disabled ? undefined : () => handlePresetChange(preset.value)
            }
            className={presetClassName}
          >
            {preset.label}
          </div>
        );
      })}
    </div>
  );
};

export default DatetimeRangePicker;
