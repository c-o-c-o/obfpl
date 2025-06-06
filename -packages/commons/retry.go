package commons

func Retry(c int, f func(c int) error, d func(c int)) error {
	for i := 0; i < c; i++ {
		err := f(i)

		if err == nil {
			break
		}

		if i+1 >= c {
			return err
		}
		d(i)
	}
	return nil
}
