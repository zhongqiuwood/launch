package gov

import (
	"time"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/ok-chain/okchain/x/token"

	"github.com/tendermint/tendermint/crypto"
)

const (
	// ModuleKey is the name of the module
	ModuleName = "gov"

	// StoreKey is the store key string for gov
	StoreKey = ModuleName

	// RouterKey is the message route for gov
	RouterKey = ModuleName

	// QuerierRoute is the querier route for gov
	QuerierRoute = ModuleName

	// Parameter store default namestore
	DefaultParamspace = ModuleName
)

// Parameter store key
var (
	ParamStoreKeyDepositParams = []byte("depositparams")
	ParamStoreKeyVotingParams  = []byte("votingparams")
	ParamStoreKeyTallyParams   = []byte("tallyparams")

	// TODO: Find another way to implement this without using accounts, or find a cleaner way to implement it using accounts.
	DepositedCoinsAccAddr        = sdk.AccAddress(crypto.AddressHash([]byte("govDepositedCoins")))
	BurnedDepositCoinsAccAddr    = sdk.AccAddress(crypto.AddressHash([]byte("govBurnedDepositCoins")))
	DexListDepositedCoinsAccAddr = sdk.AccAddress(crypto.AddressHash([]byte("govDexListDepositedCoins")))
)

// Governance Keeper
type Keeper struct {
	// The reference to the Param Keeper to get and set Global Params
	paramsKeeper params.Keeper

	// The reference to the Token Keeper to get asset issued
	tokenKeeper token.Keeper

	// The reference to the feeCollectionKeeper to collect fee
	feeCollectionKeeper auth.FeeCollectionKeeper

	// The reference to the Paramstore to get and set gov specific params
	paramSpace params.Subspace

	// The reference to the CoinKeeper to modify balances
	ck BankKeeper

	// The ValidatorSet to get information about validators
	vs sdk.ValidatorSet

	// The reference to the DelegationSet to get information about delegators
	ds sdk.DelegationSet

	// The (unexposed) keys used to access the stores from the Context.
	storeKey sdk.StoreKey

	// The codec codec for binary encoding/decoding.
	cdc *codec.Codec

	// Reserved codespace
	codespace sdk.CodespaceType
}

// NewKeeper returns a governance keeper. It handles:
// - submitting governance proposals
// - depositing funds into proposals, and activating upon sufficient funds being deposited
// - users voting on proposals, with weight proportional to stake in the system
// - and tallying the result of the vote.
func NewKeeper(cdc *codec.Codec, key sdk.StoreKey, paramsKeeper params.Keeper, tokenKeeper token.Keeper, feeKeeper auth.FeeCollectionKeeper,
	paramSpace params.Subspace, ck BankKeeper, ds sdk.DelegationSet, codespace sdk.CodespaceType) Keeper {

	return Keeper{
		storeKey:            key,
		paramsKeeper:        paramsKeeper,
		tokenKeeper:         tokenKeeper,
		feeCollectionKeeper: feeKeeper,
		paramSpace:          paramSpace.WithKeyTable(ParamKeyTable()),
		ck:                  ck,
		ds:                  ds,
		vs:                  ds.GetValidatorSet(),
		cdc:                 cdc,
		codespace:           codespace,
	}
}

// Proposals

func (keeper Keeper) NewProposal(ctx sdk.Context, title string, description string, proposalType ProposalKind, param Params, height int64) Proposal {
	switch proposalType {
	case ProposalTypeParameterChange:
		return keeper.NewParametersProposal(ctx, title, description, proposalType, param, height)
	default:
		return keeper.NewTextProposal(ctx, title, description, proposalType)
	}
	return nil
}

// =====================================================
// Proposals

// Creates a NewProposal

