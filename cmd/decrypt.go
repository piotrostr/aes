package cmd

import (
	"io/ioutil"

	"github.com/piotrostr/aes/crypto"
	"github.com/spf13/cobra"
)

// decryptCmd represents the decrypt command
var decryptCmd = &cobra.Command{
	Use:   "decrypt [file]",
	Short: "Decrypts a file using AES-256-GCM",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			cmd.Help()
			return
		}
		path := args[0]

		gcm := crypto.GCM{}
		gcm.Initialize()

		ciphertext := crypto.GetFileContents(path)
		plaintext := gcm.Decrypt(ciphertext)

		var outpath string
		if path[len(path)-4:] == ".enc" {
			outpath = path[:len(path)-4]
		} else {
			outpath = path + ".dec"
		}

		ioutil.WriteFile(outpath, plaintext, 0o644)
	},
}

func init() {
	rootCmd.AddCommand(decryptCmd)
}
