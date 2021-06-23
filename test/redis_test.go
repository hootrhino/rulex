package test

import (
	"rulenginex/x"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/google/uuid"
)

//
func TestSaveGet(t *testing.T) {
	id := uuid.NewString()
	in1 := x.InEnd{
		Id:          id,
		Type:        "UDP",
		Name:        "UDP Stream",
		Description: "UDP Input Stream",
		Config: &map[string]interface{}{
			"token":         "token",
			"packet_length": 1024,
		},
	}
	x.SaveInEnd(&in1)
	inEnd := x.GetInEnd(id)
	assert.Equal(t, in1.Id, id)
	assert.Equal(t, in1.Type, "UDP")
	assert.Equal(t, in1.Name, "UDP Stream")
	assert.Equal(t, in1.Description, "UDP Input Stream")
	assert.Equal(t, (*in1.Config)["token"], "token")
	assert.Equal(t, (*in1.Config)["packet_length"], 1024)
	t.Log(inEnd)
}
