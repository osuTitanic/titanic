package constants

import "golang.org/x/exp/constraints"

type Enum[T constraints.Signed] interface {
	// Values returns a slice of all valid enum values for the type T
	Values() []T
	// TODO: add String & Value methods maybe?
}
