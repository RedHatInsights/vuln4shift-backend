package test

import (
	"fmt"
	"time"
)

func GetFloat32PtrValue(f *float32) float32 {
	if f == nil {
		return 0
	}
	return *f
}

func GetStringPtrValue(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func GetUTC(t *time.Time) time.Time {
	if t == nil {
		return time.Time{}
	}
	return t.UTC()
}

func GetMetaKeys(ks map[string]bool) (keys []string) {
	for k := range ks {
		keys = append(keys, k)
	}
	return keys
}

func GetClusterDetailKeys(ks map[string]struct{}) (keys []string) {
	for k := range ks {
		keys = append(keys, k)
	}
	return keys
}

func GetMetaStringSlice(i interface{}) []string {
	arr := i.([]interface{})
	res := make([]string, 0, len(arr))
	for _, v := range arr {
		res = append(res, fmt.Sprint(v))
	}
	return res
}

func GetMetaTotalItems(m interface{}) float64 {
	return m.(map[string]interface{})["total_items"].(float64)
}
