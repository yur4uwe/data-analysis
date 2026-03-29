import { charting } from "../wailsjs/go/models"

/**
 * Shared coordinates for points in 2D space.
 * Used by GridDataset (Scatter, Line, Bubble).
 */
export type DataPoint = {
    x: number
    y: number | null
}

/**
 * Extended coordinate for heatmaps.
 * Includes a value 'v' for the intensity/color mapping.
 */
export type HeatmapPoint = {
    v: number | null
} & DataPoint

/**
 * Mirror of charting.BaseDataset in Go.
 * Contains common visual and behavior configuration.
 */
export type BaseDataset = {
    label: string
    type: string
    borderColor: string
    borderWidth: number
    hidden: boolean
    togglable: boolean
    dataLabels: string[]
    fields: charting.MutableField[]
}

/**
 * Represents a dataset with XY coordinates.
 * Mapped from charting.GridDataset.
 */
export type GridDataset = BaseDataset & {
    data: DataPoint[]
    backgroundColor: string
    pointRadius: number
    pointStyle: string
}

/**
 * Represents a dataset with indexed values (e.g. for Bar or Pie charts).
 * Mapped from charting.CategoricalDataset.
 * Data uses 'any[]' to allow for 'null' values (missing data).
 */
export type CategoricalDataset = BaseDataset & {
    data: (number | null)[]
    backgroundColor: string[]
}

/**
 * Represents a 3D dataset for Heatmaps.
 * Mapped from charting.HeatmapDataset.
 * Uses 'pointData' as the JSON key for coordinates.
 */
export type HeatmapDataset = BaseDataset & {
    pointData: HeatmapPoint[]
    backgroundColor: string[]
}

/**
 * Discriminated union of all possible dataset types.
 */
export type Dataset = GridDataset | CategoricalDataset | HeatmapDataset

/**
 * Type-safe extension of the generated Chart model.
 */
export interface SafeChart extends Omit<charting.Chart, 'datasets'> {
    datasets: Record<string, Dataset>
}

/**
 * Type guard for HeatmapDataset.
 */
export function isHeatmapDataset(ds: Dataset): ds is HeatmapDataset {
    return (ds as HeatmapDataset).pointData !== undefined;
}

/**
 * Type guard for GridDataset.
 * Checks if the data array contains XY point objects.
 */
export function isGridDataset(ds: Dataset): ds is GridDataset {
    if (isHeatmapDataset(ds) || !Array.isArray(ds.data) || ds.data.length === 0) return false;
    // Find the first non-null element to determine type
    const firstValid = (ds.data as any[]).find(v => v !== null);
    if (!firstValid) return true; // Treat empty/all-null as compatible
    return typeof firstValid === 'object' && 'x' in firstValid;
}

/**
 * Type guard for CategoricalDataset.
 * Checks if the data array contains primitive numbers.
 */
export function isCategoricalDataset(ds: Dataset): ds is CategoricalDataset {
    if (isHeatmapDataset(ds) || !Array.isArray(ds.data)) return false;
    if (ds.data.length === 0) return true;
    const firstValid = (ds.data as any[]).find(v => v !== null);
    if (!firstValid) return true;
    return typeof firstValid === 'number';
}
