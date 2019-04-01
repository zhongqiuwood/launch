package token

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GenesisState - all slashing state that must be provided at genesis
type GenesisState struct {
	Params Params  `json:"params"`
	Info   []Token `json:"info"`
}

// DefaultGenesisState - default GenesisState used by Cosmos Hub
func DefaultGenesisState() GenesisState {
	return GenesisState{
		Params: DefaultParams(),
		Info:   []Token{DefaultGenesisStateOKB()},
	}
}

func DefaultGenesisStateOKB() Token {
	return Token{
		Name:        "OKB",
		Symbol:      "okb",
		TotalSupply: 1000000000,
		Owner:       nil,
		Mintable:    false,
	}
}

// ValidateGenesis validates the slashing genesis parameters
func ValidateGenesis(data GenesisState) error {
	// downtime := data.Params.SlashFractionDowntime
	// if downtime.IsNegative() || downtime.GT(sdk.OneDec()) {
	// 	return fmt.Errorf("Slashing fraction downtime should be less than or equal to one and greater than zero, is %s", downtime.String())
	// }

	// dblSign := data.Params.SlashFractionDoubleSign
	// if dblSign.IsNegative() || dblSign.GT(sdk.OneDec()) {
	// 	return fmt.Errorf("Slashing fraction double sign should be less than or equal to one and greater than zero, is %s", dblSign.String())
	// }

	// minSign := data.Params.MinSignedPerWindow
	// if minSign.IsNegative() || minSign.GT(sdk.OneDec()) {
	// 	return fmt.Errorf("Min signed per window should be less than or equal to one and greater than zero, is %s", minSign.String())
	// }

	// maxEvidence := data.Params.MaxEvidenceAge
	// if maxEvidence < 1*time.Minute {
	// 	return fmt.Errorf("Max evidence age must be at least 1 minute, is %s", maxEvidence.String())
	// }

	// downtimeJail := data.Params.DowntimeJailDuration
	// if downtimeJail < 1*time.Minute {
	// 	return fmt.Errorf("Downtime unblond duration must be at least 1 minute, is %s", downtimeJail.String())
	// }

	// signedWindow := data.Params.SignedBlocksWindow
	// if signedWindow < 10 {
	// 	return fmt.Errorf("Signed blocks window must be at least 10, is %d", signedWindow)
	// }

	return nil
}

// InitGenesis initialize default parameters
// and the keeper's address to pubkey map
func InitGenesis(ctx sdk.Context, keeper Keeper, data GenesisState) {
	keeper.SetParams(ctx, data.Params)
}

// ExportGenesis writes the current store values
// to a genesis file, which can be imported again
// with InitGenesis
func ExportGenesis(ctx sdk.Context, keeper Keeper) (data GenesisState) {
	params := keeper.GetParams(ctx)

	return GenesisState{
		Params: params,
	}
}
