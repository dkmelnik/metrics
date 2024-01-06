package utils

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_FormatFloat(t *testing.T) {
	testCases := []struct {
		inputNum float64
		inputPrc int
		expected string
	}{
		{150112.000, 3, "150112"},
		{3.833856e+06, 2, "3833856"},
		{-13.0, 3, "-13"},
		{7.70766, 4, "7.7077"},
		{282, 3, "282"},
		{3485734.100, 3, "3485734.1"},
		{-3.35872e+06, 5, "-3358720"},
		{0, 2, "0"},
	}

	for _, tc := range testCases {
		result := FormatFloat(tc.inputNum, tc.inputPrc)
		assert.Equal(t, tc.expected, result, fmt.Sprintf("FormatFloat(%f, %d)", tc.inputNum, tc.inputPrc))
	}
}
