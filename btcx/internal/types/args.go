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
	polycommon "github.com/polynetwork/cosmos-poly-module/headersync/poly-utils/common"
)

type ToBTCArgs struct {
	ToBtcAddress []byte
	Amount       uint64
	RedeemScript []byte
}

func (this *ToBTCArgs) Serialization(sink *polycommon.ZeroCopySink) error {
	sink.WriteVarBytes(this.ToBtcAddress)
	sink.WriteUint64(this.Amount)
	sink.WriteVarBytes(this.RedeemScript)
	return nil
}

func (this *ToBTCArgs) Deserialization(source *polycommon.ZeroCopySource) error {
	toBtcAddress, eof := source.NextVarBytes()
	if eof {
		return fmt.Errorf("ToBTCArgs deserialize toBtcAddress error")
	}
	amt, eof := source.NextUint64()
	if eof {
		return fmt.Errorf("ToBTCArgs deserialize Amount error")
	}
	redeemScript, eof := source.NextVarBytes()
	if eof {
		return fmt.Errorf("ToBTCArgs deserialize redeemScript error")
	}

	this.ToBtcAddress = toBtcAddress
	this.Amount = amt
	this.RedeemScript = redeemScript
	return nil
}

type BTCArgs struct {
	ToBtcAddress []byte
	Amount       uint64
}

func (this *BTCArgs) Serialization(sink *polycommon.ZeroCopySink) error {
	sink.WriteVarBytes(this.ToBtcAddress)
	sink.WriteUint64(this.Amount)
	return nil
}

func (this *BTCArgs) Deserialization(source *polycommon.ZeroCopySource) error {
	toBtcAddress, eof := source.NextVarBytes()
	if eof {
		return fmt.Errorf("ToBTCArgs deserialize toBtcAddress error")
	}
	amt, eof := source.NextUint64()
	if eof {
		return fmt.Errorf("ToBTCArgs deserialize Amount error")
	}
	this.ToBtcAddress = toBtcAddress
	this.Amount = amt
	return nil
}
