package main

import (
	"github.com/orestonce/ChessGame/ChessServer/ChessGame"
	"github.com/orestonce/ChessGame/ChessServer/ChessGate"
	"github.com/orestonce/ChessGame/ymd/ymdQuickRestart"
	"github.com/spf13/cobra"
)

var root = cobra.Command{
	Use: "ChessServer",
}

func main() {
	if err := root.Execute(); err != nil {
		root.Help()
	}
}

func init() {
	var redis ymdQuickRestart.RedisInfo
	var mysqlSchema string
	game := &cobra.Command{
		Use: "Game",
		Run: func(cmd *cobra.Command, args []string) {
			ChessGame.RunChessGame(redis, mysqlSchema)
		},
	}
	game.Flags().StringVarP(&redis.RedisAddr, `raddr`, ``, `127.0.0.1:6379`, `redis地址`)
	game.Flags().StringVarP(&redis.RedisPasswd, `rpasswd`, ``, ``, `redis密码`)
	game.Flags().StringVarP(&redis.Prefix, `rprefix`, ``, `chess`, `redis前缀`)
	game.Flags().StringVarP(&mysqlSchema, `mysql_schema`, ``, `root:@tcp(127.0.0.1:3306)/chess?parseTime=true&charset=utf8`, `mysql连接`)
	root.AddCommand(game)

	var laddr string
	var wspath string

	gate := &cobra.Command{
		Use: "Gate",
		Run: func(cmd *cobra.Command, args []string) {
			ChessGate.RunChessGate(laddr, redis, wspath)
		},
	}
	gate.Flags().StringVarP(&redis.RedisAddr, `raddr`, ``, `127.0.0.1:6379`, `redis地址`)
	gate.Flags().StringVarP(&redis.RedisPasswd, `rpasswd`, ``, ``, `redis密码`)
	gate.Flags().StringVarP(&redis.Prefix, `rprefix`, ``, `chess`, `redis前缀`)
	gate.Flags().StringVarP(&laddr, `laddr`, ``, `127.0.0.1:8912`, `监听地址`)
	gate.Flags().StringVarP(&wspath, `wspath`, ``, `/ChessGame`, `websocket路径`)
	root.AddCommand(gate)
}
