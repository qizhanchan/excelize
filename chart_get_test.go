// Copyright 2016 - 2026 The excelize Authors. All rights reserved. Use of
// this source code is governed by a BSD-style license that can be found in
// the LICENSE file.

package excelize

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetCharts(t *testing.T) {
	f := NewFile()
	for cell, v := range map[string]interface{}{
		"A1": "Apple", "A2": "Orange", "A3": "Pear",
		"B1": 3, "B2": 2, "B3": 4,
	} {
		assert.NoError(t, f.SetCellValue("Sheet1", cell, v))
	}
	assert.NoError(t, f.AddChart("Sheet1", "E1", &Chart{
		Type: Col,
		Series: []ChartSeries{{
			Name:       "Sheet1!$A$1",
			Categories: "Sheet1!$A$1:$A$3",
			Values:     "Sheet1!$B$1:$B$3",
		}},
		Title:     ChartTitle{Paragraph: []RichTextRun{{Text: "Fruit"}}},
		Legend:    ChartLegend{Position: "bottom"},
		Dimension: ChartDimension{Width: 480, Height: 290},
	}))
	// Add a second, different chart type on the same sheet.
	assert.NoError(t, f.AddChart("Sheet1", "E20", &Chart{
		Type: Pie,
		Series: []ChartSeries{{
			Name:       "Sheet1!$A$1",
			Categories: "Sheet1!$A$1:$A$3",
			Values:     "Sheet1!$B$1:$B$3",
		}},
		Title: ChartTitle{Paragraph: []RichTextRun{{Text: "Share"}}},
	}))

	// Round-trip through disk so the charts are re-parsed, not served from
	// the in-memory Drawings cache.
	path := filepath.Join(t.TempDir(), "charts.xlsx")
	require.NoError(t, f.SaveAs(path))
	require.NoError(t, f.Close())

	f2, err := OpenFile(path)
	require.NoError(t, err)
	defer func() { assert.NoError(t, f2.Close()) }()

	cells, charts, err := f2.GetCharts("Sheet1")
	require.NoError(t, err)
	require.Len(t, charts, 2)
	require.Len(t, cells, 2)

	byCell := map[string]*Chart{}
	for i, c := range cells {
		byCell[c] = charts[i]
	}

	col, ok := byCell["E1"]
	require.True(t, ok, "expected a chart anchored at E1")
	assert.Equal(t, Col, col.Type)
	assert.Equal(t, "Fruit", col.Title.Paragraph[0].Text)
	assert.Equal(t, "bottom", col.Legend.Position)
	require.Len(t, col.Series, 1)
	assert.Equal(t, "Sheet1!$A$1", col.Series[0].Name)
	assert.Equal(t, "Sheet1!$A$1:$A$3", col.Series[0].Categories)
	assert.Equal(t, "Sheet1!$B$1:$B$3", col.Series[0].Values)
	assert.Greater(t, col.Dimension.Width, uint(0))
	assert.Greater(t, col.Dimension.Height, uint(0))

	pie, ok := byCell["E20"]
	require.True(t, ok, "expected a chart anchored at E20")
	assert.Equal(t, Pie, pie.Type)
	assert.Equal(t, "Share", pie.Title.Paragraph[0].Text)

	// GetChartCells is the anchor-only convenience wrapper.
	chartCells, err := f2.GetChartCells("Sheet1")
	require.NoError(t, err)
	assert.ElementsMatch(t, []string{"E1", "E20"}, chartCells)

	// A sheet with no drawing returns no charts and no error.
	_, err = f2.NewSheet("Empty")
	require.NoError(t, err)
	cells, charts, err = f2.GetCharts("Empty")
	assert.NoError(t, err)
	assert.Empty(t, cells)
	assert.Empty(t, charts)
}

func TestGetChartsStacked(t *testing.T) {
	f := NewFile()
	assert.NoError(t, f.AddChart("Sheet1", "A1", &Chart{
		Type:   BarStacked,
		Series: []ChartSeries{{Values: "Sheet1!$B$1:$B$3"}},
	}))
	path := filepath.Join(t.TempDir(), "stacked.xlsx")
	require.NoError(t, f.SaveAs(path))
	require.NoError(t, f.Close())

	f2, err := OpenFile(path)
	require.NoError(t, err)
	defer func() { assert.NoError(t, f2.Close()) }()

	_, charts, err := f2.GetCharts("Sheet1")
	require.NoError(t, err)
	require.Len(t, charts, 1)
	assert.Equal(t, BarStacked, charts[0].Type)
}
