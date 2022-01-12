package usecases

import (
	"context"

	"github.com/centrifugal/centrifuge"
	"github.com/google/uuid"
)

type NodeWrapper interface {
	SubscribeClientToChannel(ctx context.Context, userID uuid.UUID, channelID uuid.UUID) error
}

type NodeWrapperImpl struct {
	node *centrifuge.Node
}

func NewNodeWrapperImpl(node *centrifuge.Node) NodeWrapper {
	return &NodeWrapperImpl{
		node,
	}
}

func (n NodeWrapperImpl) SubscribeClientToChannel(ctx context.Context, userID uuid.UUID, channelID uuid.UUID) error {
	return n.node.Subscribe(userID.String(), channelID.String(), func(so *centrifuge.SubscribeOptions) {
		so.JoinLeave = true
		so.Presence = true
		so.Recover = true
		so.Position = true
	})
}
