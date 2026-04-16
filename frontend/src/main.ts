import "./style.css";
import "./app.css";

import { ActiveLabChangeEvent } from "./events";
import { createLabTabs, updateAllFieldLabels } from "./lab-init";

import { GetLabs } from "../wailsjs/go/main/App";
import { charting } from "../wailsjs/go/models";
import { EventsOn } from "../wailsjs/runtime";
import { renderMultiChart, renderChartInto } from "./chart-render";
import { registry } from "./registry";
import { SafeChart } from "./types";

// State

const stopLoading = () => {
  const btn = document.getElementById("rerender-btn");
  if (btn) btn.classList.remove("loading");
};

window.addEventListener("load", () => {
  GetLabs().then((labsResponse) => {
    console.log("Loaded labs:", labsResponse.labs);

    // Register all labs in the registry without initializing UI
    labsResponse.labs.forEach((lab) => {
      registry.addLab(lab.ID, lab);
    });

    // Create lab tabs
    createLabTabs(labsResponse.labs);

    // Initialize first lab by dispatching activeLabChange event
    if (labsResponse.labs.length > 0) {
      window.dispatchEvent(new ActiveLabChangeEvent(labsResponse.labs[0].ID));
    }
  });
});

EventsOn("renderComplete", (data: charting.RenderResponse) => {
  console.log("Render complete:", data);
  stopLoading();
  if (data.error) {
    console.error("Render error:", data.error);
    const errorContainer = document.getElementById("error-container");
    if (errorContainer) {
      errorContainer.innerHTML = `<div style="color: red; padding: 20px;">Error: ${data.error.Message}</div>`;
    }
  } else {
    console.log("Render completed successfully,", window.activeChartId);

    const activeChartData = data.charts[window.activeChartId!] as unknown as SafeChart;

    if (activeChartData.type.startsWith("multi-")) {
      renderMultiChart(activeChartData);
    } else {
      let container = document.getElementById("chart-container");
      if (!container) {
        container = document.createElement("div");
        container.id = "chart-container";
        document.body.appendChild(container);
      }
      renderChartInto(activeChartData, container);
    }

    // Update all field labels to reflect actual rendered values
    updateAllFieldLabels(activeChartData);

    const errorContainer = document.getElementById("error-container");
    if (errorContainer) {
      errorContainer.innerHTML = "";
    }
  }
});

// Listen for render errors
EventsOn("renderError", (data: charting.RenderResponse) => {
  console.error("Render error:", data);
  stopLoading();
  const container = document.getElementById("error-container");
  if (container) {
    container.innerHTML = `<div style="color: red; padding: 20px;">Error: ${data.error.Message}</div>`;
  }
});