func (keeper Keeper) NewParametersProposal(ctx sdk.Context, title string, description string, proposalType ProposalKind, params Params, height int64) Proposal {
	proposalID, err := keeper.getNewProposalID(ctx)
	if err != nil {
		return nil
	}
	var textProposal = BasicProposal{
		ProposalID:       proposalID,
		Title:            title,
		Description:      description,
		ProposalType:     proposalType,
		Status:           StatusDepositPeriod,
		FinalTallyResult: EmptyTallyResult(),
		TotalDeposit:     sdk.DecCoins{},
		SubmitTime:       ctx.BlockHeader().Time,
	}

	var proposal Proposal = &ParameterProposal{
		textProposal,
		params,
		height,
	}

	depositPeriod := keeper.GetDepositParams(ctx).MaxDepositPeriod
	proposal.SetDepositEndTime(proposal.GetSubmitTime().Add(depositPeriod))
	keeper.SetProposal(ctx, proposal)
	keeper.InsertInactiveProposalQueue(ctx, proposal.GetDepositEndTime(), proposalID)
	return proposal
}

// Creates a NewProposal
func (keeper Keeper) NewTextProposal(ctx sdk.Context, title string, description string, proposalType ProposalKind) Proposal {
	proposalID, err := keeper.getNewProposalID(ctx)
	if err != nil {
		return nil
	}
	var proposal Proposal = &TextProposal{
		BasicProposal{
			ProposalID:       proposalID,
			Title:            title,
			Description:      description,
			ProposalType:     proposalType,
			Status:           StatusDepositPeriod,
			FinalTallyResult: EmptyTallyResult(),
			TotalDeposit:     sdk.DecCoins{},
			SubmitTime:       ctx.BlockHeader().Time,
		},
	}

	depositPeriod := keeper.GetDepositParams(ctx).MaxDepositPeriod
	proposal.SetDepositEndTime(proposal.GetSubmitTime().Add(depositPeriod))

	keeper.SetProposal(ctx, proposal)
	keeper.InsertInactiveProposalQueue(ctx, proposal.GetDepositEndTime(), proposalID)
	return proposal
}

// Creates a NewProposal
func (keeper Keeper) NewDexListProposal(ctx sdk.Context, msg MsgDexListSubmitProposal) Proposal {
	proposalID, err := keeper.getNewProposalID(ctx)
	if err != nil {
		return nil
	}
	var proposal Proposal = &DexListProposal{
		BasicProposal: BasicProposal{
			ProposalID:       proposalID,
			Title:            msg.Title,
			Description:      msg.Description,
			ProposalType:     msg.ProposalType,
			Status:           StatusDepositPeriod,
			FinalTallyResult: EmptyTallyResult(),
			TotalDeposit:     sdk.DecCoins{},
			SubmitTime:       ctx.BlockHeader().Time,
		},
		Proposer:      msg.Proposer,
		ListAsset:     msg.ListAsset,
		QuoteAsset:    msg.QuoteAsset,
		InitPrice:     msg.InitPrice,
		BlockHeight:   msg.BlockHeight,
		MaxPriceDigit: msg.MaxPriceDigit,
		MaxSizeDigit:  msg.MaxSizeDigit,
		MinTradeSize:  msg.MinTradeSize,
	}

	depositPeriod := keeper.GetDepositParams(ctx).MaxDepositPeriod
	proposal.SetDepositEndTime(proposal.GetSubmitTime().Add(depositPeriod))

	keeper.SetProposal(ctx, proposal)
	keeper.InsertInactiveProposalQueue(ctx, proposal.GetDepositEndTime(), proposalID)
	return proposal
}

// Get Proposal from store by ProposalID
func (keeper Keeper) GetProposal(ctx sdk.Context, proposalID uint64) Proposal {
	store := ctx.KVStore(keeper.storeKey)
	bz := store.Get(KeyProposal(proposalID))
	if bz == nil {
		return nil
	}
	var proposal Proposal
	keeper.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &proposal)
	return proposal
}

// Implements sdk.AccountKeeper.
func (keeper Keeper) SetProposal(ctx sdk.Context, proposal Proposal) {
	store := ctx.KVStore(keeper.storeKey)
	bz := keeper.cdc.MustMarshalBinaryLengthPrefixed(proposal)
	store.Set(KeyProposal(proposal.GetProposalID()), bz)
}

