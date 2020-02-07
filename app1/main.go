package main

import (
    _ "fmt"
    "os"
    "io"

    "github.com/spf13/cobra"
    cloud "github.com/latelee/myproject/app1/cmd/cloud"
    edge "github.com/latelee/myproject/app1/cmd/edge"
)


var (
	keadmLongDescription = `
    long description.

`
	keadmExample = `
    example message.
`
)


// NewKubeedgeCommand returns cobra.Command to run keadm commands
func NewKubeedgeCommand(in io.Reader, out, err io.Writer) *cobra.Command {

	cmds := &cobra.Command{
		Use:     "keadm",
		Short:   "keadm: Bootstrap KubeEdge cluster",
		Long:    keadmLongDescription,
		Example: keadmExample,
	}

	cmds.ResetFlags()
	cmds.AddCommand(cloud.NewCloudInit(out))
	cmds.AddCommand(edge.NewEdgeJoin(out))

	return cmds
}


func main() {
	cmd := NewKubeedgeCommand(os.Stdin, os.Stdout, os.Stderr)

	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
