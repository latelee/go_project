package cmd

import (
	"io"
    "fmt"

	"github.com/spf13/cobra"
	_ "github.com/spf13/pflag"
)

var (
	cloudInitLongDescription = `
    This is long description.
`
	cloudInitExample = `
    keadm init example...
`
)

// NewCloudInit represents the keadm init command for cloud component
func NewCloudInit(out io.Writer) *cobra.Command {

var DockerVersion      string


	var cmd = &cobra.Command{
		Use:     "init",
		Short:   "Bootstraps cloud component. Checks and install (if required) the pre-requisites.",
		Long:    cloudInitLongDescription,
		Example: cloudInitExample,
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("init cmd....")
            fmt.Println(DockerVersion)
            return nil
		},
	}
    // note：使用子命令形式，下列可在help中展开
    // 命令参数，保存的值，参数名，默认参数，说明
	cmd.Flags().StringVar(&DockerVersion, "docker", "18.08",
		"Use this key to download and use the required Docker version")

	return cmd
}