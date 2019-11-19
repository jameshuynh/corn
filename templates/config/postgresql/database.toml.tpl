[development]
  dbname  = "{{.Module}}-dev"
  host    = "127.0.0.1"
  port    = 5432
  user    = "postgres"
  adapter = "psql"
  password = ""
  sslMode = "disable"

[test]
  dbname  = "{{.Module}}-test"
  host    = "127.0.0.1"
  port    = 5432
  user    = "postgres"
  adapter = "psql"
  password = ""
  sslMode = "disable"

[production]
  dbname  = "{{.Module}}-production"
  host    = "127.0.0.1"
  port    = 5432
  user    = "postgres"
  adapter = "psql"
  password = ""
  sslMode = "disable"
