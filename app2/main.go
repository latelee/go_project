package main

import (
    "fmt"
    "os"

    "github.com/spf13/cobra"    
)

var KubeEdgeVersion    string
var DockerVersion      string

func NewCloudCoreCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use: "app.exe",
        Short: `This is short msg`,
		Long: `This is long usage message`,
		Run: func(cmd *cobra.Command, args []string) {
            fmt.Println("test of cobra")
            fmt.Println(DockerVersion)
		},
	}

	usageFmt := "Usage:\n  %s\n"
	cmd.SetUsageFunc(func(cmd *cobra.Command) error {
		fmt.Fprintf(cmd.OutOrStderr(), usageFmt, cmd.UseLine())
		return nil
	})
	cmd.SetHelpFunc(func(cmd *cobra.Command, args []string) {
		fmt.Fprintf(cmd.OutOrStdout(), "%s\n\n"+usageFmt, cmd.Long, cmd.UseLine())
	})

    cmd.ResetFlags()

    // 命令参数，保存的值，参数名，默认参数，说明
	cmd.Flags().StringVar(&DockerVersion, "docker", "18.08",
		"Use this key to download and use the required Docker version")


	return cmd
}


func main() {
	command := NewCloudCoreCommand()

	if err := command.Execute(); err != nil {
		os.Exit(1)
	}
}
