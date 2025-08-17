using System.Text.Json.Serialization;

namespace HyperliquidSDK.Types;

/// <summary>
/// Base response from Hyperliquid API
/// </summary>
public class HyperliquidResponse<T>
{
    [JsonPropertyName("status")]
    public string Status { get; set; } = string.Empty;

    [JsonPropertyName("response")]
    public T? Response { get; set; }

    [JsonPropertyName("error")]
    public string? Error { get; set; }
}

/// <summary>
/// Order request for placing new orders
/// </summary>
public class OrderRequest
{
    [JsonPropertyName("coin")]
    public string Asset { get; set; } = string.Empty;

    [JsonPropertyName("is_buy")]
    public bool IsBuy { get; set; }

    [JsonPropertyName("sz")]
    public decimal Size { get; set; }

    [JsonPropertyName("limit_px")]
    public decimal LimitPrice { get; set; }

    [JsonPropertyName("order_type")]
    public OrderType OrderType { get; set; } = new();

    [JsonPropertyName("reduce_only")]
    public bool ReduceOnly { get; set; }

    [JsonPropertyName("cloid")]
    public string? ClientOrderId { get; set; }
}

/// <summary>
/// Order type configuration
/// </summary>
public class OrderType
{
    [JsonPropertyName("limit")]
    public LimitOrderType? Limit { get; set; }

    [JsonPropertyName("trigger")]
    public TriggerOrderType? Trigger { get; set; }
}

/// <summary>
/// Limit order configuration
/// </summary>
public class LimitOrderType
{
    [JsonPropertyName("tif")]
    public string TimeInForce { get; set; } = "Gtc"; // "Gtc", "Ioc", "Alo"
}

/// <summary>
/// Trigger order configuration
/// </summary>
public class TriggerOrderType
{
    [JsonPropertyName("triggerPx")]
    public decimal TriggerPrice { get; set; }

    [JsonPropertyName("isMarket")]
    public bool IsMarket { get; set; }

    [JsonPropertyName("tpsl")]
    public string TpSl { get; set; } = string.Empty; // "tp" or "sl"
}

/// <summary>
/// Cancel request for canceling orders
/// </summary>
public class CancelRequest
{
    [JsonPropertyName("coin")]
    public string Asset { get; set; } = string.Empty;

    [JsonPropertyName("oid")]
    public long? OrderId { get; set; }

    [JsonPropertyName("cloid")]
    public string? ClientOrderId { get; set; }
}

/// <summary>
/// Modify request for modifying existing orders
/// </summary>
public class ModifyRequest
{
    [JsonPropertyName("oid")]
    public long OrderId { get; set; }

    [JsonPropertyName("order")]
    public OrderRequest Order { get; set; } = new();
}

/// <summary>
/// User's trading state
/// </summary>
public class UserState
{
    [JsonPropertyName("marginSummary")]
    public MarginSummary MarginSummary { get; set; } = new();

    [JsonPropertyName("crossMarginSummary")]
    public CrossMarginSummary CrossMarginSummary { get; set; } = new();

    [JsonPropertyName("assetPositions")]
    public List<AssetPosition> AssetPositions { get; set; } = new();

    [JsonPropertyName("withdrawable")]
    public decimal Withdrawable { get; set; }
}

/// <summary>
/// Margin summary information
/// </summary>
public class MarginSummary
{
    [JsonPropertyName("accountValue")]
    public decimal AccountValue { get; set; }

    [JsonPropertyName("totalMarginUsed")]
    public decimal TotalMarginUsed { get; set; }

    [JsonPropertyName("totalNtlPos")]
    public decimal TotalNtlPos { get; set; }

    [JsonPropertyName("totalRawUsd")]
    public decimal TotalRawUsd { get; set; }
}

/// <summary>
/// Cross margin summary information
/// </summary>
public class CrossMarginSummary
{
    [JsonPropertyName("accountValue")]
    public decimal AccountValue { get; set; }

    [JsonPropertyName("totalMarginUsed")]
    public decimal TotalMarginUsed { get; set; }

    [JsonPropertyName("totalNtlPos")]
    public decimal TotalNtlPos { get; set; }

    [JsonPropertyName("totalRawUsd")]
    public decimal TotalRawUsd { get; set; }
}

