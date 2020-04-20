package loud

type UserLocation int

const (
	HOME UserLocation = iota
	FOREST
	SHOP
	PYLCNTRL
	SETTINGS
	DEVELOP
)

// User represents an active user in the system.
type User interface {
	SetAddress(string)
	SetGold(int)
	SetPylonAmount(int)
	SetItems([]Item)
	SetCharacters([]Character)
	SetDefaultItemIndex(idx int)
	SetDefaultCharacterIndex(idx int)
	SetLocation(UserLocation)
	SetLastTransaction(string)
	SetLatestBlockHeight(int64)
	InventoryItems() []Item
	InventoryItemIDByName(string) string
	InventoryIronSwords() []Item
	InventorySwords() []Item
	InventoryCharacters() []Character
	InventoryUpgradableItems() []Item
	InventorySellableItems() []Item
	GetLocation() UserLocation
	GetPrivKey() string
	GetDefaultItemIndex() int
	GetDefaultCharacterIndex() int
	GetDefaultCharacter() *Character
	GetAddress() string
	GetGold() int
	GetPylonAmount() int
	GetUserName() string
	GetLastTransaction() string
	GetLatestBlockHeight() int64
	Reload()
	Save()
}
