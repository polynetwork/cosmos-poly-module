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

package test

import (
	"encoding/hex"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/davecgh/go-spew/spew"
	"testing"
)

func Test_ParseCoins(t *testing.T) {

	args0 := "1000000000ont,1000000000000000000ong"
	//args1 := "1000000000stake,1000000000validatortoken"
	coins, err := sdk.ParseCoins(args0)
	if err != nil {
		t.Errorf("parsecoins error:%v", err)
	}

	spew.Printf("coins are %v\n", coins)
}

func Test_UnmarshalOperator(t *testing.T) {
	addr, err := sdk.AccAddressFromBech32("cosmos1ayc6faczpj42eu7wjsjkwcj7h0q2p2e4vrlkzf")
	fmt.Printf("addr in hex is %x\n", addr.Bytes())

	addr1, err := sdk.AccAddressFromHex(hex.EncodeToString(addr))
	if err != nil {
		t.Errorf("could not unmarshal result to sdk.AccAddress:%v", err)
	}
	fmt.Printf("opeartor are %s\n", addr1.String())

}


func Test_TestBytes(t *testing.T) {
	methodBs := make([]byte, 0)
	methodBs = returnBytes()
	fmt.Printf("methodBs = %v\n", methodBs)

}

func returnBytes() []byte {
	res := []byte{1, 2, 3}
	return res
}
