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

	"github.com/charmbracelet/bubbles/table"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

// transactionsCmd represents the accounts command
var transactionsCmd = &cobra.Command{
	Use:   "transactions [accountNumber]",
	Args:  cobra.ExactArgs(1),
	Short: "List transactions for given account",
	Run: func(cmd *cobra.Command, args []string) {
		session, err := sessionFromKeychain()
		if err != nil {
			log.Error().Err(err).Msg("failed creating session from keychain")
		}

		transactions, err := session.Transactions(args[0])
		if err != nil {
			log.Fatal().Err(err).Msg("failed to get accounts")
		}

		// save session token and remember-me token in keychain
		if err := saveToKeychain(session); err != nil {
			log.Fatal().Err(err).Msg("failed saving session token to keychain")
		}

		columns := []table.Column{
			{Title: "Date", Width: 10},
			{Title: "Type", Width: 15},
			{Title: "", Width: 15},
			{Title: "Action", Width: 15},
			{Title: "Symbol", Width: 7},
			{Title: "Quantity", Width: 8},
			{Title: "Value", Width: 15},
			{Title: "Order ID", Width: 10},
		}

		rows := make([]table.Row, 0, len(transactions))

		p := message.NewPrinter(language.English)

		for _, trx := range transactions {
			rows = append(rows, []string{
				trx.TransactionDate.Format("2006-01-02"),
				trx.TransactionType,
				trx.TransactionSubType,
				trx.Action.String(),
				trx.Symbol,
				p.Sprintf("%.0f", trx.Quantity),
				p.Sprintf("$%.2f", trx.Value),
				fmt.Sprint(trx.OrderID),
			})
		}

		showTable(columns, rows)
	},
}

func init() {
	rootCmd.AddCommand(transactionsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// transactionsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// transactionsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
