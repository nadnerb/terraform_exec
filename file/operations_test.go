package file

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestExpectsToHaveFilesWithExtension(t *testing.T) {
	hasTxtFiles, err := DirectoryContainsWithExtension("../fixtures", ".txt")

	assert.Nil(t, err)
	assert.True(t, hasTxtFiles, "should contain text files")
}

func TestExpectsToNotHaveFilesWithExtension(t *testing.T) {
	hasFooFiles, err := DirectoryContainsWithExtension("../fixtures", ".foo")

	assert.Nil(t, err)
	assert.False(t, hasFooFiles, "should not contain foo files")
}
