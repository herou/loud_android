package loud

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	pylonSDK "github.com/Pylons-tech/pylons/cmd/test"
	"github.com/Pylons-tech/pylons/x/pylons/msgs"
	"github.com/Pylons-tech/pylons/x/pylons/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func CreateCookbook(user User) (string, error) { // This is for afti develop mode automation test is only using
	t := GetTestingT()
	username := user.GetUserName()
	addr := pylonSDK.GetAccountAddr(username, t)
	sdkAddr, err := sdk.AccAddressFromBech32(addr)
	log.Println("sdkAddr, err := sdk.AccAddressFromBech32(addr)", sdkAddr, err)

	ccbMsg := msgs.NewMsgCreateCookbook(
		"tst_cookbook_name",                  // cbType.Name,
		fmt.Sprintf("%d", time.Now().Unix()), // cbType.ID,
		"addghjkllsdfdggdgjkkk",              // cbType.Description,
		"asdfasdfasdf",                       // cbType.Developer,
		"1.0.0",                              // cbType.Version,
		"a@example.com",                      // cbType.SupportEmail,
		0,                                    // cbType.Level,
		5,                                    // cbType.CostPerBlock,
		sdkAddr,                              // cbType.Sender,
	)
	txhash := pylonSDK.TestTxWithMsgWithNonce(t, ccbMsg, username, false)
	if AutomateInput {
		ok, err := CheckSignatureMatchWithAftiCli(t, txhash, user.GetPrivKey(), ccbMsg, username, false)
		if !ok || err != nil {
			log.Println("error checking afticli", ok, err)
			SomethingWentWrongMsg = "automation test failed, " + err.Error()
		}
	}
	user.SetLastTransaction(txhash)
	log.Println("ended sending transaction")
	return txhash, nil
}

func GetExtraPylons(user User) (string, error) {
	t := GetTestingT()
	username := user.GetUserName()
	addr := pylonSDK.GetAccountAddr(username, t)
	sdkAddr, err := sdk.AccAddressFromBech32(addr)
	log.Println("sdkAddr, err := sdk.AccAddressFromBech32(addr)", sdkAddr, err)
	extraPylonsMsg := msgs.NewMsgGetPylons(types.PremiumTier.Fee, sdkAddr)
	txhash := pylonSDK.TestTxWithMsgWithNonce(t, extraPylonsMsg, username, false)
	user.SetLastTransaction(txhash)
	log.Println("ended sending transaction")
	return txhash, nil
}

func GetInitialCoin(user User) (string, error) {
	rcpName := "LOUD's get initial coin recipe"
	itemIDs := []string{}

	return ExecuteRecipe(user, rcpName, itemIDs)
}

func Hunt(user User, item Item) (string, error) {

	defaultCharacter := user.GetDefaultCharacter()
	defaultCharacterID := ""
	if defaultCharacter != nil {
		defaultCharacterID = defaultCharacter.ID
	} else {
		return "", errors.New("character is required to hunt!")
	}
	rcpName := "LOUD's hunt without sword recipe"
	itemIDs := []string{defaultCharacterID}

	switch item.Name {
	case WOODEN_SWORD, COPPER_SWORD:
		rcpName = "LOUD's hunt with a sword recipe"
		itemIDs = []string{defaultCharacterID, item.ID}
	}

	return ExecuteRecipe(user, rcpName, itemIDs)
}

func BuyCharacter(user User, item Character) (string, error) {
	rcpName := ""
	switch item.Name {
	case TIGER_CHARACTER:
		rcpName = "LOUD's Get Character recipe"
	default:
		return "", errors.New("You are trying to buy character which is not in shop")
	}
	if item.Price > user.GetPylonAmount() {
		return "", errors.New("You don't have enough pylon to buy this character")
	}
	return ExecuteRecipe(user, rcpName, []string{})
}

func Buy(user User, item Item) (string, error) {
	rcpName := ""
	switch item.Name {
	case WOODEN_SWORD:
		if item.Level == 1 {
			rcpName = "LOUD's Wooden sword lv1 buy recipe"
		}
	case COPPER_SWORD:
		if item.Level == 1 {
			rcpName = "LOUD's Copper sword lv1 buy recipe"
		}
	default:
		return "", errors.New("You are trying to buy item which is not in shop")
	}
	if item.Price > user.GetGold() {
		return "", errors.New("You don't have enough gold to buy this item")
	}
	return ExecuteRecipe(user, rcpName, []string{})
}

func Sell(user User, item Item) (string, error) {
	itemIDs := []string{item.ID}

	rcpName := ""
	switch item.Name {
	case WOODEN_SWORD, COPPER_SWORD:
		rcpName = "LOUD's sword sell recipe"
	}
	return ExecuteRecipe(user, rcpName, itemIDs)
}

func Upgrade(user User, item Item) (string, error) {
	itemIDs := []string{item.ID}
	rcpName := ""
	switch item.Name {
	case WOODEN_SWORD:
		if item.Level == 1 {
			rcpName = "LOUD's Wooden sword lv1 to lv2 upgrade recipe"
		}
	case COPPER_SWORD:
		if item.Level == 1 {
			rcpName = "LOUD's Copper sword lv1 to lv2 upgrade recipe"
		}
	}
	if item.GetUpgradePrice() > user.GetGold() {
		return "", errors.New("You don't have enough gold to upgrade this item")
	}
	return ExecuteRecipe(user, rcpName, itemIDs)
}

