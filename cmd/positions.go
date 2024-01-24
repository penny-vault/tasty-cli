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

// positionsCmd represents the accounts command
var positionsCmd = &cobra.Command{
	Use:   "positions [accountNumber]",
	Args:  cobra.ExactArgs(1),
	Short: "List positions for given account",
	Run: func(cmd *cobra.Command, args []string) {
		session, err := sessionFromKeychain()
		if err != nil {
			log.Error().Err(err).Msg("failed creating session from keychain")
		}

		positions, err := session.Positions(args[0])
		if err != nil {
			log.Fatal().Err(err).Msg("failed to get accounts")
		}

		// save session token and remember-me token in keychain
		if err := saveToKeychain(session); err != nil {
			log.Fatal().Err(err).Msg("failed saving session token to keychain")
		}

		columns := []table.Column{
			{Title: "Symbol", Width: 10},
			{Title: "Quantity", Width: 9},
			{Title: "Value", Width: 15},
			{Title: "P/L", Width: 15},
			{Title: "Type", Width: 10},
			{Title: "Created At", Width: 10},
			{Title: "Realized", Width: 10},
		}

		rows := make([]table.Row, 0, len(positions))

		p := message.NewPrinter(language.English)

		for _, pos := range positions {
			rows = append(rows, []string{
				pos.Symbol,
				fmt.Sprint(pos.Quantity),
				p.Sprintf("$%.2f", pos.Quantity*pos.ClosePrice),
				p.Sprintf("$%.2f", (pos.Quantity*pos.ClosePrice)-(pos.Quantity*pos.AverageOpenPrice)),
				pos.InstrumentType,
				pos.CreatedAt.Format("2006-01-02"),
				fmt.Sprintf("%.2f", pos.RealizedDayGain),
			})
		}

		showTable(columns, rows)
	},
}

func init() {
	rootCmd.AddCommand(positionsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// positionsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// positionsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
