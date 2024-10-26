/*
Copyright © 2024 netr0m <netr0m@pm.me>
*/

package common

import (
	"fmt"
	"log/slog"
	"os"
	"strings"
)

func InitLogger(debugLogging bool) {
	lvl := new(slog.LevelVar)
	if debugLogging {
		lvl.Set(slog.LevelDebug)
	} else {
		lvl.Set(slog.LevelInfo)
	}

	_handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: lvl})
	logger := slog.New(_handler)

	slog.SetDefault(logger)
}

func (e *Error) Unwrap() error { return e.Err }

func (e *Error) Error() string {
	return fmt.Sprintf("%s failed with status %s: %s", e.Operation, e.Status, e.Message)
}

func (e *Error) Debug() string {
	var debugLines []string

	if e.Request != nil {
		debugLines = append(debugLines, fmt.Sprintf("Request:\n%v", e.Request))
	}
	if e.Response != nil {
		debugLines = append(debugLines, fmt.Sprintf("Response:\n%v", e.Response))
	}
	if e.Err != nil {
		debugLines = append(debugLines, fmt.Sprintf("Error:\n%v", e.Err.Error()))
	}

	return strings.Join(debugLines, "\n")
}
