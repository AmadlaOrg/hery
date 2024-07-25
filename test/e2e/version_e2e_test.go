package e2e

import (
	"os"
	"os/exec"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("CLI", func() {
	Context("version command", func() {
		It("should print the version", func() {
			// Construct the absolute path to main.go
			cwd, err := os.Getwd()
			Expect(err).ToNot(HaveOccurred())
			cmdPath := filepath.Join(cwd, "../../main.go")

			// Run the command
			cmd := exec.Command("go", "run", cmdPath, "version")
			output, err := cmd.CombinedOutput()
			if err != nil {
				Fail(string(output)) // Provide more context on failure
			}
			Expect(err).ToNot(HaveOccurred(), string(output))
			Expect(string(output)).To(ContainSubstring("hery version 1.0.0"))
		})
	})
})
