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
	"math/big"

	"github.com/polynetwork/cosmos-poly-module/common"
	polycommon "github.com/polynetwork/cosmos-poly-module/headersync/poly-utils/common"
)

type TxArgs struct {
	FromAssetHash []byte
	ToAssetHash   []byte
	ToAddress     []byte
	Amount        *big.Int
	FeeAmount     *big.Int
	FeeAddress    []byte
}

func (this *TxArgs) Serialization(sink *polycommon.ZeroCopySink, intLen int) error {
	sink.WriteVarBytes(this.FromAssetHash)
	sink.WriteVarBytes(this.ToAssetHash)
	sink.WriteVarBytes(this.ToAddress)
	paddedAmountBs, err := common.Pad32Bytes(this.Amount, intLen)
	if err != nil {
		return fmt.Errorf("TxArgs Serialization error:%v", err)
	}
	sink.WriteBytes(paddedAmountBs)
	paddedFeeAmountBs, err := common.Pad32Bytes(this.FeeAmount, intLen)
	if err != nil {
		return fmt.Errorf("TxArgs Serialization error:%v", err)
	}
	sink.WriteBytes(paddedFeeAmountBs)
	sink.WriteVarBytes(this.FeeAddress)
	return nil
}

func (this *TxArgs) Deserialization(source *polycommon.ZeroCopySource, intLen int) error {
	fromAssetHash, eof := source.NextVarBytes()
	if eof {
		return fmt.Errorf("TxArgs deserialize fromAssetHash error")
	}
	toAssetHash, eof := source.NextVarBytes()
	if eof {
		return fmt.Errorf("TxArgs deserialize toAssetHash error")
	}
	toAddress, eof := source.NextVarBytes()
	if eof {
		return fmt.Errorf("TxArgs deserialize ToAddress error")
	}
	paddedAmountBs, eof := source.NextBytes(uint64(intLen))
	if eof {
		return fmt.Errorf("TxArgs deserialize Amount error")
	}
	amount, err := common.Unpad32Bytes(paddedAmountBs, intLen)
	if err != nil {
		return fmt.Errorf("TxArgs Deserialization error:%v", err)
	}
	paddedFeeAmountBs, eof := source.NextBytes(uint64(intLen))
	if eof {
		return fmt.Errorf("TxArgs deserialize FeeAmount error")
	}
	feeAmount, err := common.Unpad32Bytes(paddedFeeAmountBs, intLen)
	if err != nil {
		return fmt.Errorf("TxArgs Deserialization error:%v", err)
	}
	feeAddress, eof := source.NextVarBytes()
	if eof {
		return fmt.Errorf("TxArgs deserialize FeeAddress error")
	}

	this.FromAssetHash = fromAssetHash
	this.ToAssetHash = toAssetHash
	this.ToAddress = toAddress
	this.Amount = amount
	this.FeeAmount = feeAmount
	this.FeeAddress = feeAddress
	return nil
}
