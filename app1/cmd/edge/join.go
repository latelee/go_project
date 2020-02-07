package cmd

import (
	"io"
    "fmt"

	"github.com/spf13/cobra"
	_ "github.com/spf13/pflag"
)

var (
	edgeJoinLongDescription = `
"keadm join" command bootstraps KubeEdge's worker node (at the edge) component.
It checks if the pre-requisites are installed already,
If not installed, this command will help in download,
install and execute on the host.
It will also connect with cloud component to receive 
further instructions and forward telemetry data from 
devices to cloud
`
	edgeJoinExample = `
keadm join --cloudcoreip=<ip address> --edgenodeid=<unique string as edge identifier>

  - For this command --cloudcoreip flag is a Mandatory option
  - This command will download and install the default version of pre-requisites and KubeEdge

keadm join --cloudcoreip=10.20.30.40 --edgenodeid=testing123 --kubeedge-version=0.2.1 --k8sserverip=50.60.70.80:8080 --interfacename eth0

- In case, any flag is used in a format like "--docker-version" or "--docker-version=" (without a value)
  then default versions shown in help will be chosen. 
  The versions for "--docker-version", "--kubernetes-version" and "--kubeedge-version" flags should be in the
  format "18.06.3", "1.14.0" and "0.2.1" respectively
`
)

// NewCloudInit represents the keadm init command for cloud component
func NewEdgeJoin(out io.Writer) *cobra.Command {

var DockerVersion      string

	var cmd = &cobra.Command{
		Use:     "join",
		Short:   "Bootstraps edge component. Checks and install (if required) the pre-requisites. Execute it on any edge node machine you wish to join",
		Long:    edgeJoinLongDescription,
		Example: edgeJoinExample,
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("join cmd....")
            fmt.Println(DockerVersion)
            return nil
		},
	}

    // 命令参数，保存的值，参数名，默认参数，说明
	cmd.Flags().StringVar(&DockerVersion, "docker", "18.08",
		"Use this key to download and use the required Docker version")

	return cmd
}