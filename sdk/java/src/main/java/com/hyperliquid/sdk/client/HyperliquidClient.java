package com.hyperliquid.sdk.client;

import com.hyperliquid.sdk.auth.Signer;
import com.hyperliquid.sdk.types.*;
import com.hyperliquid.sdk.websocket.WebSocketManager;
import org.springframework.web.reactive.function.client.WebClient;
import reactor.core.publisher.Mono;
import reactor.core.publisher.Flux;
import reactor.util.retry.Retry;

import java.math.BigDecimal;
import java.time.Duration;
import java.util.List;
import java.util.Map;
import java.util.concurrent.CompletableFuture;

public class HyperliquidClient {
    
    private static final String MAINNET_API = "https://api.hyperliquid.xyz";
    private static final String TESTNET_API = "https://api.hyperliquid-testnet.xyz";
    private static final String MAINNET_WS = "wss://api.hyperliquid.xyz/ws";
    private static final String TESTNET_WS = "wss://api.hyperliquid-testnet.xyz/ws";
    
    private final String apiUrl;
    private final String wsUrl;
    private final WebClient webClient;
    private final Signer signer;
    private final WebSocketManager wsManager;
    private final InfoClient infoClient;
    private final ExchangeClient exchangeClient;
    
    private HyperliquidClient(Builder builder) {
        this.apiUrl = builder.apiUrl;
        this.wsUrl = builder.wsUrl;
        this.signer = builder.signer;
        
        this.webClient = WebClient.builder()
            .baseUrl(apiUrl)
            .codecs(configurer -> configurer.defaultCodecs().maxInMemorySize(10 * 1024 * 1024))
            .build();
            
        this.wsManager = new WebSocketManager(wsUrl);
        this.infoClient = new InfoClient(this);
        this.exchangeClient = new ExchangeClient(this);
    }
    
    public static Builder builder() {
        return new Builder();
    }
    
    public static HyperliquidClient mainnet(String privateKey) {
        return builder()
            .apiUrl(MAINNET_API)
            .wsUrl(MAINNET_WS)
            .privateKey(privateKey)
            .build();
    }
    
    public static HyperliquidClient testnet(String privateKey) {
        return builder()
            .apiUrl(TESTNET_API)
            .wsUrl(TESTNET_WS)
            .privateKey(privateKey)
            .build();
    }
    
    public InfoClient info() {
        return infoClient;
    }
    
    public ExchangeClient exchange() {
        return exchangeClient;
    }
    
    public WebSocketManager websocket() {
        return wsManager;
    }
    
    public CompletableFuture<Void> connect() {
        return wsManager.connect();
    }
    
    public CompletableFuture<Void> disconnect() {
        return wsManager.disconnect();
    }
    
    protected <T> Mono<T> post(String endpoint, Object request, Class<T> responseType) {
        return webClient.post()
            .uri(endpoint)
            .bodyValue(request)
            .retrieve()
            .bodyToMono(responseType)
            .retryWhen(Retry.backoff(3, Duration.ofSeconds(1))
                .maxBackoff(Duration.ofSeconds(10)));
    }
    
    protected <T> Mono<T> signedPost(String endpoint, Map<String, Object> action, Class<T> responseType) {
        Map<String, Object> signedRequest = signer.signRequest(action);
        return post(endpoint, signedRequest, responseType);
    }
    
    protected WebClient getWebClient() {
        return webClient;
    }
    
    protected Signer getSigner() {
        return signer;
    }
    
    public static class Builder {
        private String apiUrl = MAINNET_API;
        private String wsUrl = MAINNET_WS;
        private String privateKey;
        private Signer signer;
        
        public Builder apiUrl(String apiUrl) {
            this.apiUrl = apiUrl;
            return this;
        }
        
        public Builder wsUrl(String wsUrl) {
            this.wsUrl = wsUrl;
            return this;
        }
        
        public Builder privateKey(String privateKey) {
            this.privateKey = privateKey;
            this.signer = new Signer(privateKey);
            return this;
        }
        
        public Builder signer(Signer signer) {
            this.signer = signer;
            return this;
        }
        
        public HyperliquidClient build() {
            if (signer == null && privateKey != null) {
                signer = new Signer(privateKey);
            }
            return new HyperliquidClient(this);
        }
    }
}