package ftx

const (
	SymbolBtcPerp = "BTC-PERP"
	SymbolBtcUsd  = "BTC/USD"
	SymbolBtc0626 = "BTC-0626"
	SymbolBtc0923 = "BTC-0923"
	TradeHeaders  = "receive_time,id,price,size,side,liquidation,time"
)

const (
	publicWssUrl      = "wss://ftx.com/ws" // public websocket url
	channelTicker     = "ticker"           // websocket channel: ticker
	channelTrades     = "trades"           // websocket channel: trades
	channelOrderbooks = "orderbooks"       // websocket channel: orderbooks
	msgTypeSubscribed = "subscribed"
	msgTypeUpdate     = "update"
	timeFormat        = "2006-01-02T15:04:05.000Z"
)
