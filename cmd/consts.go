package cmd

const (
	AppName         = "envars"
	ENVARS_PWD      = "ENVARS_PWD"
	ERR_PWD_ERR_MSG = "Please export database password as environment variable (export " + ENVARS_PWD + "='YOUR_DB_PASSWORD'"
	ERR_DECODE      = "could not decode the file"
	ERR_FILE_OPEN	= "could not open the db file"
	ERR_LOCK_DB		= "could not lock db"
	GENERAL_ERROR   = "Ohhh... Something went wrong"
	DB_FILE_NAME    = ".env.kdbx"
)

const (
	COLOR_RED  Color = "\033[0;31m"
	COLOR_NONE Color = "\033[0m"
)
