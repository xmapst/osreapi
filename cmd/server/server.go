package server

import (
	"os"
	"path/filepath"
	"time"

	"github.com/kardianos/service"
	"github.com/spf13/cobra"

	"github.com/xmapst/osreapi/internal/logx"
	"github.com/xmapst/osreapi/internal/server"
	"github.com/xmapst/osreapi/internal/server/config"
)

var cmd = &cobra.Command{
	Use: "server",
	Aliases: []string{
		"run",
	},
	Short: "start server",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		name, err := filepath.Abs(os.Args[0])
		if err != nil {
			logx.Errorln(err)
			return err
		}
		svc, err := service.New(server.New(), &service.Config{
			Name:        config.App.ServiceName,
			DisplayName: "OSReApi",
			Description: "OS Remote Executor Api",
			Executable:  name,
			Arguments:   os.Args[1:],
		})
		if err != nil {
			return err
		}
		err = svc.Run()
		if err != nil {
			return err
		}
		return nil
	},
}

func init() {
	cmd.Flags().BoolVarP(&config.App.Debug, "debug", "d", false, "debug mode.")
	cmd.Flags().StringVar(&config.App.DBType, "db_type", "bolt", "database type. [bolt]")
	cmd.Flags().DurationVar(&config.App.ExecTimeOut, "exec_timeout", 24*time.Hour, "set the task exec command expire time")
	cmd.Flags().StringVarP(&config.App.ListenAddress, "addr", "a", "tcp://0.0.0.0:2376", "listening address.")
	cmd.Flags().BoolVarP(&config.App.Normal, "normal", "n", false, "wait for all task execution to complete.")
	cmd.Flags().IntVarP(&config.App.PoolSize, "pool_size", "p", 30, "set the size of the execution work pool.")
	cmd.Flags().Int64Var(&config.App.MaxRequests, "max_requests", 0, "max number of concurrent requests. 0 to disable.")
	cmd.Flags().StringVarP(&config.App.SelfUpdateURL, "self_url", "s", "https://oss.yfdou.com/tools/osreapi", "self Update URL")
	cmd.Flags().DurationVarP(&config.App.WebTimeout, "timeout", "t", 120*time.Second, "maximum duration before timing out read/write/idle.")
	cmd.Flags().StringVar(&config.App.DataDir, "data_dir", filepath.Join("$TEMP", config.App.ServiceName, "data"), "database directory")
	cmd.Flags().StringVar(&config.App.LogDir, "log_dir", filepath.Join("$TEMP", config.App.ServiceName, "logs"), "log output directory")
	cmd.Flags().StringVar(&config.App.ScriptDir, "script_dir", filepath.Join("$TEMP", config.App.ServiceName, "scripts"), "task script temp directory.")
	cmd.Flags().StringVar(&config.App.WorkSpace, "workspace_dir", filepath.Join("$TEMP", config.App.ServiceName, "workspace"), "task workspace temp directory")
}

func New() *cobra.Command {
	return cmd
}
