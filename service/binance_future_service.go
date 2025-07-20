package service

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"sync"
	"time"
)

// BinanceFuturesService defines the interface for interacting with the Binance Futures API.
type BinanceFuturesService interface {
	GetPing() (interface{}, error)
	GetTime() (interface{}, error)
	GetExchangeInfo() (interface{}, error)
	GetDepth(symbol string, limit int) (interface{}, error)
	GetAggTrades(symbol string, limit int) (interface{}, error)
	GetTickerPrice(symbol string) (interface{}, error)
	GetAllTickerPrices() (interface{}, error)
	GetBookTicker(symbol string) (interface{}, error)
	GetKlines(symbol, interval string, limit int) (interface{}, error)
	GetMarkPrice(symbol string) (interface{}, error)
	GetAllForceOrders(symbol string, autoCloseType string, startTime, endTime *int64, limit int) (interface{}, error)
	Get24HrTicker(symbol string) (interface{}, error)
	GetAll24HrTickers() (interface{}, error)
	GetFundingRate(symbol string, startTime, endTime *int64, limit int) (interface{}, error)
	GetRecentTrades(symbol string, limit int, fromId *int64) (interface{}, error)
}

type binanceFuturesService struct {
	futuresURL        string
	localCacheService LocalCacheService
	cacheTTL          time.Duration
	cacheDelay        time.Duration
	lock              sync.RWMutex
}

// NewBinanceFuturesService creates and returns a new BinanceFuturesService instance.
func NewBinanceFuturesService(localCacheService LocalCacheService) BinanceFuturesService {
	return &binanceFuturesService{
		futuresURL:        "https://fapi.binance.com",
		localCacheService: localCacheService,
		cacheTTL:          1 * time.Minute,
		cacheDelay:        500 * time.Millisecond,
	}
}

// fetchData makes an HTTP GET request to the given API URL with parameters.
func (s *binanceFuturesService) fetchData(apiURL string, params map[string]string) (interface{}, error) {
	u, err := url.Parse(apiURL)
	if err != nil {
		return nil, fmt.Errorf("error parsing URL: %w", err)
	}
	q := u.Query()
	for key, value := range params {
		q.Set(key, value)
	}
	u.RawQuery = q.Encode()

	resp, err := http.Get(u.String())
	if err != nil {
		return nil, fmt.Errorf("error fetching data from %s: %w", u.String(), err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received non-OK status code %d from %s, response: %s", resp.StatusCode, u.String(), resp.Status)
	}

	var response interface{}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("error decoding response from %s: %w", u.String(), err)
	}
	return response, nil
}

// fetchAndCache fetches data from the API and stores it in the local cache.
func (s *binanceFuturesService) fetchAndCache(key, delayKey, apiURL string, params map[string]string) (interface{}, error) {
	data, err := s.fetchData(apiURL, params)
	if err != nil {
		return nil, err
	}

	s.localCacheService.Set(key, data, s.cacheTTL)
	s.localCacheService.Set(delayKey, true, s.cacheDelay)
	return data, nil
}

// refreshCache asynchronously refreshes the cache for a given key if the delay period has passed.
func (s *binanceFuturesService) refreshCache(key, delayKey, apiURL string, params map[string]string) {
	s.lock.Lock()
	defer s.lock.Unlock()

	if _, delayExists := s.localCacheService.Get(delayKey); delayExists {
		return
	}
	s.localCacheService.Set(delayKey, true, s.cacheDelay)

	data, err := s.fetchData(apiURL, params)
	if err != nil {
		log.Printf("Failed to refresh futures cache for %s: %v", key, err)
		s.localCacheService.Del(delayKey)
		return
	}

	s.localCacheService.Set(key, data, s.cacheTTL)
}

// getWithCache retrieves data from cache or fetches it from the API, caching the result.
func (s *binanceFuturesService) getWithCache(cacheName, keySuffix, apiURL string, params map[string]string) (interface{}, error) {
	key := fmt.Sprintf("futures_%s:%s", cacheName, keySuffix)
	delayKey := fmt.Sprintf("futures_%s:%s:delay", cacheName, keySuffix)

	if cachedData, found := s.localCacheService.Get(key); found {
		go s.refreshCache(key, delayKey, apiURL, params)
		return cachedData, nil
	}

	return s.fetchAndCache(key, delayKey, apiURL, params)
}

// General Endpoints

// GetPing tests connectivity to the Rest API.
func (s *binanceFuturesService) GetPing() (interface{}, error) {
	return map[string]interface{}{
		"serverTime": time.Now().UnixMilli(),
		"message":    "success",
	}, nil
}

// GetTime tests connectivity to the Rest API and get the current server time.
func (s *binanceFuturesService) GetTime() (interface{}, error) {
	return map[string]interface{}{
		"serverTime": time.Now().UnixMilli(),
	}, nil
}

// GetExchangeInfo current exchange trading rules and symbol information.
func (s *binanceFuturesService) GetExchangeInfo() (interface{}, error) {
	return s.getWithCache("exchangeinfo", "global", s.futuresURL+"/fapi/v1/exchangeInfo", nil)
}

// Market Data Endpoints

// GetDepth returns the order book for a symbol.
func (s *binanceFuturesService) GetDepth(symbol string, limit int) (interface{}, error) {
	params := map[string]string{
		"symbol": symbol,
		"limit":  fmt.Sprintf("%d", limit),
	}
	return s.getWithCache("depth", fmt.Sprintf("%s-%d", symbol, limit), s.futuresURL+"/fapi/v1/depth", params)
}

