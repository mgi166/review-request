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
    "bytes"
    "fmt"
    "os"
    "os/user"
    "path"
    "reflect"
    "text/template"
    "strings"
    "strconv"
    "time"
    "github.com/ashwanthkumar/slack-go-webhook"
    "github.com/BurntSushi/toml"
    "github.com/urfave/cli"
)

const helpTemplate = `
NAME:
   {{.Name}}{{if .Usage}} - {{.Usage}}{{end}}

USAGE:
   {{if .UsageText}}{{.UsageText}}{{else}}{{.HelpName}} {{if .VisibleFlags}}[global options]{{end}} {{if .ArgsUsage}}{{.ArgsUsage}}{{else}}[arguments...]{{end}}{{end}}{{if .Version}}{{if not .HideVersion}}

ARGUMENTS
   Requires just two arguments.

   - PULL_REQUEST_URL: github pull request url. e.g) https://github.com/rails/rails/pull/1
   - PHASE: Review phase. e.g) 1, 2

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

type Config struct {
	Text     string
	Review   ReviewConfig
	Reviewer ReviewerConfig
}

type ReviewConfig struct {
	Slack ReviewSlackConfig
}

type ReviewSlackConfig struct {
	WebhookUrl string
	UserName   string
	Icon       string
	Channel    string
}

type ReviewerConfig struct {
	Monday    ReviewPhaseConfig
	Tuesday   ReviewPhaseConfig
	Wednesday ReviewPhaseConfig
	Thursday  ReviewPhaseConfig
	Friday    ReviewPhaseConfig
	Sunday    ReviewPhaseConfig
}

type ReviewPhaseConfig struct {
	Phase1 []string
	Phase2 []string
}

func createApp(app *cli.App) *cli.App {
	app.Name = "review"

	app.Name = "review"
	app.Usage = "Requests Code review to team member from terminal."
	app.Version = "0.0.1"
	app.Author = "mgi166"

	user, _ := user.Current()

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name: "config, c",
			Value: path.Join(user.HomeDir, ".review"),
			Usage: "Load configuration from `FILE`.",
		},
		cli.BoolFlag{
			Name: "dry-run, d",
			Usage: "Dry run. If true, review request is not sent. (default: false)",
		},
		cli.IntFlag{
			Name: "phase, p",
			Usage: "Review phase. For example, specify `1` when 1 phase review",
		},
	}
	cli.AppHelpTemplate = helpTemplate

	return app
}

func createConfig(context *cli.Context) Config {
	var config Config

	if _, err := toml.DecodeFile(context.String("config"), &config); err != nil {
		panic(err)
	}

	return config
}

func main () {
	app := createApp(cli.NewApp())

	app.Action = func(context *cli.Context) error {
		if len(context.Args()) != 2 {
			fmt.Println("ERROR: Specify just two arguments. Run 'review -h' and confirm usage.")
			os.Exit(1)
		}

		config := createConfig(context)

		week := time.Now().Weekday().String()

		phase := "Phase" + strconv.Itoa(context.Int("phase"))
		reviewers := reflect.ValueOf(config.Reviewer).FieldByName(week).FieldByName(phase)
		tmpl, err := template.New("text").Parse(config.Text)

		if err != nil { panic(err) }

		var buffer bytes.Buffer

		dict := make(map[string]string)

		dict["reviewers"] = strings.Join(reviewers.Interface().([]string), " ")
		dict["url"] = context.Args().Get(0)
		dict["phase"] = context.Args().Get(1)

		if err := tmpl.Execute(&buffer, dict); err != nil {
			panic(err)
		}

		payload := slack.Payload {
			Text: buffer.String(),
			Username: config.Review.Slack.UserName,
			IconEmoji: config.Review.Slack.Icon,
			Channel: config.Review.Slack.Channel,
			LinkNames: "1",
			Attachments: []slack.Attachment{slack.Attachment {}},
		}

		if context.Bool("dry-run") {
			fmt.Println("dry-run: Requests to %s.\n", dict["reviewers"])
		} else {
			if err := slack.Send(config.Review.Slack.WebhookUrl, "", payload); err != nil {
				fmt.Println(err)
			} else {
				fmt.Printf("SUCCESS: Requests to %s.\n", dict["reviewers"])
			}
		}

		return nil
	}

	app.Run(os.Args)
}
