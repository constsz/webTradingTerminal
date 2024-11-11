package oms

import (
	"context"
	"errors"
	"fmt"
	"github.com/adshao/go-binance/v2/futures"
	"log"
	"strconv"
	"time"
	"ttk/data"
)

// ----------------------------------------
// Binance Exchange Functions

func (oms *OMS) CreateEntryLimitOrder(oo *OpenOrder) bool {
	symbolUsdt := data.SymbolUSDT(oo.Symbol)
	side := oo.Side

	price := BinancePrice(oo.EntryPrice, oms.symbolRules.GetTickSize(oo.Symbol))
	qty := fmt.Sprintf("%f", oo.Quantity)

	inPosition := oms.inventory.symbolInPosition(oo.Symbol)
	lastPrice := oms.db.GetLastPrice(oo.Symbol)

	if inPosition {
		qty = fmt.Sprintf("%f", oms.inventory.Symbol[oo.Symbol].BaseQuantity)
	}

	// ОТКЛЮЧЕН: Функционал который позволяет ставить стоп ордер на вход вместо лимитки.
	var stopLimitOrder bool
	if side == data.SideLong {
		if oo.EntryPrice > lastPrice {
			stopLimitOrder = true
			fmt.Printf("\n(!) %s LONG - Rejected: EntryPrice > LastPrice\n", oo.Symbol)
		}
	} else if side == data.SideShort {
		if oo.EntryPrice < lastPrice {
			stopLimitOrder = true
			fmt.Printf("\n(!) %s SHORT - Rejected: EntryPrice < LastPrice\n", oo.Symbol)
		}
	}

	var order *futures.CreateOrderService

	if !stopLimitOrder {
		order = oms.client.NewCreateOrderService().Symbol(symbolUsdt).
			TimeInForce(futures.TimeInForceTypeGTC).
			Side(binanceSide(side)).Type(futures.OrderTypeLimit).
			Quantity(qty).Price(price)

		res, err := order.Do(context.Background())
		if err != nil {
			log.Println("CreateEntryLimitOrder:", err)
			return false
		} else {
			oo.Active = true
			oo.BinanceId = res.OrderID

			oms.orderStack.setOrder(*oo)

			return true
		}

	}

	return false

	//else if stopLimitOrder {
	// Stop Limit
	//order = oms.client.NewCreateOrderService().Symbol(symbolUsdt).
	//	TimeInForce(futures.TimeInForceTypeGTC).
	//	Side(binanceSide(side)).Type(futures.OrderTypeStop).
	//	Quantity(qty).Price(price).StopPrice(price)
	// Stop Market
	//order = oms.client.NewCreateOrderService().Symbol(symbolUsdt).
	//	TimeInForce(futures.TimeInForceTypeGTC).
	//	Side(binanceSide(side)).Type(futures.OrderTypeStopMarket).
	//	Quantity(qty).StopPrice(price)

	//}

}

func (oms *OMS) MarketEntry(oo *OpenOrder) {
	symbolUsdt := data.SymbolUSDT(oo.Symbol)
	side := oo.Side

	inPosition := oms.inventory.symbolInPosition(oo.Symbol)

	qty := fmt.Sprintf("%f", oo.Quantity)

	if inPosition {
		qty = fmt.Sprintf("%f", oms.inventory.Symbol[oo.Symbol].BaseQuantity)
		oms.CancelAllOrders(oo.Symbol)
	}

	order := oms.client.NewCreateOrderService().Symbol(symbolUsdt).
		Side(binanceSide(side)).Type(futures.OrderTypeMarket).
		Quantity(qty)

	res, err := order.Do(context.Background())
	if err != nil {
		log.Println("CreateEntryLimitOrder:", err)
	} else {
		oo.Active = true
		oo.BinanceId = res.OrderID

		oms.orderStack.setOrder(*oo)

	}

}

func (oms *OMS) MarketExit(oo *OpenOrder) {
	oms.CancelAllOrders(oo.Symbol)

	symbolUsdt := data.SymbolUSDT(oo.Symbol)

	positionQuantity := oms.inventory.Symbol[oo.Symbol].PositionQuantity
	qty := fmt.Sprintf("%f", positionQuantity)
	side := data.SideInvert(oms.inventory.Symbol[oo.Symbol].Side)

	order := oms.client.NewCreateOrderService().Symbol(symbolUsdt).
		Side(binanceSide(side)).Type(futures.OrderTypeMarket).
		Quantity(qty)

	_, err := order.Do(context.Background())
	if err != nil {
		log.Println("CreateEntryLimitOrder:", err)
	}
}

