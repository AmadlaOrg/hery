package util

import (
	"github.com/AmadlaOrg/hery/storage"
	"github.com/spf13/cobra"
	"log"
)

type Util interface {
	Concoct(
		cmd *cobra.Command,
		args []string,
		handler func(paths *storage.AbsPaths, args []string)) error
}

type utilImpl struct {
	New storage.Storage
}

// Concoct sets up the necessary storage paths and executes the provided handler function.
func (s *utilImpl) Concoct(
	cmd *cobra.Command,
	args []string,
	handler func(paths *storage.AbsPaths, args []string)) error {
	paths, err := s.New.Paths()
	if err != nil {
		log.Println("Error getting paths:", err)
		return err
	}
	handler(paths, args)
	return nil
}
