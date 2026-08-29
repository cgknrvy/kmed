package icd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSearchICD10(t *testing.T) {
	t.Run("searching for malaria", func(t *testing.T) {
		codes, err := SearchICD11("malaria")
		assert.Nil(t, err)
		assert.True(t, len(codes) > 0)
	})

	t.Run("searching invalid disease name", func(t *testing.T) {
		codes, err := SearchICD11("invalid-disease-name")
		assert.Nil(t, err)
		assert.True(t, len(codes) == 0)
	})
}