func (oms *OMS) CancelAllOrders(symbol string) {
	symbolUsdt := data.SymbolUSDT(symbol)

	err := oms.client.NewCancelAllOpenOrdersService().Symbol(symbolUsdt).Do(context.Background())
	if err != nil {
		log.Println("CancelAllOrders:", err)
	} else {
		oms.orderStack.removeOpenOrder(symbol)
	}

}

func (oms *OMS) CreateSlTpOrders(symbol string) {
	symbolUsdt := data.SymbolUSDT(symbol)

	position := oms.inventory.Symbol[symbol]

	sideInverted := data.SideInvert(position.Side)
	tickSize := oms.symbolRules.GetTickSize(symbol)

	qty := fmt.Sprintf("%f", position.PositionQuantity)

	fmt.Println(position.Entries)

	sl, tp := data.CalculateSlTpPrices(position.Entries[len(position.Entries)-1], position.SL, position.TP, position.Side)

	slPrice := BinancePrice(sl, tickSize)
	tpPrice := BinancePrice(tp, tickSize)
	//_ = tpPrice

	slOrder := oms.client.NewCreateOrderService().Symbol(symbolUsdt).
		TimeInForce(futures.TimeInForceTypeGTC).
		Side(binanceSide(sideInverted)).Type(futures.OrderTypeStopMarket).
		Quantity(qty).StopPrice(slPrice)
	slLogOrder := LogOrder{
		orderType: "SL",
		price:     slPrice,
		qty:       qty,
	}

	tpOrder := oms.client.NewCreateOrderService().Symbol(symbolUsdt).
		TimeInForce(futures.TimeInForceTypeGTC).
		Side(binanceSide(sideInverted)).Type(futures.OrderTypeLimit).
		Quantity(qty).Price(tpPrice)
	tpLogOrder := LogOrder{
		orderType: "TP",
		price:     tpPrice,
		qty:       qty,
	}

	// Send TakeProfit order
	err := sendOrder(tpOrder, tpLogOrder, 3)
	if err != nil {
		log.Println(err)
	}

	time.Sleep(20 * time.Millisecond)

	// Send StopLoss order
	err = sendOrder(slOrder, slLogOrder, 3)
	if err != nil {
		log.Println(err)
	}

}

type LogOrder struct {
	orderType string
	price     string
	qty       string
}

func (logOrder LogOrder) Print() {
	log.Println("ORDER LOG:")
	fmt.Println("Type  :", logOrder.orderType)
	fmt.Println("Price :", logOrder.price)
	fmt.Println("Qty   :", logOrder.qty)
}

func sendOrder(order *futures.CreateOrderService, logOrder LogOrder, retryTimes int) error {
	for i := 0; i < retryTimes; i++ {
		res, err := order.Do(context.Background())

		if err != nil {
			log.Printf("\nERROR in CreateSlTpOrders | Try #%d\n%v", i, err)
			log.Println("res:\n", res)
			logOrder.Print()
			time.Sleep(100 * time.Millisecond)

		} else {
			log.Println("Order placed successfully")
			logOrder.Print()
			return nil
		}

	}

	return errors.New("binance sendOrder: Run out of Retries")
}

// ----------------------------------------
// UTIL Functions

func binanceSide(side data.Side) futures.SideType {
	switch side {
	case data.SideLong:
		return futures.SideTypeBuy
	case data.SideShort:
		return futures.SideTypeSell
	default:
		log.Println("BinanceSide conversion ERROR! SideNone is not == 0 or 1, why?!")
		return ""
	}
}

func binanceSideToInternal(binanceSide futures.SideType) data.Side {
	switch binanceSide {
	case futures.SideTypeBuy:
		return data.SideLong
	case futures.SideTypeSell:
		return data.SideShort
	default:
		log.Println("BinanceSideToInternal: something went wrong, empty futures.SideType passed!")
		return data.SideLong
	}
}

func BinancePrice(priceFloat float64, tickSize string) string {
	// Calculate number of decimals
	var numDigitsAfterDecimal int

	for i, s := range tickSize {
		if i > 1 {
			numDigitsAfterDecimal++
			if string(s) == "1" {
				break
			}
		}
	}

	// Round float to needed number of decimals
	priceString := data.RoundToPrecisionByString(priceFloat, numDigitsAfterDecimal)

	return priceString
}

// BinanceInPosition Checks if symbol has open position on Binance
func binanceInPosition(client *futures.Client, symbol string) (isInPosition bool, position *futures.PositionRisk) {
	// Check binance this symbol has any p
	res, _ := client.NewGetPositionRiskService().Symbol(symbol).Do(context.Background())

	for _, p := range res {
		positionAmt, _ := strconv.ParseFloat(p.PositionAmt, 64)
		if positionAmt != 0 {
			return true, p
		}
	}

	return false, nil
}
