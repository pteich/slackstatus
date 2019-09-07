// Copyright © 2016 Peter Teich <mail@pteich.xyz>
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
	"github.com/pteich/slackstatus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"os"
)

var cfgFile string

var RootCmd = &cobra.Command{
	Use:   "slackstatus \"your message here\"",
	Short: "Post a formatted status message to Slack",
	Long: `Post a formatted status message to a Slack channel with a given
username. An optional footer info can be provided.
You can set one of the predifined colors DANGER, WARNING or GOOD or use any hex color value.

You need to set up an incoming webhook for your Slack at https://my.slack.com/services/new/incoming-webhook/ first.
	`,
	Run: func(cmd *cobra.Command, args []string) {

		message := getPipedInput()
		if len(args) > 1 {
			message = args[0]
		}

		if len(message) == 0 {
			cmd.Help()
			log.Fatal("You need to provide a message as argument or pipe a message")
		}

		var slackmsg = slackstatus.Message{
			WebhookURL: viper.GetString("webhook"),
			Username:   viper.GetString("user"),
			Channel:    viper.GetString("channel"),
			IconEmoji:  viper.GetString("iconemoji"),
			Footer:     viper.GetString("footer"),
		}

		if err := slackmsg.Send(message, viper.GetString("color")); err != nil {
			log.Fatal(err)
		}

	},
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.slackstatus.yaml)")

	RootCmd.PersistentFlags().String("webhook", "", "Slack webhook URL (required)")
	RootCmd.MarkFlagRequired("webhook")

	RootCmd.PersistentFlags().String("user", "slackstatus", "Username to use when posting")
	RootCmd.PersistentFlags().String("channel", "#slackstatus", "Channel to post into")
	RootCmd.PersistentFlags().String("iconemoji", ":monkey_face:", "Icon emoji to use whan posting")
	RootCmd.PersistentFlags().String("footer", "", "Footer text")
	RootCmd.PersistentFlags().String("color", "good", "Color for message bar")

	viper.BindPFlags(RootCmd.PersistentFlags())

	log.SetFlags(0)
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	}

	viper.SetConfigName(".slackstatus")
	viper.AddConfigPath(os.Getenv("HOME"))
	viper.AddConfigPath(".")
	viper.SetEnvPrefix("slackstatus")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

func getPipedInput() string {

	fileInput, err := os.Stdin.Stat()
	if err != nil {
		log.Fatal(err)
	}

	var text string
	if fileInput.Mode()&os.ModeNamedPipe != 0 {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			text += scanner.Text()
		}
	}

	return text
}
