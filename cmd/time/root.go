// Copyright 2019 Gregory Doran <greg@gregorydoran.co.uk>. 
// All rights reserved.

package time

import "github.com/spf13/cobra"

var Cmd = &cobra.Command{
	Use: "time",
	Short: "Time management commands",
}

func init() {
	Cmd.AddCommand(reportCmd)
}

