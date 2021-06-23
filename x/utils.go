package x

import (
	"github.com/google/uuid"
)

// XXX_${uuid}
func MakeUUID(prefix string) string {
	return prefix + "_" + uuid.NewString()
}
