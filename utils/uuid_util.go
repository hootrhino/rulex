package utils

import (
	"strings"

	"github.com/google/uuid"
)

//
// MakeUUID
//
func InUuid() string {
	return MakeUUID("IN")
}

//
// MakeUUID
//
func OutUuid() string {
	return MakeUUID("OUT")
}

//
// MakeUUID
//
func RuleUuid() string {
	return MakeUUID("RULE")
}

//
// MakeUUID
//
func MakeUUID(prefix string) string {
	return prefix + ":" + strings.Replace(uuid.NewString(), "-", "", -1)
}
