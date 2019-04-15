package querier

import (
	"fmt"

	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	keep "github.com/ok-chain/okchain/x/staking/keeper"
	"github.com/ok-chain/okchain/x/staking/types"
	"math/big"
)

// query endpoints supported by the staking Querier
const (
	QueryValidators                    = "validators"
	QueryValidator                     = "validator"
	QueryDelegatorDelegations          = "delegatorDelegations"
	QueryDelegatorUnbondingDelegations = "delegatorUnbondingDelegations"
	QueryRedelegations                 = "redelegations"
	QueryValidatorDelegations          = "validatorDelegations"
	QueryValidatorRedelegations        = "validatorRedelegations"
	QueryValidatorUnbondingDelegations = "validatorUnbondingDelegations"
	QueryDelegator                     = "delegator"
	QueryDelegation                    = "delegation"
	QueryUnbondingDelegation           = "unbondingDelegation"
	QueryDelegatorValidators           = "delegatorValidators"
	QueryDelegatorValidator            = "delegatorValidator"
	QueryPool                          = "pool"
	QueryParameters                    = "parameters"
)

// creates a querier for staking REST endpoints
func NewQuerier(k keep.Keeper, cdc *codec.Codec) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err sdk.Error) {
		switch path[0] {
		case QueryValidators:
			return queryValidators(ctx, cdc, k)
		case QueryValidator:
			return queryValidator(ctx, cdc, req, k)
		case QueryValidatorDelegations:
			return queryValidatorDelegations(ctx, cdc, req, k)
		case QueryValidatorUnbondingDelegations:
			return queryValidatorUnbondingDelegations(ctx, cdc, req, k)
		case QueryDelegation:
			return queryDelegation(ctx, cdc, req, k)
		case QueryUnbondingDelegation:
			return queryUnbondingDelegation(ctx, cdc, req, k)
		case QueryDelegatorDelegations:
			return queryDelegatorDelegations(ctx, cdc, req, k)
		case QueryDelegatorUnbondingDelegations:
			return queryDelegatorUnbondingDelegations(ctx, cdc, req, k)
		case QueryRedelegations:
			return queryRedelegations(ctx, cdc, req, k)
		case QueryDelegatorValidators:
			return queryDelegatorValidators(ctx, cdc, req, k)
		case QueryDelegatorValidator:
			return queryDelegatorValidator(ctx, cdc, req, k)
		case QueryPool:
			return queryPool(ctx, cdc, k)
		case QueryParameters:
			return queryParameters(ctx, cdc, k)
		default:
			return nil, sdk.ErrUnknownRequest("unknown staking query endpoint")
		}
	}
}

// defines the params for the following queries:
// - 'custom/staking/delegatorDelegations'
// - 'custom/staking/delegatorUnbondingDelegations'
// - 'custom/staking/delegatorRedelegations'
// - 'custom/staking/delegatorValidators'
type QueryDelegatorParams struct {
	DelegatorAddr sdk.AccAddress
}

func NewQueryDelegatorParams(delegatorAddr sdk.AccAddress) QueryDelegatorParams {
	return QueryDelegatorParams{
		DelegatorAddr: delegatorAddr,
	}
}

// defines the params for the following queries:
// - 'custom/staking/validator'
// - 'custom/staking/validatorDelegations'
// - 'custom/staking/validatorUnbondingDelegations'
// - 'custom/staking/validatorRedelegations'
type QueryValidatorParams struct {
	ValidatorAddr sdk.ValAddress
}

func NewQueryValidatorParams(validatorAddr sdk.ValAddress) QueryValidatorParams {
	return QueryValidatorParams{
		ValidatorAddr: validatorAddr,
	}
}

// defines the params for the following queries:
// - 'custom/staking/delegation'
// - 'custom/staking/unbondingDelegation'
// - 'custom/staking/delegatorValidator'
type QueryBondsParams struct {
	DelegatorAddr sdk.AccAddress
	ValidatorAddr sdk.ValAddress
}

func NewQueryBondsParams(delegatorAddr sdk.AccAddress, validatorAddr sdk.ValAddress) QueryBondsParams {
	return QueryBondsParams{
		DelegatorAddr: delegatorAddr,
		ValidatorAddr: validatorAddr,
	}
}

// defines the params for the following queries:
// - 'custom/staking/redelegation'
type QueryRedelegationParams struct {
	DelegatorAddr    sdk.AccAddress
	SrcValidatorAddr sdk.ValAddress
	DstValidatorAddr sdk.ValAddress
}

