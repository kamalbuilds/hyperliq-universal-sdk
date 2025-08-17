using System;
using System.Net.Http;
using System.Threading;
using System.Threading.Tasks;
using Microsoft.Extensions.Logging;
using HyperliquidSDK.Auth;
using HyperliquidSDK.Types;
using HyperliquidSDK.WebSocket;
using Polly;
using Polly.Extensions.Http;

namespace HyperliquidSDK.Client
{
    public class HyperliquidClient : IDisposable
    {
        private const string MainnetApi = "https://api.hyperliquid.xyz";
        private const string TestnetApi = "https://api.hyperliquid-testnet.xyz";
        private const string MainnetWs = "wss://api.hyperliquid.xyz/ws";
        private const string TestnetWs = "wss://api.hyperliquid-testnet.xyz/ws";

        private readonly string _apiUrl;
        private readonly string _wsUrl;
        private readonly HttpClient _httpClient;
        private readonly ISigner _signer;
        private readonly ILogger<HyperliquidClient>? _logger;
        private readonly WebSocketManager _wsManager;
        private readonly InfoClient _infoClient;
        private readonly ExchangeClient _exchangeClient;
        private readonly IAsyncPolicy<HttpResponseMessage> _retryPolicy;

        private HyperliquidClient(HyperliquidClientOptions options)
        {
            _apiUrl = options.ApiUrl;
            _wsUrl = options.WsUrl;
            _signer = options.Signer ?? new Signer(options.PrivateKey!);
            _logger = options.Logger;

            _httpClient = options.HttpClient ?? new HttpClient { BaseAddress = new Uri(_apiUrl) };
            
            _retryPolicy = HttpPolicyExtensions
                .HandleTransientHttpError()
                .WaitAndRetryAsync(
                    3,
                    retryAttempt => TimeSpan.FromSeconds(Math.Pow(2, retryAttempt)),
                    onRetry: (outcome, timespan, retryCount, context) =>
                    {
                        _logger?.LogWarning($"Retry {retryCount} after {timespan} seconds");
                    });

            _wsManager = new WebSocketManager(_wsUrl, _logger);
            _infoClient = new InfoClient(this);
            _exchangeClient = new ExchangeClient(this);
        }

        public static HyperliquidClient CreateMainnet(string privateKey, ILogger<HyperliquidClient>? logger = null)
        {
            return new HyperliquidClient(new HyperliquidClientOptions
            {
                ApiUrl = MainnetApi,
                WsUrl = MainnetWs,
                PrivateKey = privateKey,
                Logger = logger
            });
        }

        public static HyperliquidClient CreateTestnet(string privateKey, ILogger<HyperliquidClient>? logger = null)
        {
            return new HyperliquidClient(new HyperliquidClientOptions
            {
                ApiUrl = TestnetApi,
                WsUrl = TestnetWs,
                PrivateKey = privateKey,
                Logger = logger
            });
        }

        public static HyperliquidClient Create(HyperliquidClientOptions options)
        {
            return new HyperliquidClient(options);
        }

        public InfoClient Info => _infoClient;
        public ExchangeClient Exchange => _exchangeClient;
        public WebSocketManager WebSocket => _wsManager;

        public async Task ConnectAsync(CancellationToken cancellationToken = default)
        {
            await _wsManager.ConnectAsync(cancellationToken);
        }

        public async Task DisconnectAsync()
        {
            await _wsManager.DisconnectAsync();
        }

        internal HttpClient HttpClient => _httpClient;
        internal ISigner Signer => _signer;
        internal IAsyncPolicy<HttpResponseMessage> RetryPolicy => _retryPolicy;
        internal ILogger<HyperliquidClient>? Logger => _logger;

        public void Dispose()
        {
            _wsManager?.Dispose();
            _httpClient?.Dispose();
        }
    }

    public class HyperliquidClientOptions
    {
        public string ApiUrl { get; set; } = MainnetApi;
        public string WsUrl { get; set; } = MainnetWs;
        public string? PrivateKey { get; set; }
        public ISigner? Signer { get; set; }
        public HttpClient? HttpClient { get; set; }
        public ILogger<HyperliquidClient>? Logger { get; set; }
    }
}