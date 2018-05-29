// +build !jsoniter
// +build !gojay

package json

import "encoding/json"

var (
	Marshal       = json.Marshal
	Unmarshal     = json.Unmarshal
	MarshalIndent = json.MarshalIndent
	NewDecoder    = json.NewDecoder
	NewEncoder    = json.NewEncoder
	Valid         = json.Valid
	Compact       = json.Compact
	HTMLEscape    = json.HTMLEscape
	Indent        = json.Indent
)

type (
	Number = json.Number
	RawMessage = json.RawMessage
)
