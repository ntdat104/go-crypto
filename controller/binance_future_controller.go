package controller

import (
	"log"
	"net/http"
	"strconv"

	// Import time package for parsing timestamps
	// Make sure this path is correct for your project
	"github.com/gin-gonic/gin"
	"github.com/ntdat104/go-crypto/service"
)

type BinanceFutureController interface {
	FuturesPing(ctx *gin.Context)
	FuturesTime(ctx *gin.Context)
	FuturesExchangeInfo(ctx *gin.Context)
	FuturesDepth(ctx *gin.Context)
	FuturesAggTrades(ctx *gin.Context)
	FuturesTickerPrice(ctx *gin.Context)
	FuturesAllTickerPrices(ctx *gin.Context) // New method
	FuturesBookTicker(ctx *gin.Context)      // New method
	FuturesKlines(ctx *gin.Context)
	FuturesMarkPrice(ctx *gin.Context)      // New method
	FuturesAllForceOrders(ctx *gin.Context) // New method
	Futures24HrTicker(ctx *gin.Context)     // New method
	FuturesAll24HrTickers(ctx *gin.Context) // New method
	FuturesFundingRate(ctx *gin.Context)    // New method
	FuturesRecentTrades(ctx *gin.Context)   // New method
}

type binanceFutureController struct {
	// IMPORTANT: Change this to service.BinanceFuturesService
	binanceService service.BinanceFuturesService
}

// NewBinanceFutureController creates and returns a new BinanceFutureController instance.
func NewBinanceFutureController(binanceService service.BinanceFuturesService) BinanceFutureController {
	return &binanceFutureController{
		binanceService: binanceService,
	}
}

// FuturesPing handles the /fapi/v1/ping endpoint.
func (c *binanceFutureController) FuturesPing(ctx *gin.Context) {
	resp, err := c.binanceService.GetPing()
	if err != nil {
		log.Printf("Error in FuturesPing: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, resp)
}

// FuturesTime handles the /fapi/v1/time endpoint.
func (c *binanceFutureController) FuturesTime(ctx *gin.Context) {
	resp, err := c.binanceService.GetTime()
	if err != nil {
		log.Printf("Error in FuturesTime: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, resp)
}

// FuturesExchangeInfo handles the /fapi/v1/exchangeInfo endpoint.
func (c *binanceFutureController) FuturesExchangeInfo(ctx *gin.Context) {
	resp, err := c.binanceService.GetExchangeInfo()
	if err != nil {
		log.Printf("Error in FuturesExchangeInfo: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, resp)
}

// FuturesDepth handles the /fapi/v1/depth endpoint.
func (c *binanceFutureController) FuturesDepth(ctx *gin.Context) {
	symbol := ctx.Query("symbol")
	if symbol == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "symbol query parameter is required"})
		return
	}
	limitStr := ctx.DefaultQuery("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid limit parameter"})
		return
	}

	resp, err := c.binanceService.GetDepth(symbol, limit)
	if err != nil {
		log.Printf("Error in FuturesDepth for symbol %s, limit %d: %v", symbol, limit, err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, resp)
}

// FuturesAggTrades handles the /fapi/v1/aggTrades endpoint.
func (c *binanceFutureController) FuturesAggTrades(ctx *gin.Context) {
	symbol := ctx.Query("symbol")
	if symbol == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "symbol query parameter is required"})
		return
	}
	limitStr := ctx.DefaultQuery("limit", "500") // Default for aggTrades might be higher
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid limit parameter"})
		return
	}

	resp, err := c.binanceService.GetAggTrades(symbol, limit)
	if err != nil {
		log.Printf("Error in FuturesAggTrades for symbol %s, limit %d: %v", symbol, limit, err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, resp)
}

// FuturesTickerPrice handles the /fapi/v1/ticker/price endpoint for a single symbol.
func (c *binanceFutureController) FuturesTickerPrice(ctx *gin.Context) {
	symbol := ctx.Query("symbol")
	if symbol == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "symbol query parameter is required"})
		return
	}
	resp, err := c.binanceService.GetTickerPrice(symbol)
	if err != nil {
		log.Printf("Error in FuturesTickerPrice for symbol %s: %v", symbol, err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, resp)
}

// FuturesAllTickerPrices handles the /fapi/v1/ticker/price endpoint for all symbols.
func (c *binanceFutureController) FuturesAllTickerPrices(ctx *gin.Context) {
	resp, err := c.binanceService.GetAllTickerPrices()
	if err != nil {
		log.Printf("Error in FuturesAllTickerPrices: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, resp)
}

// FuturesBookTicker returns the best price/qty on the order book for a symbol.
func (c *binanceFutureController) FuturesBookTicker(ctx *gin.Context) {
	symbol := ctx.Query("symbol")
	if symbol == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "symbol query parameter is required"})
		return
	}
	resp, err := c.binanceService.GetBookTicker(symbol)
	if err != nil {
		log.Printf("Error in FuturesBookTicker for symbol %s: %v", symbol, err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, resp)
}

