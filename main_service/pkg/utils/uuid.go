package utils
import (
	"github.com/google/uuid"
)

// GenerateUUID returns a new randomly generated UUID.
func GenerateUUID() string {
	return uuid.New().String()
}
