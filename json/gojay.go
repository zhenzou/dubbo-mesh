//+build gojay

package json

import (
	stdjson "encoding/json"

	"github.com/json-iterator/go"
	"github.com/francoispqt/gojay"
)

var (
	json          = jsoniter.ConfigCompatibleWithStandardLibrary
	Unmarshal     = gojay.Unmarshal
	Marshal       = gojay.Marshal
	MarshalIndent = json.MarshalIndent
	NewDecoder    = json.NewDecoder
	NewEncoder    = json.NewEncoder
	Valid         = json.Valid
	Compact       = stdjson.Compact
	HTMLEscape    = stdjson.HTMLEscape
	Indent        = stdjson.Indent
)

type (
	Number = jsoniter.Number
	RawMessage = jsoniter.RawMessage
)
