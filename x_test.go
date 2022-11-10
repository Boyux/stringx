package stringx

import "math/rand"

var elements = []rune{
	'1', '2', '3', '4', '5', '6', '7', '8', '9', '0', 'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm',
	'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z', 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J',
	'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z', '!', '@', '#', '$', '%', '^', '&',
	'*', '(', ')', '{', '}', '[', ']', '\'', '\'', '"', '"', '\r', '\n', '\v', '\t', ' ', '\\',
	'你', '好', '世', '界', '💰', '🐱',
}

func random(n int) string {
	slice := make([]rune, n)
	for i := 0; i < n/2; i++ {
		slice[i] = elements[rand.Intn(len(elements))]
	}
	for j := n / 2; j < n; j++ {
		slice[j] = elements[len(elements)-6:][rand.Intn(len(elements[len(elements)-6:]))]
	}
	return string(elements)
}
