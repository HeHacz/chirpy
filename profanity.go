package main

import "strings"

var ProfaneWords = map[string]bool{
	"kerfuffle": true,
	"sharbert":  true,
	"fornax":    true,
}

func censorMessage(payload string) string {
	censoredMessage := strings.Split(payload, " ")
	for i, word := range censoredMessage {
		_, exist := ProfaneWords[strings.ToLower(word)]
		if exist {
			censoredMessage[i] = "****"
		}
	}
	return strings.Join(censoredMessage, " ")
}
