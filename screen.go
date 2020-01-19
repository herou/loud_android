package loud

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	pylonSDK "github.com/MikeSofaer/pylons/cmd/test"
	"github.com/MikeSofaer/pylons/x/pylons/handlers"
	"github.com/ahmetb/go-cursor"
	"github.com/gliderlabs/ssh"
	"github.com/mgutz/ansi"
	"github.com/nsf/termbox-go"

	terminal "github.com/wayneashleyberry/terminal-dimensions"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

// Screen represents a UI screen.
type Screen interface {
	SaveGame()
	UpdateBlockHeight(int64)
	SetScreenSize(int, int)
	HandleInputKey(termbox.Event)
	Render()
	Reset()
}

type ScreenStatus int

const (
	SHOW_LOCATION ScreenStatus = iota
	SELECT_SELL_ITEM
	WAIT_SELL_PROCESS
	RESULT_SELL_FINISH
	SELECT_BUY_ITEM
	WAIT_BUY_PROCESS
	RESULT_BUY_FINISH
	SELECT_HUNT_ITEM
	WAIT_HUNT_PROCESS
	RESULT_HUNT_FINISH
	SELECT_UPGRADE_ITEM
	WAIT_UPGRADE_PROCESS
	RESULT_UPGRADE_FINISH
)

type GameScreen struct {
	world          World
	user           User
	screenSize     ssh.Window
	activeItem     Item
	blockHeight    int64
	txFailReason   string
	txResult       handlers.ExecuteRecipeSerialize
	refreshed      bool
	scrStatus      ScreenStatus
	colorCodeCache map[string](func(string) string)
}

const allowMouseInputAndHideCursor string = "\x1b[?1003h\x1b[?25l"
const resetScreen string = "\x1bc"
const ellipsis = "…"
const hpon = "◆"
const hpoff = "◇"
const bgcolor = 232

var gameLanguage string = "en"

var shopItems []Item = []Item{
	Item{
		ID:    "001",
		Name:  "Wooden sword",
		Level: 1,
		Price: 100,
	},
	Item{
		ID:    "002",
		Name:  "Copper sword",
		Level: 1,
		Price: 250,
	},
}

func localize(key string) string {
	bundle := i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("json", json.Unmarshal)
	bundle.MustLoadMessageFile("en.json")
	bundle.MustLoadMessageFile("es.json")

	loc := i18n.NewLocalizer(bundle, gameLanguage)

	translate := loc.MustLocalize(
		&i18n.LocalizeConfig{
			MessageID: key,
			PluralCount: 1,
		})
	return translate
}

func truncateRight(message string, width int) string {
	if utf8.RuneCountInString(message) < width {
		fmtString := fmt.Sprintf("%%-%vs", width)

		return fmt.Sprintf(fmtString, message)
	}
	return string([]rune(message)[0:width-1]) + ellipsis
}

func truncateLeft(message string, width int) string {
	if utf8.RuneCountInString(message) < width {
		fmtString := fmt.Sprintf("%%-%vs", width)

		return fmt.Sprintf(fmtString, message)
	}
	strLen := utf8.RuneCountInString(message)
	return ellipsis + string([]rune(message)[strLen-width:strLen-1])
}

func justifyRight(message string, width int) string {
	if utf8.RuneCountInString(message) < width {
		fmtString := fmt.Sprintf("%%%vs", width)

		return fmt.Sprintf(fmtString, message)
	}
	strLen := utf8.RuneCountInString(message)
	return ellipsis + string([]rune(message)[strLen-width:strLen-1])
}

func centerText(message, pad string, width int) string {
	if utf8.RuneCountInString(message) > width {
		return truncateRight(message, width)
	}
	leftover := width - utf8.RuneCountInString(message)
	left := leftover / 2
	right := leftover - left

	if pad == "" {
		pad = " "
	}

	leftString := ""
	for utf8.RuneCountInString(leftString) <= left && utf8.RuneCountInString(leftString) <= right {
		leftString += pad
	}

	return fmt.Sprintf("%s%s%s", string([]rune(leftString)[0:left]), message, string([]rune(leftString)[0:right]))
}

func (screen *GameScreen) SetScreenSize(Width, Height int) {
	screen.screenSize = ssh.Window{
		Width:  Width,
		Height: Height,
	}
	screen.refreshed = false
}

