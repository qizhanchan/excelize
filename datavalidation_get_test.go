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

func TestGetDataValidationListItems(t *testing.T) {
	f := NewFile()
	// Inline list source.
	dvInline := NewDataValidation(true)
	dvInline.Sqref = "A1:A5"
	require.NoError(t, dvInline.SetDropList([]string{"East", "West", "North"}))
	require.NoError(t, f.AddDataValidation("Sheet1", dvInline))
	// Range list source, referencing seeded cells.
	for cell, v := range map[string]interface{}{"E1": "Red", "E2": "Green", "E3": "Blue"} {
		require.NoError(t, f.SetCellValue("Sheet1", cell, v))
	}
	dvRange := NewDataValidation(true)
	dvRange.Sqref = "B1:B5"
	dvRange.SetSqrefDropList("$E$1:$E$3")
	require.NoError(t, f.AddDataValidation("Sheet1", dvRange))

	path := filepath.Join(t.TempDir(), "dvlist.xlsx")
	require.NoError(t, f.SaveAs(path))
	require.NoError(t, f.Close())

	f2, err := OpenFile(path)
	require.NoError(t, err)
	defer func() { assert.NoError(t, f2.Close()) }()

	dvs, err := f2.GetDataValidations("Sheet1")
	require.NoError(t, err)
	require.Len(t, dvs, 2)

	bySqref := map[string]*DataValidation{}
	for _, dv := range dvs {
		bySqref[dv.Sqref] = dv
	}

	items, ok, err := f2.GetDataValidationListItems("Sheet1", bySqref["A1:A5"])
	require.NoError(t, err)
	assert.True(t, ok)
	assert.Equal(t, []string{"East", "West", "North"}, items)

	items, ok, err = f2.GetDataValidationListItems("Sheet1", bySqref["B1:B5"])
	require.NoError(t, err)
	assert.True(t, ok)
	assert.Equal(t, []string{"Red", "Green", "Blue"}, items)

	// Non-list validations report ok == false.
	dvNum := NewDataValidation(true)
	dvNum.Sqref = "C1:C5"
	require.NoError(t, dvNum.SetRange(1, 10, DataValidationTypeDecimal, DataValidationOperatorBetween))
	_, ok, err = f2.GetDataValidationListItems("Sheet1", dvNum)
	require.NoError(t, err)
	assert.False(t, ok)
}
