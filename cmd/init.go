/*
Copyright © 2022 Minh Vu <minhd_vu@yahoo.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"path/filepath"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var force bool

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize the gotes config",
	Long: `Initializing the gotes config will create a file at ~/.gotes.yaml unless you set the --config flag. The config file consist of these variables:
  - root_path: The root path to gotes. The notes_dir and templates_dir will be inside of this root path.
  - notes_dir: The directory where note files will be stored.
  - templates_dir: The directory where template files will be stored.
  - editor: The default editor to use when opening notes. This should be something that is callable from path.
            Examples: vi, vim, nvim, emacs, code`,
	RunE: func(cmd *cobra.Command, args []string) error {
		questions := []*survey.Question{
			{
				Name: "root_path",
				Prompt: &survey.Input{
					Message: "Root Path:",
					Default: cfg.RootPath,
					Help:    `The root path to gotes. The notes_dir and templates_dir will be inside of this root path.`,
				},
				Transform: survey.TransformString(filepath.Clean),
			},
			{
				Name: "notes_dir",
				Prompt: &survey.Input{
					Message: "Notes Dir:",
					Default: cfg.NotesDir,
					Help:    `The directory where note files will be stored.`,
				},
				Transform: survey.TransformString(filepath.Clean),
			},
			{
				Name: "templates_dir",
				Prompt: &survey.Input{
					Message: "Templates Dir:",
					Default: cfg.TemplatesDir,
					Help:    `The directory where template files will be stored.`,
				},
				Transform: survey.TransformString(filepath.Clean),
			},
			{
				Name: "editor",
				Prompt: &survey.Input{
					Message: "Editor:",
					Default: cfg.Editor,
					Help: `The default editor to use when opening notes. This should be something that is callable from path.
Examples: vi, vim, nvim, emacs, code`},
			},
		}

		answers := make(map[string]interface{})
		if err := survey.Ask(questions, &answers); err != nil {
			return err
		}

		if err := viper.MergeConfigMap(answers); err != nil {
			return err
		}

		err := viper.SafeWriteConfig()
		if force {
			err = viper.WriteConfig()
		}

		return err
	},
}

func init() {
	rootCmd.AddCommand(initCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// initCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	initCmd.Flags().BoolVarP(&force, "force", "f", false, "Overwrite config file if it exists")
}
