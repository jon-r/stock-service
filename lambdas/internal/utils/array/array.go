package array

import "reflect"

func UnpackArray(s any) []any {
	v := reflect.ValueOf(s)
	r := make([]any, v.Len())
	for i := 0; i < v.Len(); i++ {
		r[i] = v.Index(i).Interface()
	}
	return r
}

func Find[T any](collection []T, predicate func(item T) bool) (T, bool) {
	for _, item := range collection {
		if predicate(item) {
			return item, true
		}
	}

	var result T
	return result, false
}
