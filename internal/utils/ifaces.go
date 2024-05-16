package utils

func Ref[T any](V T) *T {
	return &V
}

func AnyArrToInterfaceArr(vals ...any) (arr []interface{}) {
	return append(arr, vals...)
}

func Int64ArrToInterfaceArr(vals ...int64) (arr []interface{}) {
	for _, val := range vals {
		arr = append(arr, val)
	}

	return arr
}
