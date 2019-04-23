## storeKey

### staking 
|key|value|number(key)|value details|value size|clean up|备注|   
|---|---|-------------|--------------|---------|--------|---|
|0x00                            |x/staking/types.Pool       |1    |无数组|<1k|无需清零                                |Pool                |
|0x11+OperatorAddr               |Tokens                     |N/A  |无数组|<1k|交易清理                                |LastValidatorsPower |
|0x12                            |Total Tokens               |1    |无数组|<1k|无需清零                                |LastTotalPower      |
|0x21+OperatorAddr               |x/staking/types.Validator  |N/A  |无数组|<1k|交易清理                                |Validator           |
|0x22+ConsensusAddr              |OperatorAddr               |N/A  |无数组|<1k|交易清理                                |Validator           |
|0x23+Tokens+OperatorAddr        |OperatorAddr               |N/A  |无数组|<1k|交易清理                                |Validator           |
|0x31+DelegatorAddr+ValidatorAddr|x/staking/types.Delegation |N/A  |无数组|<1k|交易清理                                |Delegation          |
|0x32+DelegatorAddr+ValidatorAddr|x/staking/types.UnbondingDelegation|N/A  |无数组|<1k|每个区块清理到期                 |UnbondingDelegation |
|0x33+ValidatorAddr+DelegatorAddr|nil|N/A                                  |无数组|<1k|每个区块清理到期                 |UnbondingDelegation|
|0x34+DelegatorAddr+ValidatorSrcAddr+ValidatorDstAddr|x/staking/types.Redelegation|N/A|无数组|<1k|每个区块清理到期      |Redelegation|
|0x35+ValidatorSrcAddr+ValidatorDstAddr+DelegatorAddr|nil|N/A                         |无数组|<1k|每个区块清理到期      |Redelegation|
|0x36+ValidatorDstAddr+ValidatorSrcAddr+DelegatorAddr|nil|N/A                         |无数组|<1k|每个区块清理到期      |Redelegation|
|0x41+Time|x/staking/types.[]DVPair|N/A|数组长度最多为单块最大交易数量|\>1k|每个区块清理到期                                |UnbondingDelegationQueue|
|0x42+Time|x/staking/types.[]DVVTriple|N/A|数组长度最多为单块最大交易数量|\>1k|每个区块清理到期                             |RedelegationQueue|
|0x43+Time|x/staking/types.[]ValAddress|N/A|数组长度最多为validator集合总数|\>1k|每个区块清理到期                          |ValidatorQueue|
|0x91+DelegatorAddr+ValidatorSrcAddr|x/staking/types.Redelegation|N/A|无数组|<1k|每个区块清理到期                       |Redelegation|
|0x92+BlockHeight|x/staking/types.[]DVVTriple|区块高度|数组长度最多为单块最大交易数量|\>1k|每个区块清理到期                  |RedelegationActionQueue|
|0x93+DelegatorAddr+ValidatorAddr|nil|N/A|无数组|<1k|每个周期清零                                                       |DelegatorPool|
   

### params 
|key|value|number(key)|value details|value size|clean up|备注|   
|---|---|-------------|--------------|---------|--------|---|
|ParamsSubspace("staking")  |x/staking/types.Params     |1|无数组|<1k|无需清零                  |Params         |


     
