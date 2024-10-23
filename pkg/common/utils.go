/*
Copyright © 2024 netr0m <netr0m@pm.me>
*/

package common

import (
	"fmt"
	"log"
	"log/slog"
	"os"
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
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	slog.SetDefault(logger)
}

func (e *Error) Unwrap() error { return e.Err }

func (e *Error) Error() string {
	return fmt.Sprintf("%s failed with status %s: %s", e.Operation, e.Status, e.Message)
}
