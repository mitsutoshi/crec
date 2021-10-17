package gmo

const (
	WssUrl    = "wss://api.coin.z.com/ws" // websocket url
	OriginUrl = "https://api.coin.z.com"  // rest api server url

	// symbol
	SymbolBtcJpy = "BTC_JPY" // symbol: BTC_JPY
	SymbolEthJpy = "ETH_JPY" // symbol: ETH_JPY
	SymbolXrpJpy = "XRP_JPY" // symbol: XRP_JPY
	SymbolLtcJpy = "LTC_JPY" // symbol: LTC_JPY
	SymbolBchJpy = "BCH_JPY" // symbol: BCH_JPY

	// status
	StatusOpen        = "OPEN"        // status: open
	StatusPreOpen     = "PREOPEN"     // status: preopen
	StatusMaintenance = "MAINTENANCE" // status: MAINTENANCE

	// side
	SideBuy  = "BUY"  // side: BUY
	SideSell = "SELL" // side: SELL

	// order's executionType
	ExecTypeLimit  = "LIMIT"
	ExecTypeMarket = "MARKET"
	ExecTypeStop   = "STOP"

	// order's msgType
	MsgTypeNewOrder     = "NOR" // msgType: new order
	MsgTypeModifyOrder  = "ROR" // msgType: modify order
	MsgTypeModifyCancel = "COR" // msgType: cancel order

	// order's settleType
	SettleTypeOpen  = "OPEN"  // selleType: open order
	SettleTypeClose = "CLOSE" // selleType: close order

	timeFormat   = "2006-01-02T15:04:05.000Z"
	TradeHeaders = "receive_time,symbol,side,size,price,timestamp"
)
