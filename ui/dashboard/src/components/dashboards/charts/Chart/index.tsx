import ErrorPanel from "@powerpipe/components/dashboards/Error";
import has from "lodash/has";
import isEmpty from "lodash/isEmpty";
import merge from "lodash/merge";
import Placeholder from "@powerpipe/components/dashboards/Placeholder";
import React, { useEffect, useRef, useState } from "react";
import ReactEChartsCore from "echarts-for-react/lib/core";
import set from "lodash/set";
import useChartThemeColors from "@powerpipe/hooks/useChartThemeColors";
import useMediaMode from "@powerpipe/hooks/useMediaMode";
import useTemplateRender from "@powerpipe/hooks/useTemplateRender";
import {
  buildChartDataset,
  getColorOverride,
  LeafNodeData,
  Width,
} from "@powerpipe/components/dashboards/common";
import { EChartsOption } from "echarts-for-react/src/types";
import {
  ChartProperties,
  ChartProps,
  ChartSeries,
  ChartSeriesOptions,
  ChartTransform,
  ChartType,
} from "@powerpipe/components/dashboards/charts/types";
import { FlowType } from "@powerpipe/components/dashboards/flows/types";
import { getChartComponent } from "@powerpipe/components/dashboards/charts";
import { GraphType } from "@powerpipe/components/dashboards/graphs/types";
import { HierarchyType } from "@powerpipe/components/dashboards/hierarchies/types";
import { injectSearchPathPrefix } from "@powerpipe/utils/url";
import { registerComponent } from "@powerpipe/components/dashboards";
import { useDashboard } from "@powerpipe/hooks/useDashboard";
import { useNavigate } from "react-router-dom";
import {
  isDiffColumn,
  parseDiffColumn,
  tableRowDiffColumn,
} from "@powerpipe/utils/data";

const getThemeColorsWithPointOverrides = (
  type: ChartType = "column",
  series: any[],
  seriesOverrides: ChartSeries | undefined,
  dataset: any[][],
  themeColorValues,
) => {
  if (isEmpty(themeColorValues)) {
    return [];
  }
  switch (type) {
    case "donut":
    case "pie": {
      const newThemeColors: string[] = [];
      for (let rowIndex = 1; rowIndex < dataset.length; rowIndex++) {
        if (rowIndex - 1 < themeColorValues.charts.length) {
          newThemeColors.push(themeColorValues.charts[rowIndex - 1]);
        } else {
          newThemeColors.push(
            themeColorValues.charts[
              (rowIndex - 1) % themeColorValues.charts.length
            ],
          );
        }
      }
      series.forEach((seriesInfo) => {
        const seriesName = seriesInfo.name;
        const overrides = seriesOverrides
          ? seriesOverrides[seriesName] || {}
          : ({} as ChartSeriesOptions);
        const pointOverrides = overrides.points || {};
        dataset.slice(1).forEach((dataRow, dataRowIndex) => {
          const pointOverride = pointOverrides[dataRow[0]];
          if (pointOverride && pointOverride.color) {
            newThemeColors[dataRowIndex] = getColorOverride(
              pointOverride.color,
              themeColorValues,
            );
          }
        });
      });
      return newThemeColors;
    }
    default:
      const newThemeColors: string[] = [];
      for (let seriesIndex = 0; seriesIndex < series.length; seriesIndex++) {
        if (seriesIndex < themeColorValues.charts.length - 1) {
          newThemeColors.push(themeColorValues.charts[seriesIndex]);
        } else {
          newThemeColors.push(
            themeColorValues.charts[
              seriesIndex % themeColorValues.charts.length
            ],
          );
        }
      }
      return newThemeColors;
  }
};

const getCommonBaseOptions = () => ({
  animation: false,
  grid: {
    left: "5%",
    right: "5%",
    top: "7%",
    bottom: "8%",
    // bottom: 40,
    containLabel: true,
  },
  legend: {
    orient: "horizontal",
    left: "center",
    top: "10",
    textStyle: {
      fontSize: 11,
      overflow: "truncate",
    },
  },
  textStyle: {
    fontFamily:
      'ui-sans-serif, system-ui, -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, "Helvetica Neue", Arial, "Noto Sans", sans-serif, "Apple Color Emoji", "Segoe UI Emoji", "Segoe UI Symbol", "Noto Color Emoji"',
  },
  tooltip: {
    appendToBody: true,
    textStyle: {
      fontSize: 11,
    },
    trigger: "item",
  },
});

