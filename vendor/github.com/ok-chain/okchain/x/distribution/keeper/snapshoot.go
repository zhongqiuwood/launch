package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	st "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/ok-chain/okchain/x/staking/types"
)

//store the snapshoot for next round distribution
func (k Keeper) SetValidatorsSnapshoot(ctx sdk.Context) {
	//fmt.Printf("SetValidatorsSnapshoot, before: %+v, set to: %+v\n", k.GetValidators(ctx), k.stakingKeeper.GetLastValidators(ctx))
	k.clearSnapshoot(ctx)
	for _, val := range k.stakingKeeper.GetLastValidators(ctx) {
		k.SetValidator(ctx, val)
		k.SetDelegationgSnapshoot(ctx, val.OperatorAddress)
	}
}
func (k Keeper) SetDelegationgSnapshoot(ctx sdk.Context, val sdk.ValAddress) {
	for _, del := range k.stakingKeeper.GetValidatorDelegations(ctx, val) {
		k.SetDelegation(ctx, del)
	}
}

// return all delegations to a specific validator. Useful for querier.
func (k Keeper) GetValidatorDelegations(ctx sdk.Context, valAddr sdk.ValAddress) (delegations []st.Delegation) {
	store := ctx.KVStore(k.kDelegationSnapshot)
	iterator := sdk.KVStorePrefixIterator(store, append(DelegationSnapshootPrefix, valAddr.Bytes()...))
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		delegation := st.MustUnmarshalDelegation(k.cdc, iterator.Value())
		//if delegation.GetValidatorAddr().Equals(valAddr) {
		delegations = append(delegations, delegation)
		//}
	}
	return delegations
}

// return a specific delegation
func (k Keeper) GetDelegation(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) (delegation st.Delegation, found bool) {
	store := ctx.KVStore(k.kDelegationSnapshot)
	value := store.Get(GetDelegationSnapshootKey(delAddr, valAddr))
	if value == nil {
		return delegation, false
	}

	delegation = st.MustUnmarshalDelegation(k.cdc, value)
	return delegation, true
}

//return all validators in current snapshoot
func (k Keeper) GetValidators(ctx sdk.Context) (validators []st.Validator) {
	store := ctx.KVStore(k.kValidatorsSnapShot)
	iterator := sdk.KVStorePrefixIterator(store, ValidatorSnapshootPrefix)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		validator := st.MustUnmarshalValidator(k.cdc, iterator.Value())
		validators = append(validators, validator)
	}
	return validators
}

//return a validator by valaddr
func (k Keeper) GetValidator(ctx sdk.Context, val sdk.ValAddress) (validator st.Validator, found bool) {
	store := ctx.KVStore(k.kValidatorsSnapShot)
	value := store.Get(GetValidatorSnapshootKey(val))
	if value == nil {
		return validator, false
	}

	// amino bytes weren't found in cache, so amino unmarshal and add it to the cache
	validator = st.MustUnmarshalValidator(k.cdc, value)
	return validator, true
}

// iterate through the validator set and perform the provided function
func (k Keeper) IterateValidators(ctx sdk.Context, fn func(index int64, validator sdk.Validator) (stop bool)) {
	store := ctx.KVStore(k.kValidatorsSnapShot)
	iterator := sdk.KVStorePrefixIterator(store, ValidatorSnapshootPrefix)
	defer iterator.Close()
	i := int64(0)
	for ; iterator.Valid(); iterator.Next() {
		validator := st.MustUnmarshalValidator(k.cdc, iterator.Value())
		stop := fn(i, validator) // XXX is this safe will the validator unexposed fields be able to get written to?
		if stop {
			break
		}
		i++
	}
}

// iterate through all of the delegations from a delegator
func (k Keeper) IterateDelegations(ctx sdk.Context, delAddr sdk.AccAddress,
	fn func(index int64, del sdk.Delegation) (stop bool)) {
	store := ctx.KVStore(k.kDelegationSnapshot)
	iterator := sdk.KVStorePrefixIterator(store, DelegationSnapshootPrefix) //smallest to largest
	defer iterator.Close()

	for i := int64(0); iterator.Valid(); iterator.Next() {
		del := types.MustUnmarshalDelegation(k.cdc, iterator.Value())
		if del.DelegatorAddress.Equals(delAddr) {
			stop := fn(i, del)
			if stop {
				break
			}
			i++
		}
	}
}

func (k Keeper) SetValidator(ctx sdk.Context, val sdk.Validator) {
	store := ctx.KVStore(k.kValidatorsSnapShot)
	store.Set(GetValidatorSnapshootKey(val.GetOperator()), k.cdc.MustMarshalBinaryLengthPrefixed(val))
}

//Store delegations belonged to validator, key = delegationPrefix+valAddr+delAddr, critical for iteration by valAddr
func (k Keeper) SetDelegation(ctx sdk.Context, del sdk.Delegation) {
	store := ctx.KVStore(k.kDelegationSnapshot)
	store.Set(GetDelegationSnapshootKey(del.GetDelegatorAddr(), del.GetValidatorAddr()), k.cdc.MustMarshalBinaryLengthPrefixed(del))
}

//clearSnapshoot delete the previous snapshoot of validators and delegations, to be optimized
func (k Keeper) clearSnapshoot(ctx sdk.Context) {
	store := ctx.KVStore(k.kValidatorsSnapShot)
	iterator := sdk.KVStorePrefixIterator(store, ValidatorSnapshootPrefix)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		store.Delete(iterator.Key())
	}
	store = ctx.KVStore(k.kDelegationSnapshot)
	iterator = sdk.KVStorePrefixIterator(store, DelegationSnapshootPrefix)
	for ; iterator.Valid(); iterator.Next() {
		store.Delete(iterator.Key())
	}
}
