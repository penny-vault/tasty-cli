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

// ordersCmd represents the accounts command
var ordersCmd = &cobra.Command{
	Use:   "list-orders [accountNumber]",
	Args:  cobra.ExactArgs(1),
	Short: "List positions for given account",
	Run: func(cmd *cobra.Command, args []string) {
		session, err := sessionFromKeychain()
		if err != nil {
			log.Error().Err(err).Msg("failed creating session from keychain")
		}

		orders, err := session.Orders(args[0])
		if err != nil {
			log.Fatal().Err(err).Msg("failed to get accounts")
		}

		// save session token and remember-me token in keychain
		if err := saveToKeychain(session); err != nil {
			log.Fatal().Err(err).Msg("failed saving session token to keychain")
		}

		columns := []table.Column{
			{Title: "ID", Width: 10},
			{Title: "Received At", Width: 10},
			{Title: "Order Type", Width: 15},
			{Title: "Time in force", Width: 15},
			{Title: "Symbol", Width: 15},
			{Title: "Value", Width: 15},
			{Title: "Status", Width: 10},
		}

		rows := make([]table.Row, 0, len(orders))

		p := message.NewPrinter(language.English)

		for _, order := range orders {
			rows = append(rows, []string{
				order.ID,
				order.ReceivedAt.Format("2006-01-02"),
				order.OrderType.String(),
				order.TimeInForce.String(),
				order.Legs[0].Symbol,
				p.Sprintf("$%.2f", order.Value),
				order.Status,
			})
		}

		showTable(columns, rows)
	},
}

func init() {
	rootCmd.AddCommand(ordersCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// ordersCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// ordersCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
