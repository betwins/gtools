package gtools

import (
	"errors"
	"regexp"
	"strconv"
)

func SplitPeriodMeasure(measure string) (int, string, error) {
	re := regexp.MustCompile(`^(\d+)([DMY])$`)
	match := re.FindStringSubmatch(measure)
	if match != nil {
		num, _ := strconv.Atoi(match[1])
		unit := match[2]
		return num, unit, nil
	}
	return 0, "", errors.New("invalid period measure")
}
