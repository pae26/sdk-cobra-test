package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var updateFlags struct {
	apply     bool
	operation string
	vault     string
	title     string
	env       string
	region    string
	secret    string
	filename  string
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = NewRootCmd()

func NewRootCmd() *cobra.Command {
	var rootCmd = &cobra.Command{
		Use:   "update-secrets",
		Short: "Update 1password and SecretsManager token.",
		Long: `Update 1password and SecretsManager token.
	Show help with [-h] option.`,
		Run: func(cmd *cobra.Command, args []string) {
			if updateFlags.vault == "" || updateFlags.filename == "" || updateFlags.env == "" || updateFlags.secret == "" {
				fmt.Println("ERROR: You must set [-f, -v, -t, -e, -s] options.")
				fmt.Println("Show help with [-h] option.")
				os.Exit(0)
			}

			if updateFlags.env != "dev" && updateFlags.env != "stg" && updateFlags.env != "prd" {
				fmt.Println("ERROR: available env name is [dev, stg, prd]")
				os.Exit(0)
			}

			switch updateFlags.operation {
			case "edit":
				editItem(updateFlags.apply, updateFlags.vault, updateFlags.title, updateFlags.filename)
			case "create":
				createItem(updateFlags.apply, updateFlags.vault, updateFlags.title, updateFlags.filename)
			}

			updateSecretsManager(updateFlags.apply, updateFlags.env, updateFlags.region, updateFlags.secret, updateFlags.filename)
		},
	}

	return rootCmd
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolVarP(&updateFlags.apply, "apply", "a", false, "default: dry-run")
	rootCmd.Flags().StringVarP(&updateFlags.operation, "operation", "o", "edit", "[1password]select edit or create. default: edit")
	rootCmd.Flags().StringVarP(&updateFlags.vault, "vault", "v", "", "[1password]vault name")
	rootCmd.Flags().StringVarP(&updateFlags.title, "title", "t", "", "[1password]title of item")
	rootCmd.Flags().StringVarP(&updateFlags.env, "env", "e", "", "[SecretsManager]environment")
	rootCmd.Flags().StringVarP(&updateFlags.region, "region", "r", "ap-northeast-1", "aws region")
	rootCmd.Flags().StringVarP(&updateFlags.secret, "secret", "s", "", "[SecretsManager]secret name")
	rootCmd.Flags().StringVarP(&updateFlags.filename, "file", "f", "", "file path defined token information")
}
