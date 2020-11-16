package main

func s1(s string) int {
	if s == "" {
		return 0
	}·

	n := 1
	for range s {
		n++
	}
	return n
}
