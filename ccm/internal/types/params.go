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
	"github.com/cosmos/cosmos-sdk/x/params"
)

// Parameter store keys
var (
	KeyCurrentChainIdForPolyChain = []byte("ChainIdForPolyChain")
)

type Params struct {
	ChainIdInPolyNet uint64 `json:"chain_id_in_poly_net" yaml:"chain_id_in_poly_net"` // chain id of current cosmos chain for cross chain in poly chain network
}

// ParamTable for ccm module.
func ParamKeyTable() params.KeyTable {
	return params.NewKeyTable().RegisterParamSet(&Params{})
}

// default ccm module parameters
func DefaultParams() Params {
	return Params{
		ChainIdInPolyNet: 0,
	}
}

// validate params
func (p Params) Validate() error {
	if err := validateChainId(p.ChainIdInPolyNet); err != nil {
		return err
	}
	return nil
}

func validateChainId(i interface{}) error {
	v, ok := i.(uint64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	if v == 0 {
		return fmt.Errorf("ChainId for Poly Chain Network needs to be configured manually, and non-zero, got: %d", v)
	}
	return nil
}

func (p Params) String() string {
	return fmt.Sprintf(`Ccm Params:
  Current CrossChainId:             %d
`,
		p.ChainIdInPolyNet,
	)
}

// Implements params.ParamSet
func (p *Params) ParamSetPairs() params.ParamSetPairs {
	return params.ParamSetPairs{
		params.NewParamSetPair(KeyCurrentChainIdForPolyChain, &p.ChainIdInPolyNet, validateChainId),
	}
}
