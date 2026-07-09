package state

// one day... yes... one day, go will support generics inside struct functions...
// my god, why does this not work, this doesn't make any sense
// like, whats the difference?? they are both just functions
// i wrote the damn function, imma go use the damn function
// hmm hold up let me see if someone else has complained about this already...
// oh no fucking way
// https://forum.golangbridge.org/t/golang-generic-support-for-methods-defined-on-structs/35525/2
// there actually is
// apparently i have to add a [T any] to the struct itself wtf, that sucks balls
// mannnnnnn
// ok but why am i complaining, this will work just fine anyway haha harharharhehahhehaehahe

func RegisterExtension(state *State, key string, value any) {
	if state == nil || state.Extensions == nil {
		return
	}
	state.Extensions[key] = value
}

func GetExtension[T any](state *State, key string) (T, bool) {
	var zero T
	if state == nil || state.Extensions == nil {
		return zero, false
	}
	value, ok := state.Extensions[key]
	if !ok {
		return zero, false
	}
	typedValue, ok := value.(T)
	if !ok {
		return zero, false
	}
	return typedValue, true
}