func (screen *GameScreen) colorFunc(color string) func(string) string {
	_, ok := screen.colorCodeCache[color]

	if !ok {
		screen.colorCodeCache[color] = ansi.ColorFunc(color)
	}

	return screen.colorCodeCache[color]
}

func (screen *GameScreen) drawBox(x, y, width, height int) {
	color := ansi.ColorCode(fmt.Sprintf("255:%v", bgcolor))

	for i := 1; i < width; i++ {
		io.WriteString(os.Stdout, fmt.Sprintf("%s%s─", cursor.MoveTo(y, x+i), color))
		io.WriteString(os.Stdout, fmt.Sprintf("%s%s─", cursor.MoveTo(y+height, x+i), color))
	}

	for i := 1; i < height; i++ {
		midString := fmt.Sprintf("%%s%%s│%%%vs│", (width - 1))
		io.WriteString(os.Stdout, fmt.Sprintf("%s%s│", cursor.MoveTo(y+i, x), color))
		io.WriteString(os.Stdout, fmt.Sprintf("%s%s│", cursor.MoveTo(y+i, x+width), color))
		io.WriteString(os.Stdout, fmt.Sprintf(midString, cursor.MoveTo(y+i, x), color, " "))
	}

	io.WriteString(os.Stdout, fmt.Sprintf("%s%s╭", cursor.MoveTo(y, x), color))
	io.WriteString(os.Stdout, fmt.Sprintf("%s%s╰", cursor.MoveTo(y+height, x), color))
	io.WriteString(os.Stdout, fmt.Sprintf("%s%s╮", cursor.MoveTo(y, x+width), color))
	io.WriteString(os.Stdout, fmt.Sprintf("%s%s╯", cursor.MoveTo(y+height, x+width), color))
}

func (screen *GameScreen) drawFill(x, y, width, height int) {
	color := ansi.ColorCode(fmt.Sprintf("0:%v", bgcolor))

	midString := fmt.Sprintf("%%s%%s%%%vs", (width))
	for i := 0; i <= height; i++ {
		io.WriteString(os.Stdout, fmt.Sprintf(midString, cursor.MoveTo(y+i, x), color, " "))
	}
}

func (screen *GameScreen) drawProgressMeter(min, max, fgcolor, bgcolor, width uint64) string {
	var blink bool
	if min > max {
		min = max
		blink = true
	}
	proportion := float64(float64(min) / float64(max))
	if math.IsNaN(proportion) {
		proportion = 0.0
	} else if proportion < 0.05 {
		blink = true
	}
	onWidth := uint64(float64(width) * proportion)
	offWidth := uint64(float64(width) * (1.0 - proportion))

	onColor := screen.colorFunc(fmt.Sprintf("%v:%v", fgcolor, bgcolor))
	offColor := onColor

	if blink {
		onColor = screen.colorFunc(fmt.Sprintf("%v+B:%v", fgcolor, bgcolor))
	}

	if (onWidth + offWidth) > width {
		onWidth = width
		offWidth = 0
	} else if (onWidth + offWidth) < width {
		onWidth += width - (onWidth + offWidth)
	}

	on := ""
	off := ""

	for i := 0; i < int(onWidth); i++ {
		on += hpon
	}

	for i := 0; i < int(offWidth); i++ {
		off += hpoff
	}

	return onColor(on) + offColor(off)
}

func (screen *GameScreen) drawVerticalLine(x, y, height int) {
	color := ansi.ColorCode(fmt.Sprintf("255:%v", bgcolor))
	for i := 1; i < height; i++ {
		io.WriteString(os.Stdout, fmt.Sprintf("%s%s│", cursor.MoveTo(y+i, x), color))
	}

	io.WriteString(os.Stdout, fmt.Sprintf("%s%s┬", cursor.MoveTo(y, x), color))
	io.WriteString(os.Stdout, fmt.Sprintf("%s%s┴", cursor.MoveTo(y+height, x), color))
}

func (screen *GameScreen) drawHorizontalLine(x, y, width int) {
	color := ansi.ColorCode(fmt.Sprintf("255:%v", bgcolor))
	for i := 1; i < width; i++ {
		io.WriteString(os.Stdout, fmt.Sprintf("%s%s─", cursor.MoveTo(y, x+i), color))
	}

	io.WriteString(os.Stdout, fmt.Sprintf("%s%s├", cursor.MoveTo(y, x), color))
	io.WriteString(os.Stdout, fmt.Sprintf("%s%s┤", cursor.MoveTo(y, x+width), color))
}

