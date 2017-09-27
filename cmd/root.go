// Copyright © 2017 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"regexp"

	"github.com/spf13/cobra"
)

var (
	pathFlag              string
	nameFlag              string
	replaceFlag           string
	ignoreCaseFlag        bool
	showAbsolutePathsFlag bool
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "go_n_find",
	Short: "CLI for finding files and folders",
	Long:  `CLI tool for finding files and folders by name or content`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if nameFlag == "" {
			return fmt.Errorf("name flag must be defined")
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
			ReplaceWith:       replaceFlag,
			IgnoreCase:        ignoreCaseFlag,
			ShowAbsolutePaths: showAbsolutePathsFlag,
		}

		if err := findInTree(options); err != nil {
			return err
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

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	RootCmd.Flags().StringVarP(&pathFlag, "path", "p", ".", "path to directory")
	RootCmd.Flags().StringVarP(&nameFlag, "name", "n", "", "name of file or directory")
	RootCmd.Flags().StringVarP(&replaceFlag, "replace", "r", "", "replace string")
	RootCmd.Flags().BoolVarP(&ignoreCaseFlag, "ignore-case", "i", false, "ignore case flag")
	RootCmd.Flags().BoolVarP(&showAbsolutePathsFlag, "absolute-paths", "a", false, "print absolute paths in result")

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

	doAction(options, fileInfo.Name())
	return nil
}

func doAction(options *findOptions, fileName string) {
	if options.Name != "" {
		var finalPathPrint = ""
		absolutePath, err := filepath.Abs(options.Path)
		if err != nil {
			log.Fatalf("could not get absolute path: %v", err)
		}
		if options.ShowAbsolutePaths {
			finalPathPrint = absolutePath
		} else {
			finalPathPrint = options.Path
		}

		re := regexp.MustCompile(options.Name)
		if options.IgnoreCase {
			re = regexp.MustCompile("(?i)" + options.Name)
		}

		if re.MatchString(fileName) {
			if options.ReplaceWith != "" {
				pathDir := filepath.Dir(absolutePath)
				newFileName := re.ReplaceAllString(fileName, options.ReplaceWith)

				err := os.Rename(absolutePath, filepath.FromSlash(path.Join(pathDir, newFileName)))
				if err != nil {
					fmt.Printf("could not rename file: %v", err)
				}
				fmt.Printf("%s => %s\n", absolutePath, filepath.FromSlash(path.Join(pathDir, newFileName)))
			} else {
				fmt.Println(finalPathPrint)
			}
		}
	}
}

type findOptions struct {
	Path              string
	Name              string
	ReplaceWith       string
	IgnoreCase        bool
	ShowAbsolutePaths bool
}

func (o *findOptions) CreateCopy() *findOptions {
	newFindOptions := &findOptions{
		Path:              o.Path,
		Name:              o.Name,
		ReplaceWith:       o.ReplaceWith,
		IgnoreCase:        o.IgnoreCase,
		ShowAbsolutePaths: o.ShowAbsolutePaths,
	}
	return newFindOptions
}
