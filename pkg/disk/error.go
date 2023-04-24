package disk

import (
	"fmt"
	"regexp"
)

type ErrorMessage struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func parseErrorMessage(errStr string) (*ErrorMessage, error) {
	r := regexp.MustCompile(`"code": (-?\d+),\n\t"message": "([^"]+)"`)
	match := r.FindStringSubmatch(errStr)
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
