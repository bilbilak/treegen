package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	app "github.com/bilbilak/treegen/config"
	"github.com/bilbilak/treegen/internal"
)

var (
	Help    bool
	Version bool
	License bool
)

var rootCmd = &cobra.Command{
	Use:   strings.ToLower(app.Name) + " [flags]\n  " + strings.ToLower(app.Name) + " [FILE]... [STDIN]",
	Short: "ASCII Tree Directory and File Structure Generator",
	Long:  app.Name + ` is a powerful CLI tool for generating an entire file and folder structure from its ASCII tree representation.`,
	Args:  cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		if Help {
			_ = cmd.Help()
			return
		}

		if Version {
			fmt.Println(app.Version)
			return
		}

		if License {
			fmt.Println(app.License)
			return
		}

		processed := internal.ProcessInput(args)

		if !processed {
			internal.Help()
		}
	},
}

func init() {
	rootCmd.CompletionOptions.HiddenDefaultCmd = true
	rootCmd.SetHelpCommand(&cobra.Command{Hidden: true})

	rootCmd.PersistentFlags().BoolVar(&Help, "help", false, "Show usage instructions")

	rootCmd.Flags().BoolVar(&Version, "version", false, "Display the installed version number")
	rootCmd.Flags().BoolVar(&License, "license", false, "Display the license name")

	rootCmd.Flags().BoolVarP(&internal.Force, "force", "f", false, "Overwrite existing files")
	rootCmd.Flags().BoolVarP(&internal.Verbose, "verbose", "v", false, "Enable verbose output")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		internal.FatalError(err)
	}
}
