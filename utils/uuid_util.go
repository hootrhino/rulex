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
// GoodsUuid
//
func GoodsUuid() string {
	return MakeUUID("GOODS")
}

//
// MakeUUID
//
func OutUuid() string {
	return MakeUUID("OUT")
}
func DeviceUuid() string {
	return MakeUUID("DEVICE")
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
