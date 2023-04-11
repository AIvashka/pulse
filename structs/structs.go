package structs

import (
	"fmt"
	"sort"
	"time"
)

type Order struct {
	Side       OrderSide
	Instrument string
	Size       float64
	Price      float64
	CreatedAt  time.Time
}

type OrderSide int

const (
	OrderSideBuy OrderSide = iota
	OrderSideSell
)

func (side OrderSide) String() string {
	switch side {
	case OrderSideBuy:
		return "BUY"
	case OrderSideSell:
		return "SELL"
	default:
		return fmt.Sprintf("UNKNOWN (%d)", side)
	}
}

func ParseOrderSide(side float64) OrderSide {
	if side == 0 {
		return OrderSideBuy
	} else if side == 1 {
		return OrderSideSell
	} else {
		panic(fmt.Sprintf("invalid order side: %d", side))
	}
}

type OrderBook struct {
	Bids []*Order
	Asks []*Order
}

type OrderBookMap map[string]*OrderBook

func (o OrderBookMap) AddOrder(order *Order) {
	ob, ok := o[order.Instrument]
	if !ok {
		ob = &OrderBook{}
		o[order.Instrument] = ob
	}

	if order.Side == OrderSideBuy {
		ob.Bids = append(ob.Bids, &Order{
			Size:  order.Size,
			Price: order.Price,
		})
	} else {
		ob.Asks = append(ob.Asks, &Order{
			Size:  order.Size,
			Price: order.Price,
		})
	}
}

func (o OrderBookMap) Sort() {
	for _, ob := range o {
		sort.Sort(bidsByPrice(ob.Bids))
		sort.Sort(asksByPrice(ob.Asks))
	}
}

type bidsByPrice []*Order

func (b bidsByPrice) Len() int           { return len(b) }
func (b bidsByPrice) Swap(i, j int)      { b[i], b[j] = b[j], b[i] }
func (b bidsByPrice) Less(i, j int) bool { return b[i].Price > b[j].Price }

type asksByPrice []*Order

func (a asksByPrice) Len() int           { return len(a) }
func (a asksByPrice) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a asksByPrice) Less(i, j int) bool { return a[i].Price < a[j].Price }

func GroupOrdersBySymbol(orders []*Order) OrderBookMap {
	orderBooks := make(OrderBookMap)

	for _, order := range orders {
		orderBooks.AddOrder(order)
	}

	orderBooks.Sort()

	return orderBooks
}

type SpreadLevel struct {
	BaseAsset  string  `json:"base_asset"`
	QuoteAsset string  `json:"quote_asset"`
	Quantity   float64 `json:"quantity"`
}

type SpreadConfig struct {
	Levels []*SpreadLevel `json:"levels"`
}

type SpreadRecord struct {
	Symbol   string
	Quantity float64
	Spread   float64
	SaveTime time.Time
}

type Config struct {
	DBHost       string
	DBPort       string
	DBUser       string
	DBPassword   string
	DBName       string
	SpreadConfig SpreadConfig
	APIKey       string
	APISecret    string
	Interval     int
}
