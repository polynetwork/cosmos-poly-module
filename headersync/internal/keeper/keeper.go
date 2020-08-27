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

package keeper

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/polynetwork/cosmos-poly-module/headersync/internal/types"
	polycommon "github.com/polynetwork/poly/common"
	vconfig "github.com/polynetwork/poly/consensus/vbft/config"
	polysig "github.com/polynetwork/poly/core/signature"
	polytype "github.com/polynetwork/poly/core/types"
	"github.com/polynetwork/poly/merkle"
)

// Keeper of the mint store
type Keeper struct {
	cdc      *codec.Codec
	storeKey sdk.StoreKey
}

// NewKeeper creates a new mint Keeper instance
func NewKeeper(
	cdc *codec.Codec, key sdk.StoreKey) Keeper {
	return Keeper{
		cdc:      cdc,
		storeKey: key,
	}
}

func (keeper Keeper) SyncGenesisHeader(ctx sdk.Context, genesisHeaderStr string) error {
	genesisHeader := &polytype.Header{}

	genesisHeaderBytes, err := hex.DecodeString(genesisHeaderStr)
	if err != nil {
		return types.ErrSyncGenesisHeader(fmt.Sprintf("hex.DecodeString error: %s", err.Error()))
	}
	source := polycommon.NewZeroCopySource(genesisHeaderBytes)
	if err := genesisHeader.Deserialization(source); err != nil {
		return types.ErrDeserializeHeader(err)
	}
	if consensusPeer, _ := keeper.GetConsensusPeers(ctx, genesisHeader.ChainID); consensusPeer != nil {
		return types.ErrSyncGenesisHeader(fmt.Sprintf("Genesis Header already synced, ConsensusPeers exists: %s", consensusPeer.String()))
	}
	if err := keeper.UpdateConsensusPeer(ctx, genesisHeader); err != nil {
		return err
	}
	// Make sure the header contains poly.NewChainConfig info
	if _, err := keeper.GetConsensusPeers(ctx, genesisHeader.ChainID); err != nil {
		return types.ErrSyncGenesisHeader(fmt.Sprintf("After UpdteConsensusPeer, Get Consensus Peers error: %v", err))
	}
	return nil
}

func (keeper Keeper) SyncBlockHeaders(ctx sdk.Context, headerStrs []string) error {
	for _, headerStr := range headerStrs {
		header := &polytype.Header{}
		headerBs, err := hex.DecodeString(headerStr)
		if err != nil {
			return types.ErrSyncBlockHeader("Decode header string to bytes", 0, 0, err)
		}
		source := polycommon.NewZeroCopySource(headerBs)
		if err := header.Deserialization(source); err != nil {
			return types.ErrDeserializeHeader(err)
		}
		if err := keeper.ProcessHeader(ctx, header, nil, nil); err != nil {
			return types.ErrSyncBlockHeader("ProcessHeader", header.ChainID, header.Height, err)
		}
	}
	return nil
}

func (keeper Keeper) ProcessHeader(ctx sdk.Context, header *polytype.Header, headerProof []byte, curHeader *polytype.Header) error {
	// header to be checked if containing valid NewChainConfig
	var cpHeader *polytype.Header
	if curHeader == nil || headerProof == nil {
		if err := keeper.VerifyHeaderSig(ctx, header); err != nil {
			if err := keeper.VerifyHeaderByKeyHeaderHash(ctx, header); err == nil {
				return nil
			}
			return err
		}
		cpHeader = header
	} else {
		if err := keeper.VerifyHistoricalHeader(ctx, header, headerProof, curHeader); err != nil {
			return err
		}
		cpHeader = curHeader
	}

	if err := keeper.UpdateConsensusPeer(ctx, cpHeader); err != nil {
		return err
	}
	return nil
}

func (keeper Keeper) VerifyHeaderSig(ctx sdk.Context, header *polytype.Header) error {
	consensusPeer, err := keeper.GetConsensusPeers(ctx, header.ChainID)
	if err != nil {
		return types.ErrSyncBlockHeader("GetConsensusPeer", header.ChainID, header.Height, err)
	}
	if header.Height <= consensusPeer.Height {
		return types.ErrSyncBlockHeader("Compare height", header.ChainID, header.Height,
			fmt.Errorf("Stored consensus header.Height: %d, trying to sync height:%d", consensusPeer.Height, header.Height))
	}

	if len(header.Bookkeepers)*3 < len(consensusPeer.PeerMap)*2 {
		return types.ErrBookKeeperNum(len(header.Bookkeepers), len(consensusPeer.PeerMap))
	}
	for _, bookkeeper := range header.Bookkeepers {
		pubkey := vconfig.PubkeyID(bookkeeper)
		_, present := consensusPeer.PeerMap[pubkey]
		if !present {
			return types.ErrInvalidPublicKey(pubkey)
		}
	}
	hash := header.Hash()
	if e := polysig.VerifyMultiSignature(hash[:], header.Bookkeepers, len(header.Bookkeepers), header.SigData); e != nil {
		return types.ErrVerifyMultiSigFail(err, header.Height)
	}
	return nil
}
func (keeper Keeper) VerifyHeaderByKeyHeaderHash(ctx sdk.Context, header *polytype.Header) error {
	headerHash := header.Hash()
	keyHeaderHash, err := keeper.GetKeyHeaderHash(ctx, header.ChainID)
	if err != nil {
		return fmt.Errorf("VerifyHeaderByKeyHeaderHash, GetKeyHeaderHash Error: %s", err.Error())
	}
	if headerHash == *keyHeaderHash {
		return nil
	}
	return fmt.Errorf("VerifyHeaderByKeyHeaderHash, not equal, expect: %s, got: %s", keyHeaderHash.ToArray(), headerHash.ToArray())
}

