package common

type Iterator interface {
	HasNext() bool
	Next() (interface{}, bool)
}