const getXAxisLabelRotation = (number_of_rows: number) => {
  if (number_of_rows < 5) {
    return 0;
  }
  if (number_of_rows < 10) {
    return 30;
  }
  if (number_of_rows < 15) {
    return 45;
  }
  if (number_of_rows < 20) {
    return 60;
  }
  return 90;
};

const getXAxisLabelWidth = (number_of_rows: number) => {
  if (number_of_rows < 5) {
    return null;
  }
  if (number_of_rows < 10) {
    return 85;
  }
  if (number_of_rows < 15) {
    return 75;
  }
  if (number_of_rows < 20) {
    return 60;
  }
  return 50;
};

const getCommonBaseOptionsForChartType = (
  type: ChartType | undefined,
  width: Width | undefined,
  dataset: any[][],
  shouldBeTimeSeries: boolean,
  series: any[],
  seriesOverrides: ChartSeries | undefined,
  themeColors,
) => {
  switch (type) {
    case "bar":
      return {
        color: getThemeColorsWithPointOverrides(
          type,
          series,
          seriesOverrides,
          dataset,
          themeColors,
        ),
        legend: {
          show: series ? series.length > 1 : false,
          textStyle: {
            color: themeColors.foreground,
          },
        },
        // Declare an x-axis (category axis).
        // The category map the first row in the dataset by default.
        xAxis: {
          axisLabel: { color: themeColors.foreground, fontSize: 10 },
          axisLine: {
            show: true,
            lineStyle: { color: themeColors.foregroundLightest },
          },
          axisTick: { show: true },
          nameGap: 25,
          nameLocation: "center",
          nameTextStyle: { color: themeColors.foreground },
          splitLine: { show: false },
        },
        // Declare a y-axis (value axis).
        yAxis: {
          type: "category",
          axisLabel: {
            color: themeColors.foreground,
            overflow: "truncate",
          },
          axisLine: { lineStyle: { color: themeColors.foregroundLightest } },
          axisTick: { show: false },
          nameGap: 100,
          nameLocation: "center",
          nameTextStyle: { color: themeColors.foreground },
        },
      };
    case "area":
    case "line":
      return {
        color: getThemeColorsWithPointOverrides(
          type,
          series,
          seriesOverrides,
          dataset,
          themeColors,
        ),
        legend: {
          show: series ? series.length > 1 : false,
          textStyle: {
            color: themeColors.foreground,
          },
        },
        // Declare an x-axis (category or time axis, depending on the type of the first column).
        // The category/time map the first row in the dataset by default.
        xAxis: {
          type: shouldBeTimeSeries ? "time" : "category",
          boundaryGap: type !== "area",
          axisLabel: {
            color: themeColors.foreground,
            fontSize: 10,
            rotate: getXAxisLabelRotation(dataset.length - 1),
            width: getXAxisLabelWidth(dataset.length),
            overflow: "truncate",
          },
          axisLine: { lineStyle: { color: themeColors.foregroundLightest } },
          axisTick: { show: false },
          nameGap: 30,
          nameLocation: "center",
          nameTextStyle: { color: themeColors.foreground },
        },
        // Declare a y-axis (value axis).
        yAxis: {
          axisLabel: { color: themeColors.foreground, fontSize: 10 },
          axisLine: {
            show: true,
            lineStyle: { color: themeColors.foregroundLightest },
          },
          axisTick: { show: true },
          splitLine: { show: false },
          nameGap: width ? width + 42 : 50,
          nameLocation: "center",
          nameTextStyle: { color: themeColors.foreground },
        },
        tooltip: {
          trigger: "axis",
        },
      };
    case "column":
      return {
        color: getThemeColorsWithPointOverrides(
          type,
          series,
          seriesOverrides,
          dataset,
          themeColors,
        ),
        legend: {
          show: series ? series.length > 1 : false,
          textStyle: {
            color: themeColors.foreground,
          },
        },
        // Declare an x-axis (category or time axis, depending on the value of the first column).
        // The category/time map the first row in the dataset by default.
        xAxis: {
          type: shouldBeTimeSeries ? "time" : "category",
          axisLabel: {
            color: themeColors.foreground,
            fontSize: 10,
            rotate: getXAxisLabelRotation(dataset.length - 1),
            width: getXAxisLabelWidth(dataset.length),
            overflow: "truncate",
          },
          axisLine: { lineStyle: { color: themeColors.foregroundLightest } },
          axisTick: { show: false },
          nameGap: 30,
          nameLocation: "center",
          nameTextStyle: { color: themeColors.foreground },
        },
        // Declare a y-axis (value axis).
        yAxis: {
          axisLabel: { color: themeColors.foreground, fontSize: 10 },
          axisLine: {
            show: true,
            lineStyle: { color: themeColors.foregroundLightest },
          },
          axisTick: { show: true },
          splitLine: { show: false },
          nameGap: width ? width + 42 : 50,
          nameLocation: "center",
          nameTextStyle: { color: themeColors.foreground },
        },
        ...(shouldBeTimeSeries ? { tooltip: { trigger: "axis" } } : {}),
      };
    case "pie":
      return {
        color: getThemeColorsWithPointOverrides(
          type,
          series,
          seriesOverrides,
          dataset,
          themeColors,
        ),
        legend: {
          show: false,
          textStyle: {
            color: themeColors.foreground,
          },
        },
      };
    case "donut":
      return {
        color: getThemeColorsWithPointOverrides(
          type,
          series,
          seriesOverrides,
          dataset,
          themeColors,
        ),
        legend: {
          show: false,
          textStyle: {
            color: themeColors.foreground,
          },
        },
      };
    default:
      return {};
  }
};

