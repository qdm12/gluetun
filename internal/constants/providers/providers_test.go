package providers

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

func Test_All(t *testing.T) {
	t.Parallel()

	all := All()
	assert.NotContains(t, all, Custom)
	assert.NotEmpty(t, all)
}

func Test_AllWithCustom(t *testing.T) {
	t.Parallel()

	all := AllWithCustom()
	assert.Contains(t, all, Custom)
	assert.Len(t, all, len(All())+1)
}

func TestWorkflowHasAll(t *testing.T) {
	t.Parallel()

	const path = "../../../.github/workflows/update-servers-list.yml"
	file, err := os.Open(path)
	require.NoError(t, err)
	defer file.Close()

	var data struct {
		On struct {
			WorkflowDispatch struct {
				Inputs struct {
					Provider struct {
						Options []string `yaml:"options"`
					} `yaml:"provider"`
				} `yaml:"inputs"`
			} `yaml:"workflow_dispatch"`
		} `yaml:"on"`
	}
	decoder := yaml.NewDecoder(file)
	err = decoder.Decode(&data)
	require.NoError(t, err)

	providers := All()
	expected := make([]string, len(providers)+1)
	expected[0] = "all"
	copy(expected[1:], providers)
	assert.Equal(t, expected, data.On.WorkflowDispatch.Inputs.Provider.Options)
}
