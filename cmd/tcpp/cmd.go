package tcpp

import (
	//"fmt"

	"webdemo/common/conf"
	"webdemo/pkg/klog"

	"github.com/spf13/cobra"
	_ "github.com/spf13/pflag"
)

var (
	name             = `tcp`
	shortDescription = `tcp command`
	longDescription  = `tcp...
`
	example = `  example comming up...
`
)

var theCmd = []conf.UserCmdFunc{
	// fix warning: struct literal uses unkeyed fields
	{
		Name:      "server",
		ShortHelp: "server",
		Func:      TcpServer,
	},
	{
		Name:      "bar",
		ShortHelp: "foo1",
		Func:      Foo,
	},
}

func NewCmd() *cobra.Command {

	var cmd = &cobra.Command{
		Use:     name,
		Short:   shortDescription,
		Long:    longDescription + "\n", // + common.GetHelpInfo(theCmd),
		Example: example,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 { // 当没有参数时，使用默认的Help
				cmd.Help()
				return nil
			}
			// TcpServer()
			idx := -1
			for idx1, item := range theCmd {
				if args[0] == item.Name {
					idx = idx1 // why ???
					break
				}
			}
			if idx == -1 {
				klog.Printf("arg '%v' not support", args[0])
				cmd.Help()
				return nil
			}

			theCmd[idx].Func(args)

			return nil
		},
	}
	// note：使用子命令形式，下列可在help中展开
	// 命令参数，保存的值，参数名，默认参数，说明
	//cmd.Flags().StringVar(&mode, "db", "-", "set the database name")

	return cmd
}