func (screen *GameScreen) redrawBorders() {
	io.WriteString(os.Stdout, ansi.ColorCode(fmt.Sprintf("255:%v", bgcolor)))
	screen.drawBox(1, 1, screen.screenSize.Width-1, screen.screenSize.Height-1)
	screen.drawVerticalLine(screen.screenSize.Width/2-2, 1, screen.screenSize.Height)

	y := screen.screenSize.Height
	if y < 20 {
		y = 5
	} else {
		y = (y / 2) - 2
	}
	screen.drawHorizontalLine(1, y+2, screen.screenSize.Width/2-3)
	screen.drawHorizontalLine(1, screen.screenSize.Height-2, screen.screenSize.Width/2-3)
}

func (screen *GameScreen) renderUserCommands() {

	infoLines := []string{}
	switch screen.scrStatus {
	case SHOW_LOCATION:
		cmdMap := map[UserLocation]string{
			HOME:   localize("home"),
			FOREST: localize("forest"),
			SHOP:   localize("shop"),
			SETTINGS: localize("settings"),
		}
		cmdString := cmdMap[screen.user.GetLocation()]
		infoLines = strings.Split(cmdString, "\n")
	case SELECT_BUY_ITEM:
		for idx, item := range shopItems {
			infoLines = append(infoLines, fmt.Sprintf("%d) %s Lv%d", idx+1, localize(item.Name), item.Level))
		}
		infoLines = append(infoLines, localize("C)ancel"))
	case SELECT_SELL_ITEM:
		userItems := screen.user.InventoryItems()
		for idx, item := range userItems {
			infoLines = append(infoLines, fmt.Sprintf("%d) %s Lv%d", idx+1, localize(item.Name), item.Level))
		}
		infoLines = append(infoLines, localize("C)ancel"))
	case SELECT_HUNT_ITEM:
		userWeaponItems := screen.user.InventoryItems()
		infoLines = append(infoLines, localize("N)o item"))
		for idx, item := range userWeaponItems {
			infoLines = append(infoLines, fmt.Sprintf("%d) %s Lv%d", idx+1, localize(item.Name), item.Level))
		}
		infoLines = append(infoLines, localize("Get I)nitial Coin"))
		infoLines = append(infoLines, localize("C)ancel"))
	case SELECT_UPGRADE_ITEM:
		userUpgradeItems := screen.user.UpgradableItems()
		for idx, item := range userUpgradeItems {
			infoLines = append(infoLines, fmt.Sprintf("%d) %s Lv%d", idx+1, localize(item.Name), item.Level))
		}
		infoLines = append(infoLines, localize("C)ancel"))
	case RESULT_BUY_FINISH:
		fallthrough
	case RESULT_HUNT_FINISH:
		fallthrough
	case RESULT_SELL_FINISH:
		fallthrough
	case RESULT_UPGRADE_FINISH:
		infoLines = append(infoLines, localize("Go) on"))
	default:
	}

	// box start point (x, y)
	x := 2
	y := screen.screenSize.Height/2 + 1

	bgcolor := uint64(bgcolor)
	fmtFunc := screen.colorFunc(fmt.Sprintf("255:%v", bgcolor))
	for index, line := range infoLines {
		io.WriteString(os.Stdout, fmt.Sprintf("%s%s",
			cursor.MoveTo(y+index, x), fmtFunc(line)))
		if index+2 > int(screen.screenSize.Height) {
			break
		}
	}
}

