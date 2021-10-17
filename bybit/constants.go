package bybit

const (
	TradeHeaders = "receive_time,timestamp,trade_time_ms,price,side,size,tick_direction,trade_id,cross_seq"
)

const (
	publicWssUrl = "wss://stream.bytick.com/realtime" // public websocket url
	timeFormat   = "2006-01-02T15:04:05.000Z"
)