const getOptionOverridesForChartType = (
  type: ChartType = "column",
  properties: ChartProperties | undefined,
  shouldBeTimeSeries: boolean,
) => {
  if (!properties) {
    return {};
  }

  let overrides = {};

  // orient: "horizontal",
  //     left: "center",
  //     top: "top",

  if (properties.legend) {
    // Legend display
    const legendDisplay = properties.legend.display;
    if (legendDisplay === "all") {
      overrides = set(overrides, "legend.show", true);
    } else if (legendDisplay === "none") {
      overrides = set(overrides, "legend.show", false);
    }

    // Legend display position
    const legendPosition = properties.legend.position;
    if (legendPosition === "top") {
      overrides = set(overrides, "legend.orient", "horizontal");
      overrides = set(overrides, "legend.left", "center");
      overrides = set(overrides, "legend.top", 10);
      overrides = set(overrides, "legend.bottom", "auto");
    } else if (legendPosition === "right") {
      overrides = set(overrides, "legend.orient", "vertical");
      overrides = set(overrides, "legend.left", 10);
      overrides = set(overrides, "legend.top", "middle");
      overrides = set(overrides, "legend.bottom", "auto");
    } else if (legendPosition === "bottom") {
      overrides = set(overrides, "legend.orient", "horizontal");
      overrides = set(overrides, "legend.left", "center");
      overrides = set(overrides, "legend.top", "auto");
      overrides = set(overrides, "legend.bottom", 10);
    } else if (legendPosition === "left") {
      overrides = set(overrides, "legend.orient", "vertical");
      overrides = set(overrides, "legend.left", 10);
      overrides = set(overrides, "legend.top", "middle");
      overrides = set(overrides, "legend.bottom", "auto");
    }
  }

  // Axes settings
  if (properties.axes) {
    // X axis settings
    if (properties.axes.x) {
      // X axis display setting
      const xAxisDisplay = properties.axes.x.display;
      if (xAxisDisplay === "all") {
        overrides = set(overrides, "xAxis.show", true);
      } else if (xAxisDisplay === "none") {
        overrides = set(overrides, "xAxis.show", false);
      }

      // X axis min setting
      if (type === "bar" && has(properties, "axes.x.min")) {
        overrides = set(overrides, "xAxis.min", properties.axes.x.min);
      }
      // Y axis max setting
      if (type === "bar" && has(properties, "axes.x.max")) {
        overrides = set(overrides, "xAxis.max", properties.axes.x.max);
      }

      // X axis labels settings
      if (properties.axes.x.labels) {
        // X axis labels display setting
        const xAxisTicksDisplay = properties.axes.x.labels.display;
        if (xAxisTicksDisplay === "all") {
          overrides = set(overrides, "xAxis.axisLabel.show", true);
        } else if (xAxisTicksDisplay === "none") {
          overrides = set(overrides, "xAxis.axisLabel.show", false);
        }
      }

      // X axis title settings
      if (properties.axes.x.title) {
        // X axis title display setting
        const xAxisTitleDisplay = properties.axes.x.title.display;
        if (xAxisTitleDisplay === "none") {
          overrides = set(overrides, "xAxis.name", null);
        }

        // X Axis title align setting
        const xAxisTitleAlign = properties.axes.x.title.align;
        if (xAxisTitleAlign === "start") {
          overrides = set(overrides, "xAxis.nameLocation", "start");
        } else if (xAxisTitleAlign === "center") {
          overrides = set(overrides, "xAxis.nameLocation", "center");
        } else if (xAxisTitleAlign === "end") {
          overrides = set(overrides, "xAxis.nameLocation", "end");
        }

        // X Axis title value setting
        const xAxisTitleValue = properties.axes.x.title.value;
        if (xAxisTitleValue) {
          overrides = set(overrides, "xAxis.name", xAxisTitleValue);
        }
      }

      // X Axis range setting (for timeseries plots)
      // Valid chart types: column, area, line (bar, donut and pie make no sense)
      if (["column", "area", "line"].includes(type) && shouldBeTimeSeries) {
        // X axis min setting (for timeseries)
        if (has(properties, "axes.x.min")) {
          // ECharts wants millis since epoch, not seconds
          overrides = set(overrides, "xAxis.min", properties.axes.x.min * 1000);
        }
        // Y axis max setting (for timeseries)
        if (has(properties, "axes.x.max")) {
          // ECharts wants millis since epoch, not seconds
          overrides = set(overrides, "xAxis.max", properties.axes.x.max * 1000);
        }
      }
    }

    // Y axis settings
    if (properties.axes.y) {
      // Y axis display setting
      const yAxisDisplay = properties.axes.y.display;
      if (yAxisDisplay === "all") {
        overrides = set(overrides, "yAxis.show", true);
      } else if (yAxisDisplay === "none") {
        overrides = set(overrides, "yAxis.show", false);
      }

      // Y axis min setting
      if (type !== "bar" && has(properties, "axes.y.min")) {
        overrides = set(overrides, "yAxis.min", properties.axes.y.min);
      }
      // Y axis max setting
      if (type !== "bar" && has(properties, "axes.y.max")) {
        overrides = set(overrides, "yAxis.max", properties.axes.y.max);
      }

      // Y axis labels settings
      if (properties.axes.y.labels) {
        // Y axis labels display setting
        const yAxisTicksDisplay = properties.axes.y.labels.display;
        if (yAxisTicksDisplay === "all") {
          overrides = set(overrides, "yAxis.axisLabel.show", true);
        } else if (yAxisTicksDisplay === "none") {
          overrides = set(overrides, "yAxis.axisLabel.show", false);
        }
      }

      // Y axis title settings
      if (properties.axes.y.title) {
        // Y axis title display setting
        const yAxisTitleDisplay = properties.axes.y.title.display;
        if (yAxisTitleDisplay === "none") {
          overrides = set(overrides, "yAxis.name", null);
        }

        // Y Axis title align setting
        const yAxisTitleAlign = properties.axes.y.title.align;
        if (yAxisTitleAlign === "start") {
          overrides = set(overrides, "yAxis.nameLocation", "start");
        } else if (yAxisTitleAlign === "center") {
          overrides = set(overrides, "yAxis.nameLocation", "center");
        } else if (yAxisTitleAlign === "end") {
          overrides = set(overrides, "yAxis.nameLocation", "end");
        }

        // Y Axis title value setting
        const yAxisTitleValue = properties.axes.y.title.value;
        if (yAxisTitleValue) {
          overrides = set(overrides, "yAxis.name", yAxisTitleValue);
        }
      }
    }
  }

  return overrides;
};

