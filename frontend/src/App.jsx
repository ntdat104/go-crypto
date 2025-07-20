import React, { useState, useEffect } from 'react';
import { Light as SyntaxHighlighter } from 'react-syntax-highlighter';
import { docco } from 'react-syntax-highlighter/dist/esm/styles/hljs'; // Or any other style you prefer

// Main App component
const App = () => {
    // State to manage the selected API endpoint
    const [selectedApi, setSelectedApi] = useState('');
    // State to store input parameters for the selected API
    const [params, setParams] = useState({});
    // State to store the API response
    const [response, setResponse] = useState(null);
    // State for loading indicator
    const [loading, setLoading] = useState(false);
    // State for error messages
    const [error, setError] = useState(null);
    // State to store the generated curl command
    const [curlCommand, setCurlCommand] = useState('');
    // State to show copy status (e.g., "Copied!")
    const [copyStatus, setCopyStatus] = useState('');
    // State to show copy status for API response
    const [responseCopyStatus, setResponseCopyStatus] = useState('');

    // Base URL for your Go backend
    const BASE_URL = 'https://go-crypto-production.up.railway.app/api/crypto';

    // Define all API endpoints and their required/optional parameters with descriptions
    const apiEndpoints = {
        // Binance Spot API
        'Spot Ping': {
            path: '/ping', method: 'GET', params: [],
            description: 'Test connectivity to the Binance Spot API. Returns an empty object on success.'
        },
        'Spot Server Time': {
            path: '/time', method: 'GET', params: [],
            description: 'Test connectivity to the Binance Spot API and get the current server time.'
        },
        'Spot Exchange Info': {
            path: '/exchangeInfo', method: 'GET', params: [],
            description: 'Current exchange trading rules and symbol information for Binance Spot.'
        },
        'Spot Ticker Price (Single)': {
            path: '/ticker/price', method: 'GET', params: [{ name: 'symbol', type: 'text', required: true, defaultValue: 'BTCUSDT', description: 'The cryptocurrency trading pair symbol (e.g., BTCUSDT).' }],
            description: 'Latest price for a specific trading pair on Binance Spot.'
        },
        'Spot All Ticker Prices': {
            path: '/ticker/allPrices', method: 'GET', params: [],
            description: 'Latest prices for all trading pairs on Binance Spot.'
        },
        'Spot Book Ticker (Single)': {
            path: '/bookTicker', method: 'GET', params: [{ name: 'symbol', type: 'text', required: true, defaultValue: 'BTCUSDT', description: 'The cryptocurrency trading pair symbol (e.g., BTCUSDT).' }],
            description: 'Best price/quantity on the order book for a specific trading pair on Binance Spot.'
        },
        'Spot Depth': {
            path: '/depth', method: 'GET', params: [
                { name: 'symbol', type: 'text', required: true, defaultValue: 'BTCUSDT', description: 'The cryptocurrency trading pair symbol (e.g., BTCUSDT).' },
                { name: 'limit', type: 'number', required: false, defaultValue: 10, description: 'Limit the number of bids and asks to return (default: 10, max: 1000).' }
            ],
            description: 'Order book depth information for a specific trading pair on Binance Spot.'
        },
        'Spot Recent Trades': {
            path: '/trades', method: 'GET', params: [
                { name: 'symbol', type: 'text', required: true, defaultValue: 'BTCUSDT', description: 'The cryptocurrency trading pair symbol (e.g., BTCUSDT).' },
                { name: 'limit', type: 'number', required: false, defaultValue: 10, description: 'Limit the number of recent trades to return (default: 10, max: 1000).' }
            ],
            description: 'Get recent trades for a specific trading pair on Binance Spot.'
        },
        'Spot Klines': {
            path: '/klines', method: 'GET', params: [
                { name: 'symbol', type: 'text', required: true, defaultValue: 'BTCUSDT', description: 'The cryptocurrency trading pair symbol (e.e., BTCUSDT).' },
                { name: 'interval', type: 'text', required: true, defaultValue: '1h', description: 'Kline interval (e.g., 1m, 5m, 1h, 1d).' },
                { name: 'limit', type: 'number', required: false, defaultValue: 10, description: 'Limit the number of klines to return (default: 10, max: 1000).' }
            ],
            description: 'Kline/Candlestick data for a specific trading pair on Binance Spot.'
        },
        'Spot Historical Trades': {
            path: '/historicalTrades', method: 'GET', params: [
                { name: 'symbol', type: 'text', required: true, defaultValue: 'BTCUSDT', description: 'The cryptocurrency trading pair symbol (e.g., BTCUSDT).' },
                { name: 'limit', type: 'number', required: false, defaultValue: 500, description: 'Limit the number of historical trades to return (default: 500, max: 1000).' },
                { name: 'fromId', type: 'number', required: false, description: 'Trade ID to fetch from. All trades with ID >= fromId will be returned.' }
            ],
            description: 'Get historical trades for a specific trading pair on Binance Spot.'
        },
        'Spot Aggregate Trades': {
            path: '/aggregateTrades', method: 'GET', params: [
                { name: 'symbol', type: 'text', required: true, defaultValue: 'BTCUSDT', description: 'The cryptocurrency trading pair symbol (e.g., BTCUSDT).' },
                { name: 'fromId', type: 'number', required: false, description: 'Trade ID to fetch from. All trades with ID >= fromId will be returned.' },
                { name: 'startTime', type: 'number', required: false, description: 'Timestamp in milliseconds to get aggregate trades from.' },
                { name: 'endTime', type: 'number', required: false, description: 'Timestamp in milliseconds to get aggregate trades until.' },
                { name: 'limit', type: 'number', required: false, defaultValue: 500, description: 'Limit the number of aggregate trades to return (default: 500, max: 1000).' }
            ],
            description: 'Get compressed, aggregate trades for a specific trading pair on Binance Spot.'
        },
        'Spot Average Price': {
            path: '/avgPrice', method: 'GET', params: [{ name: 'symbol', type: 'text', required: true, defaultValue: 'BTCUSDT', description: 'The cryptocurrency trading pair symbol (e.g., BTCUSDT).' }],
            description: 'Current average price for a specific trading pair on Binance Spot.'
        },
        'Spot 24hr Ticker (Single)': {
            path: '/ticker/24hr', method: 'GET', params: [{ name: 'symbol', type: 'text', required: true, defaultValue: 'BTCUSDT', description: 'The cryptocurrency trading pair symbol (e.g., BTCUSDT).' }],
            description: '24-hour rolling window price change statistics for a specific trading pair on Binance Spot.'
        },
        'Spot All Book Tickers': {
            path: '/bookTicker/all', method: 'GET', params: [],
            description: 'Best price/quantity on the order book for all trading pairs on Binance Spot.'
        },

        // Binance Futures API
        'Futures Ping': {
            path: '/futures/ping', method: 'GET', params: [],
            description: 'Test connectivity to the Binance Futures API. Returns an empty object on success.'
        },
        'Futures Time': {
            path: '/futures/time', method: 'GET', params: [],
            description: 'Test connectivity to the Binance Futures API and get the current server time.'
        },
        'Futures Exchange Info': {
            path: '/futures/exchangeInfo', method: 'GET', params: [],
            description: 'Current exchange trading rules and symbol information for Binance Futures.'
        },
        'Futures Depth': {
            path: '/futures/depth', method: 'GET', params: [
                { name: 'symbol', type: 'text', required: true, defaultValue: 'BTCUSDT', description: 'The cryptocurrency trading pair symbol (e.g., BTCUSDT).' },
                { name: 'limit', type: 'number', required: false, defaultValue: 10, description: 'Limit the number of bids and asks to return (default: 10, max: 1000).' }
            ],
            description: 'Order book depth information for a specific trading pair on Binance Futures.'
        },
        'Futures Aggregate Trades': {
            path: '/futures/aggTrades', method: 'GET', params: [
                { name: 'symbol', type: 'text', required: true, defaultValue: 'BTCUSDT', description: 'The cryptocurrency trading pair symbol (e.g., BTCUSDT).' },
                { name: 'limit', type: 'number', required: false, defaultValue: 500, description: 'Limit the number of aggregate trades to return (default: 500, max: 1000).' }
            ],
            description: 'Get compressed, aggregate trades for a specific trading pair on Binance Futures.'
        },
        'Futures Ticker Price (Single)': {
            path: '/futures/ticker/price', method: 'GET', params: [{ name: 'symbol', type: 'text', required: true, defaultValue: 'BTCUSDT', description: 'The cryptocurrency trading pair symbol (e.g., BTCUSDT).' }],
            description: 'Latest price for a specific trading pair on Binance Futures.'
        },
        'Futures All Ticker Prices': {
            path: '/futures/ticker/allPrices', method: 'GET', params: [],
            description: 'Latest prices for all trading pairs on Binance Futures.'
        },
        'Futures Book Ticker': {
            path: '/futures/bookTicker', method: 'GET', params: [{ name: 'symbol', type: 'text', required: true, defaultValue: 'BTCUSDT', description: 'The cryptocurrency trading pair symbol (e.g., BTCUSDT).' }],
            description: 'Best price/quantity on the order book for a specific trading pair on Binance Futures.'
        },
        'Futures Klines': {
            path: '/futures/klines', method: 'GET', params: [
                { name: 'symbol', type: 'text', required: true, defaultValue: 'BTCUSDT', description: 'The cryptocurrency trading pair symbol (e.g., BTCUSDT).' },
                { name: 'interval', type: 'text', required: true, defaultValue: '1h', description: 'Kline interval (e.g., 1m, 5m, 1h, 1d).' },
                { name: 'limit', type: 'number', required: false, defaultValue: 500, description: 'Limit the number of klines to return (default: 500, max: 1500).' }
            ],
            description: 'Kline/Candlestick data for a specific trading pair on Binance Futures.'
        },
        'Futures Mark Price': {
            path: '/futures/markPrice', method: 'GET', params: [{ name: 'symbol', type: 'text', required: true, defaultValue: 'BTCUSDT', description: 'The cryptocurrency trading pair symbol (e.g., BTCUSDT).' }],
            description: 'Get mark price for a specific trading pair on Binance Futures.'
        },
        'Futures All Force Orders': {
            path: '/futures/allForceOrders', method: 'GET', params: [
                { name: 'symbol', type: 'text', required: true, defaultValue: 'BTCUSDT', description: 'The cryptocurrency trading pair symbol (e.g., BTCUSDT).' },
                { name: 'autoCloseType', type: 'text', required: false, description: 'Filter by auto close type (e.g., LIQUIDATION, ADL).' },
                { name: 'startTime', type: 'number', required: false, description: 'Timestamp in milliseconds to get force orders from.' },
                { name: 'endTime', type: 'number', required: false, description: 'Timestamp in milliseconds to get force orders until.' },
                { name: 'limit', type: 'number', required: false, defaultValue: 500, description: 'Limit the number of force orders to return (default: 500, max: 1000).' }
            ],
            description: 'Get all liquidation orders for a specific trading pair on Binance Futures.'
        },
        'Futures 24hr Ticker (Single)': {
            path: '/futures/24hrTicker', method: 'GET', params: [{ name: 'symbol', type: 'text', required: true, defaultValue: 'BTCUSDT', description: 'The cryptocurrency trading pair symbol (e.g., BTCUSDT).' }],
            description: '24-hour rolling window price change statistics for a specific trading pair on Binance Futures.'
        },
        'Futures All 24hr Tickers': {
            path: '/futures/all24hrTickers', method: 'GET', params: [],
            description: '24-hour rolling window price change statistics for all trading pairs on Binance Futures.'
        },
        'Futures Funding Rate': {
            path: '/futures/fundingRate', method: 'GET', params: [
                { name: 'symbol', type: 'text', required: true, defaultValue: 'BTCUSDT', description: 'The cryptocurrency trading pair symbol (e.g., BTCUSDT).' },
                { name: 'startTime', type: 'number', required: false, description: 'Timestamp in milliseconds to get funding rates from.' },
                { name: 'endTime', type: 'number', required: false, description: 'Timestamp in milliseconds to get funding rates until.' },
                { name: 'limit', type: 'number', required: false, defaultValue: 100, description: 'Limit the number of funding rates to return (default: 100, max: 1000).' }
            ],
            description: 'Get funding rate history for a specific trading pair on Binance Futures.'
        },
        'Futures Recent Trades': {
            path: '/futures/recentTrades', method: 'GET', params: [
                { name: 'symbol', type: 'text', required: true, defaultValue: 'BTCUSDT', description: 'The cryptocurrency trading pair symbol (e.g., BTCUSDT).' },
                { name: 'limit', type: 'number', required: false, defaultValue: 500, description: 'Limit the number of recent trades to return (default: 500, max: 1000).' },
                { name: 'fromId', type: 'number', required: false, description: 'Trade ID to fetch from. All trades with ID >= fromId will be returned.' }
            ],
            description: 'Get recent trades for a specific trading pair on Binance Futures.'
        },
    };

    // Function to construct the URL with query parameters
    const buildUrl = (apiName, currentParams) => {
        const apiInfo = apiEndpoints[apiName];
        if (!apiInfo) return '';
        let url = `${BASE_URL}${apiInfo.path}`;
        const queryParams = new URLSearchParams();

        apiInfo.params.forEach(param => {
            const paramValue = currentParams[param.name];
            if (paramValue !== '' && paramValue !== null && paramValue !== undefined) {
                queryParams.append(param.name, paramValue);
            }
        });

        if (queryParams.toString()) {
            url += `?${queryParams.toString()}`;
        }
        return url;
    };

    // Function to build the curl command
    const buildCurlCommand = (apiName, currentParams) => {
        const url = buildUrl(apiName, currentParams);
        const apiInfo = apiEndpoints[apiName];
        if (!apiInfo) return '';
        return `curl -X ${apiInfo.method} "${url}"`;
    };

    // Effect to reset parameters and update curl when a new API is selected
    useEffect(() => {
        if (selectedApi) {
            const apiInfo = apiEndpoints[selectedApi];
            const newParams = {};
            apiInfo.params.forEach(param => {
                newParams[param.name] = param.defaultValue !== undefined ? param.defaultValue : '';
            });
            setParams(newParams);
            setResponse(null); // Clear previous response
            setError(null); // Clear previous error
            setCurlCommand(buildCurlCommand(selectedApi, newParams)); // Update curl command
        } else {
            setParams({});
            setResponse(null);
            setError(null);
            setCurlCommand('');
        }
    }, [selectedApi]);

    // Handle input changes for parameters
    const handleParamChange = (e) => {
        const { name, value, type } = e.target;
        setParams(prevParams => {
            const updatedParams = {
                ...prevParams,
                [name]: type === 'number' ? (value === '' ? '' : Number(value)) : value,
            };
            setCurlCommand(buildCurlCommand(selectedApi, updatedParams)); // Update curl command immediately
            return updatedParams;
        });
    };

    // Handle API call submission
    const handleSubmit = async () => {
        setLoading(true);
        setResponse(null);
        setError(null);
        setResponseCopyStatus(''); // Clear copy status for response

        const url = buildUrl(selectedApi, params); // Use the updated buildUrl
        console.log("Calling API:", url); // Log the URL for debugging

        try {
            const res = await fetch(url);
            if (!res.ok) {
                const errorData = await res.json();
                throw new Error(errorData.error || `HTTP error! status: ${res.status}`);
            }
            const data = await res.json();
            setResponse(data);
        } catch (err) {
            console.error("API call error:", err);
            setError(err.message);
        } finally {
            setLoading(false);
        }
    };

    // Generic function to copy text to clipboard
    const copyToClipboard = (textToCopy, setStatus) => {
        const textarea = document.createElement('textarea');
        textarea.value = textToCopy;
        document.body.appendChild(textarea);
        textarea.select();
        try {
            document.execCommand('copy');
            setStatus('Copied!');
            setTimeout(() => setStatus(''), 2000); // Clear message after 2 seconds
        } catch (err) {
            setStatus('Failed to copy!');
            console.error('Failed to copy text: ', err);
        }
        document.body.removeChild(textarea);
    };

    return (
        <div className="min-h-screen bg-gray-100 p-4 font-sans antialiased flex flex-col items-center">
            <div className="bg-white p-6 rounded-lg shadow-lg w-full max-w-6xl md:flex md:flex-row md:space-x-6">
                {/* Left Column: Controls and Documentation */}
                <div className="md:w-1/2 flex flex-col space-y-6">
                    {/* API Selection */}
                    <div className="bg-gray-50 p-4 rounded-lg border border-gray-200">
                        <label htmlFor="api-select" className="block text-gray-700 text-sm font-bold mb-2">
                            Select API Endpoint:
                        </label>
                        <div className="relative">
                            <select
                                id="api-select"
                                className="block appearance-none w-full bg-white border border-gray-300 text-gray-700 py-3 px-4 pr-8 rounded-lg leading-tight focus:outline-none focus:border-blue-500 shadow-sm"
                                value={selectedApi}
                                onChange={(e) => setSelectedApi(e.target.value)}
                            >
                                <option value="">-- Please choose an API --</option>
                                {Object.keys(apiEndpoints).map((apiName) => (
                                    <option key={apiName} value={apiName}>
                                        {apiName}
                                    </option>
                                ))}
                            </select>
                            <div className="pointer-events-none absolute inset-y-0 right-0 flex items-center px-2 text-gray-700">
                                <svg className="fill-current h-4 w-4" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20"><path d="M9.293 12.95l.707.707L15.657 8l-1.414-1.414L10 10.828 5.757 6.586 4.343 8z"/></svg>
                            </div>
                        </div>
                    </div>

                    {selectedApi && (
                        <>
                            {/* API Documentation Block */}
                            <div className="p-4 border border-gray-200 rounded-lg bg-gray-50">
                                <h2 className="text-xl font-semibold text-gray-800 mb-2">API Documentation</h2>
                                <p className="text-gray-700 mb-4">
                                    {apiEndpoints[selectedApi].description}
                                </p>

                                <h3 className="text-lg font-semibold text-gray-800 mb-2">Parameters</h3>
                                {apiEndpoints[selectedApi].params.length > 0 ? (
                                    <div className="grid grid-cols-1 gap-2">
                                        {apiEndpoints[selectedApi].params.map((param) => (
                                            <div key={param.name} className="flex flex-col">
                                                <p className="text-gray-700 text-sm">
                                                    <span className="font-bold">{param.name}</span>: {param.type}
                                                    {param.required ? <span className="text-red-500 ml-1">(required)</span> : <span className="text-gray-500 ml-1">(optional)</span>}
                                                </p>
                                                {param.description && (
                                                    <p className="text-gray-600 text-xs ml-2 italic">{param.description}</p>
                                                )}
                                            </div>
                                        ))}
                                    </div>
                                ) : (
                                    <p className="text-gray-600 text-sm">No parameters required for this endpoint.</p>
                                )}
                            </div>

                            {/* Parameters Input */}
                            <div className="p-4 border border-gray-200 rounded-lg bg-gray-50">
                                <h2 className="text-xl font-semibold text-gray-800 mb-4">Input Parameters</h2>
                                {apiEndpoints[selectedApi].params.length > 0 ? (
                                    <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                                        {apiEndpoints[selectedApi].params.map((param) => (
                                            <div key={param.name} className="flex flex-col">
                                                <label htmlFor={param.name} className="text-gray-700 text-sm font-bold mb-1">
                                                    {param.name} {param.required && <span className="text-red-500">*</span>}
                                                </label>
                                                <input
                                                    type={param.type}
                                                    id={param.name}
                                                    name={param.name}
                                                    value={params[param.name] || ''}
                                                    onChange={handleParamChange}
                                                    placeholder={param.defaultValue ? `e.g., ${param.defaultValue}` : `Enter ${param.name}`}
                                                    className="shadow-sm appearance-none border rounded-lg w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline focus:border-blue-500"
                                                    required={param.required}
                                                />
                                            </div>
                                        ))}
                                    </div>
                                ) : (
                                    <p className="text-gray-600">No parameters to input for this endpoint.</p>
                                )}
                            </div>

                            {/* Curl Command Display */}
                            <div className="p-4 border border-gray-200 rounded-lg bg-gray-50">
                                <h2 className="text-xl font-semibold text-gray-800 mb-4 flex justify-between items-center">
                                    Curl Command
                                    {curlCommand && (
                                        <button
                                            onClick={() => copyToClipboard(curlCommand, setCopyStatus)}
                                            className="bg-gray-200 hover:bg-gray-300 text-gray-800 text-sm font-semibold py-1 px-3 rounded-lg flex items-center transition duration-200 ease-in-out"
                                        >
                                            <svg className="w-4 h-4 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M8 5H6a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2v-1M8 5a2 2 0 002 2h2a2 2 0 002-2M8 5a2 2 0 012-2h2a2 2 0 012 2m0 0h2a2 2 0 012 2v3m2 0h-2M10 18H8m4 0h2m-4 0h2"></path></svg>
                                            {copyStatus || 'Copy'}
                                        </button>
                                    )}
                                </h2>
                                {curlCommand ? (
                                    <SyntaxHighlighter
                                        language="bash"
                                        style={docco}
                                        showLineNumbers={false}
                                        wrapLines={true}
                                        customStyle={{
                                            backgroundColor: '#1f2937', // bg-gray-800
                                            color: '#f9fafb', // text-white
                                            padding: '1rem', // p-4
                                            borderRadius: '0.5rem', // rounded-lg
                                            overflow: 'auto',
                                            maxHeight: '16rem', // max-h-64
                                            fontSize: '0.875rem', // text-sm
                                        }}
                                    >
                                        {curlCommand}
                                    </SyntaxHighlighter>
                                ) : (
                                    <p className="text-gray-600">Select an API and input parameters to generate the curl command.</p>
                                )}
                            </div>

                            {/* Submit Button */}
                            <div className="text-center">
                                <button
                                    onClick={handleSubmit}
                                    disabled={loading}
                                    className="bg-blue-600 hover:bg-blue-700 text-white font-bold py-3 px-6 rounded-lg shadow-md transition duration-300 ease-in-out transform hover:scale-105 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-opacity-75 disabled:opacity-50 disabled:cursor-not-allowed"
                                >
                                    {loading ? 'Loading...' : 'Call API'}
                                </button>
                            </div>
                        </>
                    )}
                </div>

                {/* Right Column: API Response */}
                <div className="md:w-1/2 mt-6 md:mt-0 p-4 border border-gray-200 rounded-lg bg-gray-50 flex flex-col">
                    <h2 className="text-xl font-semibold text-gray-800 mb-4 flex justify-between items-center">
                        API Response
                        {response && (
                            <button
                                onClick={() => copyToClipboard(JSON.stringify(response, null, 2), setResponseCopyStatus)}
                                className="bg-gray-200 hover:bg-gray-300 text-gray-800 text-sm font-semibold py-1 px-3 rounded-lg flex items-center transition duration-200 ease-in-out"
                            >
                                <svg className="w-4 h-4 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M8 5H6a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2v-1M8 5a2 2 0 002 2h2a2 2 0 002-2M8 5a2 2 0 012-2h2a2 2 0 012 2m0 0h2a2 2 0 012 2v3m2 0h-2M10 18H8m4 0h2m-4 0h2"></path></svg>
                                {responseCopyStatus || 'Copy'}
                            </button>
                        )}
                    </h2>
                    {loading && (
                        <div className="flex items-center justify-center py-4">
                            <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-gray-900"></div>
                            <span className="ml-3 text-gray-700">Fetching data...</span>
                        </div>
                    )}
                    {error && (
                        <div className="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded-lg relative" role="alert">
                            <strong className="font-bold">Error!</strong>
                            <span className="block sm:inline ml-2">{error}</span>
                        </div>
                    )}
                    {response && (
                        <SyntaxHighlighter
                            language="json"
                            style={docco}
                            showLineNumbers={false}
                            wrapLines={true}
                            customStyle={{
                                backgroundColor: '#1f2937', // bg-gray-800
                                color: '#f9fafb', // text-white
                                padding: '1rem', // p-4
                                borderRadius: '0.5rem', // rounded-lg
                                overflow: 'auto',
                                flexGrow: 1, // flex-grow
                                maxHeight: '24rem', // max-h-96
                                fontSize: '0.875rem', // text-sm
                            }}
                        >
                            {JSON.stringify(response, null, 2)}
                        </SyntaxHighlighter>
                    )}
                    {!loading && !error && !response && (
                        <p className="text-gray-600">API response will appear here.</p>
                    )}
                </div>
            </div>
        </div>
    );
};

export default App;
