package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"path"
	"sort"
	"strconv"
	"strings"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/launch/pkg"
	"github.com/ok-chain/okchain/app"
	"github.com/ok-chain/okchain/x/token"
	"github.com/tendermint/go-amino"
	tmtypes "github.com/tendermint/tendermint/types"
)

const (
	captainJSON = "accounts/captain.json"
	adminJSON   = "accounts/admin.json"
	othersJSON  = "accounts/others.json"

	genesisTemplate = "params/genesis_template.json"
	genTxPath       = "gentx/data"
	genesisFile     = "genesis.json"

	okbDenomination     = "okb"
	okbGenesisTotal     = 1000000000
	addressGenesisTotal = 7

	timeGenesisString = "2019-03-13 23:00:00 -0000 UTC"
	chainID           = "okchain"
)

// constants but can't use `const`
var (
	timeGenesis time.Time

	// vesting times
	timeGenesisTwoMonths time.Time
	timeGenesisOneYear   time.Time
	timeGenesisTwoYears  time.Time
)

// initialize the times!
func init() {
	var err error
	timeLayoutString := "2006-01-02 15:04:05 -0700 MST"
	timeGenesis, err = time.Parse(timeLayoutString, timeGenesisString)
	if err != nil {
		panic(err)
	}
	timeGenesisTwoMonths = timeGenesis.AddDate(0, 2, 0)
	timeGenesisOneYear = timeGenesis.AddDate(1, 0, 0)
	timeGenesisTwoYears = timeGenesis.AddDate(2, 0, 0)
}

// max precision on amt is two decimals ("centi-atoms")
func atomToUAtomInt(amt float64) sdk.Int {
	// amt is specified to 2 decimals ("centi-atoms").
	// multiply by 100 to get the number of centi-atoms
	// and round to int64.
	// Multiply by remaining to get uAtoms.
	var precision float64 = 100
	var remaining int64 = 10000

	catoms := int64(math.Round(amt * precision))
	uAtoms := catoms * remaining
	return sdk.NewInt(uAtoms)
}

// convert atoms with two decimal precision to coins
func newCoins(amt float64) sdk.DecCoins {
	uAtoms := sdk.MustNewDecFromStr(strconv.FormatFloat(amt, 'f', -1, 64))

	return sdk.DecCoins{sdk.NewDecCoinFromDec(okbDenomination, uAtoms)}
}

func main() {
	// for each path, accumulate the contributors file.
	// icf addresses are in bech32, fundraiser are in hex
	contribs := make(map[string]float64)
	accumulateBechContributors(captainJSON, contribs)
	captainAccount := makeGenesisAccounts(contribs, nil, MultisigAccount{})

	if len(captainAccount) != 1 {
		panic(fmt.Errorf("Invalid captain account!"))
	}

	accumulateBechContributors(adminJSON, contribs)
	accumulateContributors(othersJSON, contribs)
	genesisAccounts := makeGenesisAccounts(contribs, nil, MultisigAccount{})

	// check totals
	checkTotals(genesisAccounts)

	fmt.Println("-----------")
	fmt.Println("TOTAL addrs", len(genesisAccounts))
	fmt.Println("TOTAL okbs", okbGenesisTotal)

	// load gentxs
	fs, err := ioutil.ReadDir(genTxPath)
	if err != nil {
		panic(err)
	}

	var genTxs []json.RawMessage
	for _, f := range fs {
		name := f.Name()
		if name == "README.md" {
			continue
		}
		bz, err := ioutil.ReadFile(path.Join(genTxPath, name))
		if err != nil {
			panic(err)
		}
		genTxs = append(genTxs, json.RawMessage(bz))
	}

	fmt.Println("-----------")
	fmt.Println("TOTAL gen txs", len(genTxs))

	// XXX: the app state is decoded using amino JSON (eg. ints are strings)
	// doesn't seem like we need to register anything though
	cdc := amino.NewCodec()

	genesisDoc := makeGenesisDoc(cdc, captainAccount[0], genesisAccounts, genTxs)
	// write the genesis file
	bz, err := cdc.MarshalJSON(genesisDoc)
	if err != nil {
		panic(err)
	}
	buf := bytes.NewBuffer([]byte{})
	err = json.Indent(buf, bz, "", "  ")
	if err != nil {
		panic(err)
	}
	err = ioutil.WriteFile(genesisFile, buf.Bytes(), 0600)
	if err != nil {
		panic(err)
	}
}

func fromBech32(address string) sdk.AccAddress {
	bech32PrefixAccAddr := "okchain"
	bz, err := sdk.GetFromBech32(address, bech32PrefixAccAddr)
	if err != nil {
		panic(err)
	}
	if len(bz) != sdk.AddrLen {
		panic("Incorrect address length")
	}
	return sdk.AccAddress(bz)
}

