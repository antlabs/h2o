package model

type KeyVal[T, U any] struct {
	Key       T
	Val       U
	RawVal    any
	IsString  bool
	IsInt     bool
	IsFloat64 bool
}

func (k *KeyVal[T, U]) SetIs() KeyVal[T, U] {
	switch k.RawVal.(type) {
	case string:
		k.IsString = true
	case int:
		k.IsInt = true
	case float64:
		k.IsFloat64 = true
	}
	return *k
}
