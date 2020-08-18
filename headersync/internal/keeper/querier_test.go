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

package keeper_test

import (
	"encoding/hex"
	"fmt"
	"github.com/polynetwork/cosmos-poly-module/headersync"
	polycommon "github.com/polynetwork/poly/common"
	polytype "github.com/polynetwork/poly/core/types"
	"github.com/stretchr/testify/assert"
	"testing"

	"github.com/stretchr/testify/require"

	keep "github.com/polynetwork/cosmos-poly-module/headersync/internal/keeper"
	"github.com/polynetwork/cosmos-poly-module/headersync/internal/types"

	abci "github.com/tendermint/tendermint/abci/types"
)

func TestN_headersync_Querier(t *testing.T) {
	app, ctx := createTestApp(true)

	h0s, _ := hex.DecodeString(header0)
	h0 := new(polytype.Header)
	err := h0.Deserialization(polycommon.NewZeroCopySource(h0s))
	assert.Nil(t, err)

	err = app.HeaderSyncKeeper.SyncGenesisHeader(ctx, header0)
	assert.Nil(t, err, "Sync genesis header fail")

	querier := keep.NewQuerier(app.HeaderSyncKeeper)
	// query synced Poly Chain consensus peer
	query := abci.RequestQuery{
		Path: fmt.Sprintf("custom/%s/%s", headersync.StoreKey, types.QueryConsensusPeers),
		Data: app.Codec().MustMarshalJSON(types.NewQueryConsensusPeersParams(h0.ChainID)),
	}
	cpBs, err := querier(ctx, []string{types.QueryConsensusPeers}, query)
	require.NoError(t, err)

	consensusPeersBs, err := ExtractChainConfig(h0)
	assert.Nil(t, err)
	require.Equal(t, consensusPeersBs, cpBs, "Synced consensus 0 is not equal to the querier result")
}
