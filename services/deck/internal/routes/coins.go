package routes

import (
	"crypto/md5"
	"encoding/hex"
	"net/http"
	"strconv"

	"github.com/osuTitanic/titanic/internal/activity"
	"github.com/osuTitanic/titanic/internal/constants"
	"github.com/osuTitanic/titanic/services/deck/internal/server"
)

const (
	initialCoinBalance  int64 = 10
	rechargeCoinBalance int64 = 99
	coinChecksumSecret        = "osuycoins"
)

// /web/coins.php -> osu! coins april fools event
//
// Learn more:
// https://osu.ppy.sh/wiki/en/History_of_osu%21/April_Fools/osu%21coin
func Coins(ctx *server.Context) {
	user, ok := ctx.HandleUserAuthenticationSimple("u", "h", false)
	if !ok {
		// Response already set by function
		return
	}

	username := ctx.QueryValue("u")
	count, err := ctx.QueryValueInt("c")
	if err != nil {
		ctx.Response.WriteHeader(http.StatusBadRequest)
		return
	}
	key := "bancho:coins:" + strconv.Itoa(user.Id)

	// Set to initial coin balance if it doesn't exist (NX)
	if err := ctx.State.Redis.SetNX(ctx.Request.Context(), key, initialCoinBalance, 0).Err(); err != nil {
		ctx.Logger.Error("Failed to initialize coin balance", "user_id", user.Id, "error", err)
		ctx.Response.WriteHeader(http.StatusInternalServerError)
		return
	}

	checksumPayload := username + strconv.Itoa(count) + coinChecksumSecret
	checksum := md5.Sum([]byte(checksumPayload))
	if hex.EncodeToString(checksum[:]) != ctx.QueryValue("cs") {
		ctx.Logger.Warn("Invalid coin checksum", "user_id", user.Id, "username", username)
		ctx.Response.WriteHeader(http.StatusBadRequest)
		return
	}

	action := ctx.QueryValue("action")
	amount := int64(0)
	var coins int64

	switch action {
	case "earn":
		amount = 1
		coins, err = ctx.State.Redis.IncrBy(ctx.Request.Context(), key, amount).Result()
	case "use":
		amount = -1
		coins, err = ctx.State.Redis.IncrBy(ctx.Request.Context(), key, amount).Result()
	case "recharge":
		amount = rechargeCoinBalance
		coins = rechargeCoinBalance
		err = ctx.State.Redis.Set(ctx.Request.Context(), key, coins, 0).Err()
	default:
		coins, err = ctx.State.Redis.Get(ctx.Request.Context(), key).Int64()
	}
	if err != nil {
		ctx.Logger.Error("Failed to update coin balance", "user_id", user.Id, "action", action, "error", err)
		ctx.Response.WriteHeader(http.StatusInternalServerError)
		return
	}

	activityType := constants.ActivityOsuCoinsReceived
	if action == "use" {
		activityType = constants.ActivityOsuCoinsUsed
	}

	err = activity.Submit(
		ctx.State,
		user.Id,
		nil, // mode independent
		activityType,
		map[string]any{
			"username": user.Name,
			"amount":   amount,
			"coins":    coins,
		},
		false, // should not be sent to #announce
		true,  // should be hidden in user profile
	)
	if err != nil {
		ctx.Logger.Warn("Failed to broadcast activity", "user_id", user.Id, "error", err)
	}

	ctx.Logger.Info("Updated osu! coin balance", "user_id", user.Id, "action", action, "amount", amount, "coins", coins)
	ctx.RenderText(http.StatusOK, strconv.FormatInt(coins, 10))
}
