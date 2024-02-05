package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/xmapst/osreapi/cmd/client"
	"github.com/xmapst/osreapi/cmd/server"
	"github.com/xmapst/osreapi/internal/info"
)

const longText = `The remote executor (OSReApi) provides API remote operation mode,batch execution of Shell, Powershell, Python and other commands, and easily completes common management tasks such as running automated operation and maintenance scripts, polling processes, installing or uninstalling software, updating applications, and installing patches.`

func main() {
	cmd := &cobra.Command{
		Use:                   os.Args[0],
		Short:                 "Operating system remote execution interface",
		Long:                  longText,
		DisableAutoGenTag:     true,
		DisableFlagsInUseLine: true,
		CompletionOptions: cobra.CompletionOptions{
			HiddenDefaultCmd: true,
		},
		Version: info.Version,
	}

	cmd.SetFlagErrorFunc(flagErrorFunc)
	cmd.SetVersionTemplate("{{.Version}}\n")
	cmd.SetHelpTemplate(`{{if or .Runnable .HasSubCommands}}{{.UsageString}}{{end}}`)
	cmd.PersistentFlags().BoolP("help", "h", false, "Print usage")
	_ = cmd.PersistentFlags().MarkShorthandDeprecated("help", "please use --help")
	cmd.Flags().BoolP("version", "v", false, "Print version information and quit")
	cmd.AddCommand(server.New())
	cmd.AddCommand(client.New())
	cmd.AddCommand(&cobra.Command{
		Use:   "version",
		Short: "print version information and quit",
		RunE: func(cmd *cobra.Command, args []string) error {
			info.PrintHeadInfo()
			return nil
		},
	})

	if err := cmd.Execute(); err != nil {
		os.Exit(128)
	}
}

func flagErrorFunc(cmd *cobra.Command, err error) error {
	if err == nil {
		return nil
	}

	usage := ""
	if cmd.HasSubCommands() {
		usage = "\n\n" + cmd.UsageString()
	}
	return fmt.Errorf("%s\nSee '%s --help'.%s", err, cmd.CommandPath(), usage)
}
