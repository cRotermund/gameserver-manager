package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/cRotermund/gameserver-manager/harness/internal/harness"
	"github.com/spf13/cobra"
)

var (
	secretsEnvFile string
	secretsRegion  string
	secretsProfile string
	secretsDryRun  bool
	secretsOnly    []string
)

var secretsCmd = &cobra.Command{
	Use:   "secrets",
	Short: "Bootstrap and rotate production secrets in AWS SSM",
}

var bootstrapSecretsCmd = newSecretsSyncCommand("bootstrap", false)
var rotateSecretsCmd = newSecretsSyncCommand("rotate", true)

func newSecretsSyncCommand(action string, overwrite bool) *cobra.Command {
	command := &cobra.Command{
		Use:   action,
		Short: fmt.Sprintf("%s production secrets in AWS SSM", action),
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			root, err := harness.ProjectRoot()
			if err != nil {
				return err
			}

			envFile := secretsEnvFile
			if !filepath.IsAbs(envFile) {
				envFile = filepath.Join(root, envFile)
			}
			parameters, err := harness.LoadProductionSecrets(envFile, secretsOnly)
			if err != nil {
				return err
			}

			fmt.Printf("--- %s production secrets from %s ---\n", action, envFile)
			return harness.PutSSMParameters(parameters, secretsRegion, secretsProfile, overwrite, secretsDryRun)
		},
	}
	return command
}

func init() {
	secretsEnvFile = harness.DefaultSecretsEnvFile()
	secretsRegion = harness.DefaultAWSRegion()

	for _, command := range []*cobra.Command{bootstrapSecretsCmd, rotateSecretsCmd} {
		command.Flags().StringVar(&secretsEnvFile, "env-file", secretsEnvFile, "Path to the gitignored secrets env file")
		command.Flags().StringVar(&secretsRegion, "region", secretsRegion, "AWS region")
		command.Flags().StringVar(&secretsProfile, "profile", "", "AWS CLI profile")
		command.Flags().BoolVar(&secretsDryRun, "dry-run", false, "Validate the env file and print parameter names without writing")
		command.Flags().StringSliceVar(&secretsOnly, "only", nil, "Only sync these secret variables (repeatable or comma-separated)")
		secretsCmd.AddCommand(command)
	}
	rootCmd.AddCommand(secretsCmd)
}
