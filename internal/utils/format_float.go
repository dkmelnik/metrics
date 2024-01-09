package utils

import (
	"fmt"
	"strconv"
	"strings"
)

func FormatFloat(num float64, prc int) string {
	var (
		zero, dot = "0", "."

		str = fmt.Sprintf("%."+strconv.Itoa(prc)+"f", num)
	)
	return strings.TrimRight(strings.TrimRight(str, zero), dot)
}
