package utils

import (
	"strings"

	"github.com/lithammer/shortuuid/v4"
)

// MakeUUID
func InUuid() string {
	return MakeUUID("IN")
}

// GoodsUuid
func GoodsUuid() string {
	return MakeUUID("GOODS")
}

// MakeUUID
func OutUuid() string {
	return MakeUUID("OUT")
}
func DeviceUuid() string {
	return MakeUUID("DEVICE")
}
func PluginUuid() string {
	return MakeUUID("PLUGIN")
}
func VisualUuid() string {
	return MakeUUID("VISUAL")
}
func GroupUuid() string {
	return MakeUUID("GROUP")
}
func AppUuid() string {
	return MakeUUID("APP")
}
func AiBaseUuid() string {
	return MakeUUID("AIBASE")
}
func DataSchemaUuid() string {
	return MakeUUID("SCHEMA")
}
func CronTaskUuid() string {
	return MakeUUID("CRONTASK")
}

// MakeUUID
func RuleUuid() string {
	return MakeUUID("RULE")
}

// MakeUUID
func MakeUUID(prefix string) string {
	return prefix + strings.ToUpper(shortuuid.New()[:6])
}
func MakeLongUUID(prefix string) string {
	return prefix + strings.ToUpper(shortuuid.New())
}
