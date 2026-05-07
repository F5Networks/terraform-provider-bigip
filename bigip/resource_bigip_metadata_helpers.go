package bigip

import (
	"fmt"
	"sort"

	bigip "github.com/f5devcentral/go-bigip"
)

func tfMetadataToBigipMetadata(tf map[string]interface{}) ([]bigip.ResourceMetadata, error) {
	if tf == nil {
		return nil, nil
	}

	keys := make([]string, 0, len(tf))
	for k := range tf {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	metadata := make([]bigip.ResourceMetadata, 0, len(tf))
	for _, k := range keys {
		v := tf[k]
		s, ok := v.(string)
		if !ok {
			return nil, fmt.Errorf("metadata value for key %q must be a string", k)
		}
		metadata = append(metadata, bigip.ResourceMetadata{
			Name:    k,
			Value:   s,
			Persist: "true",
		})
	}

	return metadata, nil
}

func bigipMetadataToTfMetadata(metadata []bigip.ResourceMetadata) map[string]string {
	tf := make(map[string]string, len(metadata))
	for _, m := range metadata {
		if m.Name == "" {
			continue
		}
		tf[m.Name] = m.Value
	}
	return tf
}
