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
	"github.com/charmbracelet/bubbles/table"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

// accountsCmd represents the accounts command
var accountsCmd = &cobra.Command{
	Use:   "accounts",
	Short: "List accounts accessible by customer",
	Run: func(cmd *cobra.Command, args []string) {
		session, err := sessionFromKeychain()
		if err != nil {
			log.Error().Err(err).Msg("failed creating session from keychain")
		}

		accounts, err := session.Accounts()
		if err != nil {
			log.Fatal().Err(err).Msg("failed to get accounts")
		}

		// save session token and remember-me token in keychain
		if err := saveToKeychain(session); err != nil {
			log.Fatal().Err(err).Msg("failed saving session token to keychain")
		}

		columns := []table.Column{
			{Title: "Account #", Width: 10},
			{Title: "Nickname", Width: 15},
			{Title: "Type", Width: 10},
			{Title: "Margin", Width: 10},
			{Title: "Access Level", Width: 12},
		}

		rows := make([]table.Row, 0, len(accounts))

		for _, acct := range accounts {
			rows = append(rows, []string{acct.AccountNumber, acct.Nickname, acct.AccountType, acct.MarginOrCash, acct.AuthorityLevel})
		}

		showTable(columns, rows)
	},
}

func init() {
	rootCmd.AddCommand(accountsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// accountsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// accountsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
