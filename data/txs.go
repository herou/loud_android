package loud

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/Pylons-tech/LOUD/log"
	pylonSDK "github.com/Pylons-tech/pylons_sdk/cmd/test_utils"
	"github.com/Pylons-tech/pylons_sdk/x/pylons/msgs"
	"github.com/Pylons-tech/pylons_sdk/x/pylons/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// TextCreatedByLOUD is used to watermark recipes and trades that are generated by loud game
const TextCreatedByLOUD = "created by loud game"

// ItemBuyReqTrdInfo is used to watermark sword buy trade by loud game
const ItemBuyReqTrdInfo = "sword buy request created by loud game"

// ChrBuyReqTrdInfo is used to watermark character buy trade by loud game
const ChrBuyReqTrdInfo = "character buy request created by loud game"

// ItemSellReqTrdInfo is used to watermark sword sell trade by loud game
const ItemSellReqTrdInfo = "sword sell request created by loud game"

// ChrSellReqTrdInfo is used to watermark character sell request by loud game
const ChrSellReqTrdInfo = "character sell request created by loud game"

// CreateCookbook is a cookbook creation function and is for afti develop mode automation test is only using
func CreateCookbook(user User) (string, error) {
	t := GetTestingT()
	username := user.GetUserName()
	sdkAddr := GetSDKAddrFromUserName(username)

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

	txhash, _ := SendTxMsg(user, ccbMsg)
	if AutomateInput {
		ok, err := CheckSignatureMatchWithAftiCli(t, txhash, user.GetPrivKey(), ccbMsg, username, false)
		if !ok || err != nil {
			log.WithFields(log.Fields{
				"check_ok": ok,
				"error":    err,
			}).Warnln("error checking afticli")
			SomethingWentWrongMsg = "automation test failed, " + err.Error()
		}
	}
	return txhash, nil
}

// GetExtraPylons is used for getting extra pylons from faucet
func GetExtraPylons(user User) (string, error) {
	sdkAddr := GetSDKAddrFromUserName(user.GetUserName())
	extraPylonsMsg := msgs.NewMsgGetPylons(types.PremiumTier.Fee, sdkAddr)
	return SendTxMsg(user, extraPylonsMsg)
}

// BuyGoldWithPylons is used for buying gold from pylons
func BuyGoldWithPylons(user User) (string, error) {
	return ExecuteRecipe(user, RcpBuyGoldWithPylon, []string{})
}

// DevGetTestItems is used for getting developer test items
func DevGetTestItems(user User) (string, error) {
	return ExecuteRecipe(user, RcpGetTestItems, []string{})
}

// RunHuntRecipe is a helper function to call hunt recipes
func RunHuntRecipe(monsterName, rcpName string, user User) (string, error) {
	activeCharacter := user.GetActiveCharacter()
	activeCharacterID := ""
	if activeCharacter != nil {
		activeCharacterID = activeCharacter.ID
	} else {
		return "", errors.New("character is required to hunt rabbits")
	}

	user.SetFightMonster(monsterName)
	activeWeapon := user.GetFightWeapon()
	itemIDs := []string{activeCharacterID}
	if activeWeapon != nil {
		itemIDs = []string{activeCharacterID, activeWeapon.ID}
	}

	return ExecuteRecipe(user, rcpName, itemIDs)
}

// HuntRabbits is a function to execute rabbit hunt recipe
func HuntRabbits(user User) (string, error) {
	return RunHuntRecipe(TextRabbit, RcpHuntRabbits, user)
}

// FightTroll is a function to execute fight troll recipe
func FightTroll(user User) (string, error) {
	return RunHuntRecipe(TextTroll, RcpFightTroll, user)
}

// FightWolf is a function to execute fight wolf recipe
func FightWolf(user User) (string, error) { // 🐺
	return RunHuntRecipe(TextWolf, RcpFightWolf, user)
}

// FightGoblin is a function to execute fight goblin recipe
func FightGoblin(user User) (string, error) { // 👺
	return RunHuntRecipe(TextGoblin, RcpFightGoblin, user)
}

