package util

import (
	"github.com/AmadlaOrg/hery/storage"
	"github.com/spf13/cobra"
	"log"
)

type IUtil interface {
	Concoct(
		cmd *cobra.Command,
		args []string,
		handler func(paths *storage.AbsPaths, args []string)) error
}

type SUtil struct {
	NewStorageService storage.IStorage
}

// Concoct sets up the necessary storage paths and executes the provided handler function.
func (s *SUtil) Concoct(
	cmd *cobra.Command,
	args []string,
	handler func(paths *storage.AbsPaths, args []string)) error {
	paths, err := s.NewStorageService.Paths()
	if err != nil {
		log.Println("Error getting paths:", err)
		return err
	}
	handler(paths, args)
	return nil
}
