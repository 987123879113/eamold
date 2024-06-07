package utils

import (
	"fmt"
	"strconv"
	"strings"
)

func GenerateListStringInt64(input []int64) string {
	var output strings.Builder

	for _, v := range input {
		output.WriteString(fmt.Sprintf("%d?", v))
	}

	return output.String()
}

func GenerateListString(input []string) string {
	if len(input) == 0 {
		return ""
	}

	return strings.Join(input, "?") + "?"
}

func SplitListStringInt64(input string) []int64 {
	musicIdsSplit := strings.Split(strings.TrimRight(input, "?"), "?")

	output := make([]int64, len(musicIdsSplit))
	for i, v := range musicIdsSplit {
		val, _ := strconv.ParseInt(v, 10, 64)
		output[i] = val
	}

	return output
}
