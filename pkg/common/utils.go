/*
Copyright © 2024 netr0m <netr0m@pm.me>
*/

package common

import "fmt"

func (e *Error) Unwrap() error { return e.Err }

func (e *Error) Error() string {
	return fmt.Sprintf("%s failed with status %s: %s", e.Operation, e.Status, e.Message)
}
