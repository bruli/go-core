package fixtures

func SetData[T any](def T, value *T) T {
	if value == nil {
		return def
	}
	return *value
}
