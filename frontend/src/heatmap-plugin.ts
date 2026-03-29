export function hexToRgb(hex: string) {
    hex = hex.replace("#", "");
    if (hex.length === 3) {
        hex = hex.split("").map((s) => s + s).join("");
    }
    const num = parseInt(hex, 16);
    return {
        r: (num >> 16) & 255,
        g: (num >> 8) & 255,
        b: num & 255,
    };
}

export const heatmapLegendPlugin = {
    id: "heatmapLegend",
    afterDraw: (chart: any) => {
        const chartType = chart.config.type;
        if (chartType !== "heatmap" && chartType !== "multi-heatmap") return;

        const ctx = chart.ctx;
        const area = chart.chartArea;
        const dataset = chart.data.datasets[0];
        if (!dataset || !dataset.data.length) return;

        // Retrieve min/max values and colors
        const values = dataset.data.map((d: any) => d.v);
        const min = Math.min(...values);
        const max = Math.max(...values);

        // Find the colors used for the heatmap (need to reach back to the processDataset logic or re-calculate)
        // For simplicity, we'll assume the interpolateColor logic is accessible or we use the colors from the dataset config if we can store them.
        // Since we can't easily import from static-config here without circular deps, 
        // we'll use a standard approach: pull from the dataset if we stored it there.
        const colors = dataset.backgroundColorList || ["#1d4ed8", "#b91c1c"];

        const legendWidth = 200;
        const legendHeight = 12;
        const x = area.left + (area.right - area.left - legendWidth) / 2;
        const y = chart.height - 40; // Relative to canvas bottom

        // Draw gradient bar
        const gradient = ctx.createLinearGradient(x, 0, x + legendWidth, 0);
        const gamma = 0.3; // Match the sensitivity in static-config.ts

        // Add more stops to simulate the power curve
        const numStops = 10;
        for (let i = 0; i <= numStops; i++) {
            const t = i / numStops;
            // We need to "invert" the gamma for the legend so the colors align
            // The cells use: normalizedColor = Math.pow(normalizedValue, gamma)
            // So at position 't' on the legend, we want the color for value 't^gamma'
            // Wait, let's just draw the colors at their actual mapped positions
            const mappedT = Math.pow(t, gamma);

            // Since our interpolateColor logic is in static-config.ts and hard to access here,
            // we interpolate between the first and last colors directly for the legend.
            // If we have more than 2 colors, we'd need a more complex loop.
            const startColor = hexToRgb(colors[0]);
            const endColor = hexToRgb(colors[colors.length - 1]);

            const r = Math.round(startColor.r + (endColor.r - startColor.r) * mappedT);
            const g = Math.round(startColor.g + (endColor.g - startColor.g) * mappedT);
            const b = Math.round(startColor.b + (endColor.b - startColor.b) * mappedT);

            gradient.addColorStop(t, `rgb(${r},${g},${b})`);
        }

        ctx.fillStyle = gradient;
        ctx.fillRect(x, y, legendWidth, legendHeight);

        // Draw labels
        ctx.fillStyle = "#000000";
        ctx.font = "bold 11px Arial";
        ctx.textAlign = "center";
        const formatValue = (v: number) => {
            if (Math.abs(v) < 0.0001 || Math.abs(v) > 10000) {
                return v.toExponential(2);
            }
            return v.toFixed(4);
        };
        ctx.fillText(formatValue(min), x, y + legendHeight + 14);
        ctx.fillText(formatValue(max), x + legendWidth, y + legendHeight + 14);
    }
};
