package maps

func MapValues[K comparable, V any](m map[K]V ) []V {
	vv := make([]V, 0, len(m))
	for _, v := range m {
		vv = append(vv, v)
	}
	return vv
}

func MapKeys[K comparable, V any](m map[K]V ) []K {
	vv := make([]K, 0, len(m))
	for k := range m {
		vv = append(vv, k)
	}
	return vv
}