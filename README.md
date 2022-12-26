# envars

## this is a simple cli tool that stores your project's environement variables in encrypted KeePass database. THIS TOOL IS FOR DEVELOPMENT ONLY - DO NOT USE IT IN PRODUCTION!

## how to build the code?
* install `go` on your machine
* run `git clone https://github.com/snirt/envars.git`
* `cd envars`
* `go build -ldflags "-s -w" -o "./bin/envars" .`
* now you can copy the file that created in the `bin` folder to `/usr/local/bin/` or create a soft link.


## how to use it?


`envars -h` will list the tool's commands 

## example of usage
* open your project's directory in terminal
* run `envars` or `envars add`
* now you should choose your database password 
* to be able to access the database you should export your master password to the session: `export ENVARS_PWD=YOUR_PASSWORD`
* now you can list your environment variables: `envars ls`
* `eval $(envars list --export)` command will export all your variables to the current session.
* create an alias by adding this line to your shell configuration: `alias export_vars='eval $(envars list -e)'`

## how to use it in VSCode?
* open your `launch.json` file
* add to your json object in `configurations` this: `"console": "integratedTerminal",`
* open the integrated terminal and export your project's master password as described above: `export ENVARS_PWD=YOUR_PASSWORD`
* run `export_vars`
* now you can run your code and your environment variables are safe. they will e available until you'll close the terminal session. 
