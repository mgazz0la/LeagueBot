package utils

func Map[A any, B any](f func(A) B, as []A) []B {
	bs := make([]B, len(as))
	for i := range as {
		bs[i] = f(as[i])
	}
	return bs
}

func Mapi[A any, B any](f func(A, int) B, as []A) []B {
	bs := make([]B, len(as))
	for i := range as {
		bs[i] = f(as[i], i)
	}
	return bs
}

func MapValues[K comparable, V any](m map[K]V) []V {
	vs := make([]V, len(m))
	for _, v := range m {
		vs = append(vs, v)
	}

	return vs
}