func (keeper Keeper) UpdateConsensusPeer(ctx sdk.Context, header *polytype.Header) error {

	blkInfo := &vconfig.VbftBlockInfo{}
	if err := json.Unmarshal(header.ConsensusPayload, blkInfo); err != nil {
		return types.ErrUnmarshalSpecificTypeFail(blkInfo, err)
	}
	if blkInfo.NewChainConfig != nil {
		consensusPeers := &types.ConsensusPeers{
			ChainID: header.ChainID,
			Height:  header.Height,
			PeerMap: make(map[string]*types.Peer),
		}
		for _, p := range blkInfo.NewChainConfig.Peers {
			consensusPeers.PeerMap[p.ID] = &types.Peer{Index: p.Index, PeerPubkey: p.ID}
		}
		if err := keeper.SetConsensusPeers(ctx, *consensusPeers); err != nil {
			return err
		}
		if err := keeper.SetKeyHeaderHash(ctx, consensusPeers.ChainID, header.Hash()); err != nil {
			return err
		}
	}
	return nil
}

func (keeper Keeper) SetConsensusPeers(ctx sdk.Context, consensusPeers types.ConsensusPeers) error {
	store := ctx.KVStore(keeper.storeKey)
	sink := polycommon.NewZeroCopySink(nil)
	consensusPeers.Serialization(sink)
	store.Set(GetConsensusPeerKey(consensusPeers.ChainID), sink.Bytes())
	return nil
}

func (keeper Keeper) GetConsensusPeers(ctx sdk.Context, chainId uint64) (*types.ConsensusPeers, error) {
	store := ctx.KVStore(keeper.storeKey)
	consensusPeerBytes := store.Get(GetConsensusPeerKey(chainId))
	if consensusPeerBytes == nil {
		return nil, types.ErrGetConsensusPeers(chainId)
	}
	consensusPeers := new(types.ConsensusPeers)
	if err := consensusPeers.Deserialization(polycommon.NewZeroCopySource(consensusPeerBytes)); err != nil {
		return nil, types.ErrDeserializeConsensusPeer(err)
	}
	return consensusPeers, nil
}

func (keeper Keeper) SetKeyHeaderHash(ctx sdk.Context, chainId uint64, keyHeaderHash polycommon.Uint256) error {
	store := ctx.KVStore(keeper.storeKey)
	store.Set(GetKeyHeaderHashKey(chainId), keyHeaderHash.ToArray())
	return nil
}

func (keeper Keeper) GetKeyHeaderHash(ctx sdk.Context, chainId uint64) (*polycommon.Uint256, error) {
	store := ctx.KVStore(keeper.storeKey)
	headerHashBs := store.Get(GetKeyHeaderHashKey(chainId))
	if headerHashBs == nil {
		return nil, types.ErrGetKeyHeaderHash(fmt.Sprintf("Empty key header hash with chainId: %d", chainId))
	}
	headerHash, err := polycommon.Uint256ParseFromBytes(headerHashBs)
	if err != nil {
		return nil, types.ErrGetKeyHeaderHash(fmt.Sprintf("Error: Uint256 from bytes: %x with chainId: %d", headerHashBs, err))
	}
	return &headerHash, nil
}

func (keeper Keeper) VerifyHistoricalHeader(ctx sdk.Context, header *polytype.Header, headerProof []byte, curHeader *polytype.Header) error {
	if err := keeper.VerifyHeaderSig(ctx, curHeader); err != nil {
		return err
	}
	value, err := merkle.MerkleProve(headerProof, curHeader.BlockRoot[:])
	if err != nil {
		return fmt.Errorf("VerifyHistoricalHeader, MerkleProve error: %s", err.Error())
	}
	hashToBeVerified := header.Hash()
	if !bytes.Equal(value, hashToBeVerified[:]) {
		return fmt.Errorf("VerifyHistoricalHeader error, historical header height: %d, current epoch header height: %d, expect: %x, got: %x", header.Height, curHeader.Height, hashToBeVerified[:], value)
	}
	return nil
}