// FuturesKlines handles the /fapi/v1/klines endpoint.
func (c *binanceFutureController) FuturesKlines(ctx *gin.Context) {
	symbol := ctx.Query("symbol")
	interval := ctx.Query("interval")
	if symbol == "" || interval == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "symbol and interval query parameters are required"})
		return
	}
	limitStr := ctx.DefaultQuery("limit", "500")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid limit parameter"})
		return
	}

	resp, err := c.binanceService.GetKlines(symbol, interval, limit)
	if err != nil {
		log.Printf("Error in FuturesKlines for symbol %s, interval %s, limit %d: %v", symbol, interval, limit, err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, resp)
}

// FuturesMarkPrice handles the /fapi/v1/premiumIndex endpoint.
func (c *binanceFutureController) FuturesMarkPrice(ctx *gin.Context) {
	symbol := ctx.Query("symbol")
	if symbol == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "symbol query parameter is required"})
		return
	}
	resp, err := c.binanceService.GetMarkPrice(symbol)
	if err != nil {
		log.Printf("Error in FuturesMarkPrice for symbol %s: %v", symbol, err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, resp)
}

// FuturesAllForceOrders handles the /fapi/v1/allForceOrders endpoint.
func (c *binanceFutureController) FuturesAllForceOrders(ctx *gin.Context) {
	symbol := ctx.Query("symbol")
	if symbol == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "symbol query parameter is required"})
		return
	}

	autoCloseType := ctx.Query("autoCloseType") // Optional

	var startTime, endTime *int64
	if s := ctx.Query("startTime"); s != "" {
		t, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid startTime parameter"})
			return
		}
		startTime = &t
	}
	if s := ctx.Query("endTime"); s != "" {
		t, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid endTime parameter"})
			return
		}
		endTime = &t
	}

	limitStr := ctx.DefaultQuery("limit", "500")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid limit parameter"})
		return
	}

	resp, err := c.binanceService.GetAllForceOrders(symbol, autoCloseType, startTime, endTime, limit)
	if err != nil {
		log.Printf("Error in FuturesAllForceOrders for symbol %s: %v", symbol, err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, resp)
}

// Futures24HrTicker handles the /fapi/v1/ticker/24hr endpoint for a single symbol.
func (c *binanceFutureController) Futures24HrTicker(ctx *gin.Context) {
	symbol := ctx.Query("symbol")
	if symbol == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "symbol query parameter is required"})
		return
	}
	resp, err := c.binanceService.Get24HrTicker(symbol)
	if err != nil {
		log.Printf("Error in Futures24HrTicker for symbol %s: %v", symbol, err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, resp)
}

// FuturesAll24HrTickers handles the /fapi/v1/ticker/24hr endpoint for all symbols.
func (c *binanceFutureController) FuturesAll24HrTickers(ctx *gin.Context) {
	resp, err := c.binanceService.GetAll24HrTickers()
	if err != nil {
		log.Printf("Error in FuturesAll24HrTickers: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, resp)
}

// FuturesFundingRate handles the /fapi/v1/fundingRate endpoint.
func (c *binanceFutureController) FuturesFundingRate(ctx *gin.Context) {
	symbol := ctx.Query("symbol")
	if symbol == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "symbol query parameter is required"})
		return
	}

	var startTime, endTime *int64
	if s := ctx.Query("startTime"); s != "" {
		t, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid startTime parameter"})
			return
		}
		startTime = &t
	}
	if s := ctx.Query("endTime"); s != "" {
		t, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid endTime parameter"})
			return
		}
		endTime = &t
	}

	limitStr := ctx.DefaultQuery("limit", "100") // Default limit for funding rate
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid limit parameter"})
		return
	}

	resp, err := c.binanceService.GetFundingRate(symbol, startTime, endTime, limit)
	if err != nil {
		log.Printf("Error in FuturesFundingRate for symbol %s: %v", symbol, err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, resp)
}

// FuturesRecentTrades handles the /fapi/v1/trades endpoint.
func (c *binanceFutureController) FuturesRecentTrades(ctx *gin.Context) {
	symbol := ctx.Query("symbol")
	if symbol == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "symbol query parameter is required"})
		return
	}
	limitStr := ctx.DefaultQuery("limit", "500")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid limit parameter"})
		return
	}

	var fromId *int64
	if s := ctx.Query("fromId"); s != "" {
		id, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid fromId parameter"})
			return
		}
		fromId = &id
	}

	resp, err := c.binanceService.GetRecentTrades(symbol, limit, fromId)
	if err != nil {
		log.Printf("Error in FuturesRecentTrades for symbol %s, limit %d, fromId %v: %v", symbol, limit, fromId, err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, resp)
}