const getSeriesForChartType = (
  type: ChartType = "column",
  data: LeafNodeData | undefined,
  properties: ChartProperties | undefined,
  rowSeriesLabels: string[],
  transform: ChartTransform,
  shouldBeTimeSeries: boolean,
  themeColors,
  dataset,
) => {
  if (!data) {
    return [];
  }

  // Keep a map of the series names with their index and configured color override
  const seriesMap = {};

  const series: any[] = [];
  const seriesNames =
    transform === "crosstab"
      ? rowSeriesLabels
      : data.columns
          .slice(1)
          .filter((col) => col.name !== "__diff")
          .map((col) => col.name);
  const seriesNamesWithoutDiffColumns = seriesNames.filter(
    (s) => !isDiffColumn(s),
  );
  const seriesLength = seriesNames.length;
  const hasDiffCol = !!data.columns.find((col) => col.name === "__diff");

  for (let seriesIndex = 0; seriesIndex < seriesLength; seriesIndex++) {
    let seriesName = seriesNames[seriesIndex];
    const diff = parseDiffColumn(seriesName);
    let seriesColor = "auto";
    const seriesMapSettings = {
      index: seriesIndex,
      // Don't set if a diff - see if we set it below via override for paired column
      // and if not, then we'll try to look up the colour from the paired column index
      color: diff.isDiff ? undefined : seriesColor,
      isDiff: diff.isDiff,
      pairedColumn: diff.pairedColumn,
      pairedColumnIndex: diff.isDiff
        ? seriesNamesWithoutDiffColumns.indexOf(diff.pairedColumn)
        : -1,
    };
    let seriesOverrides;

    if (diff.isDiff) {
      seriesName = `${seriesName.split("_")[0]} (Diff)`;
    }

    if (properties) {
      if (
        properties.series &&
        properties.series[diff.pairedColumn || seriesName]
      ) {
        seriesOverrides = properties.series[seriesName];
      }
      if (seriesOverrides && seriesOverrides.title) {
        seriesName = seriesOverrides.title;
      }
      if (seriesOverrides && seriesOverrides.color) {
        seriesColor = getColorOverride(seriesOverrides.color, themeColors);
        seriesMapSettings.color = seriesColor;
      }

      seriesMap[seriesName] = seriesMapSettings;
    }

    if (diff.isDiff && seriesMapSettings.color === undefined) {
      seriesMapSettings.color =
        themeColors.charts[
          seriesMapSettings.pairedColumnIndex % themeColors.charts.length
        ];
    }

    switch (type) {
      case "bar":
      case "column":
        series.push({
          name: seriesName,
          type: "bar",
          ...(hasDiffCol || (properties && properties.grouping === "compare")
            ? {}
            : { stack: "total" }),
          itemStyle: {
            borderRadius:
              // Only round the last series and take into account bar vs chart e.g. orientation
              seriesIndex + 1 === seriesLength
                ? type === "bar"
                  ? [0, 5, 5, 0]
                  : [5, 5, 0, 0]
                : undefined,
            color: !diff.isDiff
              ? seriesMapSettings.color
              : {
                  type: "pattern",
                  image: (() => {
                    // Create a canvas for the pattern
                    const canvas = document.createElement("canvas");
                    const ctx = canvas.getContext("2d");
                    canvas.width = 8; // Pattern size
                    canvas.height = 8;

                    // Set base color
                    const baseColor = seriesMapSettings.color;

                    // Define the colors
                    const lightColor = lightenColor(baseColor, 0.3);
                    // const darkColor = darkenColor(baseColor, 0.2);

                    // Draw light background
                    ctx.fillStyle = lightColor;
                    ctx.fillRect(0, 0, 8, 8);

                    // Draw wider diagonal dark lines
                    ctx.strokeStyle = baseColor;
                    ctx.lineWidth = 4; // Increase line width for wider lines
                    ctx.beginPath();
                    ctx.moveTo(-2, 6); // Adjust start and end points for better alignment
                    ctx.lineTo(6, -2);
                    ctx.moveTo(2, 10);
                    ctx.lineTo(10, 2);
                    ctx.stroke();

                    return canvas;
                  })(),
                  repeat: "repeat",
                },
            borderColor: themeColors.dashboardPanel,
            borderWidth: 1,
          },
          emphasis: {
            itemStyle: {
              borderRadius: [5, 5],
            },
          },
          barMaxWidth: 75,
          // Per https://stackoverflow.com/a/56116442, when using time series you have to manually encode each series
          // We assume that the first dimension/column is the timestamp
          ...(shouldBeTimeSeries ? { encode: { x: 0, y: seriesName } } : {}),
          // label: {
          //   show: true,
          //   position: 'outside'
          // },
        });
        break;
      case "donut":
        series.push({
          name: seriesName,
          type: "pie",
          center: ["50%", "50%"],
          radius: ["30%", "50%"],
          label: { color: themeColors.foreground, fontSize: 10 },
          itemStyle: {
            borderRadius: 5,
            borderColor: themeColors.dashboardPanel,
            borderWidth: 2,
          },
          emphasis: {
            itemStyle: {
              color: "inherit",
            },
          },
        });
        break;
      case "pie":
        if (diff.isDiff) {
          // Multi-series pie (donut) logic
          series.push({
            name: seriesName + " (Original)",
            type: "pie",
            center: ["50%", "50%"],
            radius: ["20%", "40%"], // Inner radius for donut effect
            data: dataset.slice(1).map((item) => {
              return {
                value: item[2], // 'Count'
                name: item[0], // 'Type'
              };
            }),
            label: {
              position: "outside",
              formatter: "{b}: {c}",
              color: themeColors.foreground,
              fontSize: 10,
            },
            emphasis: {
              itemStyle: {
                color: "inherit",
                shadowBlur: 5,
                shadowOffsetX: 0,
                shadowColor: "rgba(0, 0, 0, 0.5)",
              },
            },
            itemStyle: {
              borderRadius: 5,
              borderColor: themeColors.dashboardPanel,
              borderWidth: 2,
            },
          });

          series.push({
            name: seriesName + " (Diff)",
            type: "pie",
            center: ["50%", "50%"],
            radius: ["45%", "65%"], // Outer radius for donut effect
            data: dataset.slice(1).map((item, index) => {
              // Get the matching color from the inner series
              const matchingColor =
                themeColors.charts[index % themeColors.charts.length];
              return {
                value: item[1], // 'Count_diff'
                name: item[0], // 'Type'
                itemStyle: {
                  color: {
                    type: "pattern",
                    image: (() => {
                      // Create a canvas for the hatching pattern
                      const canvas = document.createElement("canvas");
                      const ctx = canvas.getContext("2d");
                      canvas.width = 8; // Pattern size
                      canvas.height = 8;

                      // Fill the background with the matching color
                      ctx.fillStyle = matchingColor;
                      ctx.fillRect(0, 0, 8, 8);

                      // Draw the hatching pattern (e.g., diagonal lines)
                      ctx.strokeStyle = lightenColor(matchingColor, 0.3); // Use a slightly lighter color for contrast
                      ctx.lineWidth = 2; // Adjust line width for desired effect
                      ctx.beginPath();
                      ctx.moveTo(0, 8); // Start from bottom left
                      ctx.lineTo(8, 0); // Draw diagonal line to top right
                      ctx.stroke();

                      return canvas;
                    })(),
                    repeat: "repeat",
                  },
                  borderColor: themeColors.dashboardPanel,
                  borderWidth: 2,
                  borderRadius: 5,
                },
              };
            }),
            label: {
              position: "outside",
              formatter: "{b}: {c}",
              color: themeColors.foreground,
              fontSize: 10,
            },
            emphasis: {
              itemStyle: {
                color: "inherit",
                shadowBlur: 5,
                shadowOffsetX: 0,
                shadowColor: "rgba(0, 0, 0, 0.5)",
              },
            },
            itemStyle: {
              borderRadius: 5,
              borderColor: themeColors.dashboardPanel,
              borderWidth: 2,
            },
          });
        } else {
          // Single-series pie (standard pie logic)
          series.push({
            name: seriesName,
            type: "pie",
            center: ["50%", "50%"],
            radius: ["30%", "50%"],
            label: { color: themeColors.foreground, fontSize: 10 },
            itemStyle: {
              borderRadius: 5,
              borderColor: themeColors.dashboardPanel,
              borderWidth: 2,
            },
            emphasis: {
              itemStyle: {
                color: "inherit",
              },
            },
          });
        }
        break;

      case "area":
        series.push({
          name: seriesName,
          type: "line",
          ...(properties && properties.grouping === "compare"
            ? {}
            : { stack: "total" }),
          // Per https://stackoverflow.com/a/56116442, when using time series you have to manually encode each series
          // We assume that the first dimension/column is the timestamp
          ...(shouldBeTimeSeries ? { encode: { x: 0, y: seriesName } } : {}),
          areaStyle: {},
          emphasis: {
            focus: "series",
          },
          itemStyle: {
            color: !diff.isDiff
              ? seriesMapSettings.color
              : {
                  type: "pattern",
                  image: (() => {
                    // Create a canvas for the pattern
                    const canvas = document.createElement("canvas");
                    const ctx = canvas.getContext("2d");
                    canvas.width = 8; // Pattern size
                    canvas.height = 8;

                    // Set base color
                    const baseColor = seriesMapSettings.color;

                    // Define the colors
                    const lightColor = lightenColor(baseColor, 0.3);
                    // const darkColor = darkenColor(baseColor, 0.2);

                    // Draw light background
                    ctx.fillStyle = lightColor;
                    ctx.fillRect(0, 0, 8, 8);

                    // Draw wider diagonal dark lines
                    ctx.strokeStyle = baseColor;
                    ctx.lineWidth = 4; // Increase line width for wider lines
                    ctx.beginPath();
                    ctx.moveTo(-2, 6); // Adjust start and end points for better alignment
                    ctx.lineTo(6, -2);
                    ctx.moveTo(2, 10);
                    ctx.lineTo(10, 2);
                    ctx.stroke();

                    return canvas;
                  })(),
                  repeat: "repeat",
                },
            borderColor: themeColors.dashboardPanel,
            borderWidth: 1,
          },
        });
        break;
      case "line":
        series.push({
          name: seriesName,
          type: "line",
          itemStyle: { color: seriesMapSettings.color },
          lineStyle: {
            type: diff.isDiff ? "dashed" : "solid",
          },
          // Per https://stackoverflow.com/a/56116442, when using time series you have to manually encode each series
          // We assume that the first dimension/column is the timestamp
          ...(shouldBeTimeSeries ? { encode: { x: 0, y: seriesName } } : {}),
        });
        break;
    }
  }
  return series;
};

