package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"

	"github.com/cRotermund/gameserver-manager/harness/internal/harness"
	"github.com/spf13/cobra"
)

var applyCmd = &cobra.Command{
	Use:   "apply [overlay]",
	Short: "Apply a Kustomize overlay to the cluster",
	Long:  "Apply a Kustomize overlay (local or prod) to the connected Kubernetes cluster.",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		overlay := "local"
		if len(args) > 0 {
			overlay = args[0]
		}
		return applyOverlay(overlay)
	},
}

func init() {
	rootCmd.AddCommand(applyCmd)
}

func applyOverlay(overlay string) error {
	if err := validateOverlay(overlay); err != nil {
		return err
	}

	root, err := harness.ProjectRoot()
	if err != nil {
		return err
	}

	overlayDir := filepath.Join(root, "infra", "k8s", overlay)
	if st, err := os.Stat(overlayDir); err != nil || !st.IsDir() {
		return fmt.Errorf("overlay directory not found: %s", overlayDir)
	}

	fmt.Printf("\n--- Applying overlay %q ---\n", overlay)

	if harness.KustomizeAvailable() {
		manifest, err := harness.KustomizeBuild(overlayDir)
		if err != nil {
			return err
		}
		return harness.RunWithInput(manifest, "kubectl", "apply", "-f", "-")
	}

	return harness.Run("kubectl", "apply", "-k", overlayDir)
}

func validateOverlay(overlay string) error {
	valid := []string{"local", "prod"}
	if !slices.Contains(valid, overlay) {
		return fmt.Errorf("invalid overlay %q (valid: %v)", overlay, valid)
	}
	return nil
}