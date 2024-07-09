package cmd

import (
	collectionPkgCmd "github.com/AmadlaOrg/hery/collection/cmd"
	"github.com/spf13/cobra"
)

var CollectionCmd = &cobra.Command{
	Use:   "collection",
	Short: "Collections",
	//Run: func(cmd *cobra.Command, args []string) {
	/*entityDir, err := storage.Path()
	if err != nil {
		fmt.Println("could not get the root storage directory:", err)
		return
	}

	// TODO: Maybe display that it is downloading
	err = entity.Get(entityDir)
	if err != nil {
		fmt.Println("Error crawling directories:", err)
		return
	}*/
	//},
}

func init() {
	CollectionCmd.AddCommand(collectionPkgCmd.ListCmd)
	CollectionCmd.AddCommand(collectionPkgCmd.InitCmd)
}
