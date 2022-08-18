package cmd

import (
	"io/ioutil"

	"github.com/piotrostr/aes/crypto"
	"github.com/spf13/cobra"
)

// decryptCmd represents the decrypt command
var decryptCmd = &cobra.Command{
	Use:   "decrypt",
	Short: "A brief description of your command",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			cmd.Help()
			return
		}
		path := args[0]
		ciphertext := crypto.GetFileContents(path)
		plaintext := crypto.Decrypt(ciphertext)
		outpath := path + ".dec"
		ioutil.WriteFile(outpath, plaintext, 0o644)
	},
}

func init() {
	rootCmd.AddCommand(decryptCmd)
}