/// <summary>
/// Asset position information
/// </summary>
public class AssetPosition
{
    [JsonPropertyName("position")]
    public Position Position { get; set; } = new();

    [JsonPropertyName("type")]
    public string Type { get; set; } = string.Empty;
}

/// <summary>
/// Detailed position information
/// </summary>
public class Position
{
    [JsonPropertyName("coin")]
    public string Coin { get; set; } = string.Empty;

    [JsonPropertyName("entryPx")]
    public decimal? EntryPrice { get; set; }

    [JsonPropertyName("leverage")]
    public Leverage Leverage { get; set; } = new();

    [JsonPropertyName("liquidationPx")]
    public decimal? LiquidationPrice { get; set; }

    [JsonPropertyName("marginUsed")]
    public decimal MarginUsed { get; set; }

    [JsonPropertyName("positionValue")]
    public decimal PositionValue { get; set; }

    [JsonPropertyName("returnOnEquity")]
    public decimal ReturnOnEquity { get; set; }

    [JsonPropertyName("szi")]
    public decimal Size { get; set; }

    [JsonPropertyName("unrealizedPnl")]
    public decimal UnrealizedPnl { get; set; }
}

/// <summary>
/// Leverage configuration
/// </summary>
public class Leverage
{
    [JsonPropertyName("type")]
    public string Type { get; set; } = string.Empty; // "cross" or "isolated"

    [JsonPropertyName("value")]
    public int Value { get; set; }

    [JsonPropertyName("rawUsd")]
    public decimal? RawUsd { get; set; }
}

/// <summary>
/// Open order information
/// </summary>
public class OpenOrder
{
    [JsonPropertyName("coin")]
    public string Coin { get; set; } = string.Empty;

    [JsonPropertyName("limitPx")]
    public decimal LimitPrice { get; set; }

    [JsonPropertyName("oid")]
    public long OrderId { get; set; }

    [JsonPropertyName("side")]
    public string Side { get; set; } = string.Empty; // "A" for ask/sell, "B" for bid/buy

    [JsonPropertyName("sz")]
    public decimal Size { get; set; }

    [JsonPropertyName("timestamp")]
    public long Timestamp { get; set; }

    [JsonPropertyName("origSz")]
    public decimal OriginalSize { get; set; }

    [JsonPropertyName("cloid")]
    public string? ClientOrderId { get; set; }

    [JsonPropertyName("reduceOnly")]
    public bool ReduceOnly { get; set; }

    [JsonPropertyName("orderType")]
    public string OrderType { get; set; } = string.Empty;
}

/// <summary>
/// Filled order information
/// </summary>
public class Fill
{
    [JsonPropertyName("coin")]
    public string Coin { get; set; } = string.Empty;

    [JsonPropertyName("px")]
    public decimal Price { get; set; }

    [JsonPropertyName("sz")]
    public decimal Size { get; set; }

    [JsonPropertyName("side")]
    public string Side { get; set; } = string.Empty;

    [JsonPropertyName("time")]
    public long Time { get; set; }

    [JsonPropertyName("startPosition")]
    public decimal StartPosition { get; set; }

    [JsonPropertyName("dir")]
    public string Direction { get; set; } = string.Empty;

    [JsonPropertyName("closedPnl")]
    public decimal ClosedPnl { get; set; }

    [JsonPropertyName("hash")]
    public string Hash { get; set; } = string.Empty;

    [JsonPropertyName("oid")]
    public long OrderId { get; set; }

    [JsonPropertyName("crossed")]
    public bool Crossed { get; set; }

    [JsonPropertyName("fee")]
    public decimal Fee { get; set; }

    [JsonPropertyName("tid")]
    public long TradeId { get; set; }

    [JsonPropertyName("feeToken")]
    public string FeeToken { get; set; } = string.Empty;
}

/// <summary>
/// L2 order book data
/// </summary>
public class L2Book
{
    [JsonPropertyName("coin")]
    public string Coin { get; set; } = string.Empty;

    [JsonPropertyName("time")]
    public long Time { get; set; }

    [JsonPropertyName("levels")]
    public List<List<L2Level>> Levels { get; set; } = new();
}

/// <summary>
/// L2 order book level
/// </summary>
public class L2Level
{
    [JsonPropertyName("px")]
    public decimal Price { get; set; }

