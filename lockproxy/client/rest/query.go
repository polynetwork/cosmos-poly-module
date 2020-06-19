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

package rest

import (
	"encoding/hex"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/polynetwork/cosmos-poly-module/lockproxy/client/common"
	"strconv"
)

func registerQueryRoutes(cliCtx context.CLIContext, r *mux.Router, queryRoute string) {
	r.HandleFunc(
		fmt.Sprintf("/lockproxy/proxy_hash_by_operator/{%s}", Operator),
		queryProxyHashByOperatorHandlerFn(cliCtx, queryRoute),
	).Methods("GET")

	r.HandleFunc(
		fmt.Sprintf("/lockproxy/proxy_hash/{%s}/{%s}", LockProxyHash, ToChainId),
		queryProxyHashHandlerFn(cliCtx, queryRoute),
	).Methods("GET")

	r.HandleFunc(
		fmt.Sprintf("/lockproxy/asset_hash/{%s}/{%s}/{%s}", LockProxyHash, AssetDenom, ToChainId),
		queryAssetHashHandlerFn(cliCtx, queryRoute),
	).Methods("GET")

}

func queryProxyHashByOperatorHandlerFn(cliCtx context.CLIContext, queryRoute string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}
		operatorAddr, err := sdk.AccAddressFromBech32(mux.Vars(r)[Operator])
		if err != nil {
			return
		}
		res, ok := checkResponseQueryProxyHashByOperatorResponse(w, cliCtx, queryRoute, operatorAddr)
		if !ok {
			return
		}

		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func checkResponseQueryProxyHashByOperatorResponse(
	w http.ResponseWriter, cliCtx context.CLIContext, queryRoute string, operator sdk.AccAddress) (res []byte, ok bool) {

	res, err := common.QueryProxyByOperator(cliCtx, queryRoute, operator)
	if err != nil {
		rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
		return nil, false
	}

	return res, true
}

func queryProxyHashHandlerFn(cliCtx context.CLIContext, queryRoute string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}
		// query for rewards from a particular delegator
		vars := mux.Vars(r)
		lockproxy, err := hex.DecodeString(vars[LockProxyHash])
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		chainId, err := strconv.ParseUint(vars[ToChainId], 10, 64)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		res, ok := checkResponseQueryProxyHashResponse(w, cliCtx, queryRoute, lockproxy, chainId)
		if !ok {
			return
		}

		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func checkResponseQueryProxyHashResponse(
	w http.ResponseWriter, cliCtx context.CLIContext, queryRoute string, lockproxy []byte, chainId uint64) (res []byte, ok bool) {

	res, err := common.QueryProxyHash(cliCtx, queryRoute, lockproxy, chainId)
	if err != nil {
		rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
		return nil, false
	}

	return res, true
}

func queryAssetHashHandlerFn(cliCtx context.CLIContext, queryRoute string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}
		// query for rewards from a particular delegator
		vars := mux.Vars(r)
		lockproxy, err := hex.DecodeString(vars[LockProxyHash])
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		chainId, err := strconv.ParseUint(vars[ToChainId], 10, 64)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		res, ok := checkResponseQueryAssetHashResponse(w, cliCtx, queryRoute, lockproxy, vars[AssetDenom], chainId)
		if !ok {
			return
		}

		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func checkResponseQueryAssetHashResponse(
	w http.ResponseWriter, cliCtx context.CLIContext, queryRoute string, lockproxy []byte, denom string, chainId uint64) (res []byte, ok bool) {

	res, err := common.QueryAssetHash(cliCtx, queryRoute, lockproxy, denom, chainId)
	if err != nil {
		rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
		return nil, false
	}

	return res, true
}