func NewQueryRedelegationParams(delegatorAddr sdk.AccAddress, srcValidatorAddr sdk.ValAddress, dstValidatorAddr sdk.ValAddress) QueryRedelegationParams {
	return QueryRedelegationParams{
		DelegatorAddr:    delegatorAddr,
		SrcValidatorAddr: srcValidatorAddr,
		DstValidatorAddr: dstValidatorAddr,
	}
}

func queryValidators(ctx sdk.Context, cdc *codec.Codec, k keep.Keeper) (res []byte, err sdk.Error) {
	// stakingParams := k.GetParams(ctx)
	validators := k.GetAllValidators(ctx)

	var outputValidators []types.FormattedValidator

	for _, validator := range validators {
		var formattedValidator types.FormattedValidator
		formattedValidator.OperatorAddress = validator.OperatorAddress
		formattedValidator.ConsPubKey = validator.ConsPubKey
		formattedValidator.Jailed = validator.Jailed
		formattedValidator.UnbondingHeight = validator.UnbondingHeight
		formattedValidator.UnbondingCompletionTime = validator.UnbondingCompletionTime
		formattedValidator.DelegatorShares = validator.DelegatorShares
		formattedValidator.Description = validator.Description
		formattedValidator.Commission = validator.Commission
		formattedValidator.Status = validator.Status
		formattedValidator.Tokens = validator.Tokens.ToDec().QuoInt(sdk.NewIntFromBigInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(sdk.Precision), nil)))
		formattedValidator.MinSelfDelegation = validator.MinSelfDelegation.ToDec().QuoInt(sdk.NewIntFromBigInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(sdk.Precision), nil)))

		outputValidators = append(outputValidators, formattedValidator)
	}

	res, errRes := codec.MarshalJSONIndent(cdc, outputValidators)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", errRes.Error()))
	}
	return res, nil
}

func queryValidator(ctx sdk.Context, cdc *codec.Codec, req abci.RequestQuery, k keep.Keeper) (res []byte, err sdk.Error) {
	var params QueryValidatorParams

	errRes := cdc.UnmarshalJSON(req.Data, &params)
	if errRes != nil {
		return []byte{}, sdk.ErrInternal(fmt.Sprintf("failed to parse params: %s", err))
	}

	validator, found := k.GetValidator(ctx, params.ValidatorAddr)
	if !found {
		return []byte{}, types.ErrNoValidatorFound(types.DefaultCodespace)
	}
	var formattedValidator types.FormattedValidator
	formattedValidator.OperatorAddress = validator.OperatorAddress
	formattedValidator.ConsPubKey = validator.ConsPubKey
	formattedValidator.Jailed = validator.Jailed
	formattedValidator.UnbondingHeight = validator.UnbondingHeight
	formattedValidator.UnbondingCompletionTime = validator.UnbondingCompletionTime
	formattedValidator.DelegatorShares = validator.DelegatorShares
	formattedValidator.Description = validator.Description
	formattedValidator.Commission = validator.Commission
	formattedValidator.Status = validator.Status
	formattedValidator.Tokens = validator.Tokens.ToDec().QuoInt(sdk.NewIntFromBigInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(sdk.Precision), nil)))
	formattedValidator.MinSelfDelegation = validator.MinSelfDelegation.ToDec().QuoInt(sdk.NewIntFromBigInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(sdk.Precision), nil)))

	res, errRes = codec.MarshalJSONIndent(cdc, formattedValidator)
	if errRes != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", errRes.Error()))
	}
	return res, nil
}

func queryValidatorDelegations(ctx sdk.Context, cdc *codec.Codec, req abci.RequestQuery, k keep.Keeper) (res []byte, err sdk.Error) {
	var params QueryValidatorParams

	errRes := cdc.UnmarshalJSON(req.Data, &params)
	if errRes != nil {
		return []byte{}, sdk.ErrInternal(fmt.Sprintf("failed to parse params: %s", err))
	}

	delegations := k.GetValidatorDelegations(ctx, params.ValidatorAddr)

	res, errRes = codec.MarshalJSONIndent(cdc, delegations)
	if errRes != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", errRes.Error()))
	}
	return res, nil
}

