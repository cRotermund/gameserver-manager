package harness

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const (
	defaultSecretsEnvFile = "infra/k8s/prod/secrets.env"
	defaultAWSRegion      = "us-east-1"
)

type SSMParameter struct {
	Name  string
	Value string
}

type ssmPutParameterInput struct {
	Name      string `json:"Name"`
	Type      string `json:"Type"`
	Value     string `json:"Value"`
	Overwrite bool   `json:"Overwrite"`
}

func DefaultSecretsEnvFile() string {
	return defaultSecretsEnvFile
}

func DefaultAWSRegion() string {
	return defaultAWSRegion
}

func LoadProductionSecrets(path string) ([]SSMParameter, error) {
	values, err := loadEnvFile(path)
	if err != nil {
		return nil, err
	}

	required := []string{
		"AWS_ACCESS_KEY_ID",
		"AWS_SECRET_ACCESS_KEY",
		"AWS_ROLE_ARN",
	}
	for _, key := range required {
		if strings.TrimSpace(values[key]) == "" {
			return nil, fmt.Errorf("required variable %q is missing or empty in %s", key, path)
		}
		if strings.ContainsAny(values[key], "\r\n") {
			return nil, fmt.Errorf("variable %q must be a single-line value", key)
		}
	}

	return []SSMParameter{
		{
			Name:  "/kubernetes/appsecrets/gameserver-manager/control-plane-api/aws-access-key-id",
			Value: values["AWS_ACCESS_KEY_ID"],
		},
		{
			Name:  "/kubernetes/appsecrets/gameserver-manager/control-plane-api/aws-secret-access-key",
			Value: values["AWS_SECRET_ACCESS_KEY"],
		},
		{
			Name:  "/kubernetes/appsecrets/gameserver-manager/control-plane-api/aws-role-arn",
			Value: values["AWS_ROLE_ARN"],
		},
	}, nil
}

func PutSSMParameters(parameters []SSMParameter, region, profile string, overwrite, dryRun bool) error {
	if strings.TrimSpace(region) == "" {
		return fmt.Errorf("AWS region must not be empty")
	}

	for _, parameter := range parameters {
		if strings.TrimSpace(parameter.Name) == "" || parameter.Value == "" {
			return fmt.Errorf("parameter name and value must not be empty")
		}
		if dryRun {
			fmt.Printf("would %s %s\n", actionName(overwrite), parameter.Name)
			continue
		}
		if err := putSSMParameter(parameter, region, profile, overwrite); err != nil {
			return fmt.Errorf("write %s: %w", parameter.Name, err)
		}
		fmt.Printf("wrote %s\n", parameter.Name)
	}

	return nil
}

func putSSMParameter(parameter SSMParameter, region, profile string, overwrite bool) error {
	input := ssmPutParameterInput{
		Name:      parameter.Name,
		Type:      "SecureString",
		Value:     parameter.Value,
		Overwrite: overwrite,
	}

	temporary, err := os.CreateTemp("", "gameserver-manager-ssm-*.json")
	if err != nil {
		return fmt.Errorf("create temporary AWS CLI input: %w", err)
	}
	temporaryPath := temporary.Name()
	defer os.Remove(temporaryPath)

	if err := temporary.Chmod(0600); err != nil {
		_ = temporary.Close()
		return fmt.Errorf("restrict temporary AWS CLI input: %w", err)
	}
	if err := json.NewEncoder(temporary).Encode(input); err != nil {
		_ = temporary.Close()
		return fmt.Errorf("encode AWS CLI input: %w", err)
	}
	if err := temporary.Close(); err != nil {
		return fmt.Errorf("close temporary AWS CLI input: %w", err)
	}

	args := []string{
		"ssm", "put-parameter",
		"--cli-input-json", "file://" + filepath.ToSlash(temporaryPath),
		"--region", region,
		"--query", "Version",
		"--output", "text",
	}
	if profile != "" {
		args = append(args, "--profile", profile)
	}
	return Run("aws", args...)
}

func loadEnvFile(path string) (map[string]string, error) {
	contents, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read env file %s: %w", path, err)
	}

	values := make(map[string]string)
	for lineNumber, line := range strings.Split(string(contents), "\n") {
		line = strings.TrimSpace(strings.TrimSuffix(line, "\r"))
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		line = strings.TrimPrefix(line, "export ")
		key, value, found := strings.Cut(line, "=")
		if !found || strings.TrimSpace(key) == "" {
			return nil, fmt.Errorf("invalid env file %s line %d: expected KEY=VALUE", path, lineNumber+1)
		}
		key = strings.TrimSpace(key)
		value = strings.TrimSpace(value)
		if len(value) >= 2 && ((value[0] == '"' && value[len(value)-1] == '"') || (value[0] == '\'' && value[len(value)-1] == '\'')) {
			if value[0] == '"' {
				unquoted, err := strconv.Unquote(value)
				if err != nil {
					return nil, fmt.Errorf("invalid quoted value for %s line %d: %w", key, lineNumber+1, err)
				}
				value = unquoted
			} else {
				value = value[1 : len(value)-1]
			}
		}
		if _, exists := values[key]; exists {
			return nil, fmt.Errorf("duplicate variable %q in %s line %d", key, path, lineNumber+1)
		}
		values[key] = value
	}

	return values, nil
}

func actionName(overwrite bool) string {
	if overwrite {
		return "rotate"
	}
	return "create"
}
