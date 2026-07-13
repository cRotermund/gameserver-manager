package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/cRotermund/gameserver-manager/harness/internal/harness"
	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:   "delete [overlay]",
	Short: "Delete a Kustomize overlay from the cluster",
	Long:  "Delete a Kustomize overlay (local or prod) from the connected Kubernetes cluster.",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		overlay := "local"
		if len(args) > 0 {
			overlay = args[0]
		}

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

		fmt.Printf("\n--- Deleting overlay %q ---\n", overlay)

		if harness.KustomizeAvailable() {
			manifest, err := harness.KustomizeBuild(overlayDir)
			if err != nil {
				return err
			}
			return harness.RunWithInput(manifest, "kubectl", "delete", "-f", "-")
		}

		return harness.Run("kubectl", "delete", "-k", overlayDir)
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
}