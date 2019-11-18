[development]
  dbname  = "{{.AppName}}-dev"
  host    = "127.0.0.1"
  port    = 5432
  user    = "postgres"
  adapter = "psql"
  password = ""
  sslMode = "disable"

[test]
  dbname  = "{{.AppName}}-test"
  host    = "127.0.0.1"
  port    = 5432
  user    = "postgres"
  adapter = "psql"
  password = ""
  sslMode = "disable"

[production]
  dbname  = "{{.AppName}}-production"
  host    = "127.0.0.1"
  port    = 5432
  user    = "postgres"
  adapter = "psql"
  password = ""
  sslMode = "disable"
