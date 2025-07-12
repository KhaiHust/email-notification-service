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
