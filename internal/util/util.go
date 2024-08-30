package util

import (
	"fmt"
)

func PairsToOptions(prefix string, kvs []string) ([]string, error) {
	if prefix != "" {
		prefix = prefix + ":"
	}

	if len(kvs) == 0 {
		return nil, nil // TODO: erase the options?
	}
	if len(kvs)%2 != 0 {
		return nil, fmt.Errorf("final args should be in k1,v1,k2,v2,... pairs")
	}
	out := make([]string, 0, len(kvs)/2)
	for i := 0; i < len(kvs); i += 2 {
		k := kvs[i]
		v := kvs[i+1]

		out = append(out, prefix+k+"="+v)
	}
	return out, nil
}