// FightGiant is a function to execute fight giant recipe
func FightGiant(user User, tarBonus int) (string, error) { // 🗿
	rcp := RcpFightGiant
	switch tarBonus {
	case FireSpecial:
		rcp = RcpFightFireGiant
	case IceSpecial:
		rcp = RcpFightIceGiant
	case AcidSpecial:
		rcp = RcpFightAcidGiant
	}
	return RunHuntRecipe(TextGiant, rcp, user)
}

// FightDragonFire is a function to execute fight fire dragon recipe
func FightDragonFire(user User) (string, error) { // 🦐
	return RunHuntRecipe(TextDragonFire, RcpFightDragonFire, user)
}

// FightDragonIce is a function to execute fight ice dragon recipe
func FightDragonIce(user User) (string, error) { // 🦈
	return RunHuntRecipe(TextDragonIce, RcpFightDragonIce, user)
}

// FightDragonAcid is a function to execute fight acid dragon recipe
func FightDragonAcid(user User) (string, error) { // 🐊
	return RunHuntRecipe(TextDragonAcid, RcpFightDragonAcid, user)
}

// FightDragonUndead is a function to execute fight undead dragon recipe
func FightDragonUndead(user User) (string, error) { // 🐉
	return RunHuntRecipe(TextDragonUndead, RcpFightDragonUndead, user)
}

// BuyCharacter is a function to execute buy character recipe
func BuyCharacter(user User, ch Character) (string, error) {
	rcpName := ""
	switch ch.Name {
	case TextTigerChr:
		rcpName = RcpBuyCharacter
	default:
		return "", errors.New("You are trying to buy character which is not in shop")
	}
	if ch.Price > user.GetPylonAmount() {
		return "", errors.New("You don't have enough pylon to buy this character")
	}
	return ExecuteRecipe(user, rcpName, []string{})
}

// RenameCharacter is a function to execute rename character recipe
func RenameCharacter(user User, ch Character, newName string) (string, error) {
	t := GetTestingT()
	addr := pylonSDK.GetAccountAddr(user.GetUserName(), GetTestingT())
	sdkAddr, _ := sdk.AccAddressFromBech32(addr)
	renameMsg := msgs.NewMsgUpdateItemString(ch.ID, "Name", newName, sdkAddr)
	txhash, err := pylonSDK.TestTxWithMsgWithNonce(t, renameMsg, user.GetUserName(), false)
	if err != nil {
		return "", fmt.Errorf("error sending transaction; %+v", err)
	}
	user.SetLastTransaction(txhash, Sprintf("rename character from %s to %s", ch.Name, newName))
	return txhash, nil
}

// BuyItem is a function to execute buy item recipe
func BuyItem(user User, item Item) (string, error) {
	rcpName := ""
	itemIDs := []string{}
	switch item.Name {
	case WoodenSword:
		if item.Level == 1 {
			rcpName = RcpBuyWoodenSword
		}
	case CopperSword:
		if item.Level == 1 {
			rcpName = RcpBuyCopperSword
		}
	case SilverSword:
		if item.Level == 1 {
			rcpName = RcpBuySilverSword
			itemIDs = []string{user.InventoryItemIDByName(GoblinEar)}
		}
	case BronzeSword:
		if item.Level == 1 {
			rcpName = RcpBuyBronzeSword
			itemIDs = []string{user.InventoryItemIDByName(WolfTail)}
		}
	case IronSword:
		if item.Level == 1 {
			rcpName = RcpBuyIronSword
			itemIDs = []string{user.InventoryItemIDByName(TrollToes)}
		}
	case AngelSword:
		if item.Level == 1 {
			rcpName = RcpBuyAngelSword
			itemIDs = []string{
				user.InventoryItemIDByName(DropDragonFire),
				user.InventoryItemIDByName(DropDragonIce),
				user.InventoryItemIDByName(DropDragonAcid),
			}
		}
	default:
		return "", errors.New("You are trying to buy item which is not in shop")
	}
	if item.Price > user.GetGold() {
		return "", errors.New("You don't have enough gold to buy this item")
	}
	return ExecuteRecipe(user, rcpName, itemIDs)
}

