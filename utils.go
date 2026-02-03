package main

import "strings"

func removeProfane(p string) string {
	words1 := strings.Split(p, " ")

	for i, word := range words1 {
		lowerWord := strings.ToLower(word)
		if lowerWord == "kerfuffle" || lowerWord == "sharbert" || lowerWord == "fornax" {
			words1[i] = "****"
		}

	}

	return strings.Join(words1, " ")
}
