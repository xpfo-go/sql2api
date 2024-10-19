package cmd

import (
	"context"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/xpfo-go/logs"
	"github.com/xpfo-go/sql2api/api"
	"github.com/xpfo-go/sql2api/server"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/go-sql-driver/mysql"
)

func init() {
	rootCmd.Flags().IntP("port", "P", 8000, "server port, default 8000")
	rootCmd.Flags().StringP("host", "H", "0.0.0.0", "server host, default 0.0.0.0")
}

var rootCmd = &cobra.Command{
	Use:   "sql2api",
	Short: "sql to api",
	Long:  "sql to api, Code by Go",
	Run: func(cmd *cobra.Command, args []string) {
		Start(cmd)
	},
}

func Start(cmd *cobra.Command) {
	// 1. init
	initLogs()

	// 2. watch the signal
	ctx, cancelFunc := context.WithCancel(context.Background())

	// 3. default database
	//database.DefaultDBClient = database.NewMysqlClient(&database.MysqlConfig{
	//	Database: "demo",
	//	Host:     "127.0.0.1",
	//	User:     "root",
	//	Password: "123456",
	//	Port:     3306,
	//})

	//if err := database.DefaultDBClient.Connect(); err != nil {
	//	panic(err.Error())
	//}

	api.RegisterRouter()

	// 4. start the server
	port, err := cmd.Flags().GetInt("port")
	if err != nil {
		panic(err)
	}

	host, err := cmd.Flags().GetString("host")
	if err != nil {
		panic(err)
	}

	httpServer := server.NewServer(server.Config{
		Host: host,
		Port: port,
	})
	go httpServer.Run(ctx)

	interrupt(cancelFunc)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

func interrupt(onSignal func()) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)

	for s := range c {
		logs.Infof("Caught signal %s. Exiting.", s)
		onSignal()
		close(c)
	}
}
