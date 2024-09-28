package cmd

import (
	//"fmt"

	web "webdemo/module/gin"

	"github.com/spf13/cobra"
	_ "github.com/spf13/pflag"
)

var (
	name             = `website`
	shortDescription = name + ` command`
	longDescription  = name + `  ...
`
	example = `  example comming up...
`
)

func NewCmd() *cobra.Command {

	var cmd = &cobra.Command{
		Use:     name,
		Short:   shortDescription,
		Long:    longDescription + "\n",
		Example: example,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 { // 当没有参数时，使用默认的Help
				obj := web.NewGinServer()
				// obj.Start(web.WEB_WEBSITE)
				obj.RunAll()
				// obj.RunStaticWebSite()
			}

			return nil
		},
	}
	// note：使用子命令形式，下列可在help中展开
	// 命令参数，保存的值，参数名，默认参数，说明
	//cmd.Flags().StringVar(&mode, "db", "-", "set the database name")

	return cmd
}
