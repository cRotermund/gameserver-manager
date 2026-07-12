package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"

	"github.com/cRotermund/gameserver-manager/harness/internal/harness"
	"github.com/spf13/cobra"
)

var buildTag string

var buildCmd = &cobra.Command{
	Use:   "build [service...]",
	Short: "Build container images for one or more services",
	Long:  "Build container images for one or more services (api, web, bot). Defaults to all services if none specified.",
	RunE: func(cmd *cobra.Command, args []string) error {
		root, err := harness.ProjectRoot()
		if err != nil {
			return err
		}

		cli, err := harness.ContainerCLI()
		if err != nil {
			return err
		}

		chosen, err := validateServices(args)
		if err != nil {
			return err
		}

		var images []string
		for _, name := range chosen {
			img := fmt.Sprintf("localhost/%s:%s", harness.Images[name], buildTag)
			fmt.Printf("\n--- Building %s (%s) ---\n", name, img)

			svcDir := filepath.Join(root, harness.ServiceDirs[name])
			if err := harness.Run(cli, "build", "-t", img, svcDir); err != nil {
				return err
			}
			images = append(images, img)
		}

		if len(images) > 0 {
			fmt.Println()
			if err := harness.LoadImages(images); err != nil {
				fmt.Fprintf(os.Stderr, "warning: image load failed: %v\n", err)
			}
		}

		return nil
	},
}

func init() {
	buildCmd.Flags().StringVar(&buildTag, "tag", "latest", "Image tag")
	rootCmd.AddCommand(buildCmd)
}

func validateServices(args []string) ([]string, error) {
	if len(args) == 0 {
		chosen := make([]string, 0, len(harness.ServiceDirs))
		for name := range harness.ServiceDirs {
			chosen = append(chosen, name)
		}
		slices.Sort(chosen)
		return chosen, nil
	}

	for _, svc := range args {
		if _, ok := harness.ServiceDirs[svc]; !ok {
			valid := make([]string, 0, len(harness.ServiceDirs))
			for name := range harness.ServiceDirs {
				valid = append(valid, name)
			}
			slices.Sort(valid)
			return nil, fmt.Errorf("invalid service %q (valid: %v)", svc, valid)
		}
	}

	return args, nil
}