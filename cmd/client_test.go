package cmd

/*import (
	clientCmdPkg "github.com/AmadlaOrg/hery/client/cmd"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestClientCmd(t *testing.T) {
	// Check that the root command is properly set up
	assert.Equal(t, "client", ClientCmd.Use)
	assert.Equal(t, "HERY client", ClientCmd.Short)

	// Check that the ConnectToServerCmd has been added
	found := false
	for _, c := range ClientCmd.Commands() {
		if c == clientCmdPkg.ConnectToServerCmd {
			found = true
			break
		}
	}
	assert.True(t, found, "ConnectToServerCmd should be registered with ClientCmd")
}

// Optionally, you can test the full command execution structure
func TestClientCmd_Execute(t *testing.T) {
	rootCmd := &cobra.Command{Use: "hery"}
	rootCmd.AddCommand(ClientCmd)

	_, err := executeCommand(rootCmd, "client")
	assert.NoError(t, err, "ClientCmd execution should not return an error")
}

// Utility function to execute a command for testing
func executeCommand(root *cobra.Command, args ...string) (output string, err error) {
	_, err = root.ExecuteC()
	return
}*/
