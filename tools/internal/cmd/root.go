// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "website-docs",
	Short: "Doc-generation tooling for the Gemara website",
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(newOpenAPI2MDCmd())
	rootCmd.AddCommand(newLexicon2MDCmd())
	rootCmd.AddCommand(newTermLinkerCmd())
}
