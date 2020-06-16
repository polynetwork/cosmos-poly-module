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
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/polynetwork/cosmos-poly-module/ft/client/common"
)

func registerQueryRoutes(cliCtx context.CLIContext, r *mux.Router, queryRoute string) {
	r.HandleFunc(
		fmt.Sprintf("/ft/denom_info/{%s}", Denom),
		queryDemonHandlerFn(cliCtx, queryRoute),
	).Methods("GET")

	r.HandleFunc(
		fmt.Sprintf("/ft/denom_cc_info/{%s}/{%s}", Denom, ChainId),
		queryDemonCrossChainInfoHandlerFn(cliCtx, queryRoute),
	).Methods("GET")

}

func queryDemonHandlerFn(cliCtx context.CLIContext, queryRoute string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}
		// query for rewards from a particular delegator
		vars := mux.Vars(r)
		denom := vars[Denom]

		res, ok := checkResponseQueryDenomInfoResponse(w, cliCtx, queryRoute, denom)
		if !ok {
			return
		}

		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func queryDemonCrossChainInfoHandlerFn(cliCtx context.CLIContext, queryRoute string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}
		// query for rewards from a particular delegator
		vars := mux.Vars(r)
		denom := vars[Denom]
		chainId, ok := rest.ParseUint64OrReturnBadRequest(w, vars[ChainId])
		if !ok {
			return
		}
		res, ok := checkResponseQueryDenomCCInfoResponse(w, cliCtx, queryRoute, denom, chainId)
		if !ok {
			return
		}

		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func checkResponseQueryDenomInfoResponse(
	w http.ResponseWriter, cliCtx context.CLIContext, queryRoute string, denom string) (res []byte, ok bool) {

	res, err := common.QueryDenomInfo(cliCtx, queryRoute, denom)
	if err != nil {
		rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
		return nil, false
	}

	return res, true
}
func checkResponseQueryDenomCCInfoResponse(
	w http.ResponseWriter, cliCtx context.CLIContext, queryRoute string, denom string, chainId uint64) (res []byte, ok bool) {

	res, err := common.QueryDenomCrossChainInfo(cliCtx, queryRoute, denom, chainId)
	if err != nil {
		rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
		return nil, false
	}

	return res, true
}
