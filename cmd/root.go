// Copyright © 2017 Mladen Popadic <mladen.popadic.4@gmail.com>

package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/mpopadic/go_n_find/colors"
	"github.com/spf13/cobra"
)

var (
	pathFlag              string
	nameFlag              string
	replaceFlag           string
	ignoreCaseFlag        bool
	showAbsolutePathsFlag bool
	forceReplaceFlag      bool
	contentFlag           string
)

var (
	_numberOfResults     int
	_renameMap           map[string]string
	_replaceContentFiles []string
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "go_n_find",
	Short: "CLI for finding files and folders",
	Long:  `CLI tool for finding files and folders by name or content`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if pathFlag == "" {
			return fmt.Errorf("path flag is required")
		}
		if nameFlag == "" && contentFlag == "" {
			return fmt.Errorf("name flag or content flag are required")
		}
		return nil
	},
	// Uncomment the following line if your bare application
	// has an action associated with it:
	RunE: func(cmd *cobra.Command, args []string) error {

		// Set findOptions
		options := &findOptions{
			Path:              pathFlag,
			Name:              nameFlag,
			Content:           contentFlag,
			ReplaceWith:       replaceFlag,
			IgnoreCase:        ignoreCaseFlag,
			ShowAbsolutePaths: showAbsolutePathsFlag,
			ForceReplace:      forceReplaceFlag,
		}

		_numberOfResults = 0

		if options.ReplaceWith != "" && !options.ForceReplace {
			_renameMap = make(map[string]string)
		}

		if err := findInTree(options); err != nil {
			return err
		}

		colors.CYAN.Printf("Number of results: %d\n", _numberOfResults)
		if options.ReplaceWith != "" && options.Content == "" {
			if !options.ForceReplace {
				response := waitResponse("Are you sure? [Yes/No] ", map[string][]string{
					"Yes": []string{"Yes", "Y", "y"},
					"No":  []string{"No", "N", "n"},
				})
				switch response {
				case "Yes":
					renamePaths(_renameMap)
				case "No":
					colors.RED.Print(response)
				}
			} else {
				renamePaths(_renameMap)
			}
		}

		if options.ReplaceWith != "" && options.Content != "" {
			if !options.ForceReplace {
				response := waitResponse("Are you sure? [Yes/No] ", map[string][]string{
					"Yes": []string{"Yes", "Y", "y"},
					"No":  []string{"No", "N", "n"},
				})
				switch response {
				case "Yes":
					reg, err := regexp.Compile(options.Content)
					if err != nil {
						colors.RED.Printf("invalid content regular expresion\n")
						os.Exit(1)
					}
					replaceContent(_replaceContentFiles, reg, options.ReplaceWith)
				case "No":
					colors.RED.Print(response)
				}
			} else {
				reg, err := regexp.Compile(options.Content)
				if err != nil {
					colors.RED.Printf("invalid content regular expresion\n")
					os.Exit(1)
				}
				replaceContent(_replaceContentFiles, reg, options.ReplaceWith)
			}
		}

		return nil
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	colors.InitColors()

	RootCmd.Flags().StringVarP(&pathFlag, "path", "p", "", "path to directory")
	RootCmd.Flags().StringVarP(&nameFlag, "name", "n", "", "regular expression for matching file or directory name; This flag filter files if content flag is used")
	RootCmd.Flags().StringVarP(&replaceFlag, "replace", "r", "", "replaces mached regular expression parts with given value")
	RootCmd.Flags().BoolVarP(&ignoreCaseFlag, "ignore-case", "i", false, "ignore case for all regular expresions; Add '(?i)' in front of specific regex for ignore case")
	RootCmd.Flags().BoolVarP(&showAbsolutePathsFlag, "absolute-paths", "a", false, "print absolute paths in result")
	RootCmd.Flags().BoolVarP(&forceReplaceFlag, "force-replace", "f", false, "Force replace without responding")
	RootCmd.Flags().StringVarP(&contentFlag, "content", "c", "", "regular expression for matching file content")

}

func findInTree(options *findOptions) error {
	fileInfo, err := os.Stat(options.Path)
	if err != nil {
		return fmt.Errorf("could not get fileInfo for %s: %v", options.Path, err)
	}

	if fileInfo.IsDir() {
		files, err := ioutil.ReadDir(options.Path)
		if err != nil {
			return fmt.Errorf("could not read directory %s: %v", options.Path, err)
		}
		for _, file := range files {
			childOptions := options.CreateCopy()
			childOptions.Path = path.Join(options.Path, file.Name())
			findInTree(childOptions)
		}
	}

	doAction(options, fileInfo)
	return nil
}

