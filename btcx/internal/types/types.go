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
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type DenomInfo struct {
	Creator          string
	Denom            string
	AssetHash        string
	TotalSupply      sdk.Int
	RedeemScipt      string
	RedeemScriptHash string
}

func (msg DenomInfo) String() string {
	return fmt.Sprintf(`
  Creator:        	 			%s
  Denom: 						%s
  AssetHash:					%s
  TotalSupply:		 			%s
  RedeemScriptHash(AssetHash):  %s
  RedeemScipt: 					%s
`, msg.Creator, msg.Denom, msg.AssetHash, msg.TotalSupply.String(), msg.RedeemScriptHash, msg.RedeemScipt)
}

type DenomCrossChainInfo struct {
	DenomInfo
	ToChainId   uint64
	ToAssetHash string
}

func (msg DenomCrossChainInfo) String() string {
	return msg.DenomInfo.String() + fmt.Sprintf(`
  ToChainId:       	 			%d
  ToAssetHash:		 			%s
`, msg.ToChainId, msg.ToAssetHash)
}
