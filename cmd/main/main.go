package main

import (
	"fback/internal/config"
	"fback/internal/modules/account"
	"github.com/zedisdog/ty/application"
)

func main() {
	application.Init(config.NewConfig())
	application.RegisterModule(new(account.Module))
	//application.RegisterModule(new(customer.Module))
	//application.RegisterModule(new(frontend.Module))
	application.Boot()
	application.Run()
	application.Wait()
}
