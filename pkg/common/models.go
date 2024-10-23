/*
Copyright © 2024 netr0m <netr0m@pm.me>
*/
package common

/*
Common model for errors
*/
type Error struct {
	Operation string
	Message   string
	Status    string
	Err       error
	Response  interface{}
	Request   interface{}
}