// SellItem is a function to execute sell item recipe
func SellItem(user User, item Item) (string, error) {
	itemIDs := []string{item.ID}

	rcpName := ""
	if item.Value > 0 {
		rcpName = RcpSellSword
	}
	return ExecuteRecipe(user, rcpName, itemIDs)
}

// UpgradeItem is a function to execute upgrade item recipe
func UpgradeItem(user User, item Item) (string, error) {
	itemIDs := []string{item.ID}
	rcpName := ""
	switch item.Name {
	case WoodenSword:
		if item.Level == 1 {
			rcpName = RcpWoodenSwordUpgrade
		}
	case CopperSword:
		if item.Level == 1 {
			rcpName = RcpCopperSwordUpgrade
		}
	}
	if item.GetUpgradePrice() > user.GetGold() {
		return "", errors.New("You don't have enough gold to upgrade this item")
	}
	return ExecuteRecipe(user, rcpName, itemIDs)
}

// CreateBuyGoldTrdReq is a function to create gold buying trade request
func CreateBuyGoldTrdReq(user User, goldEnterValue string, pylonEnterValue string) (string, error) {
	loudValue, err := strconv.Atoi(goldEnterValue)
	if err != nil {
		return "", err
	}
	if loudValue == 0 {
		return "", errors.New("gold amount shouldn't be zero to be a valid trading")
	}
	pylonValue, err := strconv.Atoi(pylonEnterValue)
	if err != nil {
		return "", err
	}
	if pylonValue == 0 {
		return "", errors.New("pylon amount shouldn't be zero to be a valid trading")
	}

	sdkAddr := GetSDKAddrFromUserName(user.GetUserName())

	inputCoinList := types.GenCoinInputList("loudcoin", int64(loudValue))

	outputCoins := sdk.Coins{sdk.NewInt64Coin("pylon", int64(pylonValue))}
	extraInfo := TextCreatedByLOUD

	createTrdMsg := msgs.NewMsgCreateTrade(
		inputCoinList,
		nil,
		outputCoins,
		nil,
		extraInfo,
		sdkAddr)
	return SendTxMsg(user, createTrdMsg)
}

// CreateSellGoldTrdReq is a function to create sell gold trade request
func CreateSellGoldTrdReq(user User, goldEnterValue string, pylonEnterValue string) (string, error) {
	loudValue, err := strconv.Atoi(goldEnterValue)
	if err != nil {
		return "", err
	}
	if loudValue == 0 {
		return "", errors.New("gold amount shouldn't be zero to be a valid trading")
	}
	pylonValue, err := strconv.Atoi(pylonEnterValue)
	if err != nil {
		return "", err
	}
	if pylonValue == 0 {
		return "", errors.New("pylon amount shouldn't be zero to be a valid trading")
	}

	sdkAddr := GetSDKAddrFromUserName(user.GetUserName())

	inputCoinList := types.GenCoinInputList("pylon", int64(pylonValue))

	outputCoins := sdk.Coins{sdk.NewInt64Coin("loudcoin", int64(loudValue))}
	extraInfo := TextCreatedByLOUD

	createTrdMsg := msgs.NewMsgCreateTrade(
		inputCoinList,
		nil,
		outputCoins,
		nil,
		extraInfo,
		sdkAddr)
	return SendTxMsg(user, createTrdMsg)
}

// CreateBuyItemTrdReq is a function to create buy item trade request
func CreateBuyItemTrdReq(user User, itspec ItemSpec, pylonEnterValue string) (string, error) {
	// trade creator will get sword from pylon

	itemInputs := GetItemInputsFromItemSpec(itspec)

	pylonValue, err := strconv.Atoi(pylonEnterValue)
	if err != nil {
		return "", err
	}
	if pylonValue == 0 {
		return "", errors.New("pylon amount shouldn't be zero to be a valid trading")
	}

	sdkAddr := GetSDKAddrFromUserName(user.GetUserName())

	outputCoins := sdk.Coins{sdk.NewInt64Coin("pylon", int64(pylonValue))}
	extraInfo := ItemBuyReqTrdInfo

	createTrdMsg := msgs.NewMsgCreateTrade(
		nil,
		itemInputs,
		outputCoins,
		nil,
		extraInfo,
		sdkAddr)
	return SendTxMsg(user, createTrdMsg)
}

