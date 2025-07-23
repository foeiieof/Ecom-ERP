package infrastructure

import "github.com/jinzhu/copier"

func MapStruct[T any, U any](from T) (U, error) {
	var to U
	err := copier.Copy(&to, &from)
	return to, err
}

