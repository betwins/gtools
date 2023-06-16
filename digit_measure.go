package gtools

import (
	"errors"
	"regexp"
	"strconv"
)

func SplitDigitMeasure(measure string) (int, string, error) {
	re := regexp.MustCompile(`(\d+)([KMG])`)
	match := re.FindStringSubmatch(measure)
	if match != nil {
		num, _ := strconv.Atoi(match[1])
		unit := match[2]
		return num, unit, nil
	}
	return 0, "", errors.New("invalid digit measure")
}
