package utils

import (
	"net/url"
	"strconv"
	"strings"
	"time"
)

func ToTimePointer(t int64) *time.Time {
	if t == 0 {
		return nil
	}
	tm := time.Unix(t, 0)
	return &tm
}
func ToUnixTimeToPointer(t *time.Time) *int64 {
	if t == nil {
		return nil
	}
	ux := t.Unix()
	return &ux
}
func FromUnixPointerToTime(t *int64) *time.Time {
	if t == nil {
		return nil
	}
	tm := time.Unix(*t, 0)
	return &tm
}
func ToInt64Pointer(t int64) *int64 {
	return &t
}
func ToStringPointer(t string) *string {
	return &t
}

func ParseInt64Pointer(value string) (*int64, error) {
	if value == "" {
		return nil, nil
	}
	v, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return nil, err
	}
	return &v, nil
}

func GetQueryInt64Pointer(values url.Values, key string) (*int64, error) {
	return ParseInt64Pointer(values.Get(key))
}
func GetInt64PointerWithDefault(values url.Values, key string, defaultVal int64) (*int64, error) {
	val := values.Get(key)
	if val == "" {
		return &defaultVal, nil
	}
	parsed, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		return nil, err
	}
	return &parsed, nil
}
func GetQueryStringPointer(values url.Values, key string) *string {
	if v := values.Get(key); v != "" {
		return &v
	}
	return nil
}

func GetQueryStringArray(values url.Values, key string) []string {
	if v := values.Get(key); v != "" {
		return strings.Split(v, ",")
	}
	return nil
}
