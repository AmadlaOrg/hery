package e2e

import (
	"os"
	"os/exec"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("CLI", func() {
	Context("help command", func() {
		It("should print the help message", func() {
			// Construct the absolute path to main.go
			cwd, err := os.Getwd()
			Expect(err).ToNot(HaveOccurred())
			cmdPath := filepath.Join(cwd, "../../main.go")

			// Run the command
			cmd := exec.Command("go", "run", cmdPath, "help")
			output, err := cmd.CombinedOutput()
			if err != nil {
				Fail(string(output)) // Provide more context on failure
			}
			Expect(err).ToNot(HaveOccurred(), string(output))

			expectedOutput := `HERY CLI application

Usage:
  hery [command]

Available Commands:
  client      HERY client
  collection  Collections
  completion  Generate the autocompletion script for the specified shell
  compose     Compose the specified entity
  entity      Entity commands
  help        Help about any command
  query       Query entities
  server      HERY Server
  settings    List the paths and other environment variables for HERY
  version     Print the version number of hery

Flags:
  -h, --help      help for hery
  -v, --version   version for hery

Use "hery [command] --help" for more information about a command.`

			Expect(string(output)).To(ContainSubstring(expectedOutput))
		})
	})
})
