package orderbook

// import (
// 	"os"
// 	"sync"

// 	"github.com/VictorLowther/btree"
// )

// type Limit struct{
// 	orders []*Limit
// }



// type DataFeed struct{
// 	Provider string
// 	Symbol string
// 	BestAsk float64
// 	BestBid float64
// 	Spread float64
// }


// type BestPrice struct{
// 	Provider string
// 	Price    float64
// 	Size     float64
// }

// type CrossSpread struct{
// 	Symbol  string
// 	BestAsk BestPrice
// 	BestBid BestPrice
// 	Spread  float64
// }

// type Provider interface{
// 	Start() error
// 	GetOrderbooks() Orderbooks
// 	Name() string
// }

// type Orderbooks map[string]*Book

// type Book struct{
// 	Symbol string
// 	Asks *Limits
// 	Bids *Limits
// }

// func NewBook(symbol string) *Book{
// 	return &Book{
// 		Symbol: symbol,
// 		Asks: NewLimits(false),
// 		Bids: NewLimits(true),
// 	}
// }

// func (b *Book) Spread() float64{
// 	if b.Asks.data.Len() == 0 || b.Bids.data.Len() == 0 {
// 		return 0.0
// 	}
// 	bestAsk := b.Asks.Best().Price
// 	bestBid := b.Bids.Best().Price
// 	return bestAsk - bestBid
// }

// func (b *Book) BestBid() *Limit{
// 	return b.Bids.Best()
// }

// func (b *Book) BestAsk() *Limit{
// 	return b.Asks.Best()
// }

// func getBidByPrice(price float64) btree.CompareAgainst[*Limit]{
// 	return func(l *Limit) int {
// 		switch {
// 		case l.Price > price:
// 			return -1
// 		case l.Price < price:
// 			return 1
// 		default:
// 			return 0
// 		}
// 	}
// }

// func getAskByPrice(price float64) btree.CompareAgainst[*Limit]{
// 	return func(l *Limit) int {
// 		switch {
// 		case l.Price < price:
// 			return -1
// 		case l.Price > price:
// 			return 1
// 		default:
// 			return 0
// 		}
// 	}
// }

// func sortByBestBid(a, b *Limit) bool {
// 	return a.Price > b.Price
// }

// func sortByBestAsk(a, b *Limit) bool {
// 	return a.Price < b.Price
// }

// type Limits struct{
// 	isBids      bool
// 	lock        sync.RWMutex
// 	data        *btree.Tree[*Limit]
// 	totalVolume float64
// }

// func NewLimits(isBids bool) *Limits{
// 	f := sortByBestAsk
// 	if isBids {
// 		f = sortByBestBid
// 	}
// 	return &Limits{
// 		isBids: isBids,
// 		data: btree.New(f),
// 	}
// }

// func (l *Limits) Len() int {
// 	return l.data.Len()
// }

// func (l *Limits) Best() *Limit{
// 	l.lock.RLock()
// 	defer l.lock.RUnlock()

// 	if l.data.Len() == 0 {
// 		return nil
// 	}
// 	iter := l.data.Iterator(nil, nil)
// 	iter.Next()
// 	return iter.Item()
// }

// func (l *Limits) Update(price float64, size float64) {
// 	l.lock.Lock()
// 	defer l.lock.Unlock()

// 	getFunc := getAskByPrice(price)
// 	if l.isBids {
// 		getFunc = getBidByPrice(price)
// 	}

// 	if limit, ok := l.data.Get(getFunc); ok {
// 		if size == 0.0 {
// 			l.data.Delete(limit)
// 			return
// 		}
// 		limit.totalVolume = size
// 		return
// 	}
// 	if size == 0.0 {
// 		return
// 	}

// 	limit := NewLimit(price)
// 	limit.TotalVolume = size
// 	l.data.Insert(limit)
// }

// func (l *Limits) addOrder(price float64, o *Order){
// 	if o.isBids != l.isBids{
// 		panic("the side of the limits does not match the side of the odrer")
// 	}

// 	f := getAskByPrice(price)
// 	if l.isBids{
// 		f = getBidByPrice(price)
// 	}

// 	var (
// 		limit *Limit
// 		ok bool
// 	)

// 	limit, ok = l.data.Get(f)
// 	if !ok {
// 		limit = NewLimit(price)
// 		l.data.Insert(limit)
// 	}

// 	l.totalVolume += o.size

// 	limit.addOrder(o)
// }

// func loadFromFile(src string) error{
// 	f, err := os.Open(src)
// 	if err != nil{

// 	}
// }