// Utility functions for color adjustments
function lightenColor(color, amount) {
  const [r, g, b] = hexToRgb(color);
  return rgbToHex(
    Math.min(255, r + amount * 255),
    Math.min(255, g + amount * 255),
    Math.min(255, b + amount * 255),
  );
}

function darkenColor(color, amount) {
  const [r, g, b] = hexToRgb(color);
  return rgbToHex(
    Math.max(0, r - amount * 255),
    Math.max(0, g - amount * 255),
    Math.max(0, b - amount * 255),
  );
}

function hexToRgb(hex) {
  const bigint = parseInt(hex.slice(1), 16);
  return [(bigint >> 16) & 255, (bigint >> 8) & 255, bigint & 255];
}

function rgbToHex(r, g, b) {
  return `#${((1 << 24) + (r << 16) + (g << 8) + b).toString(16).slice(1)}`;
}

const adjustGridConfig = (
  config: EChartsOption,
  properties: ChartProperties | undefined,
) => {
  let newConfig = { ...config };
  if (!!newConfig?.xAxis?.name) {
    newConfig = set(newConfig, "grid.containLabel", false);
    newConfig = set(newConfig, "grid.bottom", "20%");
  }
  if (!!newConfig?.yAxis?.name) {
    newConfig = set(newConfig, "grid.containLabel", false);
    newConfig = set(newConfig, "grid.left", "25%");
    newConfig = set(newConfig, "grid.bottom", "25%");
  }
  if (newConfig?.legend?.show) {
    const configuredPosition = properties?.legend?.position || "top";
    switch (configuredPosition) {
      case "top":
        newConfig = set(newConfig, "grid.top", "20%");
        newConfig = set(newConfig, "grid.top", "20%");
        break;
      case "right":
        newConfig = set(newConfig, "grid.right", "35%");
        break;
      case "bottom":
        newConfig = set(newConfig, "grid.bottom", "25%");
        break;
      case "left":
        newConfig = set(newConfig, "grid.left", "50%");
        break;
    }
  }
  return newConfig;
};

