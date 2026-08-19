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

// TestGetConditionalFormatThreeColorScale guards the three-color-scale
// read-back: the middle stop's value must not leak into MaxValue.
func TestGetConditionalFormatThreeColorScale(t *testing.T) {
	f := NewFile()
	require.NoError(t, f.SetConditionalFormat("Sheet1", "D1:D5", []ConditionalFormatOptions{
		{
			Type: "3_color_scale", Criteria: "=",
			MinType: "min", MinColor: "#FF0000",
			MidType: "percentile", MidValue: "50", MidColor: "#FFFF00",
			MaxType: "max", MaxColor: "#00FF00",
		},
	}))
	path := filepath.Join(t.TempDir(), "cf3.xlsx")
	require.NoError(t, f.SaveAs(path))
	require.NoError(t, f.Close())

	f2, err := OpenFile(path)
	require.NoError(t, err)
	defer func() { assert.NoError(t, f2.Close()) }()

	cfs, err := f2.GetConditionalFormats("Sheet1")
	require.NoError(t, err)
	require.Len(t, cfs["D1:D5"], 1)
	opt := cfs["D1:D5"][0]
	assert.Equal(t, "3_color_scale", opt.Type)
	assert.Equal(t, "50", opt.MidValue)
	// The maximum stop is "max" with no explicit value; MaxValue must be
	// empty rather than the middle stop's "50".
	assert.Empty(t, opt.MaxValue)
	assert.Equal(t, "#FF0000", opt.MinColor)
	assert.Equal(t, "#FFFF00", opt.MidColor)
	assert.Equal(t, "#00FF00", opt.MaxColor)
}
