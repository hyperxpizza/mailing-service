package utils

import "github.com/google/uuid"

func GenerateOneTimeToken() string {
	var token uuid.UUID
	token = uuid.New()

	return token.String()
}
