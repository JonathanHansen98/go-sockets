package channel_manager

import (
	"errors"
	"fmt"
)

type channelManager struct {
	Channels map[string]*Channel
}

func NewChannelManager() *channelManager {
	return &channelManager{Channels: make(map[string]*Channel)}
}

func (cm *channelManager) AddChannel(name string) *Channel {
	p := NewChannel(name)

	cm.Channels[name] = p

	return p
}

func (cm *channelManager) RemoveChannel(name string) {
	fmt.Printf("Removing channel: %s", name)

	c, ok := cm.Channels[name]

	if !ok {
		fmt.Printf("Channel does not exist: %s\n", name)
		return
	}

	c.DisconnectClients()

	delete(cm.Channels, name)
}

func (cm *channelManager) GetChannel(name string) (*Channel, error) {
	c, ok := cm.Channels[name]

	if !ok {
		s := fmt.Sprintf("Channel does not exist: %s", name)
		return nil, errors.New(s)
	}

	return c, nil
}