func CreateBuyLoudTradeRequest(user User, loudEnterValue string, pylonEnterValue string) (string, error) {
	t := GetTestingT()
	loudValue, err := strconv.Atoi(loudEnterValue)
	if err != nil {
		return "", err
	}
	pylonValue, err := strconv.Atoi(pylonEnterValue)
	if err != nil {
		return "", err
	}

	eugenAddr := pylonSDK.GetAccountAddr(user.GetUserName(), nil)
	sdkAddr, err := sdk.AccAddressFromBech32(eugenAddr)

	inputCoinList := types.GenCoinInputList("loudcoin", int64(loudValue))

	outputCoins := sdk.Coins{sdk.NewInt64Coin("pylon", int64(pylonValue))}
	extraInfo := "created by loud game"

	createTrdMsg := msgs.NewMsgCreateTrade(
		inputCoinList,
		nil,
		outputCoins,
		nil,
		extraInfo,
		sdkAddr)
	log.Println("started sending transaction", user.GetUserName(), createTrdMsg)
	txhash := pylonSDK.TestTxWithMsgWithNonce(t, createTrdMsg, user.GetUserName(), false)
	user.SetLastTransaction(txhash)
	log.Println("ended sending transaction")
	return txhash, nil
}

func CreateSellLoudTradeRequest(user User, loudEnterValue string, pylonEnterValue string) (string, error) {
	t := GetTestingT()
	loudValue, err := strconv.Atoi(loudEnterValue)
	if err != nil {
		return "", err
	}
	pylonValue, err := strconv.Atoi(pylonEnterValue)
	if err != nil {
		return "", err
	}

	eugenAddr := pylonSDK.GetAccountAddr(user.GetUserName(), nil)
	sdkAddr, _ := sdk.AccAddressFromBech32(eugenAddr)

	inputCoinList := types.GenCoinInputList("pylon", int64(pylonValue))

	outputCoins := sdk.Coins{sdk.NewInt64Coin("loudcoin", int64(loudValue))}
	extraInfo := "created by loud game"

	createTrdMsg := msgs.NewMsgCreateTrade(
		inputCoinList,
		nil,
		outputCoins,
		nil,
		extraInfo,
		sdkAddr)
	log.Println("started sending transaction", user.GetUserName(), createTrdMsg)
	txhash := pylonSDK.TestTxWithMsgWithNonce(t, createTrdMsg, user.GetUserName(), false)
	user.SetLastTransaction(txhash)
	log.Println("ended sending transaction")
	return txhash, nil
}

func CreateBuySwordTradeRequest(user User, activeItem Item, pylonEnterValue string) (string, error) {
	// trade creator will get sword from pylon
	t := GetTestingT()

	itemInputs := GetItemInputsFromActiveItem(activeItem)

	pylonValue, err := strconv.Atoi(pylonEnterValue)
	if err != nil {
		return "", err
	}

	eugenAddr := pylonSDK.GetAccountAddr(user.GetUserName(), nil)
	sdkAddr, err := sdk.AccAddressFromBech32(eugenAddr)

	outputCoins := sdk.Coins{sdk.NewInt64Coin("pylon", int64(pylonValue))}
	extraInfo := "sword buy request created by loud game"

	createTrdMsg := msgs.NewMsgCreateTrade(
		nil,
		itemInputs,
		outputCoins,
		nil,
		extraInfo,
		sdkAddr)
	log.Println("started sending transaction", user.GetUserName(), createTrdMsg)
	txhash := pylonSDK.TestTxWithMsgWithNonce(t, createTrdMsg, user.GetUserName(), false)
	user.SetLastTransaction(txhash)
	log.Println("ended sending transaction")
	return txhash, nil
}

func CreateSellSwordTradeRequest(user User, activeItem Item, pylonEnterValue string) (string, error) {
	// trade creator will get pylon from sword
	t := GetTestingT()

	pylonValue, err := strconv.Atoi(pylonEnterValue)
	if err != nil {
		return "", err
	}

	eugenAddr := pylonSDK.GetAccountAddr(user.GetUserName(), nil)
	sdkAddr, _ := sdk.AccAddressFromBech32(eugenAddr)

	inputCoinList := types.GenCoinInputList("pylon", int64(pylonValue))
	itemOutputList, err := GetItemOutputFromActiveItem(activeItem)
	if err != nil {
		return "", err
	}

	extraInfo := "sword sell request created by loud game"

	createTrdMsg := msgs.NewMsgCreateTrade(
		inputCoinList,
		nil,
		nil,
		itemOutputList,
		extraInfo,
		sdkAddr)
	log.Println("started sending transaction", user.GetUserName(), createTrdMsg)
	txhash := pylonSDK.TestTxWithMsgWithNonce(t, createTrdMsg, user.GetUserName(), false)
	user.SetLastTransaction(txhash)
	log.Println("ended sending transaction")
	return txhash, nil
}

func FulfillTrade(user User, tradeID string) (string, error) {
	t := GetTestingT()
	eugenAddr := pylonSDK.GetAccountAddr(user.GetUserName(), nil)
	sdkAddr, _ := sdk.AccAddressFromBech32(eugenAddr)
	ffTrdMsg := msgs.NewMsgFulfillTrade(tradeID, sdkAddr, []string{})

	log.Println("started sending transaction", user.GetUserName(), ffTrdMsg)
	txhash := pylonSDK.TestTxWithMsgWithNonce(t, ffTrdMsg, user.GetUserName(), false)
	user.SetLastTransaction(txhash)
	log.Println("ended sending transaction")
	return txhash, nil
}