    [JsonPropertyName("sz")]
    public decimal Size { get; set; }

    [JsonPropertyName("n")]
    public int Orders { get; set; }
}

/// <summary>
/// Candlestick data
/// </summary>
public class Candle
{
    [JsonPropertyName("T")]
    public long CloseTime { get; set; }

    [JsonPropertyName("c")]
    public decimal Close { get; set; }

    [JsonPropertyName("h")]
    public decimal High { get; set; }

    [JsonPropertyName("i")]
    public string Interval { get; set; } = string.Empty;

    [JsonPropertyName("l")]
    public decimal Low { get; set; }

    [JsonPropertyName("n")]
    public int Trades { get; set; }

    [JsonPropertyName("o")]
    public decimal Open { get; set; }

    [JsonPropertyName("s")]
    public string Symbol { get; set; } = string.Empty;

    [JsonPropertyName("t")]
    public long OpenTime { get; set; }

    [JsonPropertyName("v")]
    public decimal Volume { get; set; }
}

/// <summary>
/// Asset metadata
/// </summary>
public class AssetInfo
{
    [JsonPropertyName("name")]
    public string Name { get; set; } = string.Empty;

    [JsonPropertyName("szDecimals")]
    public int SizeDecimals { get; set; }

    [JsonPropertyName("maxLeverage")]
    public int MaxLeverage { get; set; }

    [JsonPropertyName("onlyIsolated")]
    public bool OnlyIsolated { get; set; }
}

/// <summary>
/// Exchange metadata
/// </summary>
public class Meta
{
    [JsonPropertyName("universe")]
    public List<AssetInfo> Universe { get; set; } = new();
}

/// <summary>
/// Order response from exchange
/// </summary>
public class OrderResponse
{
    [JsonPropertyName("status")]
    public string Status { get; set; } = string.Empty;

    [JsonPropertyName("response")]
    public OrderResponseData Response { get; set; } = new();
}

/// <summary>
/// Order response data
/// </summary>
public class OrderResponseData
{
    [JsonPropertyName("type")]
    public string Type { get; set; } = string.Empty;

    [JsonPropertyName("data")]
    public OrderResponseDataInner Data { get; set; } = new();
}

/// <summary>
/// Inner order response data
/// </summary>
public class OrderResponseDataInner
{
    [JsonPropertyName("statuses")]
    public List<OrderStatus> Statuses { get; set; } = new();
}

/// <summary>
/// Order status information
/// </summary>
public class OrderStatus
{
    [JsonPropertyName("resting")]
    public RestingOrder? Resting { get; set; }

    [JsonPropertyName("filled")]
    public FilledOrder? Filled { get; set; }

    [JsonPropertyName("error")]
    public string? Error { get; set; }
}

/// <summary>
/// Resting order details
/// </summary>
public class RestingOrder
{
    [JsonPropertyName("oid")]
    public long OrderId { get; set; }

    [JsonPropertyName("cloid")]
    public string? ClientOrderId { get; set; }
}

/// <summary>
/// Filled order details
/// </summary>
public class FilledOrder
{
    [JsonPropertyName("totalSz")]
    public decimal TotalSize { get; set; }

    [JsonPropertyName("avgPx")]
    public decimal AveragePrice { get; set; }

    [JsonPropertyName("oid")]
    public long OrderId { get; set; }
}

/// <summary>
/// Transfer request
/// </summary>
public class TransferRequest
{
    [JsonPropertyName("destination")]
    public string Destination { get; set; } = string.Empty;

    [JsonPropertyName("amount")]
    public decimal Amount { get; set; }

    [JsonPropertyName("time")]
    public long Time { get; set; }
}

/// <summary>
/// Withdrawal request
/// </summary>
public class WithdrawRequest
{
    [JsonPropertyName("destination")]
    public string Destination { get; set; } = string.Empty;

    [JsonPropertyName("amount")]
    public decimal Amount { get; set; }

    [JsonPropertyName("time")]
    public long Time { get; set; }
}

/// <summary>
/// Builder fee information for MEV protection
/// </summary>
public class BuilderInfo
{
    [JsonPropertyName("b")]
    public string Builder { get; set; } = string.Empty;

    [JsonPropertyName("f")]
    public int Fee { get; set; } // Fee in tenths of basis points
}