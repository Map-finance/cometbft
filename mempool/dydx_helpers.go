package mempool

import (
	"github.com/cometbft/cometbft/libs/log"
	"github.com/cometbft/cometbft/proto/dydxcometbft/clob"
	"github.com/cometbft/cometbft/types"
	cosmostx "github.com/cosmos/cosmos-sdk/types/tx"
)

// IsShortTermClobOrderTransaction returns true if the provided `tx` is a
// Cosmos transaction containing a short-term `MsgPlaceOrder` or
// short-term `MsgCancelOrder` or `MsgBatchCancel` message.
// 支持 dydxprotocol 和 h2x 两种 proto 包名
func IsShortTermClobOrderTransaction(
	tx types.Tx,
	mempoolLogger log.Logger,
) bool {
	cosmosTx := &cosmostx.Tx{}
	err := cosmosTx.Unmarshal(tx)
	if err != nil {
		mempoolLogger.Error("isClobOrderTransaction error. Invalid Cosmos Transaction.")
		return false
	}
	if cosmosTx.Body != nil && len(cosmosTx.Body.Messages) == 1 {
		bytes := cosmosTx.Body.Messages[0].Value
		typeUrl := cosmosTx.Body.Messages[0].TypeUrl

		// 检查 MsgPlaceOrder（支持 dydxprotocol 和 h2x 两种包名）
		if typeUrl == "/dydxprotocol.clob.MsgPlaceOrder" || typeUrl == "/h2x.clob.MsgPlaceOrder" {
			msgPlaceOrder := &clob.MsgPlaceOrder{}
			err := msgPlaceOrder.Unmarshal(bytes)
			// In the case of an unmarshalling error, panic.
			// Chances are, the protos are out of sync with the dydx v4 repo.
			if err != nil {
				panic(
					"Failed to unmarshal MsgPlaceOrder from Cosmos transaction in CometBFT mempool.",
				)
			}
			return msgPlaceOrder.Order.OrderId.IsShortTermOrder()
		}

		// 检查 MsgCancelOrder（支持 dydxprotocol 和 h2x 两种包名）
		if typeUrl == "/dydxprotocol.clob.MsgCancelOrder" || typeUrl == "/h2x.clob.MsgCancelOrder" {
			msgCancelOrder := &clob.MsgCancelOrder{}
			err := msgCancelOrder.Unmarshal(bytes)
			// In the case of an unmarshalling error, panic.
			// Chances are, the protos are out of sync with the dydx v4 repo.
			if err != nil {
				panic("Failed to unmarshal MsgCancelOrder from Cosmos transaction.")
			}
			return msgCancelOrder.OrderId.IsShortTermOrder()
		}

		// 检查 MsgBatchCancel（支持 dydxprotocol 和 h2x 两种包名）
		if typeUrl == "/dydxprotocol.clob.MsgBatchCancel" || typeUrl == "/h2x.clob.MsgBatchCancel" {
			// MsgBatchCancel only processes short term order cancellations as of right now.
			return true
		}
	}

	return false
}
