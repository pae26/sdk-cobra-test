package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

var onepasswordFlags struct {
	apply    bool
	create   bool
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

		updateOnepassword(onepasswordFlags.apply, onepasswordFlags.create, onepasswordFlags.vault, onepasswordFlags.title, onepasswordFlags.filename)
	},
}

func updateOnepassword(apply bool, create bool, vault string, title string, filename string) {
	if apply {
		if create {
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
			if title == "" {
				fmt.Println("ERROR: Set title of item with [-t] option.")
				os.Exit(0)
			}

			vault_ary := strings.Split(vault, ",")
			for _, v := range vault_ary {
				output, err := exec.Command("op", "edit", "document", title, filename, "--vault", v).CombinedOutput()
				if err != nil {
					fmt.Println(string(output))
				}
			}
		}

	} else {
		fmt.Println("DRY-RUN finished. Use -a option to apply.")
		if create {
			fmt.Printf("%-11s: %s\n", "operation", "create")
		} else {
			fmt.Printf("%-11s: %s\n", "operation", "edit")
		}
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
	onepasswordCmd.Flags().BoolVarP(&onepasswordFlags.create, "create", "c", false, "default: edit")
	onepasswordCmd.Flags().StringVarP(&onepasswordFlags.vault, "vault", "v", "", "vault name")
	onepasswordCmd.Flags().StringVarP(&onepasswordFlags.title, "title", "t", "", "title of item")
	onepasswordCmd.Flags().StringVarP(&onepasswordFlags.filename, "file", "f", "", "file path defined token information")
}
