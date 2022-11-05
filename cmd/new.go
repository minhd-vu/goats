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
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var newFlags newCmdFlags

type newCmdFlags struct {
	template string
	output   string
	force    bool
}

// newCmd represents the new command
var newCmd = &cobra.Command{
	Use:   "new [filename] [-t template] [-o output]",
	Short: "Create a new note",
	Long:  `Create a new note with a given filename. If no filename is given, the filename will be the timestamp.`,
	RunE: func(_ *cobra.Command, args []string) error {
		if newFlags.output == "" {
			newFlags.output = filepath.Join(cfg.RootPath, cfg.NotesDir)
		}

		if err := os.MkdirAll(newFlags.output, os.ModePerm); err != nil {
			return err
		}

		var filename string
		if len(args) > 0 {
			filename = args[0]
		} else {
			filename = time.Now().Format(time.RFC3339)
		}

		filename = fmt.Sprintf("%s.%s", filename, cfg.FileType)
		path := filepath.Join(newFlags.output, filename)
		file, err := os.Create(path)
		if err != nil {
			return err
		}

		if newFlags.template != "" {
			templateFilename := fmt.Sprintf("%s.%s", newFlags.template, cfg.FileType)
			templateFilepath := filepath.Join(cfg.RootPath, cfg.TemplatesDir, templateFilename)
			templateFile, err := os.Open(templateFilepath)
			if err != nil {
				return err
			}

			_, err = io.Copy(file, templateFile)
			if err != nil {
				return err
			}
		}

		cmd := exec.Command(cfg.Editor, path)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		if err := cmd.Run(); err != nil {
			log.Error(err)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(newCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// newCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	newCmd.Flags().StringVarP(&newFlags.template, "template", "t", "", "The template to use")
	newCmd.Flags().StringVarP(&newFlags.output, "output", "o", "", "Where the new note will be created")
	newCmd.Flags().BoolVarP(&newFlags.force, "force", "f", false, "Overwrite note if already exists")
}