func accumulateBechContributors(fileName string, contribs map[string]float64) error {

	allocations := pkg.ListToMap(fileName)

	for addr, amt := range allocations {
		if _, ok := contribs[addr]; ok {
			fmt.Println("Duplicate addr", addr)
		}
		contribs[addr] += amt
	}
	return nil
}

// load a map of hex addresses and convert them to bech32
func accumulateContributors(fileName string, contribs map[string]float64) error {
	allocations := pkg.ObjToMap(fileName)

	for addr, amt := range allocations {
		if _, ok := contribs[addr]; ok {
			fmt.Println("Duplicate addr", addr)
		}
		contribs[addr] += amt
	}
	return nil
}

//----------------------------------------------------------
// AiB Data

type Account struct {
	Address string  `json:"addr"`
	Amount  float64 `json:"amount"`
	Lock    string  `json:"lock"`
}

type MultisigAccount struct {
	Address   string   `json:"addr"`
	Threshold int      `json:"threshold"`
	Pubs      []string `json:"pubs"`
	Amount    float64  `json:"amount"`
}

//---------------------------------------------------------------
// gaia accounts and genesis doc

// compose the gaia genesis accounts from the inputs,
// check total and for duplicates,
// sort by address
func makeGenesisAccounts(
	contribs map[string]float64,
	employees []Account,
	multisig MultisigAccount) []app.GenesisAccount {

	var genesisAccounts []app.GenesisAccount
	{
		// public, private, and icf contribs
		for addr, amt := range contribs {
			acc := app.GenesisAccount{
				Address: fromBech32(addr),
				Coins:   newCoins(amt),
			}
			genesisAccounts = append(genesisAccounts, acc)
		}

	}

	// sort the accounts
	sort.SliceStable(genesisAccounts, func(i, j int) bool {
		return strings.Compare(
			genesisAccounts[i].Address.String(),
			genesisAccounts[j].Address.String(),
		) < 0
	})

	return genesisAccounts
}

// check total atoms and no duplicates
func checkTotals(genesisAccounts []app.GenesisAccount) {
	// check uAtom total
	uAtomTotal := sdk.NewDec(0)
	for _, account := range genesisAccounts {
		uAtomTotal = uAtomTotal.Add(account.Coins[0].Amount)
	}

	if len(genesisAccounts) != addressGenesisTotal {
		panicStr := fmt.Sprintf("expected %d addresses, got %d addresses allocated in genesis", addressGenesisTotal, len(genesisAccounts))
		panic(panicStr)
	}

	// ensure no duplicates
	checkdupls := make(map[string]struct{})
	for _, acc := range genesisAccounts {
		if _, ok := checkdupls[acc.Address.String()]; ok {
			panic(fmt.Sprintf("Got duplicate: %v", acc.Address))
		}
		checkdupls[acc.Address.String()] = struct{}{}
	}
	if len(checkdupls) != len(genesisAccounts) {
		panic("length mismatch!")
	}
}

// json marshal the initial app state (accounts and gentx) and add them to the template
func makeGenesisDoc(cdc *amino.Codec, captainAccounts app.GenesisAccount, genesisAccounts []app.GenesisAccount, genTxs []json.RawMessage) *tmtypes.GenesisDoc {
	// read the template with the params
	genesisDoc, err := tmtypes.GenesisDocFromFile(genesisTemplate)
	if err != nil {
		panic(err)
	}
	// set genesis time
	genesisDoc.GenesisTime = timeGenesis

	// read the gaia state from the generic tendermint app state bytes
	// and populate with the accounts and gentxs
	var genesisState app.GenesisState
	err = cdc.UnmarshalJSON(genesisDoc.AppState, &genesisState)
	if err != nil {
		panic(err)
	}

	genesisState.Accounts = genesisAccounts
	genesisState.GenTxs = genTxs

	if len(genesisState.Token.Info) != 1 {
		panic(fmt.Errorf("No genesis denom!"))
	}
	genesisState.Token.Info[0].Owner = captainAccounts.Address

	// fix staking data
	genesisState.StakingData.Pool.NotBondedTokens = token.ToUnit(okbGenesisTotal) //atomToUAtomInt(okbGenesisTotal)
	genesisState.StakingData.Params.BondDenom = okbDenomination

	// marshal the gaia app state back to json and update the genesisDoc
	genesisStateJSON, err := cdc.MarshalJSON(genesisState)
	if err != nil {
		panic(err)
	}
	genesisDoc.AppState = genesisStateJSON

	return genesisDoc
}
