package files

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestCloudFileSerialisation(t *testing.T) {
	cf := CloudFile{"TestName"}
	serialised := cf.toBytes()

	deserialised := NewCloudFileFromByte(serialised)
	assert.Equal(t, deserialised.GoogleId, cf.GoogleId, "Titles are identical.")
}