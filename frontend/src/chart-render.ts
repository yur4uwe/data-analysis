import * as Plotly from "plotly.js-dist-min";
import { charting } from "../wailsjs/go/models";
import { defaultPlotlyLayout, newPlotlyAxes, processDatasetToPlotly } from "./static-config";
import { SafeChart } from "./types";

/**
 * Destroys all currently active Plotly instances.
 */
export function destroyAllCharts() {
  const containers = document.querySelectorAll(".plotly-graph-div");
  containers.forEach((container) => {
    Plotly.purge(container as HTMLElement);
  });
}

/**
 * Destroys a specific chart instance.
 */
export function destroyPreviousChart(containerId: string) {
  const container = document.getElementById(containerId);
  if (container) {
    Plotly.purge(container);
    container.innerHTML = "";
  }
}

export function renderChartInto(chartConfig: SafeChart, container: HTMLElement) {
  if (!chartConfig) {
    console.error("renderChartInto: chartConfig is null or undefined!");
    return;
  }

  const chartType = chartConfig.type || "line";
  const plotlyContainerId = `plotly-${chartConfig.id}`;
  
  // Reset container to block (in case it was a grid from multi-chart mode)
  container.style.display = "block";
  container.style.gridTemplateColumns = "none";
  
  // Clear previous content and create a dedicated plotly div
  container.innerHTML = "";
  const plotlyDiv = document.createElement("div");
  plotlyDiv.id = plotlyContainerId;
  plotlyDiv.style.width = "100%";
  plotlyDiv.style.height = "100%";
  plotlyDiv.className = "plotly-graph-div";
  container.appendChild(plotlyDiv);

  console.log("Rendering Plotly chart ID:", chartConfig.id, "Type:", chartType);

  const labels = Array.isArray(chartConfig.labels) ? chartConfig.labels : [];
  const traces = Object.values(chartConfig.datasets || {}).map(processDatasetToPlotly(chartType, labels));

  const layout = {
    ...defaultPlotlyLayout(chartConfig.title || "", chartType, chartConfig as any),
    ...newPlotlyAxes(chartConfig as any)
  };

  const config: Partial<Plotly.Config> = {
    responsive: true,
    displaylogo: false,
    modeBarButtonsToRemove: ["select2d", "lasso2d"]
  };

  // Use requestAnimationFrame to ensure the div is in the DOM and has dimensions
  requestAnimationFrame(() => {
    Plotly.newPlot(plotlyDiv, traces as any, layout as any, config).then(() => {
      // One more resize to be absolutely sure layout changes during render didn't break it
      Plotly.Plots.resize(plotlyDiv);
    });
  });
}

export function renderMultiChart(chartConfig: SafeChart) {
  if (!chartConfig || !chartConfig.datasets) {
    console.error("renderMultiChart: chartConfig or datasets is missing");
    return;
  }

  const container = document.getElementById("chart-container")!;
  if (!container) return;

  container.innerHTML = "";
  container.style.display = "grid";
  container.style.gridTemplateColumns = "repeat(auto-fit, minmax(450px, 1fr))";
  container.style.gap = "20px";

  const singleType = (chartConfig.type || "").replace("multi-", "");

  // 1. Create all wrappers first so the grid layout is fully established
  const tasks = Object.entries(chartConfig.datasets).map(([datasetId, dataset]) => {
    if (!dataset) return null;

    const wrapper = document.createElement("div");
    wrapper.className = "chart-wrapper";
    wrapper.style.minHeight = "450px";
    wrapper.style.minWidth = "0"; 
    wrapper.style.overflow = "hidden";
    container.appendChild(wrapper);

    const syntheticChart = charting.Chart.createFrom({
      ...chartConfig,
      type: singleType,
      id: `${chartConfig.id}-${datasetId}`,
      title: dataset.label || datasetId,
      datasets: { [datasetId]: dataset },
    }) as unknown as SafeChart;

    return { chart: syntheticChart, wrapper };
  }).filter(t => t !== null);

  // 2. Render all charts after the browser has had a chance to reflow the grid
  requestAnimationFrame(() => {
    tasks.forEach(task => {
      if (task) renderChartInto(task.chart, task.wrapper);
    });
  });
}

// Helper functions for dataset control
export function toggleDatasetVisibility(chartId: string, datasetIndex: number) {
  // Plotly containers in our app have IDs like `plotly-${chartId}`
  const plotlyContainerId = `plotly-${chartId}`;
  const container = document.getElementById(plotlyContainerId);
  
  if (container) {
    const data = (container as any).data;
    if (data && data[datasetIndex]) {
      const currentVisible = data[datasetIndex].visible;
      const nextVisible = currentVisible === true ? "legendonly" : true;
      
      Plotly.restyle(container, { visible: nextVisible } as any, [datasetIndex]);
    }
  }
}
