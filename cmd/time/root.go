// Copyright 2019 Gregory Doran <greg@gregorydoran.co.uk>.
// All rights reserved.

package time

import "github.com/spf13/cobra"

// Cmd contains commands related to time management
var Cmd = &cobra.Command{
	Use:   "time",
	Short: "Time management commands",
}

func init() {
	Cmd.AddCommand(reportCmd)
}
