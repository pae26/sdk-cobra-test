package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/spf13/cobra"
)

var secretFlags struct {
	apply    bool
	env      string
	region   string
	secret   string
	filename string
}

// secretCmd represents the secret command
var secretCmd = &cobra.Command{
	Use:     "aws-secrets-manager",
	Aliases: []string{"aws"},
	Short:   "Update SecretsManager token.",
	Long: `Update SecretsManager token.
You must set [-e, -f, -s] options.`,
	Example: "  update-secrets aws-secrets-manager -e dev -f testtoken.yaml -s arn:aws:secretsmanager:ap-northeast-1:xxx",
	Run: func(cmd *cobra.Command, args []string) {
		if secretFlags.env == "" || secretFlags.filename == "" || secretFlags.secret == "" {
			fmt.Println("ERROR: You must set [-e, -f, -s] options.")
			fmt.Println("Show help with [-h] option.")
			os.Exit(0)
		}

		if secretFlags.env != "dev" && secretFlags.env != "stg" && secretFlags.env != "prd" {
			fmt.Println("ERROR: available env name is [dev, stg, prd]")
			os.Exit(0)
		}

		updateSecretsManager(secretFlags.apply, secretFlags.env, secretFlags.region, secretFlags.secret, secretFlags.filename)

	},
}

func updateSecretsManager(apply bool, env string, region string, secret string, filename string) {
	profile := os.Getenv("AWS_PROFILE")
	if !(strings.Contains(profile, env)) {
		fmt.Println("ERROR: env name inconsistent with AWS profile")
		os.Exit(0)
	}

	sess := session.Must(session.NewSessionWithOptions(session.Options{Profile: profile}))
	svc := secretsmanager.New(
		sess,
		aws.NewConfig().WithRegion(region),
	)

	input := &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(secret),
	}
	secretValue, err := svc.GetSecretValue(input)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(0)
	}
	arn := aws.StringValue(secretValue.ARN)

	tokenText, err := os.ReadFile(filename)
	if err != nil {
		fmt.Println(err)
	}
	secretString := string(tokenText)

	if apply {
		input := &secretsmanager.UpdateSecretInput{
			SecretId:     aws.String(arn),
			SecretString: aws.String(secretString),
		}

		_, err := svc.UpdateSecret(input)
		if err != nil {
			/* error handling refer to official sample code. */
			/* https://docs.aws.amazon.com/sdk-for-go/api/service/secretsmanager/#example_SecretsManager_UpdateSecret_shared00 */
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
				case secretsmanager.ErrCodeMalformedPolicyDocumentException:
					fmt.Println(secretsmanager.ErrCodeMalformedPolicyDocumentException, aerr.Error())
				case secretsmanager.ErrCodeInternalServiceError:
					fmt.Println(secretsmanager.ErrCodeInternalServiceError, aerr.Error())
				case secretsmanager.ErrCodePreconditionNotMetException:
					fmt.Println(secretsmanager.ErrCodePreconditionNotMetException, aerr.Error())
				default:
					fmt.Println(aerr.Error())
				}
			} else {
				fmt.Println(err.Error())
			}
			return
		}
		fmt.Println("SecretsManager updated.")
	} else {
		fmt.Println("[aws-secrets-manager]")
		fmt.Println("DRY-RUN finished. Use -a option to apply.")
		fmt.Printf("%-11s: %s\n", "env", env)
		fmt.Printf("%-11s: %s\n", "file path", filename) //TODO: add option display or hidden file contents.
		fmt.Printf("%-11s: %s\n", "region", region)
		fmt.Printf("%-11s: %s\n", "secret name", secret) //TODO: add option display or hidden secret ARN.
		fmt.Println()
	}
}

func init() {
	rootCmd.AddCommand(secretCmd)

	secretCmd.Flags().BoolVarP(&secretFlags.apply, "apply", "a", false, "default: dry-run")
	secretCmd.Flags().StringVarP(&secretFlags.env, "env", "e", "", "environment")
	secretCmd.Flags().StringVarP(&secretFlags.region, "region", "r", "ap-northeast-1", "aws region")
	secretCmd.Flags().StringVarP(&secretFlags.secret, "secret", "s", "", "secret name")
	secretCmd.Flags().StringVarP(&secretFlags.filename, "file", "f", "", "file path defined token information")
}
