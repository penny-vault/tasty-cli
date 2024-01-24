// Copyright 2021-2023
// SPDX-License-Identifier: Apache-2.0
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"

	gotasty "github.com/penny-vault/go-tasty"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var useSandbox bool

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:   "login <username>",
	Short: "Create a new Session in the tastytrade Open API",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// look up credentials in keychain
		password := accountCredentials(args[0], useSandbox)

		if password == "" {
			var err error
			password, err = promptPassword()
			if err != nil {
				log.Fatal().Err(err).Msg("prompt for password failed")
			}
		}

		// create a session
		session, err := gotasty.NewSession(args[0], password, gotasty.SessionOpts{
			RememberMe: true,
			Sandbox:    useSandbox,
		})
		if err != nil {
			log.Fatal().Err(err).Msg("create session failed")
		}

		// save session token and remember-me token in keychain
		if err := saveToKeychain(session); err != nil {
			log.Fatal().Err(err).Msg("failed saving session token to keychain")
		}

		fmt.Println("You are now logged in!")
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)

	loginCmd.Flags().BoolVar(&useSandbox, "sandbox", false, "Use the sandbox environment for testing")
}
