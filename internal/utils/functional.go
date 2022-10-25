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

func Keys[K comparable, V any](m map[K]V) []K {
	keys := make([]K, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

func MapKeys[K comparable, L any, V any](f func(K) L, m map[K]V) []L {
	ls := make([]L, 0, len(m))
	for k := range m {
		ls = append(ls, f(k))
	}
	return ls
}

func Values[K comparable, V any](m map[K]V) []V {
	values := make([]V, 0, len(m))
	for _, v := range m {
		values = append(values, v)
	}

	return values
}

func MapValues[K comparable, V any, W any](f func(V) W, m map[K]V) []W {
	ws := make([]W, 0, len(m))
	for _, v := range m {
		ws = append(ws, f(v))
	}
	return ws
}

func First[A any, B any](a A, _ B) A {
	return a
}

func Second[A any, B any](_ A, b B) B {
	return b
}

func FindFirstI[A any](as []A, f func(A) bool) (int, bool) {
	for i := range as {
		if f(as[i]) {
			return i, true
		}
	}

	return 0, false
}