func (screen *GameScreen) renderUserSituation() {
	infoLines := []string{}
	desc := ""
	waitProcessEnd := localize("wait process to end")
	switch screen.scrStatus {
	case SHOW_LOCATION:
		locationDescMap := map[UserLocation]string{
			HOME:   localize("home desc"),
			FOREST: localize("forest desc"),
			SHOP:   localize("shop desc"),
			SETTINGS : localize("settings desc"),
		}
		desc = locationDescMap[screen.user.GetLocation()]
	case SELECT_BUY_ITEM:
		desc = localize("select buy item desc")
	case SELECT_SELL_ITEM:
		desc = localize("select sell item desc")
	case SELECT_HUNT_ITEM:
		desc = localize("select hunt item desc")
	case SELECT_UPGRADE_ITEM:
		desc = localize("select upgrade item desc")
	case WAIT_BUY_PROCESS:
		desc = fmt.Sprintf("%s %s Lv%d.\n%s", localize("wait buy process desc"), localize(screen.activeItem.Name), screen.activeItem.Level, waitProcessEnd)
	case WAIT_HUNT_PROCESS:
		if len(screen.activeItem.Name) > 0 {
			desc = fmt.Sprintf("%s %s Lv%d.\n%s", localize("wait hunt process desc"), localize(screen.activeItem.Name), screen.activeItem.Level, waitProcessEnd)
		} else {
			desc = fmt.Sprintf("%s\n%s", localize("hunting without weapon"), waitProcessEnd)
		}
	case WAIT_SELL_PROCESS:
		desc = fmt.Sprintf("%s %s Lv%d.\n%s", localize("wait sell process desc"), localize(screen.activeItem.Name), screen.activeItem.Level, waitProcessEnd)
	case WAIT_UPGRADE_PROCESS:
		desc = fmt.Sprintf("%s %s.\n%s", localize("wait upgrade process desc"), localize(screen.activeItem.Name), waitProcessEnd)
	case RESULT_BUY_FINISH:
		if screen.txFailReason != "" {
			desc = screen.txFailReason
		} else {
			desc = fmt.Sprintf("%s %s Lv%d.\n%s", localize("result buy finish desc"), localize(screen.activeItem.Name), screen.activeItem.Level, localize("use for hunting"))
		}
	case RESULT_HUNT_FINISH:
		desc = fmt.Sprintf("%s %d.", localize("result hunt finish desc"), screen.txResult.Amount)
	case RESULT_SELL_FINISH:
		desc = fmt.Sprintf("%s %s Lv%d.", localize("result sell finish desc"), localize(screen.activeItem.Name), screen.activeItem.Level)
	case RESULT_UPGRADE_FINISH:
		desc = fmt.Sprintf("%s: %s.", localize("result upgrade finish desc"), localize(screen.activeItem.Name))
	}

	basicLines := strings.Split(desc, "\n")

	for _, line := range basicLines {
		infoLines = append(infoLines, ChunkString(line, screen.screenSize.Width/2-4)...)
	}

	// box start point (x, y)
	x := 2
	y := 2

	bgcolor := uint64(bgcolor)
	fmtFunc := screen.colorFunc(fmt.Sprintf("255:%v", bgcolor))
	for index, line := range infoLines {
		io.WriteString(os.Stdout, fmt.Sprintf("%s%s", cursor.MoveTo(y+index, x), fmtFunc(line)))
		if index+2 > int(screen.screenSize.Height) {
			break
		}
	}
}

func (screen *GameScreen) renderCharacterSheet() {
	var HP uint64 = 10
	var MaxHP uint64 = 10
	bgcolor := uint64(bgcolor)
	warning := ""
	if float32(HP) < float32(MaxHP)*.25 {
		bgcolor = 124
		warning = localize("health low warning")
	} else if float32(HP) < float32(MaxHP)*.1 {
		bgcolor = 160
		warning = localize("health critical warning")
	}

	x := screen.screenSize.Width/2 - 1
	width := (screen.screenSize.Width - x)
	fmtFunc := screen.colorFunc(fmt.Sprintf("255:%v", bgcolor))

	infoLines := []string{
		centerText(fmt.Sprintf("%v", screen.user.GetUserName()), " ", width),
		centerText(warning, "─", width),
		screen.drawProgressMeter(1, 1, 208, bgcolor, 1) + fmtFunc(truncateRight(fmt.Sprintf(" %s: %v", localize("gold"), screen.user.GetGold()), width-1)),
		screen.drawProgressMeter(HP, MaxHP, 196, bgcolor, 10) + fmtFunc(truncateRight(fmt.Sprintf(" HP: %v/%v", HP, MaxHP), width-10)),
		// screen.drawProgressMeter(HP, MaxHP, 225, bgcolor, 10) + fmtFunc(truncateRight(fmt.Sprintf(" XP: %v/%v", HP, 10), width-10)),
		// screen.drawProgressMeter(HP, MaxHP, 208, bgcolor, 10) + fmtFunc(truncateRight(fmt.Sprintf(" AP: %v/%v", HP, MaxHP), width-10)),
		// screen.drawProgressMeter(HP, MaxHP, 117, bgcolor, 10) + fmtFunc(truncateRight(fmt.Sprintf(" RP: %v/%v", HP, MaxHP), width-10)),
		// screen.drawProgressMeter(HP, MaxHP, 76, bgcolor, 10) + fmtFunc(truncateRight(fmt.Sprintf(" MP: %v/%v", HP, MaxHP), width-10)),
	}

	infoLines = append(infoLines, centerText(localize("inventory items"), "─", width))
	items := screen.user.InventoryItems()
	for _, item := range items {
		infoLines = append(infoLines, truncateRight(fmt.Sprintf("%s Lv%d", localize(item.Name), item.Level), width))
	}
	infoLines = append(infoLines, centerText(" ❦ ", "─", width))

	for index, line := range infoLines {
		io.WriteString(os.Stdout, fmt.Sprintf("%s%s", cursor.MoveTo(2+index, x), fmtFunc(line)))
		if index+2 > int(screen.screenSize.Height) {
			break
		}
	}

	nodeLines := []string{
		centerText(localize("pylons network status"), " ", width),
		centerText(screen.user.GetLastTransaction(), " ", width),
		centerText("BlockHeight: "+strconv.FormatInt(screen.blockHeight, 10), " ", width),
		centerText(" ❦ ", "─", width),
	}

	for index, line := range nodeLines {
		io.WriteString(os.Stdout, fmt.Sprintf("%s%s", cursor.MoveTo(2+len(infoLines)+index, x), fmtFunc(line)))
		if index+2 > int(screen.screenSize.Height) {
			break
		}
	}

	lastLine := len(infoLines) + len(nodeLines) + 1
	screen.drawFill(x, lastLine+1, width, screen.screenSize.Height-(lastLine+2))
}