func queryValidatorUnbondingDelegations(ctx sdk.Context, cdc *codec.Codec, req abci.RequestQuery, k keep.Keeper) (res []byte, err sdk.Error) {
	var params QueryValidatorParams

	errRes := cdc.UnmarshalJSON(req.Data, &params)
	if errRes != nil {
		return []byte{}, sdk.ErrInternal(fmt.Sprintf("failed to parse params: %s", err))
	}

	unbonds := k.GetUnbondingDelegationsFromValidator(ctx, params.ValidatorAddr)
	var outputUnbonds []types.FormattedUnbondingDelegation
	for _, ubd := range unbonds {
		var formattedUbd types.FormattedUnbondingDelegation
		var formattedUbdEntries []types.FormattedUnbondingDelegationEntry
		formattedUbd.DelegatorAddress = ubd.DelegatorAddress
		formattedUbd.ValidatorAddress = ubd.ValidatorAddress
		for i := 0; i < len(ubd.Entries); i++ {
			var formattedUbdEntry types.FormattedUnbondingDelegationEntry
			formattedUbdEntry.CreationHeight = ubd.Entries[i].CreationHeight
			formattedUbdEntry.CompletionTime = ubd.Entries[i].CompletionTime
			formattedUbdEntry.InitialBalance = ubd.Entries[i].InitialBalance.ToDec().QuoInt(sdk.NewIntFromBigInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(sdk.Precision), nil)))
			formattedUbdEntry.Balance = ubd.Entries[i].Balance.ToDec().QuoInt(sdk.NewIntFromBigInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(sdk.Precision), nil)))

			formattedUbdEntries = append(formattedUbdEntries, formattedUbdEntry)
		}
		formattedUbd.Entries = formattedUbdEntries

		outputUnbonds = append(outputUnbonds, formattedUbd)
	}

	res, errRes = codec.MarshalJSONIndent(cdc, outputUnbonds)
	if errRes != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", errRes.Error()))
	}
	return res, nil
}

func queryDelegatorDelegations(ctx sdk.Context, cdc *codec.Codec, req abci.RequestQuery, k keep.Keeper) (res []byte, err sdk.Error) {
	var params QueryDelegatorParams

	errRes := cdc.UnmarshalJSON(req.Data, &params)
	if errRes != nil {
		return []byte{}, sdk.ErrInternal(fmt.Sprintf("failed to parse params: %s", err))
	}

	delegations := k.GetAllDelegatorDelegations(ctx, params.DelegatorAddr)

	res, errRes = codec.MarshalJSONIndent(cdc, delegations)
	if errRes != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", errRes.Error()))
	}
	return res, nil
}

func queryDelegatorUnbondingDelegations(ctx sdk.Context, cdc *codec.Codec, req abci.RequestQuery, k keep.Keeper) (res []byte, err sdk.Error) {
	var params QueryDelegatorParams

	errRes := cdc.UnmarshalJSON(req.Data, &params)
	if errRes != nil {
		return []byte{}, sdk.ErrInternal(fmt.Sprintf("failed to parse params: %s", err))
	}

	unbondingDelegations := k.GetAllUnbondingDelegations(ctx, params.DelegatorAddr)
	var outputUnbonds []types.FormattedUnbondingDelegation
	for _, ubd := range unbondingDelegations {
		var formattedUbd types.FormattedUnbondingDelegation
		var formattedUbdEntries []types.FormattedUnbondingDelegationEntry
		formattedUbd.DelegatorAddress = ubd.DelegatorAddress
		formattedUbd.ValidatorAddress = ubd.ValidatorAddress
		for i := 0; i < len(ubd.Entries); i++ {
			var formattedUbdEntry types.FormattedUnbondingDelegationEntry
			formattedUbdEntry.CreationHeight = ubd.Entries[i].CreationHeight
			formattedUbdEntry.CompletionTime = ubd.Entries[i].CompletionTime
			formattedUbdEntry.InitialBalance = ubd.Entries[i].InitialBalance.ToDec().QuoInt(sdk.NewIntFromBigInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(sdk.Precision), nil)))
			formattedUbdEntry.Balance = ubd.Entries[i].Balance.ToDec().QuoInt(sdk.NewIntFromBigInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(sdk.Precision), nil)))

			formattedUbdEntries = append(formattedUbdEntries, formattedUbdEntry)
		}
		formattedUbd.Entries = formattedUbdEntries

		outputUnbonds = append(outputUnbonds, formattedUbd)
	}

	res, errRes = codec.MarshalJSONIndent(cdc, unbondingDelegations)
	if errRes != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", errRes.Error()))
	}
	return res, nil
}

