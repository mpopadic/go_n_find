# go_n_find

CLI tool for finding files and folders and renaming them



```
Usage:
  go_n_find [flags]
  go_n_find [command]

Available Commands:
  help        Help about any command
  version     Prints version number

Flags:
  -a, --absolute-paths   print absolute paths in result
  -c, --content string   regular expression for matching file content
  -f, --force-replace    Force replace without responding
  -h, --help             help for go_n_find
  -i, --ignore-case      ignore case for all regular expresions; Add '(?i)' in front of specific regex for ignore case
  -n, --name string      regular expression for matching file or directory name; This flag filter files if content flag is used
  -p, --path string      path to directory
  -r, --replace string   replaces mached regular expression parts with given value

Use "go_n_find [command] --help" for more information about a command.
```


