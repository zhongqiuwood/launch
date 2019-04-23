package version

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"sync"
)

const (
	versionBitCoin    = "VersionBitCoin" // okdex v1
	versionOKCoin     = "VersionOKCoin"  // okdex v2

	// ...
	VersionXXB     = "VersionXXB"        // okdex vn
)

type VersionController struct {}

var controller *VersionController
var once sync.Once

func GetVersionController() (*VersionController) {
	once.Do(func() {
		controller = &VersionController{}
	})
	return controller
}

type AppModule interface {
	NewHandlerVersion1(ctx sdk.Context, msg sdk.Msg, keeper interface{}) sdk.Result
	EndBlockerVersion1(ctx sdk.Context, keeper interface{}) sdk.Tags

	NewHandlerVersion2(ctx sdk.Context, msg sdk.Msg, keeper interface{}) sdk.Result
	EndBlockerVersion2(ctx sdk.Context, keeper interface{}) sdk.Tags
}

func (v *VersionController) NewHandler(ctx sdk.Context, msg sdk.Msg,
	keeper interface{}, module AppModule) sdk.Result {

	blockeight := ctx.BlockHeight()
	version := v.getVersion(blockeight)

	switch version {
	case versionBitCoin:
		return module.NewHandlerVersion1(ctx, msg, keeper)
	case versionOKCoin:
		return module.NewHandlerVersion2(ctx, msg, keeper)
	default:
		errMsg := fmt.Sprintf("Unrecognized version: %s", version)
		return sdk.ErrUnknownRequest(errMsg).Result()
	}
}


func (v *VersionController) EndBlocker(ctx sdk.Context,
	keeper interface{}, module AppModule) sdk.Tags {

	blockeight := ctx.BlockHeight()
	version := v.getVersion(blockeight)

	switch version {
	case versionBitCoin:
		return module.EndBlockerVersion1(ctx, keeper)
	case versionOKCoin:
		return module.EndBlockerVersion2(ctx, keeper)
	default:
		//errMsg := fmt.Sprintf("Unrecognized version: %s", version)
		return sdk.NewTags()
	}
}

func (v *VersionController) getVersion(blockeight int64) string {

	if blockeight < 1000000 {
		return versionBitCoin
	}

	return versionOKCoin
}