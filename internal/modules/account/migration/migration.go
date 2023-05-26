package migration

import "embed"

//go:embed *.sql
var Migration embed.FS
