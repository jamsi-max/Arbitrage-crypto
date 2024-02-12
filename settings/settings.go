package settings

var ServerSymbols = map[string]bool{
	// "BTCUSD": true,
	"SOLUSD":  true,
	"XLMUSD":  true,
	"WAXPUSD": true,
	// "ETHUSD",
	// "DOGEUSD",
	// "ADAUSD",
	// "LTCUSD",
	// "SOLUSD",
}

var Symbols = []string{
	// "BTCUSD",
	// "ETHUSD",
	// "DOGEUSD",
	// "ADAUSD",
	// "LTCUSD",
	"SOLUSD",
	"XLMUSD",
	"WAXPUSD",
}

var Pairs = map[string]map[string]string{
	"XLMUSD": {
		// "Finex":    "XLMUSDTPERP", // не поддерживается
		// "Kraken":   "XLM/USD", // не работает в России
		// "Coinbase": "XLM-USD", // не работает в России

		"Binance": "XLMUSDT",
		"Bybit":   "orderbook.1.XLMUSDT",
		"Cucoin":  "/market/level2:XLM-USDT",
		"OKX":     "XLM-USDT",
		"MEXC":    "spot@public.increase.depth.v3.api@XLMUSDT",
		"Gateio":  "XLM_USDT",
	},
	"SOLUSD": {
		// "Finex":    "SOLUSDTPERP", // высокие комиссии
		// "Kraken":   "SOL/USD", // не работает в России
		// "Coinbase": "SOL-USD", // не работает в России

		"Binance": "SOLUSDT",
		"Bybit":   "orderbook.1.SOLUSDT",
		"Cucoin":  "/market/level2:SOL-USDT",
		"OKX":     "SOL-USDT",
		"MEXC":    "spot@public.increase.depth.v3.api@SOLUSDT",
		"Gateio":  "SOL_USDT",
	},
	"WAXPUSD": {
		// "Kraken":   "WAX/USD", // не работает в России
		// "Coinbase": "WAXT-USD", // не поддерживается
		// "Finex":  "WAXPUSDTPERP", // не поддерживается

		"Binance": "WAXPUSDT",
		"Bybit":   "orderbook.1.WAXPUSDT",
		"Cucoin":  "/market/level2:WAX-USDT",
		"OKX":     "WAXP-USDT",
		"MEXC":    "spot@public.increase.depth.v3.api@WAXPUSDT",
		"Gateio":  "WAXP_USDT",
	},
	// "ADAUSD": {
	// 	"Binance":  "ADAUSDT",
	// 	"Kraken":   "ADA/USD",
	// 	"Coinbase": "ADA-USD",
	// 	"Bybit":    "orderbook.1.ADAUSDT",
	// 	"Cucoin":   "/market/ticker:ADA-USDT",
	// },
	// "DOGEUSD": {
	// 	"Binance":  "DOGEUSDT",
	// 	"Kraken":   "XDG/USD",
	// 	"Coinbase": "DOGE-USD",
	// 	"Bybit":    "orderbook.1.DOGEUSDT",
	// 	"Cucoin":   "/market/ticker:DOGE-USDT",
	// },
	// "BTCUSD": {
	// 	"Binance":  "BTCUSDT",
	// 	// "Kraken":   "XBT/USD",
	// 	"Coinbase": "BTC-USD",
	// 	"Bybit":    "orderbook.1.BTCUSDT",
	// 	"Cucoin":   "/market/ticker:BTC-USDT",
	// },
	// "ETHUSD": {
	// 	"Binance":  "ETHUSDT",
	// 	"Kraken":   "ETH/USD",
	// 	"Coinbase": "ETH-USD",
	// 	"Bybit":    "orderbook.1.ETHUSDT",
	// 	"Cucoin":   "/market/ticker:ETH-USDT",
	// },
	// "LTCUSD": {
	// 	"Binance":  "LTCUSDT",
	// 	"Kraken":   "LTC/USD",
	// 	"Coinbase": "LTC-USD",
	// 	"Bybit":    "orderbook.1.LTCUSDT",
	// 	"Cucoin":   "/market/ticker:LTC-USDT",
	// },
}
