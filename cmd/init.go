// Copyright © 2018 Andrew Fields <andy@andybug.com>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package cmd

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh/terminal"
	log "github.com/sirupsen/logrus"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a new abakus repository",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		InitLogging()

		cwd, _ := os.Getwd()
		home := filepath.Join(cwd, ".abakus")
		objects := filepath.Join(home, "objects")
		snapshots := filepath.Join(home, "snapshots")

		log.WithFields(log.Fields{
			"cmd": "init",
			"home": home,
		}).Info("Initializing abakus repository")

		createDir(home)
		createDir(objects)
		createDir(snapshots)

		configPath := filepath.Join(home, "config")
		config, pw := promptForConfiguration()
		config.write(configPath, pw)
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
	// initCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func createDir(path string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		log.WithFields(log.Fields{"dir": path,}).Debug("Creating directory")
		err := os.Mkdir(path, 0755)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func promptForConfiguration() (config Config, password string) {
	aws := promptForAWSConfiguration()
	password = promptSecret("Password")
	
	config = Config{
		Aws: aws,
	}

	return
}

func promptForAWSConfiguration() (aws *AWSConfig) {
	wantAws := promptString("Enter AWS configuration? (y/n)")
	if wantAws != "y" {
		return nil
	}
	
	bucket := promptString("S3 bucket")
	prefix := promptString("Prefix")
	accessKey := promptSecret("Access Key ID")
	secretKey := promptSecret("Secret Access Key")

	log.WithFields(log.Fields{
		"bucket": bucket,
		"prefix": prefix,
	}).Debug("AWS configuration entered")

	aws = &AWSConfig{
		Bucket: bucket,
		Prefix: prefix,
		AccessKey: accessKey,
		SecretKey: secretKey,
	}

	return
}

func promptString(msg string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("%s: ", msg)
	retval, _ := reader.ReadString('\n')
	return strings.TrimSpace(retval)
}

func promptSecret(msg string) string {
	fmt.Printf("%s: ", msg)
	bSecret, _ := terminal.ReadPassword(int(syscall.Stdin))
	retval := string(bSecret)
	fmt.Print("\n")
	return retval
}