func doAction(options *findOptions, fileInfo os.FileInfo) {
	absolutePath, err := filepath.Abs(options.Path)
	if err != nil {
		log.Fatalf("could not get absolute path: %v", err)
	}
	finalPathPrint := getPathPrintFormat(options.Path, absolutePath, options.ShowAbsolutePaths)

	if options.Name != "" && options.Content == "" {
		re := createRegex(options.Name, options.IgnoreCase)

		if re.MatchString(fileInfo.Name()) {
			_numberOfResults++
			if options.ReplaceWith != "" {
				pathDir := filepath.Dir(absolutePath)
				newFileName := re.ReplaceAllString(fileInfo.Name(), options.ReplaceWith)

				if options.ForceReplace {
					err := os.Rename(absolutePath, filepath.FromSlash(path.Join(pathDir, newFileName)))
					if err != nil {
						fmt.Printf("could not rename file: %v", err)
					}
					colors.RED.Print(absolutePath)
					colors.CYAN.Print(" => ")
					colors.GREEN.Println(filepath.FromSlash(path.Join(pathDir, newFileName)))
				} else {
					_renameMap[absolutePath] = filepath.FromSlash(path.Join(pathDir, newFileName))

					fmt.Print(absolutePath)
					colors.CYAN.Print(" => ")
					fmt.Println(filepath.FromSlash(path.Join(pathDir, newFileName)))
				}
			} else {
				fmt.Println(filepath.FromSlash(finalPathPrint))
			}
		}
	}
	if options.Content != "" {
		if options.Name != "" {
			regName := createRegex(options.Name, options.IgnoreCase)
			if regName.MatchString(fileInfo.Name()) && !fileInfo.IsDir() {
				_replaceContentFiles = append(_replaceContentFiles, absolutePath)
				re := createRegex(options.Content, options.IgnoreCase)

				fileBytes, err := ioutil.ReadFile(absolutePath)
				if err != nil {
					log.Fatalf("could not read file content: %v", err)
				}
				fileString := string(fileBytes)

				fileLines := strings.Split(fileString, "\n")

				printedFileName := false
				for lineNumber, line := range fileLines {
					if re.MatchString(line) {
						_numberOfResults++
						if !printedFileName {
							colors.CYAN.Printf("%s:\n", finalPathPrint)
							printedFileName = !printedFileName
						}
						allIndexes := re.FindAllStringIndex(line, -1)

						colors.YELLOW.Printf("%v:", lineNumber+1)
						location := 0
						for _, match := range allIndexes {
							fmt.Printf("%s", line[location:match[0]])
							colors.GREEN.Printf("%s", line[match[0]:match[1]])
							location = match[1]
						}
						fmt.Printf("%s", line[location:])
						fmt.Println()
					}
				}
			}
		} else {
			if !fileInfo.IsDir() {
				_replaceContentFiles = append(_replaceContentFiles, absolutePath)
				re := createRegex(options.Content, options.IgnoreCase)

				fileBytes, err := ioutil.ReadFile(absolutePath)
				if err != nil {
					log.Fatalf("could not read file content: %v", err)
				}
				fileString := string(fileBytes)

				fileLines := strings.Split(fileString, "\n")

				printedFileName := false
				for lineNumber, line := range fileLines {
					if re.MatchString(line) {
						_numberOfResults++
						if !printedFileName {
							colors.CYAN.Printf("%s:\n", finalPathPrint)
							printedFileName = !printedFileName
						}
						allIndexes := re.FindAllStringIndex(line, -1)

						colors.YELLOW.Printf("%v:", lineNumber+1)
						location := 0
						for _, match := range allIndexes {
							fmt.Printf("%s", line[location:match[0]])
							colors.GREEN.Printf("%s", line[match[0]:match[1]])
							location = match[1]
						}
						fmt.Printf("%s", line[location:])
						fmt.Println()
					}
				}
			}
		}
	}
}

type findOptions struct {
	Path              string
	Name              string
	Content           string
	ReplaceWith       string
	IgnoreCase        bool
	ShowAbsolutePaths bool
	ForceReplace      bool
}

func (o *findOptions) CreateCopy() *findOptions {
	newFindOptions := &findOptions{
		Path:              o.Path,
		Name:              o.Name,
		Content:           o.Content,
		ReplaceWith:       o.ReplaceWith,
		IgnoreCase:        o.IgnoreCase,
		ShowAbsolutePaths: o.ShowAbsolutePaths,
		ForceReplace:      o.ForceReplace,
	}
	return newFindOptions
}

func waitResponse(question string, responseAliases map[string][]string) string {
	colors.YELLOW.Printf("%s ", question)
	var respond string

	for {
		fmt.Scanf("%s\n", &respond)

		for response, aliases := range responseAliases {
			for _, alias := range aliases {
				if respond == alias {
					return response
				}
			}
		}

		colors.YELLOW.Printf("%s ", question)
	}
}

func renamePaths(paths map[string]string) error {
	for oldPath, newPath := range paths {
		err := os.Rename(oldPath, newPath)
		if err != nil {
			return fmt.Errorf("could not rename file: %v", err)
		}
		colors.RED.Print(oldPath)
		colors.CYAN.Print(" => ")
		colors.GREEN.Println(newPath)
	}
	return nil
}

func getPathPrintFormat(filePath, absolutePath string, showAbsolute bool) string {
	var result = ""
	if showAbsolute {
		result = absolutePath
	} else {
		result = filePath
	}
	return filepath.Clean(result)
}

func createRegex(text string, ignoreCase bool) *regexp.Regexp {
	re, err := regexp.Compile(text)
	if err != nil {
		colors.RED.Printf("regular expresion for name flag is not valid\n")
		os.Exit(1)
	}
	if ignoreCase {
		re, err = regexp.Compile("(?i)" + text)
		if err != nil {
			colors.RED.Printf("regular expresion for name flag is not valid\n")
			os.Exit(1)
		}
	}
	return re
}

func replaceContent(filePaths []string, oldContent *regexp.Regexp, newContent string) {
	for _, filePath := range filePaths {

		fileInfo, err := os.Stat(filePath)
		if err != nil {
			log.Fatalf("could not get file info; %v", err)
		}
		fileBytes, err := ioutil.ReadFile(filePath)
		if err != nil {
			log.Fatalf("could not read file content: %v", err)
		}
		fileString := string(fileBytes)

		newFileString := oldContent.ReplaceAllString(fileString, newContent)

		err = ioutil.WriteFile(filePath, []byte(newFileString), fileInfo.Mode())
		if err != nil {
			fmt.Printf("could not write to file: %v", err)
		} else {
			colors.GREEN.Printf("%s\n", filePath)
		}
	}
}
