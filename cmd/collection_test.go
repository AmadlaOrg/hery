package cmd

import (
	collectionPkgCmd "github.com/AmadlaOrg/hery/collection/cmd"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCollectionCmd(t *testing.T) {
	// Check that the root command is properly set up
	assert.Equal(t, "collection", CollectionCmd.Use)
	assert.Equal(t, "Collections", CollectionCmd.Short)

	// Check that the ListCmd has been added
	foundListCmd := false
	foundInitCmd := false
	for _, c := range CollectionCmd.Commands() {
		if c == collectionPkgCmd.ListCmd {
			foundListCmd = true
		}
		if c == collectionPkgCmd.InitCmd {
			foundInitCmd = true
		}
	}

	assert.True(t, foundListCmd, "ListCmd should be registered with CollectionCmd")
	assert.True(t, foundInitCmd, "InitCmd should be registered with CollectionCmd")
}

// Optionally, you can test the full command execution structure
func TestCollectionCmd_Execute(t *testing.T) {
	rootCmd := &cobra.Command{Use: "hery"}
	rootCmd.AddCommand(CollectionCmd)

	_, err := executeCommand(rootCmd, "collection")
	assert.NoError(t, err, "CollectionCmd execution should not return an error")
}
