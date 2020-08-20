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
	polycommon "github.com/polynetwork/poly/common"
)

type TxArgs struct {
	FromAssetHash []byte
	ToAssetHash   []byte
	ToAddress     []byte
	Amount        *big.Int
	FeeAmount     *big.Int
	FeeAddress    []byte
	FromAddress   []byte
	Nonce         *big.Int
}

func (txargs *TxArgs) Serialization(sink *polycommon.ZeroCopySink, intLen int) error {
	sink.WriteVarBytes(txargs.FromAssetHash)
	sink.WriteVarBytes(txargs.ToAssetHash)
	sink.WriteVarBytes(txargs.ToAddress)
	paddedAmountBs, err := common.PadFixedBytes(txargs.Amount, intLen)
	if err != nil {
		return fmt.Errorf("serialization error:%v", err)
	}
	sink.WriteBytes(paddedAmountBs)
	paddedFeeAmountBs, err := common.PadFixedBytes(txargs.FeeAmount, intLen)
	if err != nil {
		return fmt.Errorf("serialization error:%v", err)
	}
	sink.WriteBytes(paddedFeeAmountBs)
	sink.WriteVarBytes(txargs.FeeAddress)
	sink.WriteVarBytes(txargs.FromAddress)
	paddedNonceBs, err := common.PadFixedBytes(txargs.Nonce, intLen)
	if err != nil {
		return fmt.Errorf("serialization error:%v", err)
	}
	sink.WriteBytes(paddedNonceBs)
	return nil
}

func (txargs *TxArgs) Deserialization(source *polycommon.ZeroCopySource, intLen int) error {
	fromAssetHash, eof := source.NextVarBytes()
	if eof {
		return fmt.Errorf("deserialize fromAssetHash error")
	}
	toAssetHash, eof := source.NextVarBytes()
	if eof {
		return fmt.Errorf("deserialize toAssetHash error")
	}
	toAddress, eof := source.NextVarBytes()
	if eof {
		return fmt.Errorf("deserialize ToAddress error")
	}
	paddedAmountBs, eof := source.NextBytes(uint64(intLen))
	if eof {
		return fmt.Errorf("deserialize Amount error")
	}
	amount, err := common.UnpadFixedBytes(paddedAmountBs, intLen)
	if err != nil {
		return fmt.Errorf("deserialization error:%v", err)
	}
	paddedFeeAmountBs, eof := source.NextBytes(uint64(intLen))
	if eof {
		return fmt.Errorf("deserialize FeeAmount error")
	}
	feeAmount, err := common.UnpadFixedBytes(paddedFeeAmountBs, intLen)
	if err != nil {
		return fmt.Errorf("deserialization error:%v", err)
	}
	feeAddress, eof := source.NextVarBytes()
	if eof {
		return fmt.Errorf("deserialize FeeAddress error")
	}
	fromAddress, eof := source.NextVarBytes()
	if eof {
		return fmt.Errorf("deserialize FromAddress error")
	}
	paddedNonceBs, eof := source.NextBytes(uint64(intLen))
	if eof {
		return fmt.Errorf("deserialize Nonce error")
	}
	nonce, err := common.UnpadFixedBytes(paddedNonceBs, intLen)
	if err != nil {
		return fmt.Errorf("deserialization error:%v", err)
	}

	txargs.FromAssetHash = fromAssetHash
	txargs.ToAssetHash = toAssetHash
	txargs.ToAddress = toAddress
	txargs.Amount = amount
	txargs.FeeAmount = feeAmount
	txargs.FeeAddress = feeAddress
	txargs.FromAddress = fromAddress
	txargs.Nonce = nonce
	return nil
}