const injectDiffColumns = (data: LeafNodeData) => {
  const diffSeriesToAdd = {};
  let newColumns = [...data.columns];
  for (const row of data.rows) {
    const keys = Object.keys(row);
    for (const key of keys) {
      if (key === "__diff" || key.endsWith("_diff") || !!diffSeriesToAdd[key]) {
        continue;
      }
      const diff = tableRowDiffColumn(row, key);
      if (
        diff.hasDiffColumn &&
        !data.columns.find((c) => c.name === `${key}_diff`)
      ) {
        diffSeriesToAdd[`${key}_diff`] = key;
      }

      // if (hasDiffCol(key)) const series = seriesWithDiffs[key] || [];
      // series.push(row[key]);
      // seriesWithDiffs[key] = series;
    }
  }

  const diffCol = newColumns.find((c) => c.name === "__diff");
  for (const diffSeries of Object.keys(diffSeriesToAdd)) {
    const matchingColumnIndex = newColumns.findIndex(
      (c) => c.name === diffSeriesToAdd[diffSeries],
    );
    const matchingColumn = newColumns.find(
      (c) => c.name === diffSeriesToAdd[diffSeries],
    );
    if (!matchingColumn) {
      continue;
    }
    newColumns = [
      ...newColumns.slice(0, matchingColumnIndex),
      { ...matchingColumn, name: `${matchingColumn.name}_diff` },
      matchingColumn,
      ...newColumns.slice(matchingColumnIndex + 1),
    ];
  }

  if (!diffCol) {
    newColumns = [...newColumns, { name: "__diff" }];
  }

  return { columns: newColumns, rows: data.rows };
};

