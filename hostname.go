package lfe

import (
	"fmt"
	"os"
	"strings"
)

func runtimeHostname() (string, error) {
	h, e := os.Hostname()
	if e != nil {
		return "", fmt.Errorf("lfe: hostname: %w", e)
	}
	h = strings.TrimSpace(h)
	if h == "" {
		return "", fmt.Errorf("lfe: hostname is empty")
	}
	return h, nil
}
