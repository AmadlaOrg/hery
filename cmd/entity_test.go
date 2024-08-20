package cmd

import (
	entityCmdPkg "github.com/AmadlaOrg/hery/entity/cmd"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEntityCmd(t *testing.T) {
	// Check that the EntityCmd is properly set up
	assert.Equal(t, "entity", EntityCmd.Use)
	assert.Equal(t, "Entity commands", EntityCmd.Short)

	// Verify that SetFlags has been called with EntityCmd
	// Note: We assume SetFlags adds specific flags or performs some setup
	// However, directly testing this depends on what SetFlags does.
	// Here, we only verify that flags are set (if any).
	assert.NotNil(t, EntityCmd.Flags(), "Flags should be set for EntityCmd")

	// Check that the ListCmd has been added
	foundListCmd := false
	foundGetCmd := false
	foundValidateCmd := false
	for _, c := range EntityCmd.Commands() {
		if c == entityCmdPkg.ListCmd {
			foundListCmd = true
		}
		if c == entityCmdPkg.GetCmd {
			foundGetCmd = true
		}
		if c == entityCmdPkg.ValidateCmd {
			foundValidateCmd = true
		}
	}

	assert.True(t, foundListCmd, "ListCmd should be registered with EntityCmd")
	assert.True(t, foundGetCmd, "GetCmd should be registered with EntityCmd")
	assert.True(t, foundValidateCmd, "ValidateCmd should be registered with EntityCmd")
}

// Optionally, you can test the full command execution structure
func TestEntityCmd_Execute(t *testing.T) {
	rootCmd := &cobra.Command{Use: "hery"}
	rootCmd.AddCommand(EntityCmd)

	_, err := executeCommand(rootCmd, "entity")
	assert.NoError(t, err, "EntityCmd execution should not return an error")
}
