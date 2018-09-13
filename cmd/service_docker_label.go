package cmd

import (
	"github.com/spf13/cobra"
)

var serviceDockerLabelCmd = &cobra.Command{
	Use:   "docker-label",
	Short: "Manage docker labels",
}

func init() {
	serviceCmd.AddCommand(serviceDockerLabelCmd)
}
