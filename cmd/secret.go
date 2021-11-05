package cmd

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/spf13/cobra"
)

var secretFlags struct {
	env      string
	profile  string
	filename string
	secret   string
}

// secretCmd represents the secret command
var secretCmd = &cobra.Command{
	Use:   "secret",
	Short: "Update SecretsManager token.",
	Long: `Update SecretsManager token.
You must set [-e, -p, -f, -s] options.`,
	Example: "  update secret -e dev -p zozo-replace-dev-powerUser -f testtoken.yaml -s arn:aws:secretsmanager:ap-northeast-1:xxx",
	Run: func(cmd *cobra.Command, args []string) {
		updateSecretsManager(secretFlags.profile, secretFlags.filename, secretFlags.secret)
	},
}

func updateSecretsManager(profile string, filename string, secret string) {
	if secretFlags.env == "" || secretFlags.profile == "" || secretFlags.filename == "" || secretFlags.secret == "" {
		fmt.Println("You must set [-e, -p, -f, -s] options.")
		fmt.Println("Show help with [-h] option.")
		os.Exit(0)
	}

	if secretFlags.env != "dev" && secretFlags.env != "stg" && secretFlags.env != "prd" {
		fmt.Println("available env names: dev, stg, prd")
		os.Exit(0)
	}

	region := "ap-northeast-1"

	sess := session.Must(session.NewSessionWithOptions(session.Options{Profile: profile}))
	svc := secretsmanager.New(
		sess,
		aws.NewConfig().WithRegion(region),
	)

	tokenText, err := os.ReadFile(filename)
	if err != nil {
		fmt.Println(err)
	}
	secretString := string(tokenText)

	input := &secretsmanager.UpdateSecretInput{
		SecretId:     aws.String(secret),
		SecretString: aws.String(secretString),
	}

	result, err := svc.UpdateSecret(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case secretsmanager.ErrCodeInvalidParameterException:
				fmt.Println(secretsmanager.ErrCodeInvalidParameterException, aerr.Error())
			case secretsmanager.ErrCodeInvalidRequestException:
				fmt.Println(secretsmanager.ErrCodeInvalidRequestException, aerr.Error())
			case secretsmanager.ErrCodeLimitExceededException:
				fmt.Println(secretsmanager.ErrCodeLimitExceededException, aerr.Error())
			case secretsmanager.ErrCodeEncryptionFailure:
				fmt.Println(secretsmanager.ErrCodeEncryptionFailure, aerr.Error())
			case secretsmanager.ErrCodeResourceExistsException:
				fmt.Println(secretsmanager.ErrCodeResourceExistsException, aerr.Error())
			case secretsmanager.ErrCodeResourceNotFoundException:
				fmt.Println(secretsmanager.ErrCodeResourceNotFoundException, aerr.Error())
			case secretsmanager.ErrCodeInternalServiceError:
				fmt.Println(secretsmanager.ErrCodeInternalServiceError, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			fmt.Println(err.Error())
		}
		return
	}

	fmt.Println(result)
}

func init() {
	rootCmd.AddCommand(secretCmd)

	secretCmd.Flags().StringVarP(&secretFlags.env, "env", "e", "", "environment")
	secretCmd.Flags().StringVarP(&secretFlags.profile, "profile", "p", "", "AWS profile name")
	secretCmd.Flags().StringVarP(&secretFlags.filename, "file", "f", "", "file name defined token information")
	secretCmd.Flags().StringVarP(&secretFlags.secret, "secret", "s", "", "secret ARN")
}
