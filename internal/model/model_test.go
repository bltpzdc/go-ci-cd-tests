package model


import "testing"
import "github.com/stretchr/testify/assert"

func TestEq(t *testing.T) {
	assert.Equal(t, 4, 2 + 2, "they should be equal")
}
