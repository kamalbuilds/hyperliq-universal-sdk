using System.Text.Json.Serialization;

namespace HyperliquidSDK.Types;

/// <summary>
/// Base WebSocket subscription request
/// </summary>
public abstract class WebSocketSubscription
{
    [JsonPropertyName("type")]
    public abstract string Type { get; }
}

/// <summary>
/// All mids subscription
/// </summary>
public class AllMidsSubscription : WebSocketSubscription
{
    [JsonPropertyName("type")]
    public override string Type => "allMids";
}

/// <summary>
/// L2 order book subscription
/// </summary>
public class L2BookSubscription : WebSocketSubscription
{
    [JsonPropertyName("type")]
    public override string Type => "l2Book";

    [JsonPropertyName("coin")]
    public string Coin { get; set; } = string.Empty;
}

/// <summary>
/// Best bid/offer subscription
/// </summary>
public class BboSubscription : WebSocketSubscription
{
    [JsonPropertyName("type")]
    public override string Type => "bbo";

    [JsonPropertyName("coin")]
    public string Coin { get; set; } = string.Empty;
}

/// <summary>
/// Trades subscription
/// </summary>
public class TradesSubscription : WebSocketSubscription
{
    [JsonPropertyName("type")]
    public override string Type => "trades";

    [JsonPropertyName("coin")]
    public string Coin { get; set; } = string.Empty;
}

/// <summary>
/// User events subscription
/// </summary>
public class UserEventsSubscription : WebSocketSubscription
{
    [JsonPropertyName("type")]
    public override string Type => "userEvents";

    [JsonPropertyName("user")]
    public string User { get; set; } = string.Empty;
}

/// <summary>
/// User fills subscription
/// </summary>
public class UserFillsSubscription : WebSocketSubscription
{
    [JsonPropertyName("type")]
    public override string Type => "userFills";

    [JsonPropertyName("user")]
    public string User { get; set; } = string.Empty;
}

/// <summary>
/// Candle subscription
/// </summary>
public class CandleSubscription : WebSocketSubscription
{
    [JsonPropertyName("type")]
    public override string Type => "candle";

    [JsonPropertyName("coin")]
    public string Coin { get; set; } = string.Empty;

    [JsonPropertyName("interval")]
    public string Interval { get; set; } = string.Empty; // "1m", "5m", "15m", "1h", "4h", "1d"
}

/// <summary>
/// Order updates subscription
/// </summary>
public class OrderUpdatesSubscription : WebSocketSubscription
{
    [JsonPropertyName("type")]
    public override string Type => "orderUpdates";

    [JsonPropertyName("user")]
    public string User { get; set; } = string.Empty;
}

/// <summary>
/// Base WebSocket message
/// </summary>
public abstract class WebSocketMessage
{
    [JsonPropertyName("channel")]
    public abstract string Channel { get; }
}

/// <summary>
/// All mids WebSocket message
/// </summary>
public class AllMidsMessage : WebSocketMessage
{
    [JsonPropertyName("channel")]
    public override string Channel => "allMids";

    [JsonPropertyName("data")]
    public AllMidsData Data { get; set; } = new();
}

/// <summary>
/// All mids data
/// </summary>
public class AllMidsData
{
    [JsonPropertyName("mids")]
    public Dictionary<string, decimal> Mids { get; set; } = new();
}

/// <summary>
/// L2 order book WebSocket message
/// </summary>
public class L2BookMessage : WebSocketMessage
{
    [JsonPropertyName("channel")]
    public override string Channel => "l2Book";

    [JsonPropertyName("data")]
    public L2Book Data { get; set; } = new();
}

/// <summary>
/// Best bid/offer WebSocket message
/// </summary>
public class BboMessage : WebSocketMessage
{
    [JsonPropertyName("channel")]
    public override string Channel => "bbo";

    [JsonPropertyName("data")]
    public BboData Data { get; set; } = new();
}

/// <summary>
/// Best bid/offer data
/// </summary>
public class BboData
{
    [JsonPropertyName("coin")]
    public string Coin { get; set; } = string.Empty;

    [JsonPropertyName("time")]
    public long Time { get; set; }

    [JsonPropertyName("bbo")]
    public List<L2Level?> Bbo { get; set; } = new(); // [bid, ask]
}

/// <summary>
/// Trades WebSocket message
/// </summary>
public class TradesMessage : WebSocketMessage
{
    [JsonPropertyName("channel")]
    public override string Channel => "trades";

    [JsonPropertyName("data")]
    public List<Trade> Data { get; set; } = new();
}

/// <summary>
/// Trade data
/// </summary>
public class Trade
{
    [JsonPropertyName("coin")]
    public string Coin { get; set; } = string.Empty;

    [JsonPropertyName("side")]
    public string Side { get; set; } = string.Empty; // "A" for ask/sell, "B" for bid/buy

    [JsonPropertyName("px")]
    public decimal Price { get; set; }

    [JsonPropertyName("sz")]
    public decimal Size { get; set; }

    [JsonPropertyName("hash")]
    public string Hash { get; set; } = string.Empty;

    [JsonPropertyName("time")]
    public long Time { get; set; }
}

/// <summary>
/// User events WebSocket message
/// </summary>
public class UserEventsMessage : WebSocketMessage
{
    [JsonPropertyName("channel")]
    public override string Channel => "user";

    [JsonPropertyName("data")]
    public UserEventsData Data { get; set; } = new();
}

/// <summary>
/// User events data
/// </summary>
public class UserEventsData
{
    [JsonPropertyName("fills")]
    public List<Fill> Fills { get; set; } = new();
}

/// <summary>
/// User fills WebSocket message
/// </summary>
public class UserFillsMessage : WebSocketMessage
{
    [JsonPropertyName("channel")]
    public override string Channel => "userFills";

    [JsonPropertyName("data")]
    public UserFillsData Data { get; set; } = new();
}

/// <summary>
/// User fills data
/// </summary>
public class UserFillsData
{
    [JsonPropertyName("user")]
    public string User { get; set; } = string.Empty;

    [JsonPropertyName("isSnapshot")]
    public bool IsSnapshot { get; set; }

    [JsonPropertyName("fills")]
    public List<Fill> Fills { get; set; } = new();
}

/// <summary>
/// Pong message from server
/// </summary>
public class PongMessage : WebSocketMessage
{
    [JsonPropertyName("channel")]
    public override string Channel => "pong";
}

/// <summary>
/// Generic WebSocket message for other types
/// </summary>
public class GenericWebSocketMessage : WebSocketMessage
{
    [JsonPropertyName("channel")]
    public override string Channel { get; set; } = string.Empty;

    [JsonPropertyName("data")]
    public object? Data { get; set; }
}