// GetAggTrades Get compressed, aggregate trades.
func (s *binanceFuturesService) GetAggTrades(symbol string, limit int) (interface{}, error) {
	params := map[string]string{
		"symbol": symbol,
		"limit":  fmt.Sprintf("%d", limit),
	}
	return s.getWithCache("aggtrades", fmt.Sprintf("%s-%d", symbol, limit), s.futuresURL+"/fapi/v1/aggTrades", params)
}

// GetTickerPrice returns the latest price for a symbol or all symbols.
func (s *binanceFuturesService) GetTickerPrice(symbol string) (interface{}, error) {
	params := map[string]string{"symbol": symbol}
	return s.getWithCache("tickerprice", symbol, s.futuresURL+"/fapi/v1/ticker/price", params)
}

// GetAllTickerPrices returns the latest price for all symbols.
func (s *binanceFuturesService) GetAllTickerPrices() (interface{}, error) {
	return s.getWithCache("alltickerprices", "global", s.futuresURL+"/fapi/v1/ticker/price", nil)
}

// GetBookTicker returns the best price/qty on the order book for a symbol.
func (s *binanceFuturesService) GetBookTicker(symbol string) (interface{}, error) {
	params := map[string]string{"symbol": symbol}
	return s.getWithCache("bookticker", symbol, s.futuresURL+"/fapi/v1/ticker/bookTicker", params)
}

// GetKlines returns candlestick data for a symbol.
func (s *binanceFuturesService) GetKlines(symbol, interval string, limit int) (interface{}, error) {
	params := map[string]string{
		"symbol":   symbol,
		"interval": interval,
		"limit":    fmt.Sprintf("%d", limit),
	}
	return s.getWithCache("klines", fmt.Sprintf("%s-%s-%d", symbol, interval, limit), s.futuresURL+"/fapi/v1/klines", params)
}

// GetMarkPrice returns the Mark Price and Funding Rate.
func (s *binanceFuturesService) GetMarkPrice(symbol string) (interface{}, error) {
	params := map[string]string{"symbol": symbol}
	return s.getWithCache("markprice", symbol, s.futuresURL+"/fapi/v1/premiumIndex", params)
}

// GetAllForceOrders returns current or historical user's force orders.
func (s *binanceFuturesService) GetAllForceOrders(symbol string, autoCloseType string, startTime, endTime *int64, limit int) (interface{}, error) {
	params := map[string]string{
		"symbol": symbol,
	}
	if autoCloseType != "" {
		params["autoCloseType"] = autoCloseType
	}
	if startTime != nil {
		params["startTime"] = fmt.Sprintf("%d", *startTime)
	}
	if endTime != nil {
		params["endTime"] = fmt.Sprintf("%d", *endTime)
	}
	params["limit"] = fmt.Sprintf("%d", limit)

	keySuffix := fmt.Sprintf("%s-%d", symbol, limit)
	if autoCloseType != "" {
		keySuffix += "-" + autoCloseType
	}
	if startTime != nil {
		keySuffix += fmt.Sprintf("-s%d", *startTime)
	}
	if endTime != nil {
		keySuffix += fmt.Sprintf("-e%d", *endTime)
	}
	return s.getWithCache("allforceorders", keySuffix, s.futuresURL+"/fapi/v1/allForceOrders", params)
}

// Get24HrTicker 24hr Ticker Price Change Statistics.
func (s *binanceFuturesService) Get24HrTicker(symbol string) (interface{}, error) {
	params := map[string]string{"symbol": symbol}
	return s.getWithCache("ticker24hr", symbol, s.futuresURL+"/fapi/v1/ticker/24hr", params)
}

// GetAll24HrTickers 24hr Ticker Price Change Statistics for all symbols.
func (s *binanceFuturesService) GetAll24HrTickers() (interface{}, error) {
	return s.getWithCache("allticker24hr", "global", s.futuresURL+"/fapi/v1/ticker/24hr", nil)
}

// GetFundingRate returns the funding rate history.
func (s *binanceFuturesService) GetFundingRate(symbol string, startTime, endTime *int64, limit int) (interface{}, error) {
	params := map[string]string{
		"symbol": symbol,
	}
	if startTime != nil {
		params["startTime"] = fmt.Sprintf("%d", *startTime)
	}
	if endTime != nil {
		params["endTime"] = fmt.Sprintf("%d", *endTime)
	}
	params["limit"] = fmt.Sprintf("%d", limit)

	keySuffix := fmt.Sprintf("%s-%d", symbol, limit)
	if startTime != nil {
		keySuffix += fmt.Sprintf("-s%d", *startTime)
	}
	if endTime != nil {
		keySuffix += fmt.Sprintf("-e%d", *endTime)
	}
	return s.getWithCache("fundingrate", keySuffix, s.futuresURL+"/fapi/v1/fundingRate", params)
}

// GetRecentTrades returns recent trades for a symbol.
func (s *binanceFuturesService) GetRecentTrades(symbol string, limit int, fromId *int64) (interface{}, error) {
	params := map[string]string{
		"symbol": symbol,
		"limit":  fmt.Sprintf("%d", limit),
	}
	if fromId != nil {
		params["fromId"] = fmt.Sprintf("%d", *fromId)
	}

	keySuffix := fmt.Sprintf("%s-%d", symbol, limit)
	if fromId != nil {
		keySuffix += fmt.Sprintf("-%d", *fromId)
	}
	return s.getWithCache("recenttrades", keySuffix, s.futuresURL+"/fapi/v1/trades", params)
}
