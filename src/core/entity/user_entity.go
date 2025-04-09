package entity

import (
	"bytes"
	"encoding/gob"
)

type UserEntity struct {
	BaseEntity
	FullName string
	Email    string
	Password string
}

func (e *UserEntity) UnmarshalBinary(data []byte) error {
	decoder := gob.NewDecoder(bytes.NewReader(data))
	return decoder.Decode(e)
}