const buildChartOptions = (props: ChartProps, themeColors: any) => {
  const updatedData = injectDiffColumns(props.data);
  // props.data = updatedData;

  const { dataset, rowSeriesLabels, transform } = buildChartDataset(
    updatedData,
    props.properties,
  );
  const treatAsTimeSeries = ["timestamp", "timestamptz", "date"].includes(
    updatedData?.columns[0].data_type.toLowerCase() || "",
  );
  const series = getSeriesForChartType(
    props.display_type || "column",
    updatedData,
    props.properties,
    rowSeriesLabels,
    transform,
    treatAsTimeSeries,
    themeColors,
    dataset,
  );
  const config = merge(
    getCommonBaseOptions(),
    getCommonBaseOptionsForChartType(
      props.display_type || "column",
      props.width,
      dataset,
      treatAsTimeSeries,
      series,
      props.properties?.series,
      themeColors,
    ),
    getOptionOverridesForChartType(
      props.display_type || "column",
      props.properties,
      treatAsTimeSeries,
    ),
    { series },
    {
      dataset: {
        source: dataset,
      },
    },
  );
  return adjustGridConfig(config, props.properties);
};

type ChartComponentProps = {
  options: EChartsOption;
  searchPathPrefix: string[];
  type: ChartType | FlowType | GraphType | HierarchyType;
};