// Implements sdk.AccountKeeper.
func (keeper Keeper) DeleteProposal(ctx sdk.Context, proposalID uint64) {
	store := ctx.KVStore(keeper.storeKey)
	proposal := keeper.GetProposal(ctx, proposalID)
	keeper.RemoveFromInactiveProposalQueue(ctx, proposal.GetDepositEndTime(), proposalID)
	keeper.RemoveFromActiveProposalQueue(ctx, proposal.GetVotingEndTime(), proposalID)
	store.Delete(KeyProposal(proposalID))
}

// Get Proposal from store by ProposalID
// voterAddr will filter proposals by whether or not that address has voted on them
// depositorAddr will filter proposals by whether or not that address has deposited to them
// status will filter proposals by status
// numLatest will fetch a specified number of the most recent proposals, or 0 for all proposals
func (keeper Keeper) GetProposalsFiltered(ctx sdk.Context, voterAddr sdk.AccAddress, depositorAddr sdk.AccAddress, status ProposalStatus, numLatest uint64) []Proposal {

	maxProposalID, err := keeper.peekCurrentProposalID(ctx)
	if err != nil {
		return nil
	}

	matchingProposals := []Proposal{}

	if numLatest == 0 {
		numLatest = maxProposalID
	}

	for proposalID := maxProposalID - numLatest; proposalID < maxProposalID; proposalID++ {
		if voterAddr != nil && len(voterAddr) != 0 {
			_, found := keeper.GetVote(ctx, proposalID, voterAddr)
			if !found {
				continue
			}
		}

		if depositorAddr != nil && len(depositorAddr) != 0 {
			_, found := keeper.GetDeposit(ctx, proposalID, depositorAddr)
			if !found {
				continue
			}
		}

		proposal := keeper.GetProposal(ctx, proposalID)
		if proposal == nil {
			continue
		}

		if validProposalStatus(status) {
			if proposal.GetStatus() != status {
				continue
			}
		}

		matchingProposals = append(matchingProposals, proposal)
	}
	return matchingProposals
}

// Set the initial proposal ID
func (keeper Keeper) setInitialProposalID(ctx sdk.Context, proposalID uint64) sdk.Error {
	store := ctx.KVStore(keeper.storeKey)
	bz := store.Get(KeyNextProposalID)
	if bz != nil {
		return ErrInvalidGenesis(keeper.codespace, "Initial ProposalID already set")
	}
	bz = keeper.cdc.MustMarshalBinaryLengthPrefixed(proposalID)
	store.Set(KeyNextProposalID, bz)
	return nil
}

// Get the last used proposal ID
func (keeper Keeper) GetLastProposalID(ctx sdk.Context) (proposalID uint64) {
	proposalID, err := keeper.peekCurrentProposalID(ctx)
	if err != nil {
		return 0
	}
	proposalID--
	return
}

// Gets the next available ProposalID and increments it
func (keeper Keeper) getNewProposalID(ctx sdk.Context) (proposalID uint64, err sdk.Error) {
	store := ctx.KVStore(keeper.storeKey)
	bz := store.Get(KeyNextProposalID)
	if bz == nil {
		return 0, ErrInvalidGenesis(keeper.codespace, "InitialProposalID never set")
	}
	keeper.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &proposalID)
	bz = keeper.cdc.MustMarshalBinaryLengthPrefixed(proposalID + 1)
	store.Set(KeyNextProposalID, bz)
	return proposalID, nil
}

// Peeks the next available ProposalID without incrementing it
func (keeper Keeper) peekCurrentProposalID(ctx sdk.Context) (proposalID uint64, err sdk.Error) {
	store := ctx.KVStore(keeper.storeKey)
	bz := store.Get(KeyNextProposalID)
	if bz == nil {
		return 0, ErrInvalidGenesis(keeper.codespace, "InitialProposalID never set")
	}
	keeper.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &proposalID)
	return proposalID, nil
}

