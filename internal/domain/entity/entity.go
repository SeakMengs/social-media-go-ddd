package entity

import (
	"github.com/google/uuid"
)

func StringToUUID(s string) (uuid.UUID, error) {
	return uuid.Parse(s)
}
