package state

func (state *State) RegisterExtension(key string, value any) {
	if state == nil || state.Extensions == nil {
		return
	}
	state.Extensions[key] = value
}

func (state *State) GetExtension[T any](key string) (T, bool) {
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
