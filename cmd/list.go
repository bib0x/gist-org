package cmd

import (
	"fmt"
	"os"
	"sync"

	"github.com/bib0x/gist-org/internal/tag"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var listCmd = &cobra.Command{
	Use: "list",
	Short: "list different resources",
}

var tagListCmd = &cobra.Command{
	Use: "tags",
	Short: "list all tags from Org files",
	Run: func(cmd *cobra.Command, args []string) {
		path := viper.GetString("path")

		if path == "" {
			msg := fmt.Errorf("Undefined Org files path")
			cobra.CheckErr(msg)
		}

		files, err := os.ReadDir(path)
		if err != nil {
			msg := fmt.Errorf("Error gathering Org Files.\n error: %s\n", err)
			cobra.CheckErr(msg)
		}

		results := make(chan string)
		var wg sync.WaitGroup

		for _, file := range files {
			wg.Add(1)
			go tag.GetTags(path + file.Name(), results, &wg)
		}

		go func() {
			wg.Wait()
			close(results)
		}()

		tag.PrintUnsortedTags(results)
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.AddCommand(tagListCmd)
}
