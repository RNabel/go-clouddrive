package files

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestCloudFileSerialisation(t *testing.T) {
	cf := CloudFile{"TestID", "TestName", []string{"ID1", "ID2"}}
	serialised := cf.ToBytes()

	deserialised := NewCloudFileFromByte(serialised)
	assert.Equal(t, deserialised.GoogleId, cf.GoogleId, "Titles are identical.")
	assert.Equal(t, deserialised.Name, cf.Name, "Name are identical.")
	assert.Equal(t, deserialised.Parents, cf.Parents, "Name are identical.")
}