package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTagSortNil(t *testing.T) {
	assert.NotPanics(t, func() {
		SortTags(nil)
	})
}

func TestTagSortSingle(t *testing.T) {
	tags := []string{"latest"}
	SortTags(&tags)
	assert.Equal(t, "latest", tags[0])
}

func TestTagSortMultiple(t *testing.T) {
	tags := []string{"1.0", "1.0.0-1", "1.0.0"}
	SortTags(&tags)
	assert.Equal(t, "1.0.0-1", tags[0])
}

func TestTagSortMultipleLatest(t *testing.T) {
	tags := []string{"1.0", "latest"} // latest is longest
	SortTags(&tags)
	assert.Equal(t, "1.0", tags[0])
}
