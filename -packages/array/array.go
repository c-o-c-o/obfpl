package array

func In[T comparable](array []T, val T) bool {
	for _, v := range array {
		if v == val {
			return true
		}
	}
	return false
}

func Map[T, R any](array []T, f func(val T) R) []R {
	r := make([]R, 0, len(array))

	for _, v := range array {
		r = append(r, f(v))
	}

	return r
}

func MapE[T, R any](array []T, f func(val T) (R, error)) ([]R, error) {
	l := make([]R, 0, len(array))

	for _, v := range array {
		r, err := f(v)
		if err != nil {
			return nil, err
		}

		l = append(l, r)
	}

	return l, nil
}

func FilterE[T any](array []T, f func(val T) (bool, error)) ([]T, error) {
	r := make([]T, 0, len(array))

	for _, v := range array {
		ok, err := f(v)
		if err != nil {
			return nil, err
		}

		if ok {
			r = append(r, v)
		}
	}

	return r, nil
}
