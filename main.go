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

package main

import (
    "fmt"
    "os"
    "github.com/ashwanthkumar/slack-go-webhook"
    "github.com/urfave/cli"
)

const helpTemplate = `
NAME:
   {{.Name}}{{if .Usage}} - {{.Usage}}{{end}}

USAGE:
   {{if .UsageText}}{{.UsageText}}{{else}}{{.HelpName}} {{if .VisibleFlags}}[global options]{{end}} {{if .ArgsUsage}}{{.ArgsUsage}}{{else}}[arguments...]{{end}}{{end}}{{if .Version}}{{if not .HideVersion}}

ARGUMENTS
   Requires just two arguments.

   * PULL_REQUEST_URL: github pull request url. e.g) https://github.com/rails/rails/pull/1
   * PHASE: Review phase. e.g) 1, 2

EXAMPLE:
   review https://github.com/mgi166/rails-showcase/pull/1 1

VERSION:
   {{.Version}}{{end}}{{end}}{{if .Description}}

DESCRIPTION:
   {{.Description}}{{end}}{{if len .Authors}}

AUTHOR{{with $length := len .Authors}}{{if ne 1 $length}}S{{end}}{{end}}:
   {{range $index, $author := .Authors}}{{if $index}}
   {{end}}{{$author}}{{end}}{{end}}{{if .VisibleCommands}}

GLOBAL OPTIONS:
   {{range $index, $option := .VisibleFlags}}{{if $index}}
   {{end}}{{$option}}{{end}}{{end}}{{if .Copyright}}

COPYRIGHT:
   {{.Copyright}}{{end}}
`

func main () {
	app := cli.NewApp()
	app.Name = "review"
	app.Usage = "Requests Code review to team member from terminal."
	app.Version = "0.0.1"
	app.Author = "mgi166"
	cli.AppHelpTemplate = helpTemplate

	app.Action = func(context *cli.Context) error {
		webhookUrl := "https://xxxx"
		if len(context.Args()) != 2 {
			fmt.Println("ERROR: Specify just two arguments. Run 'review -h' and confirm usage.")
			os.Exit(1)
		}

		attachment := slack.Attachment {}
		payload := slack.Payload(
			"test",
			"username",
			"icon",
			"channel",
			[]slack.Attachment{attachment},
		)

		if err := slack.Send(webhookUrl, "", payload); err != nil {
			fmt.Println(err)
		}

		return nil
	}

	app.Run(os.Args)
}
