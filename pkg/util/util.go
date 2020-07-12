package util

import (
	"github.com/jinzhu/copier"
)

// Struct Map To other Struct
func StructMapToStruct(s, ts interface{}) error {
	return copier.Copy(ts, s)
}
