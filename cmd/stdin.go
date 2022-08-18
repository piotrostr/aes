/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// stdinCmd represents the stdin command
var stdinCmd = &cobra.Command{
	Use:   "stdin",
	Short: "Encrypt input from stdin",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("stdin called")
	},
}

func init() {
	rootCmd.AddCommand(stdinCmd)
}
