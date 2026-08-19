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

func TestGetAutoFilter(t *testing.T) {
	f := NewFile()
	for cell, v := range map[string]interface{}{
		"A1": "Region", "B1": "Sales",
		"A2": "East", "B2": 1000,
		"A3": "West", "B3": 3000,
		"A4": "North", "B4": 5000,
	} {
		assert.NoError(t, f.SetCellValue("Sheet1", cell, v))
	}
	require.NoError(t, f.AutoFilter("Sheet1", "A1:B4", []AutoFilterOptions{
		{Column: "B", Expression: "x > 2000"},
	}))

	path := filepath.Join(t.TempDir(), "autofilter.xlsx")
	require.NoError(t, f.SaveAs(path))
	require.NoError(t, f.Close())

	f2, err := OpenFile(path)
	require.NoError(t, err)
	defer func() { assert.NoError(t, f2.Close()) }()

	result, ok, err := f2.GetAutoFilter("Sheet1")
	require.NoError(t, err)
	require.True(t, ok)
	assert.Equal(t, "A1:B4", result.Range)
	require.Len(t, result.Columns, 1)
	col := result.Columns[0]
	assert.Equal(t, 1, col.ColID) // B is offset 1 from A
	require.Len(t, col.CustomFilters, 1)
	assert.Equal(t, ">", col.CustomFilters[0].Operator)
	assert.Equal(t, "2000", col.CustomFilters[0].Value)

	// A sheet without an AutoFilter reports ok == false.
	_, err = f2.NewSheet("Empty")
	require.NoError(t, err)
	_, ok, err = f2.GetAutoFilter("Empty")
	assert.NoError(t, err)
	assert.False(t, ok)
}

func TestGetAutoFilterDiscrete(t *testing.T) {
	f := NewFile()
	require.NoError(t, f.AutoFilter("Sheet1", "A1:A4", []AutoFilterOptions{
		{Column: "A", Expression: "x == East or x == West"},
	}))
	path := filepath.Join(t.TempDir(), "discrete.xlsx")
	require.NoError(t, f.SaveAs(path))
	require.NoError(t, f.Close())

	f2, err := OpenFile(path)
	require.NoError(t, err)
	defer func() { assert.NoError(t, f2.Close()) }()

	result, ok, err := f2.GetAutoFilter("Sheet1")
	require.NoError(t, err)
	require.True(t, ok)
	require.Len(t, result.Columns, 1)
	assert.ElementsMatch(t, []string{"East", "West"}, result.Columns[0].Filters)
}