func queryDelegatorValidators(ctx sdk.Context, cdc *codec.Codec, req abci.RequestQuery, k keep.Keeper) (res []byte, err sdk.Error) {
	var params QueryDelegatorParams

	stakingParams := k.GetParams(ctx)

	errRes := cdc.UnmarshalJSON(req.Data, &params)
	if errRes != nil {
		return []byte{}, sdk.ErrInternal(fmt.Sprintf("failed to parse params: %s", err))
	}

	validators := k.GetDelegatorValidators(ctx, params.DelegatorAddr, stakingParams.MaxValidators)
	var outputValidators []types.FormattedValidator
	for _, validator := range validators {
		var formattedValidator types.FormattedValidator
		formattedValidator.OperatorAddress = validator.OperatorAddress
		formattedValidator.ConsPubKey = validator.ConsPubKey
		formattedValidator.Jailed = validator.Jailed
		formattedValidator.UnbondingHeight = validator.UnbondingHeight
		formattedValidator.UnbondingCompletionTime = validator.UnbondingCompletionTime
		formattedValidator.DelegatorShares = validator.DelegatorShares
		formattedValidator.Description = validator.Description
		formattedValidator.Commission = validator.Commission
		formattedValidator.Status = validator.Status
		formattedValidator.Tokens = validator.Tokens.ToDec().QuoInt(sdk.NewIntFromBigInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(sdk.Precision), nil)))
		formattedValidator.MinSelfDelegation = validator.MinSelfDelegation.ToDec().QuoInt(sdk.NewIntFromBigInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(sdk.Precision), nil)))

		outputValidators = append(outputValidators, formattedValidator)
	}

	res, errRes = codec.MarshalJSONIndent(cdc, outputValidators)
	if errRes != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", errRes.Error()))
	}
	return res, nil
}

func queryDelegatorValidator(ctx sdk.Context, cdc *codec.Codec, req abci.RequestQuery, k keep.Keeper) (res []byte, err sdk.Error) {
	var params QueryBondsParams

	errRes := cdc.UnmarshalJSON(req.Data, &params)
	if errRes != nil {
		return []byte{}, sdk.ErrInternal(fmt.Sprintf("failed to parse params: %s", err))
	}

	validator, err := k.GetDelegatorValidator(ctx, params.DelegatorAddr, params.ValidatorAddr)
	if err != nil {
		return
	}

	// for a high resolution display
	var formattedValidator types.FormattedValidator
	formattedValidator.OperatorAddress = validator.OperatorAddress
	formattedValidator.ConsPubKey = validator.ConsPubKey
	formattedValidator.Jailed = validator.Jailed
	formattedValidator.UnbondingHeight = validator.UnbondingHeight
	formattedValidator.UnbondingCompletionTime = validator.UnbondingCompletionTime
	formattedValidator.DelegatorShares = validator.DelegatorShares
	formattedValidator.Description = validator.Description
	formattedValidator.Commission = validator.Commission
	formattedValidator.Status = validator.Status
	formattedValidator.Tokens = validator.Tokens.ToDec().QuoInt(sdk.NewIntFromBigInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(sdk.Precision), nil)))
	formattedValidator.MinSelfDelegation = validator.MinSelfDelegation.ToDec().QuoInt(sdk.NewIntFromBigInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(sdk.Precision), nil)))

	res, errRes = codec.MarshalJSONIndent(cdc, formattedValidator)
	if errRes != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", errRes.Error()))
	}
	return res, nil
}

func queryDelegation(ctx sdk.Context, cdc *codec.Codec, req abci.RequestQuery, k keep.Keeper) (res []byte, err sdk.Error) {
	var params QueryBondsParams

	errRes := cdc.UnmarshalJSON(req.Data, &params)
	if errRes != nil {
		return []byte{}, sdk.ErrInternal(fmt.Sprintf("failed to parse params: %s", err))
	}

	delegation, found := k.GetDelegation(ctx, params.DelegatorAddr, params.ValidatorAddr)
	if !found {
		return []byte{}, types.ErrNoDelegation(types.DefaultCodespace)
	}

	res, errRes = codec.MarshalJSONIndent(cdc, delegation)
	if errRes != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", errRes.Error()))
	}
	return res, nil
}

