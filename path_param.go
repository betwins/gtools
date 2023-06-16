package gtools

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
)

func ExtractInt(pattern, path string, pos int) (int, error) {
	re := regexp.MustCompile(pattern)
	match := re.FindStringSubmatch(path)
	if len(match) < pos+1 {
		errStr := fmt.Sprintf("no match was found #v", match)
		return 0, errors.New(errStr)
	}

	param, err := strconv.Atoi(match[pos])
	if err != nil {
		errStr := fmt.Sprintf("invalid match data type #v", match[pos])
		return 0, errors.New(errStr)
	}
	return param, nil
}

func ExtractString(pattern, path string, pos int) (string, error) {
	re := regexp.MustCompile(pattern)
	match := re.FindStringSubmatch(path)
	if len(match) < pos+1 {
		errStr := fmt.Sprintf("no match was found #v", match)
		return "", errors.New(errStr)
	}

	return match[pos], nil
}
