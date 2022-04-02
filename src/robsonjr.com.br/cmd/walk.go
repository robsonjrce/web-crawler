package cmd

import (
	"errors"
	"github.com/spf13/cobra"
)

var walkCmd = &cobra.Command{
	Use:   "walk",
	Short: "Walk through the root url",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		url, err := cmd.Flags().GetString("url")
		if err != nil {
			return err
		}
		if url == "" {
			return errors.New("url can't be null")
		}

		dest, err := cmd.Flags().GetString("dest")
		if err != nil {
			return err
		}
		if dest == "" {
			return errors.New("url can't be null")
		}

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
}

func init() {
	walkCmd.Flags().String("url", "", "the url to be crawled")
	walkCmd.MarkFlagRequired("url")
	walkCmd.Flags().String("dest", "", "the destination path to store the data")
	walkCmd.MarkFlagRequired("dest")
}