package datatypes

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	t.Run("read-file", func(t *testing.T) {
		data, err := LoadConfig("testdata/config_example.json")
		assert.NoError(t, err)
		assert.NotNil(t, data)
	})
}
