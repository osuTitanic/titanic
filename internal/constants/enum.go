package constants

type Enum[T any] interface {
	Value() T
	String() string
	Valid() bool
}
