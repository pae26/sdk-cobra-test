package cmd

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

var onepasswordFlags struct {
	apply     bool
	operation string
	vault     string
	title     string
	filename  string
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

		switch onepasswordFlags.operation {
		case "edit":
			err := editItem(onepasswordFlags.apply, onepasswordFlags.vault, onepasswordFlags.title, onepasswordFlags.filename)
			if err != nil {
				log.Fatalln(err)
			}
		case "create":
			createItem(onepasswordFlags.apply, onepasswordFlags.vault, onepasswordFlags.title, onepasswordFlags.filename)
		}
		fmt.Println("1password updated.")
	},
}

func editItem(apply bool, vault string, title string, filename string) error {
	if title == "" {
		return fmt.Errorf("ERROR: Set title of item with [-t] option.")
	}

	if apply {
		vault_ary := strings.Split(vault, ",")
		for _, v := range vault_ary {
			output, err := exec.Command("op", "edit", "document", title, filename, "--vault", v).CombinedOutput()
			if err != nil {
				return fmt.Errorf(string(output))
			}
		}
	} else {
		dryRunOnepassword("edit", vault, title, filename)
	}
	return nil
}

func createItem(apply bool, vault string, title string, filename string) {
	var cmd *exec.Cmd
	if title == "" {
		title = "(not specify)"
		cmd = exec.Command("op", "create", "document", filename, "--vault", vault)
	} else {
		cmd = exec.Command("op", "create", "document", filename, "--vault", vault, "--title", title)
	}

	if apply {
		output, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Println(string(output))
			os.Exit(0)
		}
	} else {
		dryRunOnepassword("create", vault, title, filename)
	}
}

func dryRunOnepassword(operation string, vault string, title string, filename string) {
	fmt.Println("[onepassword]")
	fmt.Println("DRY-RUN finished. Use -a option to apply.")
	fmt.Printf("%-11s: %s\n", "operation", operation)
	fmt.Printf("%-11s: %s\n", "file path", filename)
	fmt.Printf("%-11s: %s\n", "vault", vault)
	fmt.Printf("%-11s: %s\n", "title", title)
	fmt.Println()
}

func init() {
	rootCmd.AddCommand(onepasswordCmd)

	onepasswordCmd.Flags().BoolVarP(&onepasswordFlags.apply, "apply", "a", false, "default: dry-run")
	onepasswordCmd.Flags().StringVarP(&onepasswordFlags.operation, "operation", "o", "edit", "select edit or create. default: edit")
	onepasswordCmd.Flags().StringVarP(&onepasswordFlags.vault, "vault", "v", "", "vault name")
	onepasswordCmd.Flags().StringVarP(&onepasswordFlags.title, "title", "t", "", "title of item")
	onepasswordCmd.Flags().StringVarP(&onepasswordFlags.filename, "file", "f", "", "file path defined token information")
}
