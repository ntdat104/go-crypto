package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/ntdat104/go-crypto/controller" // Assuming this path is correct
	"github.com/ntdat104/go-crypto/service"    // Assuming this path is correct
)

func main() {
	// Initialize local cache service
	localCacheService := service.NewLocalCacheService()

	// Initialize Binance Spot Service and Controller
	binanceSpotService := service.NewBinanceSpotService(localCacheService)           // Assuming this is your Spot service
	binanceSpotController := controller.NewBinanceSpotController(binanceSpotService) // Assuming this is your Spot controller

	// Initialize Binance Futures Service and Controller
	binanceFuturesService := service.NewBinanceFuturesService(localCacheService)            // Assuming this is your Futures service
	binanceFutureController := controller.NewBinanceFutureController(binanceFuturesService) // Assuming this is your Futures controller

	gin.SetMode(gin.ReleaseMode) // Set Gin to release mode for production
	router := gin.Default()      // Create a new Gin router (without default middleware)

	// Allow all CORS
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization")
		c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		// Respond to OPTIONS requests and stop further processing
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// Define API routes
	apiGroup := router.Group("/api/crypto")
	{
		// Binance Spot Endpoints
		apiGroup.GET("/ping", binanceSpotController.Ping)
		apiGroup.GET("/time", binanceSpotController.ServerTime)
		apiGroup.GET("/exchangeInfo", binanceSpotController.ExchangeInfo)
		apiGroup.GET("/ticker/price", binanceSpotController.TickerPrice)
		apiGroup.GET("/ticker/allPrices", binanceSpotController.AllPrices)
		apiGroup.GET("/bookTicker", binanceSpotController.BookTicker)
		apiGroup.GET("/depth", binanceSpotController.Depth)
		apiGroup.GET("/trades", binanceSpotController.RecentTrades)
		apiGroup.GET("/klines", binanceSpotController.Klines)
		apiGroup.GET("/historicalTrades", binanceSpotController.HistoricalTrades) // Added
		apiGroup.GET("/aggregateTrades", binanceSpotController.AggregateTrades)   // Added
		apiGroup.GET("/avgPrice", binanceSpotController.AvgPrice)                 // Added
		apiGroup.GET("/ticker/24hr", binanceSpotController.Ticker24Hr)            // Added
		apiGroup.GET("/bookTicker/all", binanceSpotController.AllBookTickers)     // Added

		// Binance Futures Endpoints
		apiGroup.GET("/futures/ping", binanceFutureController.FuturesPing)
		apiGroup.GET("/futures/time", binanceFutureController.FuturesTime)
		apiGroup.GET("/futures/exchangeInfo", binanceFutureController.FuturesExchangeInfo)
		apiGroup.GET("/futures/depth", binanceFutureController.FuturesDepth)
		apiGroup.GET("/futures/aggTrades", binanceFutureController.FuturesAggTrades)
		apiGroup.GET("/futures/ticker/price", binanceFutureController.FuturesTickerPrice)
		apiGroup.GET("/futures/ticker/allPrices", binanceFutureController.FuturesAllTickerPrices)
		apiGroup.GET("/futures/bookTicker", binanceFutureController.FuturesBookTicker)
		apiGroup.GET("/futures/klines", binanceFutureController.FuturesKlines)
		apiGroup.GET("/futures/markPrice", binanceFutureController.FuturesMarkPrice)
		apiGroup.GET("/futures/allForceOrders", binanceFutureController.FuturesAllForceOrders)
		apiGroup.GET("/futures/24hrTicker", binanceFutureController.Futures24HrTicker)
		apiGroup.GET("/futures/all24hrTickers", binanceFutureController.FuturesAll24HrTickers)
		apiGroup.GET("/futures/fundingRate", binanceFutureController.FuturesFundingRate)
		apiGroup.GET("/futures/recentTrades", binanceFutureController.FuturesRecentTrades)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Run the server
	log.Println("Starting server on :" + port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
