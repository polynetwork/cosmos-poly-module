/*
 * Copyright (C) 2020 The poly network Authors
 * This file is part of The poly network library.
 *
 * The  poly network  is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Lesser General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * The  poly network  is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Lesser General Public License for more details.
 * You should have received a copy of the GNU Lesser General Public License
 * along with The poly network .  If not, see <http://www.gnu.org/licenses/>.
 */

package types

import (
	"fmt"
)

// QueryBalanceParams defines the params for querying an account balance.
type QueryContainToContractParam struct {
	KeyStore       string
	ToContractAddr []byte
	FromChainId    uint64
}

func NewQueryContainToContractParam(keystore string, toContractAddr []byte, fromChainId uint64) QueryContainToContractParam {
	return QueryContainToContractParam{keystore, toContractAddr, fromChainId}
}

type QueryContainToContractRes struct {
	KeyStore string
	Exist    bool
	Info     string
}

func (this QueryContainToContractRes) String() string {
	return fmt.Sprintf(`
  KeyStore:				%s,
  Exist:				%t,
  Info:					%s,
`, this.KeyStore, this.Exist, this.Info)
}

type QueryCcmParams struct {
}

func NewQueryCcmParams() QueryCcmParams {
	return QueryCcmParams{}
}