func (screen *GameScreen) UpdateBlockHeight(blockHeight int64) {
	screen.blockHeight = blockHeight
	screen.refreshed = false
	screen.Render()
}

func (screen *GameScreen) HandleInputKey(input termbox.Event) {
	Key := string(input.Ch)
	log.Println("Handling Key \"", Key, "\"")
	switch Key {
	case "H": // HOME
		fallthrough
	case "h":
		screen.user.SetLocation(HOME)
		screen.refreshed = false
	case "F": // FOREST
		fallthrough
	case "f":
		screen.user.SetLocation(FOREST)
		screen.refreshed = false
	case "S": // SHOP
		fallthrough
	case "s":
		screen.user.SetLocation(SHOP)
		screen.refreshed = false
	case "T": // SETTINGS
		fallthrough
	case "t":
		screen.user.SetLocation(SETTINGS)
		screen.refreshed = false
	case "G":
		fallthrough
	case "g":
		gameLanguage = "en"
		screen.refreshed = false
	case "A":
		fallthrough
	case "a":
		gameLanguage = "es"
		screen.refreshed = false
	case "C": // CANCEL
		fallthrough
	case "c":
		screen.scrStatus = SHOW_LOCATION
		screen.refreshed = false
	case "O": // GO ON
		fallthrough
	case "o":
		screen.scrStatus = SHOW_LOCATION
		screen.refreshed = false
	case "U": // HUNT
		fallthrough
	case "u":
		screen.scrStatus = SELECT_HUNT_ITEM
		screen.refreshed = false
	case "B": // BUY
		fallthrough
	case "b": // BUY
		screen.scrStatus = SELECT_BUY_ITEM
		screen.refreshed = false
	case "E": // SELL
		fallthrough
	case "e":
		screen.scrStatus = SELECT_SELL_ITEM
		screen.refreshed = false
	case "P": // UPGRADE ITEM
		fallthrough
	case "p":
		screen.scrStatus = SELECT_UPGRADE_ITEM
		screen.refreshed = false
	case "N": // Go hunt with no weapon
		fallthrough
	case "n":
		fallthrough
	case "I":
		fallthrough
	case "i":
		fallthrough
	case "1": // SELECT 1st item
		fallthrough
	case "2": // SELECT 2nd item
		fallthrough
	case "3": // SELECT 3rd item
		fallthrough
	case "4": // SELECT 4th item
		fallthrough
	case "5": // SELECT 5rd item
		fallthrough
	case "6": // SELECT 6rd item
		fallthrough
	case "7": // SELECT 7rd item
		fallthrough
	case "8": // SELECT 8rd item
		fallthrough
	case "9": // SELECT 9rd item
		screen.refreshed = false
		switch screen.scrStatus {
		case SELECT_BUY_ITEM:
			screen.activeItem = GetToBuyItemFromKey(Key)
			screen.scrStatus = WAIT_BUY_PROCESS
			screen.refreshed = false
			screen.Render()
			log.Println("started buying item")
			txhash, err := Buy(screen.user, Key)
			log.Println("ended buying item")
			if err != nil {
				screen.txFailReason = err.Error()
				screen.scrStatus = RESULT_BUY_FINISH
				screen.refreshed = false
				screen.Render()
			} else {
				time.AfterFunc(1*time.Second, func() {
					pylonSDK.WaitForNextBlock()
					screen.txResult = ProcessTxResult(screen.user, txhash)
					screen.scrStatus = RESULT_BUY_FINISH
					screen.refreshed = false
					screen.Render()
				})
			}
		case SELECT_HUNT_ITEM:
			screen.activeItem = GetWeaponItemFromKey(screen.user, Key)
			screen.scrStatus = WAIT_HUNT_PROCESS
			screen.refreshed = false
			screen.Render()
			log.Println("started hunting item")
			txhash := Hunt(screen.user, Key)
			log.Println("ended hunting item")
			time.AfterFunc(1*time.Second, func() {
				pylonSDK.WaitForNextBlock()
				screen.txResult = ProcessTxResult(screen.user, txhash)
				screen.scrStatus = RESULT_HUNT_FINISH
				screen.refreshed = false
				screen.Render()
			})
		case SELECT_SELL_ITEM:
			screen.activeItem = GetToSellItemFromKey(screen.user, Key)
			screen.scrStatus = WAIT_SELL_PROCESS
			screen.refreshed = false
			screen.Render()
			log.Println("started selling item")
			txhash := Sell(screen.user, Key)
			log.Println("ended selling item")
			time.AfterFunc(1*time.Second, func() {
				pylonSDK.WaitForNextBlock()
				screen.txResult = ProcessTxResult(screen.user, txhash)
				screen.scrStatus = RESULT_SELL_FINISH
				screen.refreshed = false
				screen.Render()
			})
		case SELECT_UPGRADE_ITEM:
			screen.activeItem = GetToUpgradeItemFromKey(screen.user, Key)
			screen.scrStatus = WAIT_UPGRADE_PROCESS
			screen.refreshed = false
			screen.Render()
			log.Println("started upgrading item")
			txhash := Upgrade(screen.user, Key)
			log.Println("ended upgrading item")
			time.AfterFunc(1*time.Second, func() {
				pylonSDK.WaitForNextBlock()
				screen.txResult = ProcessTxResult(screen.user, txhash)
				screen.scrStatus = RESULT_UPGRADE_FINISH
				screen.refreshed = false
				screen.Render()
			})
		}
	}
	screen.Render()
}

