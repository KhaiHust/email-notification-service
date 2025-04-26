package utils

import (
	"github.com/jinzhu/copier"
	"time"
)

func TimeToInt64Converter() copier.TypeConverter {
	return copier.TypeConverter{
		SrcType: time.Time{},
		DstType: int64(0),
		Fn: func(src interface{}) (interface{}, error) {
			return src.(time.Time).Unix(), nil
		},
	}
}
func CopyStruct(toValue, fromValue interface{}) error {
	return copier.CopyWithOption(toValue, fromValue, copier.Option{
		Converters: []copier.TypeConverter{
			TimeToInt64Converter(),
		},
	})
}
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
