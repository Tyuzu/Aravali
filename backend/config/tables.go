package config

type Tables struct {
	UserTable    string
	ProfileTable string
}

var TablesConfig = Tables{
	UserTable:    "users",
	ProfileTable: "profiles",
}
