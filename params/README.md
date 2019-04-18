## Genesis Parameters

launch 用于生成okchain网络使用的genesis file

`genesis_template.json`是genesis file的模板，launch先解析`genesis_template.json`中的数据结构，然后根据`launch/account/`和`launch/gentx/`中的配置文件修改`genesisDoc.AppState`生成`genesis file`

所以，***如果修改`AppState`中的字段（包括添加字段、删除字段、修改字段名）或者修改参数初始值***，需要在`genesis_template.json`中对应修改。

`AppState`的结构：
```go
type GenesisState struct {
	AuthData     auth.GenesisState     `json:"auth"`
	BankData     bank.GenesisState     `json:"bank"`
	Accounts     []GenesisAccount      `json:"accounts"`
	DistrData    distr.GenesisState    `json:"distr"`
	StakingData  staking.GenesisState  `json:"staking"`
	SlashingData slashing.GenesisState `json:"slashing"`
	GovData      gov.GenesisState      `json:"gov"`
	MintData     mint.GenesisState     `json:"mint"`
	GenTxs       []json.RawMessage     `json:"gentxs"`
	Order        order.GenesisState    `json:"order"`
	Token        token.GenesisState    `json:"token"`
}
```
