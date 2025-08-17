package examples;

import com.hyperliquid.sdk.client.HyperliquidClient;
import com.hyperliquid.sdk.types.*;
import reactor.core.publisher.Mono;
import java.math.BigDecimal;
import java.time.Duration;
import java.util.concurrent.CompletableFuture;

public class TestSDK {
    public static void main(String[] args) {
        System.out.println("🚀 Testing Hyperliquid Java SDK");
        System.out.println("=================================");
        
        try {
            // Create testnet client
            HyperliquidClient client = HyperliquidClient.testnet("");
            
            // Test 1: Connect to WebSocket
            System.out.println("\n🔌 Test 1: WebSocket Connection...");
            CompletableFuture<Void> connectFuture = client.connect();
            connectFuture.get();
            System.out.println("✅ WebSocket connected successfully");
            
            // Test 2: Get user state (example)
            System.out.println("\n👤 Test 2: Testing Info API...");
            String testAddress = "0x0000000000000000000000000000000000000000";
            
            Mono<UserState> userStateMono = client.info()
                .getUserState(testAddress)
                .timeout(Duration.ofSeconds(5));
            
            userStateMono.subscribe(
                state -> {
                    System.out.println("✅ User state retrieved");
                    System.out.println("   Account Value: " + state.getAccountValue());
                },
                error -> {
                    System.out.println("ℹ️ User state not available (expected for test address)");
                }
            );
            
            // Test 3: Subscribe to market data
            System.out.println("\n📊 Test 3: Market Data Subscription...");
            client.websocket().subscribeToAllMids(mids -> {
                System.out.println("✅ Received market prices update");
                System.out.println("   Markets: " + mids.size());
                return null;
            });
            
            // Keep running for demo
            Thread.sleep(5000);
            
            // Disconnect
            client.disconnect().get();
            System.out.println("\n✅ Disconnected successfully");
            
            System.out.println("\n✨ Java SDK tests completed!");
            
        } catch (Exception e) {
            System.err.println("Error: " + e.getMessage());
            e.printStackTrace();
        }
    }
}