const handleClick = async (
  params: any,
  navigate,
  renderTemplates,
  searchPathPrefix,
) => {
  const componentType = params.componentType;
  if (componentType !== "series") {
    return;
  }
  const dataType = params.dataType;

  switch (dataType) {
    case "node":
      if (!params.data.href) {
        return;
      }
      const renderedResults = await renderTemplates(
        { graph_node: params.data.href as string },
        [params.data],
      );
      let rowRenderResult = renderedResults[0];
      const withSearchPathPrefix = injectSearchPathPrefix(
        rowRenderResult.graph_node.result,
        searchPathPrefix,
      );
      navigate(withSearchPathPrefix);
  }
};

const Chart = ({ options, searchPathPrefix, type }: ChartComponentProps) => {
  const [echarts, setEcharts] = useState<any | null>(null);
  const navigate = useNavigate();
  const chartRef = useRef<ReactEChartsCore>(null);
  const [imageUrl, setImageUrl] = useState<string | null>(null);
  const mediaMode = useMediaMode();
  const { ready: templateRenderReady, renderTemplates } = useTemplateRender();

  // Dynamically import echarts from its own bundle
  useEffect(() => {
    import("./echarts").then((m) => setEcharts(m.echarts));
  }, []);

  useEffect(() => {
    if (!chartRef.current || !options) {
      return;
    }

    const echartInstance = chartRef.current.getEchartsInstance();
    const dataURL = echartInstance.getDataURL({});
    if (dataURL === imageUrl) {
      return;
    }
    setImageUrl(dataURL);
  }, [chartRef, imageUrl, options]);

  if (!options) {
    return null;
  }

  const eventsDict = {
    click: (params) =>
      handleClick(params, navigate, renderTemplates, searchPathPrefix),
  };

  const PlaceholderComponent = Placeholder.component;

  return (
    <PlaceholderComponent ready={!!echarts && templateRenderReady}>
      <>
        {mediaMode !== "print" && (
          <div className="relative">
            <ReactEChartsCore
              ref={chartRef}
              echarts={echarts}
              className="chart-canvas"
              onEvents={eventsDict}
              option={options}
              notMerge={true}
              lazyUpdate={true}
              style={
                type === "pie" || type === "donut" ? { height: "250px" } : {}
              }
            />
          </div>
        )}
        {mediaMode === "print" && imageUrl && (
          <div>
            <img alt="Chart" className="max-w-full max-h-full" src={imageUrl} />
          </div>
        )}
      </>
    </PlaceholderComponent>
  );
};

const ChartWrapper = (props: ChartProps) => {
  const {
    searchPathPrefix,
    themeContext: { wrapperRef },
  } = useDashboard();
  const themeColors = useChartThemeColors();

  if (!wrapperRef) {
    return null;
  }

  if (!props.data) {
    return null;
  }

  return (
    <Chart
      options={buildChartOptions(props, themeColors)}
      searchPathPrefix={searchPathPrefix}
      type={props.display_type || "column"}
    />
  );
};

const renderChart = (definition: ChartProps) => {
  // We default to column charts if not specified
  const { display_type = "column" } = definition;

  const chart = getChartComponent(display_type);

  if (!chart) {
    return <ErrorPanel error={`Unknown chart type ${display_type}`} />;
  }

  const Component = chart.component;
  return <Component {...definition} />;
};

const RenderChart = (props: ChartProps) => {
  return renderChart(props);
};

registerComponent("chart", RenderChart);

export default ChartWrapper;

export { Chart };
