## StoreKey

### 1. tokenStoreKey 

KVStoreKey的name是"token"

|key                              | value                      | number(key)          | value details | Value size | Clean up |  备注   |
|:-------------------------------:|:-------------------------:|:------:|:------:|:------:|:------:|:------:|
| symbol | struct x/token.Token | 发行币的数量 | Token具体结构体如下 | <1k | 无 | 存的token的信息 |

token的value是经过Codec.MustMarshalBinaryBare序列化后的[]byte，token的结构体如下:

```go
type Token struct {
	Name           string         `json:"name"`							// token的名字
	Symbol         string         `json:"symbol"`						// token的唯一标识
	OriginalSymbol string         `json:"original_symbol"`	// token的原始标识
	TotalSupply    int64          `json:"total_supply"`			// token的总量
	Owner          sdk.AccAddress `json:"owner"`						// token的所有者
	Mintable       bool           `json:"mintable"`					// token是否可以增发
}
```

### 2. freezeStoreKey 

KVStoreKey的name是"freeze"

|   key   |      value       |   number(key)    | value details |               Value size               | Clean up |      备注       |
| :-----: | :--------------: | :--------------: | :-----------: | :------------------------------------: | :------: | :-------------: |
| address | struct sdk.Coins | 有冻结币的用户数 |  Coins结构体  | 可能会超过1K，取决于用户冻结的币的数量 |    无    | 存的token的信息 |

token的value是Coins经过Codec.MustMarshalBinaryBare序列化后的[]byte.

## 3. lockStoreKey

KVStoreKey的name是"lock"

|   key   |      value       |    number(key)     | value details |               Value size               | Clean up |      备注       |
| :-----: | :--------------: | :----------------: | :-----------: | :------------------------------------: | :------: | :-------------: |
| address | struct sdk.Coins | 有锁定的币的用户数 |  Coins结构体  | 可能会超过1K，取决于用户锁定的币的数量 |    无    | 存的token的信息 |

token的value是Coins经过Codec.MustMarshalBinaryBare序列化后的[]byte.

## 4. tokenPairStoreKey

KVStoreKey的name是"token_pair"

|                 key                  |          value           |   number(key)    |    value details    | Value size |         Clean up         |      备注       |
| :----------------------------------: | :----------------------: | :--------------: | :-----------------: | :--------: | :----------------------: | :-------------: |
| BaseAssetSymbol+"_"+QuoteAssetSymbol | struct x/token.TokenPair | 上交易所的币对数 | TokenPair结构体如下 |    <1k     | 有接口删除上交易所的币对 | 存的token的信息 |

tokenPair的value是经过Codec.MustMarshalBinaryBare序列化后的[]byte，tokenPair的结构体如下:

```go
type TokenPair struct {
	BaseAssetSymbol  string  `json:"base_asset_symbol"`		// 基础货币
	QuoteAssetSymbol string  `json:"quote_asset_symbol"`	// 报价货币
	InitPrice        sdk.Dec `json:"price"`							  // 价格
	MaxPriceDigit    int64   `json:"max_price_digit"`	 	  // 最大交易价格的小数点位数
	MaxQuantityDigit int64   `json:"max_size_digit"`		  // 最大交易数量的小数点位数
	MinQuantity      sdk.Dec `json:"min_trade_size"`		  // 最小交易数量
}
```

## 5. feeDetailStoreKey

KVStoreKey的name是"fee_detail"

|            key            |       value        |  number(key)   | value details | Value size | Clean up |   备注    |
| :-----------------------: | :----------------: | :------------: | :-----------: | :--------: | :------: | :-------: |
| feeDetails:$(blockHeight) | struct []FeeDetail | 与区块高度对应 | FeeDetail如下 |    <1K     |    无    | feeDetail |

feeDetailStoreKey的value是[]FeeDetail经过Codec.MustMarshalBinaryBare序列化后的[]byte，FeeDetail的结构体如下:

```go
type FeeDetail struct {
	Address   string `gorm:"index;type:varchar(80)" json:"address"`								// 地址
	Fee       string `json:"fee"`																									// fee
	FeeType   string `json:"feeType"` // transfer, deal, etc. see common/const.go	// fee的类型
	Timestamp int64  `gorm:"index;type:int64" json:"timestamp"`										// 时间
}
```

## Http api

|       Url       | Method |      读key       |
| :-------------: | :----: | :--------------: |
|     /products     |  GET   | token_pair(遍历) |
|     /tokens      |  GET   |   token(遍历)    |
| /token/{symbol} |  GET   |  token: symbol   |

