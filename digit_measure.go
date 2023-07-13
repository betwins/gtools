package gtools

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
)

func SplitDigitMeasure(measure string) (int, string, error) {
	re := regexp.MustCompile(`^(\d+)([KMG]*)$`)
	match := re.FindStringSubmatch(measure)
	if match != nil {
		num, _ := strconv.Atoi(match[1])
		unit := match[2]
		return num, unit, nil
	}
	return 0, "", errors.New("invalid digit measure")
}

func multi1024(before float64) float64 {
	return before * 1024
}

func multi1024twice(before float64) float64 {
	return before * 1024 * 1024
}

func multi1024thrice(before float64) float64 {
	return before * 1024 * 1024 * 1024
}

func divid1024(before float64) float64 {
	return before / 1024
}

func divid1024twice(before float64) float64 {
	return before / 1024 / 1024
}

func divid1024thrice(before float64) float64 {
	return before / 1024 / 1024
}

var unitConversionMap = map[string]func(before float64) float64{
	"K2M": divid1024,
	"K2G": divid1024twice,
	"M2K": multi1024,
	"M2G": divid1024,
	"G2K": multi1024twice,
	"G2M": multi1024,
	"K2":  multi1024,
	"M2":  multi1024twice,
	"G2":  multi1024thrice,
	"2K":  divid1024,
	"2M":  divid1024twice,
	"2G":  divid1024thrice,
}

func ConvertDigitMeasure(before float64, beforeUnit, afterUnit string) (float64, error) {

	convertFuncKey := fmt.Sprintf("%s2%s", beforeUnit, afterUnit)
	if convertFunc, exist := unitConversionMap[convertFuncKey]; exist {
		return convertFunc(before), nil
	} else {
		return 0, errors.New("invalid unit(not K、M、G")
	}
}
