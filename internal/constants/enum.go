package constants

type HasValidator interface {
	Valid() bool
}

type HasString interface {
	String() string
}

type HasValue[T any] interface {
	Value() T
}

type Enum[T any] interface {
	HasValue[T]
	HasString
	HasValidator
}
