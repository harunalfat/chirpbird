package usecases

import (
	"context"

	"github.com/centrifugal/centrifuge"
)

type NodeWrapper interface {
	SubscribeClientToChannel(ctx context.Context, userID string, channelID string) error
}

type NodeWrapperImpl struct {
	node *centrifuge.Node
}

func NewNodeWrapperImpl(node *centrifuge.Node) NodeWrapper {
	return &NodeWrapperImpl{
		node,
	}
}

func (n NodeWrapperImpl) SubscribeClientToChannel(ctx context.Context, userID string, channelID string) error {
	return n.node.Subscribe(userID, channelID, func(so *centrifuge.SubscribeOptions) {
		so.JoinLeave = true
		so.Presence = true
		so.Recover = true
		so.Position = true
	})
}
