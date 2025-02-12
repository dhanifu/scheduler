package config

import "github.com/google/uuid"

type KeyConstant string

const (
	RequestIDKey KeyConstant = "requestid"
)

func GenerateRequestID() string {
	return uuid.New().String()
}
