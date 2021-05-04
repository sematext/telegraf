package tags

import (
	"bytes"
	"sort"
)

const (
	SematextProcessedTag = "sematext.processed"
)

func GetTagsKey(tags map[string]string) string {
	var output bytes.Buffer

	// sort tags to ensure the order
	sortedTagKeys := make([]string, 0, len(tags))
	for t := range tags {
		sortedTagKeys = append(sortedTagKeys, t)
	}
	sort.Strings(sortedTagKeys)

	for _, tagKey := range sortedTagKeys {
		output.WriteString(tagKey)
		output.WriteString("=")
		output.WriteString(tags[tagKey])
		output.WriteString(",")
	}

	return output.String()
}