func (screen *GameScreen) Render() {
	var HP uint64 = 10

	if screen.screenSize.Height < 20 || screen.screenSize.Width < 60 {
		clear := cursor.ClearEntireScreen()
		move := cursor.MoveTo(1, 1)
		io.WriteString(os.Stdout,
			fmt.Sprintf("%s%s%s", clear, move, localize("screen size warning")))
		return
	} else if HP == 0 {
		clear := cursor.ClearEntireScreen()
		dead := localize("dead")
		move := cursor.MoveTo(screen.screenSize.Height/2, screen.screenSize.Width/2-utf8.RuneCountInString(dead)/2)
		io.WriteString(os.Stdout, clear+move+dead)
		screen.refreshed = false
		return
	}

	if !screen.refreshed {
		clear := cursor.ClearEntireScreen() + allowMouseInputAndHideCursor
		io.WriteString(os.Stdout, clear)
		screen.redrawBorders()
		screen.refreshed = true
	}

	screen.renderUserCommands()
	screen.renderUserSituation()
	screen.renderCharacterSheet()
}

func (screen *GameScreen) Reset() {
	io.WriteString(os.Stdout, fmt.Sprintf("%s👋\n", resetScreen))
}

func (screen *GameScreen) SaveGame() {
	screen.user.Save()
}

// NewScreen manages the window rendering for game
func NewScreen(world World, user User) Screen {
	width, _ := terminal.Width()
	height, _ := terminal.Height()

	window := ssh.Window{
		Width:  int(width),
		Height: int(height),
	}

	screen := GameScreen{
		world:          world,
		user:           user,
		screenSize:     window,
		colorCodeCache: make(map[string](func(string) string))}

	return &screen
}
