package dict

func Keys[K comparable, V any](dict map[K]V) []K {
	r := make([]K, 0, len(dict))

	for k := range dict {
		r = append(r, k)
	}

	return r
}

func Values[K comparable, V any](dict map[K]V) []V {
	r := make([]V, 0, len(dict))

	for _, v := range dict {
		r = append(r, v)
	}

	return r
}

func Map[K, RK comparable, V comparable, RV any](dict map[K]V, f func(k K, v V) (RK, RV)) map[RK]RV {
	r := make(map[RK]RV)

	for k, v := range dict {
		rk, rv := f(k, v)
		r[rk] = rv
	}

	return r
}

func FromArray[T, K, V comparable](array []T, f func(v T) (K, V)) map[K]V {
	r := make(map[K]V)

	for _, v := range array {
		key, val := f(v)
		r[key] = val
	}

	return r
}

func FromArrayE[T, K, V comparable](array []T, f func(v T) (K, V, error)) (map[K]V, error) {
	r := make(map[K]V)

	for _, v := range array {
		key, val, err := f(v)
		if err != nil {
			return nil, err
		}

		r[key] = val
	}

	return r, nil
}
