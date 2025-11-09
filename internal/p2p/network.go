package p2p

import (
	"context"
	"encoding/json"
	"log"
	"sync"
	"time"

	"MeshComm/internal/types"

	"github.com/google/uuid"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p/core/host"
)

type MeshNetwork struct {
	ctx       context.Context
	host      host.Host
	pubsub    *pubsub.PubSub
	topics    map[types.MessageCategory]*pubsub.Topic
	subs      map[types.MessageCategory]*pubsub.Subscription
	messages  []types.Message
	user      types.User
	msgMutex  sync.RWMutex
	callbacks []func(types.Message)
}

func NewMeshNetwork(ctx context.Context, host host.Host, nickname string) (*MeshNetwork, error) {
	// Create a new PubSub service using GossipSub
	ps, err := pubsub.NewGossipSub(ctx, host)
	if err != nil {
		return nil, err
	}

	mn := &MeshNetwork{
		ctx:      ctx,
		host:     host,
		pubsub:   ps,
		topics:   make(map[types.MessageCategory]*pubsub.Topic),
		subs:     make(map[types.MessageCategory]*pubsub.Subscription),
		messages: make([]types.Message, 0),
		user: types.User{
			ID:       uuid.New().String(),
			Nickname: nickname,
			NodeID:   host.ID().String(),
		},
		callbacks: make([]func(types.Message), 0),
	}

	// Join all message categories
	categories := []types.MessageCategory{
		types.Emergency,
		types.General,
		types.Help,
	}

	for _, category := range categories {
		if err := mn.joinTopic(category); err != nil {
			return nil, err
		}
	}

	return mn, nil
}

func (mn *MeshNetwork) joinTopic(category types.MessageCategory) error {
	topic, err := mn.pubsub.Join(string(category))
	if err != nil {
		return err
	}
	mn.topics[category] = topic

	sub, err := topic.Subscribe()
	if err != nil {
		return err
	}
	mn.subs[category] = sub

	// Start message handler for this category
	go mn.handleMessages(category, sub)
	return nil
}

func (mn *MeshNetwork) handleMessages(category types.MessageCategory, sub *pubsub.Subscription) {
	for {
		msg, err := sub.Next(mn.ctx)
		if err != nil {
			log.Printf("Error receiving message in %s: %v", category, err)
			return
		}

		// Skip messages from ourselves
		if msg.ReceivedFrom == mn.host.ID() {
			continue
		}

		var message types.Message
		if err := json.Unmarshal(msg.Data, &message); err != nil {
			log.Printf("Error unmarshaling message: %v", err)
			continue
		}

		mn.msgMutex.Lock()
		mn.messages = append(mn.messages, message)
		mn.msgMutex.Unlock()

		// Notify callbacks
		for _, cb := range mn.callbacks {
			cb(message)
		}
	}
}

func (mn *MeshNetwork) PublishMessage(content string, category types.MessageCategory) error {
	message := types.Message{
		ID:        uuid.New().String(),
		Content:   content,
		Category:  category,
		Sender:    mn.user,
		Timestamp: time.Now(),
	}

	data, err := json.Marshal(message)
	if err != nil {
		return err
	}

	topic, ok := mn.topics[category]
	if !ok {
		return nil
	}

	return topic.Publish(mn.ctx, data)
}

func (mn *MeshNetwork) GetMessages() []types.Message {
	mn.msgMutex.RLock()
	defer mn.msgMutex.RUnlock()

	// Return a copy of messages
	messages := make([]types.Message, len(mn.messages))
	copy(messages, mn.messages)
	return messages
}

func (mn *MeshNetwork) OnMessage(callback func(types.Message)) {
	mn.callbacks = append(mn.callbacks, callback)
}
