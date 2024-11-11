package data

import (
	"fmt"
	"math"
	"strings"
	"ttk/ta"
)

func SymbolUSDT(symbol string) string {
	return symbol + "USDT"
}

func SymbolUSDTList(symbols []string) []string {
	symbolsUSDT := make([]string, len(symbols))

	for i, s := range symbols {
		symbolsUSDT[i] = SymbolUSDT(s)
	}

	return symbolsUSDT

}

func SymbolNoUSDT(symbolUSDT string) string {
	return strings.ReplaceAll(symbolUSDT, "USDT", "")
}

func SymbolNoUSDTList(symbols []string) []string {
	var symbolsNoUsdtList []string

	for _, s := range symbols {
		if strings.Contains(s, "USDT") && !strings.Contains(s, "USDT_") {
			symbolsNoUsdtList = append(symbolsNoUsdtList, SymbolNoUSDT(s))
		}
	}

	return symbolsNoUsdtList

}

func ConvertQuantityFromUSDT(qtyUSD float64, symbolPrice float64, qtyPrecision int) float64 {
	var qtySymbol float64
	if symbolPrice != 0 {
		qtySymbol = RoundToPrecision(qtyUSD/symbolPrice, qtyPrecision)
	} else {
		qtySymbol = 0
	}

	return qtySymbol
}

func RoundToPrecision(x float64, precision int) float64 {
	mult := 1.0

	if precision != 0 {
		for i := 0; i < precision; i++ {
			mult *= 10
		}
	}

	x = math.Round(x*mult) / mult

	return x
}

func RoundToPrecisionByString(x float64, precision int) string {
	switch precision {
	case 0:
		return fmt.Sprintf("%.0f", x)
	case 1:
		return fmt.Sprintf("%.1f", x)
	case 2:
		return fmt.Sprintf("%.2f", x)
	case 3:
		return fmt.Sprintf("%.3f", x)
	case 4:
		return fmt.Sprintf("%.4f", x)
	case 5:
		return fmt.Sprintf("%.5f", x)
	case 6:
		return fmt.Sprintf("%.6f", x)
	case 7:
		return fmt.Sprintf("%.7f", x)
	case 8:
		return fmt.Sprintf("%.8f", x)
	case 9:
		return fmt.Sprintf("%.9f", x)
	case 10:
		return fmt.Sprintf("%.10f", x)
	default:
		return fmt.Sprintf("%.4f", x)
	}
}

func CalculateSlTpPrices(p, pctSL, pctTP float64, side Side) (sl, tp float64) {

	// LONG SL, TP
	if side == SideLong {
		sl = p / ta.Pct(pctSL)
		tp = p * ta.Pct(pctTP)
	} else {
		// SHORT SL, TP
		sl = p * ta.Pct(pctSL)
		tp = p / ta.Pct(pctTP)
	}

	return
}

func Average(prices []float64) float64 {
	var sum float64

	for _, p := range prices {
		sum += p
	}

	return sum / float64(len(prices))
}
