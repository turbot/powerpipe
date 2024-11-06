import React, { useState, useRef } from "react";
import { DayPicker } from "react-day-picker";
import { format } from "date-fns";
import "react-day-picker/dist/style.css";
import "react-time-picker/dist/TimePicker.css";
import "./DateRangePicker.css";

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

const DateRangePicker = () => {
  const [selectedPreset, setSelectedPreset] = useState("1h");
  const [startDate, setStartDate] = useState(null);
  const [endDate, setEndDate] = useState(null);
  const [startTime, setStartTime] = useState("00:00:00");
  const [endTime, setEndTime] = useState("23:59:59");
  const [tab, setTab] = useState("absolute");
  const [duration, setDuration] = useState(1);
  const [unitOfTime, setUnitOfTime] = useState("hours");
  const [showCustomPanel, setShowCustomPanel] = useState(false);
  const [savedDateTime, setSavedDateTime] = useState({});
  const customButtonRef = useRef(null);

  const handlePresetChange = (preset) => {
    setSelectedPreset(preset);
    if (preset === "custom") {
      setShowCustomPanel(!showCustomPanel); // Toggle visibility
    } else {
      setShowCustomPanel(false);
    }
  };

  const handleApply = () => {
    const urlParams = new URLSearchParams(window.location.search);

    if (tab === "absolute") {
      // Handle absolute case
      const formattedStartDate = startDate
        ? format(startDate, "yyyy-MM-dd")
        : "";
      const formattedEndDate = endDate ? format(endDate, "yyyy-MM-dd") : "";
      const detectionFromValue = `${formattedStartDate} ${startTime}`;
      const detectionToValue = `${formattedEndDate} ${endTime}`;

      // Add/Update URL parameters for absolute values
      urlParams.set("detection_from", detectionFromValue);
      urlParams.set("detection_to", detectionToValue);

      // Clear relative parameter if it was previously set
      urlParams.delete("detection_to_relative");

      // Save the selected date/time
      setSavedDateTime({
        type: "absolute",
        detection_from: detectionFromValue,
        detection_to: detectionToValue,
      });

      alert(
        `Saved Absolute Date/Time:\nDetection From: ${detectionFromValue}\nDetection To: ${detectionToValue}`,
      );
    } else if (tab === "relative") {
      // Handle relative case
      const relativeValue = `T-${duration}${unitOfTime[0].toUpperCase()}`; // E.g., "T-5M" for 5 minutes

      // Add/Update URL parameter for relative value
      urlParams.set("detection_to", relativeValue);

      // Clear absolute parameters if they were previously set
      urlParams.delete("detection_from");

      // Save the selected relative time
      setSavedDateTime({
        type: "relative",
        detection_to: relativeValue,
      });

      alert(`Saved Relative Date/Time:\nDetection To: ${relativeValue}`);
    }

    // Update the URL without refreshing the page
    window.history.replaceState(null, "", "?" + urlParams.toString());

    // Hide the custom panel
    setShowCustomPanel(false);
  };

  const handleCancel = () => {
    setSelectedPreset("1h");
    setStartDate(null);
    setEndDate(null);
    setStartTime("00:00:00");
    setEndTime("23:59:59");
    setDuration(1);
    setUnitOfTime("hours");
    setShowCustomPanel(false);
  };

  const handleTimeOptionClick = (value, unit) => {
    setDuration(value);
    setUnitOfTime(unit);
  };

  return (
    <div className="date-range-picker">
      <div className="presets">
        {presets.map((preset) => (
          <button
            key={preset.value}
            onClick={() => handlePresetChange(preset.value)}
            className={`preset-button ${selectedPreset === preset.value ? "active" : ""}`}
            ref={preset.value === "custom" ? customButtonRef : null}
          >
            {preset.label}
          </button>
        ))}
      </div>

      {showCustomPanel && selectedPreset === "custom" && (
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
              className={`tab-button ${tab === "absolute" ? "active" : ""}`}
              onClick={() => setTab("absolute")}
            >
              Absolute
            </button>
            <button
              className={`tab-button ${tab === "relative" ? "active" : ""}`}
              onClick={() => setTab("relative")}
            >
              Relative
            </button>
          </div>

          {tab === "absolute" ? (
            <div className="absolute-panel">
              <div className="calendar-container">
                <DayPicker
                  mode="range"
                  selected={{ from: startDate, to: endDate }}
                  onSelect={({ from, to }) => {
                    setStartDate(from);
                    setEndDate(to);
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
                    value={startDate ? format(startDate, "yyyy-MM-dd") : ""}
                    onChange={(e) => setStartDate(new Date(e.target.value))}
                  />
                  <label>Start time</label>
                  <input
                    type="time"
                    value={startTime}
                    step="1"
                    onChange={(e) => setStartTime(e.target.value)}
                  />
                </div>
                <div className="time-input-container">
                  <label>End date</label>
                  <input
                    type="date"
                    value={endDate ? format(endDate, "yyyy-MM-dd") : ""}
                    onChange={(e) => setEndDate(new Date(e.target.value))}
                  />
                  <label>End time</label>
                  <input
                    type="time"
                    value={endTime}
                    step="1"
                    onChange={(e) => setEndTime(e.target.value)}
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
                      onClick={() => handleTimeOptionClick(min, "minutes")}
                      className={`time-option ${duration === min && unitOfTime === "minutes" ? "active" : ""}`}
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
                      onClick={() => handleTimeOptionClick(hour, "hours")}
                      className={`time-option ${duration === hour && unitOfTime === "hours" ? "active" : ""}`}
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
                      onClick={() => handleTimeOptionClick(day, "days")}
                      className={`time-option ${duration === day && unitOfTime === "days" ? "active" : ""}`}
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
                      onClick={() => handleTimeOptionClick(week, "weeks")}
                      className={`time-option ${duration === week && unitOfTime === "weeks" ? "active" : ""}`}
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
                  min="1"
                  max="9999"
                  value={duration}
                  onChange={(e) => setDuration(Number(e.target.value))}
                />
                <select
                  value={unitOfTime}
                  onChange={(e) => setUnitOfTime(e.target.value)}
                >
                  <option value="minutes">Minutes</option>
                  <option value="hours">Hours</option>
                  <option value="days">Days</option>
                  <option value="weeks">Weeks</option>
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

export default DateRangePicker;
