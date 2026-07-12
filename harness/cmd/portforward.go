package cmd

import (
	"fmt"

	"github.com/cRotermund/gameserver-manager/harness/internal/harness"
	"github.com/spf13/cobra"
)

var (
	pfLocalPort  int
	pfRemotePort int
)

var portForwardCmd = &cobra.Command{
	Use:   "port-forward <api|web>",
	Short: "Forward a service port to localhost",
	Long:  "Forward a Kubernetes service port to localhost. Use Ctrl+C to stop.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		service := args[0]

		if service == "bot" {
			return fmt.Errorf("bot has no service port to forward")
		}

		defaultPort, ok := harness.ServicePorts[service]
		if !ok {
			return fmt.Errorf("invalid service %q (valid: api, web)", service)
		}

		lp := defaultPort
		if cmd.Flags().Changed("local-port") {
			lp = pfLocalPort
		}

		rp := defaultPort
		if cmd.Flags().Changed("remote-port") {
			rp = pfRemotePort
		}

		name := fmt.Sprintf("svc/control-plane-%s", service)
		fmt.Printf("\n--- Port-forward %s localhost:%d -> :%d ---\n", name, lp, rp)
		fmt.Println("  Press Ctrl+C to stop.")

		return harness.Run("kubectl", "port-forward", name, fmt.Sprintf("%d:%d", lp, rp))
	},
}

func init() {
	portForwardCmd.Flags().IntVar(&pfLocalPort, "local-port", 0, "Local port (default: 8080 for api, 3000 for web)")
	portForwardCmd.Flags().IntVar(&pfRemotePort, "remote-port", 0, "Service port (default: same as service default)")
	rootCmd.AddCommand(portForwardCmd)
}