// CreateSellItemTrdReq is a function to create sell item trade request
func CreateSellItemTrdReq(user User, activeItem Item, pylonEnterValue string) (string, error) {
	// trade creator will get pylon from sword

	pylonValue, err := strconv.Atoi(pylonEnterValue)
	if err != nil {
		return "", err
	}
	if pylonValue == 0 {
		return "", errors.New("pylon amount shouldn't be zero to be a valid trading")
	}

	sdkAddr := GetSDKAddrFromUserName(user.GetUserName())

	inputCoinList := types.GenCoinInputList("pylon", int64(pylonValue))
	itemOutputList, err := GetItemOutputFromActiveItem(activeItem)
	if err != nil {
		return "", err
	}

	extraInfo := ItemSellReqTrdInfo

	createTrdMsg := msgs.NewMsgCreateTrade(
		inputCoinList,
		nil,
		nil,
		itemOutputList,
		extraInfo,
		sdkAddr)
	return SendTxMsg(user, createTrdMsg)
}

// CreateBuyCharacterTrdReq is a function to create buy character trade request
func CreateBuyCharacterTrdReq(user User, chspec CharacterSpec, pylonEnterValue string) (string, error) {
	// trade creator will get character from pylon

	itemInputs := GetItemInputsFromCharacterSpec(chspec)

	pylonValue, err := strconv.Atoi(pylonEnterValue)
	if err != nil {
		return "", err
	}
	if pylonValue == 0 {
		return "", errors.New("pylon amount shouldn't be zero to be a valid trading")
	}

	sdkAddr := GetSDKAddrFromUserName(user.GetUserName())

	outputCoins := sdk.Coins{sdk.NewInt64Coin("pylon", int64(pylonValue))}
	extraInfo := ChrBuyReqTrdInfo

	createTrdMsg := msgs.NewMsgCreateTrade(
		nil,
		itemInputs,
		outputCoins,
		nil,
		extraInfo,
		sdkAddr)
	return SendTxMsg(user, createTrdMsg)
}

// CreateSellCharacterTrdReq is a function to create sell character trade request
func CreateSellCharacterTrdReq(user User, activeCharacter Character, pylonEnterValue string) (string, error) {
	// trade creator will get pylon from character

	pylonValue, err := strconv.Atoi(pylonEnterValue)
	if err != nil {
		return "", err
	}
	if pylonValue == 0 {
		return "", errors.New("pylon amount shouldn't be zero to be a valid trading")
	}

	sdkAddr := GetSDKAddrFromUserName(user.GetUserName())

	inputCoinList := types.GenCoinInputList("pylon", int64(pylonValue))
	itemOutputList, err := GetItemOutputFromActiveCharacter(activeCharacter)
	if err != nil {
		return "", err
	}

	extraInfo := ChrSellReqTrdInfo

	createTrdMsg := msgs.NewMsgCreateTrade(
		inputCoinList,
		nil,
		nil,
		itemOutputList,
		extraInfo,
		sdkAddr)
	return SendTxMsg(user, createTrdMsg)
}

// FulfillTrade is a function to create fulfill trade by id
func FulfillTrade(user User, tradeID string, itemIDs []string) (string, error) {
	sdkAddr := GetSDKAddrFromUserName(user.GetUserName())
	ffTrdMsg := msgs.NewMsgFulfillTrade(tradeID, sdkAddr, itemIDs)

	return SendTxMsg(user, ffTrdMsg)
}

// CancelTrade is a function to cancel trade by id
func CancelTrade(user User, tradeID string) (string, error) {
	sdkAddr := GetSDKAddrFromUserName(user.GetUserName())
	ccTrdMsg := msgs.NewMsgDisableTrade(tradeID, sdkAddr)

	return SendTxMsg(user, ccTrdMsg)
}
