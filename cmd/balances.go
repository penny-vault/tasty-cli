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
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

// balanceCmd represents the accounts command
var balanceCmd = &cobra.Command{
	Use:   "balances",
	Short: "List balances for given account",
	Run: func(cmd *cobra.Command, args []string) {
		session, err := sessionFromKeychain()
		if err != nil {
			log.Error().Err(err).Msg("failed creating session from keychain")
		}

		accounts, err := session.Accounts()
		if err != nil {
			log.Fatal().Err(err).Msg("failed to get accounts")
		}

		columns := []table.Column{
			{Title: "Updated At", Width: 10},
			{Title: "Cash Balance", Width: 15},
			{Title: "Net Liquidating Value", Width: 15},
			{Title: "Long Equity Value", Width: 15},
			{Title: "Equity Buying Power", Width: 15},
			{Title: "Pending Cash", Width: 15},
		}

		rows := make([]table.Row, 0, len(accounts))
		for _, acct := range accounts {
			balance, err := session.Balance(acct.AccountNumber)
			if err != nil {
				log.Fatal().Err(err).Msg("failed to get accounts")
			}

			p := message.NewPrinter(language.English)

			rows = append(rows, []string{
				balance.UpdatedAt.Format("2006-01-02"),
				p.Sprintf("$%.2f", balance.CashBalance),
				p.Sprintf("$%.2f", balance.NetLiquidatingValue),
				p.Sprintf("$%.2f", balance.LongEquityValue),
				p.Sprintf("$%.2f", balance.EquityBuyingPower),
				p.Sprintf("$%.2f", balance.PendingCash),
			})

		}

		// save session token and remember-me token in keychain
		if err := saveToKeychain(session); err != nil {
			log.Fatal().Err(err).Msg("failed saving session token to keychain")
		}

		showTable(columns, rows)
	},
}

func init() {
	rootCmd.AddCommand(balanceCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// balanceCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// balanceCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
