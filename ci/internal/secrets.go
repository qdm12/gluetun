package internal

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"
)

type Logger interface {
	Info(msg string)
	Infof(format string, args ...any)
}

func readSecrets(ctx context.Context, expectedSecrets []string,
	logger Logger,
) (lines []string, err error) {
	scanner := bufio.NewScanner(os.Stdin)
	lines = make([]string, 0, len(expectedSecrets))

	for i := range expectedSecrets {
		logger.Infof("🤫 reading %s from Stdin...", expectedSecrets[i])
		if !scanner.Scan() {
			break
		}
		lines = append(lines, strings.TrimSpace(scanner.Text()))
		logger.Infof("🤫 %s secret read successfully", expectedSecrets[i])
		if ctx.Err() != nil {
			return nil, ctx.Err()
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("reading secrets from stdin: %w", err)
	}

	if len(lines) < len(expectedSecrets) {
		return nil, fmt.Errorf("expected %d secrets via Stdin, but only received %d",
			len(expectedSecrets), len(lines))
	}
	for i, line := range lines {
		if line == "" {
			return nil, fmt.Errorf("secret on line %d/%d was empty", i+1, len(lines))
		}
	}

	return lines, nil
}
