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
	"errors"
	"fmt"
	"strings"
	"syscall"

	"github.com/keybase/go-keychain"
	gotasty "github.com/penny-vault/go-tasty"
	"github.com/rs/zerolog/log"
	"golang.org/x/term"
)

func accountCredentials(account string, sandbox bool) string {
	log.Info().Bool("Sandbox", sandbox).Str("Username", account).Msg("fetching credentials")
	query := keychain.NewItem()
	query.SetSecClass(keychain.SecClassGenericPassword)
	if sandbox {
		query.SetService("tastytrade-api-sandbox")
	} else {
		query.SetService("tastytrade-api")
	}
	query.SetAccount(account)
	query.SetMatchLimit(keychain.MatchLimitOne)
	query.SetReturnData(true)
	results, err := keychain.QueryItem(query)
	if err != nil {
		// Error
		log.Error().Err(err).Msg("error encountered when querying keychain")
	} else {
		for _, r := range results {
			return string(r.Data)
		}
	}

	return ""
}

func sessionFromKeychain() (*gotasty.Session, error) {
	query := keychain.NewItem()
	query.SetSecClass(keychain.SecClassGenericPassword)
	query.SetService("tastytrade-api")
	query.SetAccount("session")
	query.SetMatchLimit(keychain.MatchLimitOne)
	query.SetReturnAttributes(true)
	query.SetReturnData(true)
	results, err := keychain.QueryItem(query)
	if err != nil {
		// Error
		return nil, err
	}

	for _, r := range results {
		sessionData := r.Data
		return gotasty.NewSessionFromBytes(sessionData)
	}

	return nil, errors.New("session data doesn't exist in keychain -- call login first")
}

func promptPassword() (string, error) {
	fmt.Print("Enter Password: ")
	bytePassword, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return "", err
	}

	password := string(bytePassword)
	return strings.TrimSpace(password), nil
}

func saveToKeychain(session *gotasty.Session) error {
	// Save session token
	item := keychain.NewItem()
	item.SetSecClass(keychain.SecClassGenericPassword)
	item.SetService("tastytrade-api")
	item.SetAccount("session")
	data, err := session.Marshal()
	if err != nil {
		return err
	}

	item.SetData(data)
	item.SetSynchronizable(keychain.SynchronizableNo)
	item.SetAccessible(keychain.AccessibleWhenUnlocked)
	err = keychain.AddItem(item)
	if err == keychain.ErrorDuplicateItem {
		if err2 := keychain.DeleteItem(item); err2 != nil {
			return err2
		}
		if err2 := keychain.AddItem(item); err2 != nil {
			return err2
		}
	}

	return nil
}
