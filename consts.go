package main

const (
	AppName             = "envars"
	ENVARS_PWD          = "ENVARS_PWD"
	ENVARS_PWD_ERR_MSG  = "Please export db password as environment variable (export " + ENVARS_PWD + "='YOUR_DB_PASSWORD'"
	GENERAL_ERROR		= "Ohhh... Something went wrong"
	DB_FILE_NAME		= ".env.kdbx"
)

const (
	COLOR_RED  Color = "\033[0;31m"
	COLOR_NONE Color = "\033[0m"
)
