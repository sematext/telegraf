// logic copied from plugins/serializers/influx/escape.go
package sematext

import "strings"

const (
	escapes            = "\t\n\f\r ,="
	nameEscapes        = "\t\n\f\r ,"
	stringFieldEscapes = "\t\n\f\r\\\""
)

var (
	escaper = strings.NewReplacer(
		"\t", `\t`,
		"\n", `\n`,
		"\f", `\f`,
		"\r", `\r`,
		`,`, `\,`,
		` `, `\ `,
		`=`, `\=`,
	)

	nameEscaper = strings.NewReplacer(
		"\t", `\t`,
		"\n", `\n`,
		"\f", `\f`,
		"\r", `\r`,
		`,`, `\,`,
		` `, `\ `,
	)

	stringFieldEscaper = strings.NewReplacer(
		`"`, `\"`,
		`\`, `\\`,
	)
)

// Escape a tagkey, tagvalue, or fieldkey
func escape(s string) string {
	if strings.ContainsAny(s, escapes) {
		return escaper.Replace(s)
	}
	return s
}

// Escape a measurement name
func nameEscape(s string) string {
	if strings.ContainsAny(s, nameEscapes) {
		return nameEscaper.Replace(s)
	}
	return s
}

// Escape a string field
func stringFieldEscape(s string) string {
	if strings.ContainsAny(s, stringFieldEscapes) {
		return stringFieldEscaper.Replace(s)
	}
	return s
}
