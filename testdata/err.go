package demo

type ErrNo int32

const (
	ENo ErrNo = 1003 // 号码出错 @TakeCodeToMap(InfoMap)

	ENotFound ErrNo = 1004 // 找不到 @TakeCodeToMap(InfoMap)
)
