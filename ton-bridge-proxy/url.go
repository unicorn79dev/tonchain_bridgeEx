package main

func getTargetURL(bridgeType string) string {
	switch bridgeType {
	case "telegram-wallet":
		return "https://bridge.ton.space/bridge"
	case "tonkeeper":
		return "https://bridge.tonapi.io/bridge"
	case "mytonwallet":
		return "https://tonconnectbridge.mytonwallet.org/bridge"
	case "tonhub":
		return "https://connect.tonhubapi.com/tonconnect"
	case "dewallet":
		return "https://bridge.dewallet.pro/bridge"
	case "bitgetTonWallet":
		return "https://ton-connect-bridge.bgwapi.io/bridge"
	case "safepalwallet":
		return "https://ton-bridge.safepal.com/tonbridge/v1/bridge"
	case "okxTonWallet":
		return "https://www.okx.com/tonbridge/discover/rpc/bridge"
	case "hot":
		return "https://sse-bridge.hot-labs.org"
	case "bybitTonWallet":
		return "https://api-node.bybit.com/spot/api/web3/bridge/ton/bridge"
	case "GateWallet":
		return "https://dapp.gateio.services/tonbridge_api/bridge/v1"
	case "binanceWeb3TonWallet":
		return "https://wallet.binance.com/tonbridge/bridge"
	default:
		return ""
	}
}