func (keeper Keeper) activateVotingPeriod(ctx sdk.Context, proposal Proposal) {
	proposal.SetVotingStartTime(ctx.BlockHeader().Time)
	votingPeriod := keeper.GetVotingParams(ctx).VotingPeriod
	proposal.SetVotingEndTime(proposal.GetVotingStartTime().Add(votingPeriod))
	proposal.SetStatus(StatusVotingPeriod)
	keeper.SetProposal(ctx, proposal)

	keeper.RemoveFromInactiveProposalQueue(ctx, proposal.GetDepositEndTime(), proposal.GetProposalID())
	keeper.InsertActiveProposalQueue(ctx, proposal.GetVotingEndTime(), proposal.GetProposalID())
}

// Params

// Returns the current DepositParams from the global param store
func (keeper Keeper) GetDepositParams(ctx sdk.Context) DepositParams {
	params := keeper.GetParams(ctx)
	return DepositParams{
		MinDeposit:       params.MinDeposit,
		MaxDepositPeriod: params.MaxDepositPeriod,
	}
	// var depositParams DepositParams
	// keeper.paramSpace.Get(ctx, ParamStoreKeyDepositParams, &depositParams)
	// return depositParams
}

// Returns the current DexListDepositParams from the global param store
func (keeper Keeper) GetDexListDepositParams(ctx sdk.Context) DexListDepositParams {
	params := keeper.GetParams(ctx)
	return DexListDepositParams{
		MinDeposit:       params.DexListMinDeposit,
		MaxDepositPeriod: params.DexListMaxDepositPeriod,
	}
	// var depositParams DepositParams
	// keeper.paramSpace.Get(ctx, ParamStoreKeyDepositParams, &depositParams)
	// return depositParams
}

// Returns the current VotingParams from the global param store
func (keeper Keeper) GetVotingParams(ctx sdk.Context) VotingParams {
	params := keeper.GetParams(ctx)
	return VotingParams{
		VotingPeriod: params.VotingPeriod,
	}
	// 	var votingParams VotingParams
	// keeper.paramSpace.Get(ctx, ParamStoreKeyVotingParams, &votingParams)
	// return votingParams
}

// Returns the current DexListVotingParams from the global param store
func (keeper Keeper) GetDexListVotingParams(ctx sdk.Context) DexListVotingParams {
	params := keeper.GetParams(ctx)
	return DexListVotingParams{
		VotingPeriod: params.VotingPeriod,
		VotingFee:    params.DexListVoteFee,
	}
}

// Returns the current DexListParams from the global param store
func (keeper Keeper) GetDexListParams(ctx sdk.Context) DexListParams {
	params := keeper.GetParams(ctx)
	return DexListParams{
		MaxBlockHeight: params.DexListMaxBlockHeight,
		Fee:            params.DexListFee,
		ExpireTime:     params.DexListExpireTime,
	}
}

// Returns the current TallyParam from the global param store
func (keeper Keeper) GetTallyParams(ctx sdk.Context) TallyParams {
	params := keeper.GetParams(ctx)
	return TallyParams{
		Threshold: params.Threshold,
		Veto:      params.Veto,
		Quorum:    params.Quorum,
	}
	// var tallyParams TallyParams
	// keeper.paramSpace.Get(ctx, ParamStoreKeyTallyParams, &tallyParams)
	// return tallyParams
}

func (keeper Keeper) setDepositParams(ctx sdk.Context, depositParams DepositParams) {
	keeper.paramSpace.Set(ctx, ParamStoreKeyDepositParams, &depositParams)
}

func (keeper Keeper) setVotingParams(ctx sdk.Context, votingParams VotingParams) {
	keeper.paramSpace.Set(ctx, ParamStoreKeyVotingParams, &votingParams)
}

func (keeper Keeper) setTallyParams(ctx sdk.Context, tallyParams TallyParams) {
	keeper.paramSpace.Set(ctx, ParamStoreKeyTallyParams, &tallyParams)
}

// Votes

