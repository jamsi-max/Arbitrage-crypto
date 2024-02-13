package settings

var ServerSymbols = map[string]bool{
	"ALGOUSD": true,
	"XLMUSD":  true,
	"WAXPUSD": true,


	// "BTCUSD": true,
	// "SOLUSD":  true,
	// "ETHUSD",
	// "DOGEUSD",
	// "ADAUSD",
	// "LTCUSD",
	// "SOLUSD",
}

var Symbols = []string{
	"ALGOUSD",
	"XLMUSD",
	"WAXPUSD",


	// "BTCUSD",
	// "ETHUSD",
	// "DOGEUSD",
	// "ADAUSD",
	// "LTCUSD",
	// "SOLUSD",
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
	"ALGOUSD": {
		// "Finex":    "SOLUSDTPERP", // высокие комиссии
		// "Kraken":   "SOL/USD", // не работает в России
		// "Coinbase": "SOL-USD", // не работает в России

		"Binance": "ALGOUSDT",
		"Bybit":   "orderbook.1.ALGOUSDT",
		"Cucoin":  "/market/level2:ALGO-USDT",
		"OKX":     "ALGO-USDT",
		"MEXC":    "spot@public.increase.depth.v3.api@ALGOUSDT",
		"Gateio":  "ALGO_USDT",
	},


	// "SOLUSD": {
	// 	// "Finex":    "SOLUSDTPERP", // высокие комиссии
	// 	// "Kraken":   "SOL/USD", // не работает в России
	// 	// "Coinbase": "SOL-USD", // не работает в России

	// 	"Binance": "SOLUSDT",
	// 	"Bybit":   "orderbook.1.SOLUSDT",
	// 	"Cucoin":  "/market/level2:SOL-USDT",
	// 	"OKX":     "SOL-USDT",
	// 	"MEXC":    "spot@public.increase.depth.v3.api@SOLUSDT",
	// 	"Gateio":  "SOL_USDT",
	// },
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

var Fee = map[string]map[string]float64{
	"XLMUSD": {
		// "Finex":    "XLMUSDTPERP", // не поддерживается
		// "Kraken":   "XLM/USD", // не работает в России
		// "Coinbase": "XLM-USD", // не работает в России

		"Bybit":   0.02,
		"Binance": 0.02,
		"OKX":     0.016,
		"MEXC":    0.1,
		"Cucoin":  3,
		"Gateio":  4.25,
		"HTX":     0.02,
	},
	"WAXPUSD": {
		// "Kraken":   "WAX/USD", // не работает в России
		// "Coinbase": "WAXT-USD", // не поддерживается
		// "Finex":  "WAXPUSDTPERP", // не поддерживается

		"Bybit":   2,
		"Binance": 2,
		"OKX":     0.1,
		"MEXC":    2,
		"Cucoin":  9.158,
		"Gateio":  8.95,
		"HTX":     1.5,
	},
	"ALGOUSD": {
		// "Finex":    "SOLUSDTPERP", // высокие комиссии
		// "Kraken":   "SOL/USD", // не работает в России
		// "Coinbase": "SOL-USD", // не работает в России

		"Bybit":   0.01,
		"Binance": 0.008,
		"OKX":     0.008,
		"MEXC":    0.1,
		"Cucoin":  0.1,
		"Gateio":  2.94,
		"HTX":     0.01,
	},
}
