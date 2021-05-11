/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"log"
	"sumaclouder/utils"

	"github.com/spf13/cobra"
)

// imgupdateCmd represents the imgupdate command
var (
	imgupdateCmd = &cobra.Command{
		Use:   "imgupdate",
		Short: "A brief description of your command",
		Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
		Run: func(cmd *cobra.Command, args []string) {
			imgupdateRun()
		},
	}
	downloadSuseLink = "http://download.suse.de/ibs/Devel:/PubCloud:/Stable:/CrossCloud:/SLE15-SP3/images/" //hardcoded - to be improved
)

func init() {
	rootCmd.AddCommand(imgupdateCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// imgupdateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// imgupdateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func imgupdateRun() {
	//line := `<tr><td valign="top"><a href="SUSE-Manager-Server-BYOS.x86_64-0.9.0-GCE-Build1.23.tar.gz.sha256.asc"><img src="/icons/encrypted.png" alt="[   ]" width="16" height="16" /></a></td><td><a href="SUSE-Manager-Server-BYOS.x86_64-0.9.0-GCE-Build1.23.tar.gz.sha256.asc">SUSE-Manager-Server-BYOS.x86_64-0.9.0-GCE-Build1.23.tar.gz.sha256.asc</a></td><td align="right">07-May-2021 20:12  </td><td align="right">189  </td><td><a href="SUSE-Manager-Server-BYOS.x86_64-0.9.0-GCE-Build1.23.tar.gz.sha256.asc.mirrorlist">`
	//utils.ParseWebHTMLLine(line)

	imglist, err := utils.ListObjectsInBucket(conf.GCEAuthPath, projectID, bucketName, "")
	if err != nil {
		log.Fatalf("Error %v", err)
	}
	err = utils.CheckifImgUpdated(imglist, downloadSuseLink)
	if err != nil {
		log.Fatalf("Error %v", err)
	}
}
