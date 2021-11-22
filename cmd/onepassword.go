package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

var onepasswordFlags struct {
	apply    bool
	vault    string
	title    string
	filename string
}

// onepasswordCmd represents the onepassword command
var onepasswordCmd = &cobra.Command{
	Use:     "1password",
	Aliases: []string{"op"},
	Short:   "Update 1password item",
	Long: `Update 1password item.
Set [-v, -t, -f] options.`,
	Run: func(cmd *cobra.Command, args []string) {
		if onepasswordFlags.vault == "" || onepasswordFlags.filename == "" {
			fmt.Println("ERROR: You must set [-v, -t, -f] options.")
			fmt.Println("Show help with [-h] option.")
			os.Exit(0)
		}

		updateOnepassword(onepasswordFlags.apply, onepasswordFlags.vault, onepasswordFlags.title, onepasswordFlags.filename)
	},
}

func updateOnepassword(apply bool, vault string, title string, filename string) {
	if apply {
		if title != "" {
			output, err := exec.Command("op", "create", "document", filename, "--vault", vault, "--title", title).CombinedOutput()
			if err != nil {
				fmt.Println(string(output))
			}
		} else {
			output, err := exec.Command("op", "create", "document", filename, "--vault", vault).CombinedOutput()
			if err != nil {
				fmt.Println(string(output))
			}
		}
	} else {
		fmt.Println("DRY-RUN finished. Use -a option to apply.")
		fmt.Printf("%-11s: %s\n", "vault", vault)
		if title != "" {
			fmt.Printf("%-11s: %s\n", "title", title)
		} else {
			fmt.Printf("%-11s: %s\n", "title", "(not specify)")
		}
		fmt.Printf("%-11s: %s\n", "file path", filename)
	}
}

func init() {
	rootCmd.AddCommand(onepasswordCmd)

	onepasswordCmd.Flags().BoolVarP(&onepasswordFlags.apply, "apply", "a", false, "default: dry-run")
	onepasswordCmd.Flags().StringVarP(&onepasswordFlags.vault, "vault", "v", "", "vault name")
	onepasswordCmd.Flags().StringVarP(&onepasswordFlags.title, "title", "t", "", "title of item")
	onepasswordCmd.Flags().StringVarP(&onepasswordFlags.filename, "file", "f", "", "file path defined token information")
}
