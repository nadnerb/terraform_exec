package file

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func ExpectsToHaveFilesWithExtension(t *testing.T) {
	has, err := HasFilesWithExtension("../fixtures/some.tfvars", "txt")

	assert.Nil(t, err)
	assert.Equal(t, has, true, "should contain text files")
}

func ExpectsToNotHaveFilesWithExtension(t *testing.T) {
	has, err := HasFilesWithExtension("../fixtures/some.tfvars", "foo")

	assert.Nil(t, err)
	assert.Equal(t, has, false, "should contain foo files")
}
