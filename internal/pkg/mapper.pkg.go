package pkg

import "github.com/jinzhu/copier"

func MapStruct[T any, U any](from T) (U, error) {
	var to U
	err := copier.Copy(&to, &from)
	return to, err
}

func MapSliceStruct[T any, U any](from []T) ([]U, error){
  to := make([]U,len(from))
  for i,item := range from {
    u, err := MapStruct[T,U](item)
    if err != nil {
      return nil,err
    }
    to[i] = u
  }
  return to ,nil
}
