package utils

import "github.com/google/uuid"

//
// MakeUUID
//
func InUuid() string {
	return MakeUUID("INEND")
}

//
// MakeUUID
//
func OutUuid() string {
	return MakeUUID("OUTEND")
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
	return prefix + "_" + uuid.NewString()
}
