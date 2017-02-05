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

func main () {
	app := cli.NewApp()
	app.Name = "review"
	app.Usage = "Requests Code review to team member from terminal."
	app.Version = "0.0.1"

	app.Action = func(context *cli.Context) error {
		webhookUrl := "https://xxxx"
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