func queryUnbondingDelegation(ctx sdk.Context, cdc *codec.Codec, req abci.RequestQuery, k keep.Keeper) (res []byte, err sdk.Error) {
	var params QueryBondsParams

	errRes := cdc.UnmarshalJSON(req.Data, &params)
	if errRes != nil {
		return []byte{}, sdk.ErrInternal(fmt.Sprintf("failed to parse params: %s", err))
	}

	ubd, found := k.GetUnbondingDelegation(ctx, params.DelegatorAddr, params.ValidatorAddr)
	if !found {
		return []byte{}, types.ErrNoUnbondingDelegation(types.DefaultCodespace)
	}

	var formattedUbd types.FormattedUnbondingDelegation
	var formattedUbdEntries []types.FormattedUnbondingDelegationEntry
	formattedUbd.DelegatorAddress = ubd.DelegatorAddress
	formattedUbd.ValidatorAddress = ubd.ValidatorAddress
	for i := 0; i < len(ubd.Entries); i++ {
		var formattedUbdEntry types.FormattedUnbondingDelegationEntry
		formattedUbdEntry.CreationHeight = ubd.Entries[i].CreationHeight
		formattedUbdEntry.CompletionTime = ubd.Entries[i].CompletionTime
		formattedUbdEntry.InitialBalance = ubd.Entries[i].InitialBalance.ToDec().QuoInt(sdk.NewIntFromBigInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(sdk.Precision), nil)))
		formattedUbdEntry.Balance = ubd.Entries[i].Balance.ToDec().QuoInt(sdk.NewIntFromBigInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(sdk.Precision), nil)))

		formattedUbdEntries = append(formattedUbdEntries, formattedUbdEntry)
	}
	formattedUbd.Entries = formattedUbdEntries

	res, errRes = codec.MarshalJSONIndent(cdc, formattedUbd)
	if errRes != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", errRes.Error()))
	}
	return res, nil
}

func queryRedelegations(ctx sdk.Context, cdc *codec.Codec, req abci.RequestQuery, k keep.Keeper) (res []byte, err sdk.Error) {
	var params QueryRedelegationParams

	errRes := cdc.UnmarshalJSON(req.Data, &params)
	if errRes != nil {
		return []byte{}, sdk.ErrUnknownRequest(string(req.Data))
	}

	var redels []types.Redelegation

	if !params.DelegatorAddr.Empty() && !params.SrcValidatorAddr.Empty() && !params.DstValidatorAddr.Empty() {
		redel, found := k.GetRedelegation(ctx, params.DelegatorAddr, params.SrcValidatorAddr, params.DstValidatorAddr)
		if !found {
			return []byte{}, types.ErrNoRedelegation(types.DefaultCodespace)
		}
		redels = []types.Redelegation{redel}
	} else if params.DelegatorAddr.Empty() && !params.SrcValidatorAddr.Empty() && params.DstValidatorAddr.Empty() {
		redels = k.GetRedelegationsFromValidator(ctx, params.SrcValidatorAddr)
	} else {
		redels = k.GetAllRedelegations(ctx, params.DelegatorAddr, params.SrcValidatorAddr, params.DstValidatorAddr)
	}

	var outputRedels []types.FormattedRedelegation
	for _, red := range redels {
		var formattedRed types.FormattedRedelegation
		var formattedRedEntries []types.FormattedRedelegationEntry

		formattedRed.DelegatorAddress = red.DelegatorAddress
		formattedRed.ValidatorSrcAddress = red.ValidatorSrcAddress
		formattedRed.ValidatorDstAddress = red.ValidatorDstAddress
		for i := 0; i < len(red.Entries); i++ {
			var formattedRedEntry types.FormattedRedelegationEntry
			formattedRedEntry.CreationHeight = red.Entries[i].CreationHeight
			formattedRedEntry.CompletionTime = red.Entries[i].CompletionTime
			formattedRedEntry.SharesDst = red.Entries[i].SharesDst
			formattedRedEntry.InitialBalance = red.Entries[i].InitialBalance.ToDec().QuoInt(sdk.NewIntFromBigInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(sdk.Precision), nil)))

			formattedRedEntries = append(formattedRedEntries, formattedRedEntry)
		}
		formattedRed.Entries = formattedRedEntries

		outputRedels = append(outputRedels, formattedRed)
	}

	res, errRes = codec.MarshalJSONIndent(cdc, outputRedels)
	if errRes != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", errRes.Error()))
	}
	return res, nil
}

func queryPool(ctx sdk.Context, cdc *codec.Codec, k keep.Keeper) (res []byte, err sdk.Error) {
	pool := k.GetPool(ctx)
	var formattedPool types.FormattedPool
	formattedPool.NotBondedTokens = pool.NotBondedTokens.ToDec().QuoInt(sdk.NewIntFromBigInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(sdk.Precision), nil)))
	formattedPool.BondedTokens = pool.BondedTokens.ToDec().QuoInt(sdk.NewIntFromBigInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(sdk.Precision), nil)))

	res, errRes := codec.MarshalJSONIndent(cdc, formattedPool)
	if errRes != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", errRes.Error()))
	}
	return res, nil
}

func queryParameters(ctx sdk.Context, cdc *codec.Codec, k keep.Keeper) (res []byte, err sdk.Error) {
	params := k.GetParams(ctx)

	res, errRes := codec.MarshalJSONIndent(cdc, params)
	if errRes != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", errRes.Error()))
	}
	return res, nil
}
