package controller

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/ntdat104/go-crypto/service"
)

type BinanceSpotController interface {
	Ping(ctx *gin.Context)
	ServerTime(ctx *gin.Context)
	ExchangeInfo(ctx *gin.Context)
	TickerPrice(ctx *gin.Context)
	AllPrices(ctx *gin.Context)
	BookTicker(ctx *gin.Context)
	Depth(ctx *gin.Context)
	RecentTrades(ctx *gin.Context)
	Klines(ctx *gin.Context)
	HistoricalTrades(ctx *gin.Context)
	AggregateTrades(ctx *gin.Context)
	AvgPrice(ctx *gin.Context)
	Ticker24Hr(ctx *gin.Context)
	AllBookTickers(ctx *gin.Context)
}

type binanceSpotController struct {
	binanceService service.BinanceSpotService // Changed to BinanceSpotService based on usage
}

// NewBinanceSpotController creates and returns a new BinanceSpotController instance.
func NewBinanceSpotController(binanceService service.BinanceSpotService) BinanceSpotController {
	return &binanceSpotController{
		binanceService: binanceService,
	}
}

// Ping handles the /api/v3/ping endpoint.
func (c *binanceSpotController) Ping(ctx *gin.Context) {
	resp, err := c.binanceService.GetPing()
	if err != nil {
		log.Printf("Error in Ping: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, resp)
}

// ServerTime handles the /api/v3/time endpoint.
func (c *binanceSpotController) ServerTime(ctx *gin.Context) {
	resp, err := c.binanceService.GetServerTime()
	if err != nil {
		log.Printf("Error in ServerTime: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, resp)
}

// ExchangeInfo handles the /api/v3/exchangeInfo endpoint.
func (c *binanceSpotController) ExchangeInfo(ctx *gin.Context) {
	resp, err := c.binanceService.GetExchangeInfo()
	if err != nil {
		log.Printf("Error in ExchangeInfo: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, resp)
}

// TickerPrice handles the /api/v3/ticker/price endpoint for a single symbol.
func (c *binanceSpotController) TickerPrice(ctx *gin.Context) {
	symbol := ctx.Query("symbol")
	if symbol == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "symbol query parameter is required"})
		return
	}
	resp, err := c.binanceService.GetTickerPrice(symbol)
	if err != nil {
		log.Printf("Error in TickerPrice for symbol %s: %v", symbol, err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, resp)
}

// AllPrices handles the /api/v3/ticker/price endpoint for all symbols.
func (c *binanceSpotController) AllPrices(ctx *gin.Context) {
	resp, err := c.binanceService.GetAllTickerPrices()
	if err != nil {
		log.Printf("Error in AllPrices: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, resp)
}

// BookTicker handles the /api/v3/ticker/bookTicker endpoint for a single symbol.
func (c *binanceSpotController) BookTicker(ctx *gin.Context) {
	symbol := ctx.Query("symbol")
	if symbol == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "symbol query parameter is required"})
		return
	}
	resp, err := c.binanceService.GetBookTicker(symbol)
	if err != nil {
		log.Printf("Error in BookTicker for symbol %s: %v", symbol, err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, resp)
}

// Depth handles the /api/v3/depth endpoint.
func (c *binanceSpotController) Depth(ctx *gin.Context) {
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
		log.Printf("Error in Depth for symbol %s, limit %d: %v", symbol, limit, err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, resp)
}

// RecentTrades handles the /api/v3/trades endpoint.
func (c *binanceSpotController) RecentTrades(ctx *gin.Context) {
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

	resp, err := c.binanceService.GetRecentTrades(symbol, limit)
	if err != nil {
		log.Printf("Error in RecentTrades for symbol %s, limit %d: %v", symbol, limit, err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, resp)
}

// Klines handles the /api/v3/klines endpoint.
func (c *binanceSpotController) Klines(ctx *gin.Context) {
	symbol := ctx.Query("symbol")
	interval := ctx.Query("interval")
	if symbol == "" || interval == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "symbol and interval query parameters are required"})
		return
	}
	limitStr := ctx.DefaultQuery("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid limit parameter"})
		return
	}

	resp, err := c.binanceService.GetKlines(symbol, interval, limit)
	if err != nil {
		log.Printf("Error in Klines for symbol %s, interval %s, limit %d: %v", symbol, interval, limit, err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, resp)
}

// HistoricalTrades handles the /api/v3/historicalTrades endpoint.
func (c *binanceSpotController) HistoricalTrades(ctx *gin.Context) {
	symbol := ctx.Query("symbol")
	if symbol == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "symbol query parameter is required"})
		return
	}
	limitStr := ctx.DefaultQuery("limit", "500") // Default limit for historical trades
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid limit parameter"})
		return
	}

	var fromId *int64
	fromIdStr := ctx.Query("fromId")
	if fromIdStr != "" {
		id, err := strconv.ParseInt(fromIdStr, 10, 64)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid fromId parameter"})
			return
		}
		fromId = &id
	}

	resp, err := c.binanceService.GetHistoricalTrades(symbol, limit, fromId)
	if err != nil {
		log.Printf("Error in HistoricalTrades for symbol %s, limit %d, fromId %v: %v", symbol, limit, fromId, err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, resp)
}

// AggregateTrades handles the /api/v3/aggTrades endpoint.
func (c *binanceSpotController) AggregateTrades(ctx *gin.Context) {
	symbol := ctx.Query("symbol")
	if symbol == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "symbol query parameter is required"})
		return
	}

	var fromId, startTime, endTime *int64
	limit := 500 // Default limit

	if s := ctx.Query("fromId"); s != "" {
		id, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid fromId parameter"})
			return
		}
		fromId = &id
	}
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
	if s := ctx.Query("limit"); s != "" {
		l, err := strconv.Atoi(s)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid limit parameter"})
			return
		}
		limit = l
	}

	resp, err := c.binanceService.GetAggregateTrades(symbol, fromId, startTime, endTime, limit)
	if err != nil {
		log.Printf("Error in AggregateTrades for symbol %s, fromId %v, startTime %v, endTime %v, limit %d: %v", symbol, fromId, startTime, endTime, limit, err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, resp)
}

// AvgPrice handles the /api/v3/avgPrice endpoint.
func (c *binanceSpotController) AvgPrice(ctx *gin.Context) {
	symbol := ctx.Query("symbol")
	if symbol == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "symbol query parameter is required"})
		return
	}
	resp, err := c.binanceService.GetAvgPrice(symbol)
	if err != nil {
		log.Printf("Error in AvgPrice for symbol %s: %v", symbol, err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, resp)
}

// Ticker24Hr handles the /api/v3/ticker/24hr endpoint.
func (c *binanceSpotController) Ticker24Hr(ctx *gin.Context) {
	symbol := ctx.Query("symbol")
	if symbol == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "symbol query parameter is required"})
		return
	}
	resp, err := c.binanceService.GetTicker24Hr(symbol)
	if err != nil {
		log.Printf("Error in Ticker24Hr for symbol %s: %v", symbol, err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, resp)
}

// AllBookTickers handles the /api/v3/ticker/bookTicker endpoint for all symbols.
func (c *binanceSpotController) AllBookTickers(ctx *gin.Context) {
	resp, err := c.binanceService.GetAllBookTickers()
	if err != nil {
		log.Printf("Error in AllBookTickers: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, resp)
}
