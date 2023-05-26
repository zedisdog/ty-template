package models

import "github.com/zedisdog/ty/database"

type Account struct {
	database.CommonField
	Username string
	Password string
}
