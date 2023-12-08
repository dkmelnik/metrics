package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGenerateGuid(t *testing.T) {
	guid := GenerateGuid()

	assert.NotEmpty(t, guid, "Generated GUID should not be empty")
}