// Adds a vote on a specific proposal
func (keeper Keeper) AddVote(ctx sdk.Context, proposalID uint64, voterAddr sdk.AccAddress, option VoteOption) sdk.Error {
	proposal := keeper.GetProposal(ctx, proposalID)
	if proposal == nil {
		return ErrUnknownProposal(keeper.codespace, proposalID)
	}

	if proposal.GetStatus() != StatusVotingPeriod {
		return ErrInactiveProposal(keeper.codespace, proposalID)
	}

	if !validVoteOption(option) {
		return ErrInvalidVote(keeper.codespace, option)
	}

	vote := Vote{
		ProposalID: proposalID,
		Voter:      voterAddr,
		Option:     option,
	}

	keeper.setVote(ctx, proposalID, voterAddr, vote)
	if proposal.GetProposalType() == ProposalTypeDexList {
		keeper.AddCollectedFees(ctx, NewCoinsFromDecCoins(keeper.GetParams(ctx).DexListVoteFee), voterAddr)
	}
	return nil
}

// Gets the vote of a specific voter on a specific proposal
func (keeper Keeper) GetVote(ctx sdk.Context, proposalID uint64, voterAddr sdk.AccAddress) (Vote, bool) {
	store := ctx.KVStore(keeper.storeKey)
	bz := store.Get(KeyVote(proposalID, voterAddr))
	if bz == nil {
		return Vote{}, false
	}
	var vote Vote
	keeper.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &vote)
	return vote, true
}

func (keeper Keeper) setVote(ctx sdk.Context, proposalID uint64, voterAddr sdk.AccAddress, vote Vote) {
	store := ctx.KVStore(keeper.storeKey)
	bz := keeper.cdc.MustMarshalBinaryLengthPrefixed(vote)
	store.Set(KeyVote(proposalID, voterAddr), bz)
}

// Gets all the votes on a specific proposal
func (keeper Keeper) GetVotes(ctx sdk.Context, proposalID uint64) sdk.Iterator {
	store := ctx.KVStore(keeper.storeKey)
	return sdk.KVStorePrefixIterator(store, KeyVotesSubspace(proposalID))
}

func (keeper Keeper) deleteVote(ctx sdk.Context, proposalID uint64, voterAddr sdk.AccAddress) {
	store := ctx.KVStore(keeper.storeKey)
	store.Delete(KeyVote(proposalID, voterAddr))
}

// Deposits

// Gets the deposit of a specific depositor on a specific proposal
func (keeper Keeper) GetDeposit(ctx sdk.Context, proposalID uint64, depositorAddr sdk.AccAddress) (Deposit, bool) {
	store := ctx.KVStore(keeper.storeKey)
	bz := store.Get(KeyDeposit(proposalID, depositorAddr))
	if bz == nil {
		return Deposit{}, false
	}
	var deposit Deposit
	keeper.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &deposit)
	return deposit, true
}

func (keeper Keeper) setDeposit(ctx sdk.Context, proposalID uint64, depositorAddr sdk.AccAddress, deposit Deposit) {
	store := ctx.KVStore(keeper.storeKey)
	bz := keeper.cdc.MustMarshalBinaryLengthPrefixed(deposit)
	store.Set(KeyDeposit(proposalID, depositorAddr), bz)
}

