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
	"strconv"

	gotasty "github.com/penny-vault/go-tasty"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

// createOrderCmd represents the accounts command
var createOrderCmd = &cobra.Command{
	Use:   "create-order [accountNumber] [action buy/sell] [symbol] [quantity] [price]",
	Args:  cobra.ExactArgs(5),
	Short: "List positions for given account",
	Run: func(cmd *cobra.Command, args []string) {
		session, err := sessionFromKeychain()
		if err != nil {
			log.Error().Err(err).Msg("failed creating session from keychain")
		}

		action := gotasty.BuyToOpen
		priceEffect := gotasty.Debit
		if args[1] == "sell" {
			action = gotasty.SellToClose
			priceEffect = gotasty.Credit
		}

		price, err := strconv.ParseFloat(args[4], 64)
		if err != nil {
			log.Error().Err(err).Str("PriceRaw", args[4]).Msg("could not parse price value")
		}

		quantity, err := strconv.ParseInt(args[3], 10, 64)
		if err != nil {
			log.Error().Err(err).Str("QuantityRaw", args[3]).Msg("could not parse quantity value")
		}

		order := &gotasty.Order{
			TimeInForce: gotasty.Day,
			OrderType:   gotasty.Limit,
			Price:       price,
			PriceEffect: priceEffect,
			Legs: []*gotasty.Leg{
				{
					InstrumentType: gotasty.Equity,
					Symbol:         args[2],
					Quantity:       quantity,
					Action:         action,
				},
			},
		}

		orderResp, err := session.SubmitOrder(args[0], order)
		if err != nil {
			log.Fatal().Err(err).Msg("failed to get accounts")
		}

		fmt.Printf("%#v\n", orderResp)

		// save session token and remember-me token in keychain
		if err := saveToKeychain(session); err != nil {
			log.Fatal().Err(err).Msg("failed saving session token to keychain")
		}
	},
}

func init() {
	rootCmd.AddCommand(createOrderCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// createOrderCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// createOrderCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
