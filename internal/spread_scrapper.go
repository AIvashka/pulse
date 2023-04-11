package internal

import (
	"fmt"
	"log"
	"pulse/structs"
	"sync"
	"time"
)

type SpreadScrapper struct {
	OrderBooks structs.OrderBookMap
	Config     structs.SpreadConfig
}

func (s *SpreadScrapper) calculateEffectiveSpreadForQuantity(symbol string, baseQuantity float64) (float64, error) {
	ob, ok := s.OrderBooks[symbol]
	if !ok {
		return 0.0, fmt.Errorf("order book not found for symbol %s", symbol)
	}

	bestBid, err := getBestPriceForQuantity(ob.Bids, baseQuantity)
	if err != nil {
		return 0.0, err
	}

	bestAsk, err := getBestPriceForQuantity(ob.Asks, baseQuantity)
	if err != nil {
		return 0.0, err
	}

	midPrice := (bestBid + bestAsk) / 2
	spread := (bestAsk - bestBid) / midPrice * 100

	return spread, nil
}

func getBestPriceForQuantity(orders []*structs.Order, baseQuantity float64) (float64, error) {
	if len(orders) == 0 {
		return 0.0, fmt.Errorf("no orders found")
	}

	var totalVolume, totalValue float64
	for _, order := range orders {
		if totalVolume >= baseQuantity {
			break
		}
		if totalVolume+order.Size > baseQuantity {
			totalValue += (baseQuantity - totalVolume) * order.Price
			totalVolume = baseQuantity
		} else {
			totalValue += order.Size * order.Price
			totalVolume += order.Size
		}
	}

	if totalVolume < baseQuantity {
		return 0.0, fmt.Errorf("not enough orders to fill quantity")
	}

	return totalValue / baseQuantity, nil
}

//func (s *SpreadScrapper) Run(ch chan<- *structs.SpreadRecord) {
//	for _, level := range s.Config.Levels {
//		symbol := fmt.Sprintf("%s-%s", level.BaseAsset, level.QuoteAsset)
//		baseQuantity := level.Quantity
//
//		spread, err := s.calculateEffectiveSpreadForQuantity(symbol, baseQuantity)
//		if err != nil {
//			log.Println("failed to calculate spread for instrument %s\n", symbol)
//			continue
//		}
//
//		record := &structs.SpreadRecord{
//			Symbol:   symbol,
//			Quantity: baseQuantity,
//			Spread:   spread,
//			SaveTime: time.Now().UTC(),
//		}
//		ch <- record
//	}
//}

func (s *SpreadScrapper) Run(ch chan<- *structs.SpreadRecord) {
	var wg sync.WaitGroup
	for _, level := range s.Config.Levels {
		wg.Add(1)
		go func(level *structs.SpreadLevel) {
			defer wg.Done()
			symbol := fmt.Sprintf("%s-%s", level.BaseAsset, level.QuoteAsset)

			spread, err := s.calculateEffectiveSpreadForQuantity(symbol, level.Quantity)
			if err != nil {
				log.Printf("failed to calculate spread for instrument %s, %v", symbol, err)
				return
			}

			record := &structs.SpreadRecord{
				Symbol:   symbol,
				Quantity: level.Quantity,
				Spread:   spread,
				SaveTime: time.Now().UTC(),
			}

			ch <- record
		}(level)
	}
	wg.Wait()
	//close(ch)
}
