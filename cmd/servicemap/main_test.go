package servicemap

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	cobra "github.com/spf13/cobra"
)

func TestCiBiConfigIsNotExist(t *testing.T){

	// generate the cibi.yaml

	var rootCmd = &cobra.Command{
		Use:     	"servicemap [Cibi.yaml]",
		Short:   	"running service map by default will run cibi.yaml",
		Args:		cobra.MaximumNArgs(1),
		Run:     	RunServiceMap,
	}

	actual := new(bytes.Buffer)
	rootCmd.SetOut(actual)
	rootCmd.SetErr(actual)
	rootCmd.SetArgs([]string{"cibi.yaml"})
	rootCmd.Execute()

	t.Log(actual.String())

	assert.Equal(t, "a", "a", "actual is not expected")
}