// Adds or updates a deposit of a specific depositor on a specific proposal
// Activates voting period when appropriate
func (keeper Keeper) AddDeposit(ctx sdk.Context, proposalID uint64, depositorAddr sdk.AccAddress, depositAmount sdk.DecCoins) (sdk.Error, bool) {
	// Checks to see if proposal exists
	proposal := keeper.GetProposal(ctx, proposalID)
	if proposal == nil {
		return ErrUnknownProposal(keeper.codespace, proposalID), false
	}

	// Check if proposal is still depositable
	if (proposal.GetStatus() != StatusDepositPeriod) && (proposal.GetStatus() != StatusVotingPeriod) {
		return ErrAlreadyFinishedProposal(keeper.codespace, proposalID), false
	}

	// Send coins from depositor's account to DepositedCoinsAccAddr account
	// TODO: Don't use an account for this purpose; it's clumsy and prone to misuse.
	_, err := keeper.ck.SendCoins(ctx, depositorAddr, DepositedCoinsAccAddr, NewCoinsFromDecCoins(depositAmount))
	if err != nil {
		return err, false
	}

	// Update proposal
	proposal.SetTotalDeposit(proposal.GetTotalDeposit().Add(depositAmount))
	keeper.SetProposal(ctx, proposal)

	// Check if deposit has provided sufficient total funds to transition the proposal into the voting period
	activatedVotingPeriod := false
	var minDeposit sdk.DecCoins
	if proposal.GetProposalType() == ProposalTypeDexList {
		minDeposit = keeper.GetDexListDepositParams(ctx).MinDeposit
	} else {
		minDeposit = keeper.GetDepositParams(ctx).MinDeposit
	}
	if proposal.GetStatus() == StatusDepositPeriod && proposal.GetTotalDeposit().IsAllGTE(minDeposit) {
		keeper.activateVotingPeriod(ctx, proposal)
		activatedVotingPeriod = true
	}

	// Add or update deposit object
	currDeposit, found := keeper.GetDeposit(ctx, proposalID, depositorAddr)
	if !found {
		newDeposit := Deposit{depositorAddr, proposalID, depositAmount}
		keeper.setDeposit(ctx, proposalID, depositorAddr, newDeposit)
	} else {
		currDeposit.Amount = currDeposit.Amount.Add(depositAmount)
		keeper.setDeposit(ctx, proposalID, depositorAddr, currDeposit)
	}

	return nil, activatedVotingPeriod
}

// Gets all the deposits on a specific proposal as an sdk.Iterator
func (keeper Keeper) GetDeposits(ctx sdk.Context, proposalID uint64) sdk.Iterator {
	store := ctx.KVStore(keeper.storeKey)
	return sdk.KVStorePrefixIterator(store, KeyDepositsSubspace(proposalID))
}

// Refunds and deletes all the deposits on a specific proposal
func (keeper Keeper) RefundDeposits(ctx sdk.Context, proposalID uint64) {
	store := ctx.KVStore(keeper.storeKey)
	depositsIterator := keeper.GetDeposits(ctx, proposalID)
	defer depositsIterator.Close()
	for ; depositsIterator.Valid(); depositsIterator.Next() {
		deposit := &Deposit{}
		keeper.cdc.MustUnmarshalBinaryLengthPrefixed(depositsIterator.Value(), deposit)

		_, err := keeper.ck.SendCoins(ctx, DepositedCoinsAccAddr, deposit.Depositor, NewCoinsFromDecCoins(deposit.Amount))
		if err != nil {
			panic("should not happen")
		}

		store.Delete(depositsIterator.Key())
	}
}

// Refunds and deletes dex list the deposits on a specific proposal when proposal is rejected or not enter voting period.
func (keeper Keeper) RefundDexListDeposits(ctx sdk.Context, proposer sdk.AccAddress) {
	_, err := keeper.ck.SendCoins(ctx, DexListDepositedCoinsAccAddr, proposer, NewCoinsFromDecCoins(keeper.GetDexListParams(ctx).Fee))
	if err != nil {
		panic("should not happen")
	}
}

// Deletes all the deposits on a specific proposal without refunding them
func (keeper Keeper) DeleteDeposits(ctx sdk.Context, proposalID uint64) {
	store := ctx.KVStore(keeper.storeKey)
	depositsIterator := keeper.GetDeposits(ctx, proposalID)
	defer depositsIterator.Close()
	for ; depositsIterator.Valid(); depositsIterator.Next() {
		deposit := &Deposit{}
		keeper.cdc.MustUnmarshalBinaryLengthPrefixed(depositsIterator.Value(), deposit)

		// TODO: Find a way to do this without using accounts.
		_, err := keeper.ck.SendCoins(ctx, DepositedCoinsAccAddr, BurnedDepositCoinsAccAddr, NewCoinsFromDecCoins(deposit.Amount))
		if err != nil {
			panic("should not happen")
		}

		store.Delete(depositsIterator.Key())
	}
}

// ProposalQueues

// Returns an iterator for all the proposals in the Active Queue that expire by endTime
func (keeper Keeper) ActiveProposalQueueIterator(ctx sdk.Context, endTime time.Time) sdk.Iterator {
	store := ctx.KVStore(keeper.storeKey)
	return store.Iterator(PrefixActiveProposalQueue, sdk.PrefixEndBytes(PrefixActiveProposalQueueTime(endTime)))
}

