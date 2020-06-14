// Copyright 2019 Gregory Doran <greg@gregorydoran.co.uk>. 
// All rights reserved.

package timekeeping

import "github.com/spf13/cobra"

var TimeCmd = &cobra.Command{
	Use: "time",
	Short: "Time management commands",
}

func init() {
	TimeCmd.AddCommand(resolveCmd)
}

