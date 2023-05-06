package disk

import (
	"fmt"
	"regexp"
	"strings"
	"syscall"
)

var errorMessageRegExp = regexp.MustCompile(`"code": (-?\d+),\n\t"message": "([^"]+)"`)

type ErrorMessage struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func parseErrorMessage(errStr string) (*ErrorMessage, error) {
	match := errorMessageRegExp.FindStringSubmatch(errStr)
	if len(match) == 0 {
		return nil, fmt.Errorf("failed to parse error message")
	}

	code := 0
	if _, err := fmt.Sscanf(match[1], "%d", &code); err != nil {
		return nil, fmt.Errorf("failed to parse error code: %w", err)
	}

	em := &ErrorMessage{
		Code:    code,
		Message: match[2],
	}

	return em, nil
}

func isFileExists(message string) bool {
	return strings.EqualFold(message, syscall.Errno(syscall.EEXIST).Error())
}

func isNoSuchDevice(message string) bool {
	return strings.EqualFold(message, syscall.Errno(syscall.ENODEV).Error())
}
