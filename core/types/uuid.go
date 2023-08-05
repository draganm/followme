package types

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
)

type UUID [16]byte

func (u UUID) String() string {
	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x", u[:4], u[4:6], u[6:8], u[8:10], u[10:16])
}

func (u UUID) MarshalJSON() ([]byte, error) {
	return json.Marshal(u.String())
}

var uuidRegexp = regexp.MustCompile(`^(?P<first>[[:xdigit:]]{8})-(?P<second>[[:xdigit:]]{4})-(?P<third>[[:xdigit:]]{4})-(?P<fourth>[[:xdigit:]]{4})-(?P<fifth>[[:xdigit:]]{12})$`)

var ErrInvalidUUIDString = errors.New("invalid UUID string")

// ParseUUID parses a UUID from a given string.
// The function expects the input string in the format "123e4567-e89b-12d3-a456-426614174000",
// where each group of characters (separated by hyphens) represents hexadecimal values.
//
// The function returns a UUID object that represents the parsed UUID and an error.
// If the input string does not represent a valid UUID, the function returns an error.
//
// Example usage:
//
//	uuid, err := ParseUUID("123e4567-e89b-12d3-a456-426614174000")
//	if err != nil {
//	    log.Fatalf("failed to parse UUID: %v", err)
//	}
//
// Parameters:
//
//	s (string): The input string representing the UUID.
//
// Returns:
//
//	UUID: The UUID object that represents the parsed UUID.
//	error: An error object that describes the error, if any occurred. nil otherwise.
func ParseUUID(s string) (UUID, error) {

	matches := uuidRegexp.FindStringSubmatch(s)

	if matches == nil {
		return UUID{}, ErrInvalidUUIDString
	}

	var uuidBytes UUID
	pos := 0
	for i, match := range matches[1:] {
		hexBytes, err := hex.DecodeString(match)
		if err != nil {
			return UUID{}, fmt.Errorf("%w: could not parse part %d: %w", ErrInvalidUUIDString, i, err)
		}
		copy(uuidBytes[pos:], hexBytes)
		pos += len(hexBytes)
	}

	return uuidBytes, nil

}

func (u UUID) Compare(o UUID) int {
	return bytes.Compare(u[:], o[:])
}

func UUIDFromBytes(b []byte) (UUID, error) {
	if len(b) < 16 {
		return UUID{}, errors.New("byte slice too short for UUID")
	}

	id := UUID{}
	copy(id[:], b)
	return id, nil
}

var ZeroUUID = UUID{}