// Inserts a ProposalID into the active proposal queue at endTime
func (keeper Keeper) InsertActiveProposalQueue(ctx sdk.Context, endTime time.Time, proposalID uint64) {
	store := ctx.KVStore(keeper.storeKey)
	bz := keeper.cdc.MustMarshalBinaryLengthPrefixed(proposalID)
	store.Set(KeyActiveProposalQueueProposal(endTime, proposalID), bz)
}

// removes a proposalID from the Active Proposal Queue
func (keeper Keeper) RemoveFromActiveProposalQueue(ctx sdk.Context, endTime time.Time, proposalID uint64) {
	store := ctx.KVStore(keeper.storeKey)
	store.Delete(KeyActiveProposalQueueProposal(endTime, proposalID))
}

// Returns an iterator for all the proposals in the Inactive Queue that expire by endTime
func (keeper Keeper) InactiveProposalQueueIterator(ctx sdk.Context, endTime time.Time) sdk.Iterator {
	store := ctx.KVStore(keeper.storeKey)
	return store.Iterator(PrefixInactiveProposalQueue, sdk.PrefixEndBytes(PrefixInactiveProposalQueueTime(endTime)))
}

// Inserts a ProposalID into the inactive proposal queue at endTime
func (keeper Keeper) InsertInactiveProposalQueue(ctx sdk.Context, endTime time.Time, proposalID uint64) {
	store := ctx.KVStore(keeper.storeKey)
	bz := keeper.cdc.MustMarshalBinaryLengthPrefixed(proposalID)
	store.Set(KeyInactiveProposalQueueProposal(endTime, proposalID), bz)
}

// removes a proposalID from the Inactive Proposal Queue
func (keeper Keeper) RemoveFromInactiveProposalQueue(ctx sdk.Context, endTime time.Time, proposalID uint64) {
	store := ctx.KVStore(keeper.storeKey)
	store.Delete(KeyInactiveProposalQueueProposal(endTime, proposalID))
}

// Returns an iterator for all the proposals in the Waiting Queue that expire by endTime
func (keeper Keeper) WaitingProposalQueueIterator(ctx sdk.Context, blockHeight uint64) sdk.Iterator {
	store := ctx.KVStore(keeper.storeKey)
	return store.Iterator(PrefixWaitingProposalQueue, sdk.PrefixEndBytes(PrefixWaitingProposalQueueBlockHeight(blockHeight)))
}

// Inserts a ProposalID into the waiting proposal queue at endTime
func (keeper Keeper) InsertWaitingProposalQueue(ctx sdk.Context, blockHeight, proposalID uint64) {
	store := ctx.KVStore(keeper.storeKey)
	bz := keeper.cdc.MustMarshalBinaryLengthPrefixed(proposalID)
	store.Set(KeyWaitingProposalQueueProposal(blockHeight, proposalID), bz)
}

// removes a proposalID from the waiting Proposal Queue
func (keeper Keeper) RemoveFromWaitingProposalQueue(ctx sdk.Context, blockHeight, proposalID uint64) {
	store := ctx.KVStore(keeper.storeKey)
	store.Delete(KeyWaitingProposalQueueProposal(blockHeight, proposalID))
}

// get gov tx fee
func (keeper Keeper) GetGovTxFee() sdk.Coins {
	return sdk.Coins{sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewIntFromBigInt(sdk.MustNewDecFromStr(GovTxFee).Int))}
}

// use feeCollectionKeeper
// AddCollectedFees - add to the fee pool
func (k Keeper) AddCollectedFees(ctx sdk.Context, coins sdk.Coins, from sdk.AccAddress) sdk.Error {
	_, _, err := k.ck.SubtractCoins(ctx, from, coins)
	if err != nil {
		return sdk.ErrInsufficientCoins("Owner does not have enough okbs")
	}
	k.feeCollectionKeeper.AddCollectedFees(ctx, coins)
	return nil
}
