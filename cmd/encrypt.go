/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"io/ioutil"

	"github.com/piotrostr/aes/crypto"
	"github.com/spf13/cobra"
)

// encryptCmd represents the encrypt command
var encryptCmd = &cobra.Command{
	Use:   "encrypt [path]",
	Short: "Encrypts a file using AES-256-GCM",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			cmd.Help()
			return
		}
		path := args[0]

		gcm := crypto.GCM{}
		gcm.Initialize()

		plaintext := crypto.GetFileContents(path)
		cipertext := gcm.Encrypt(plaintext)

		outpath := path + ".enc"
		ioutil.WriteFile(outpath, cipertext, 0o644)
	},
}

func init() {
	rootCmd.AddCommand(encryptCmd)
}
