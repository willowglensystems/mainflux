// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

package producer

import (
	"time"

	"github.com/go-redis/redis"
	"github.com/mainflux/mainflux/bootstrap"
)

const (
	streamID  = "mainflux.bootstrap"
	streamLen = 1000
)

var _ bootstrap.Service = (*eventStore)(nil)

type eventStore struct {
	svc    bootstrap.Service
	client *redis.Client
}

// NewEventStoreMiddleware returns wrapper around bootstrap service that sends
// events to event store.
func NewEventStoreMiddleware(svc bootstrap.Service, client *redis.Client) bootstrap.Service {
	return eventStore{
		svc:    svc,
		client: client,
	}
}

func (es eventStore) Add(token string, cfg bootstrap.Config) (bootstrap.Config, error) {
	saved, err := es.svc.Add(token, cfg)
	if err != nil {
		return saved, err
	}

	var channels []string
	for _, ch := range saved.MFChannels {
		channels = append(channels, ch.ID)
	}

	ev := createConfigEvent{
		mfThing:    saved.MFThing,
		owner:      saved.Owner,
		name:       saved.Name,
		mfChannels: channels,
		externalID: saved.ExternalID,
		content:    saved.Content,
		timestamp:  time.Now(),
	}

	es.add(ev)

	return saved, err
}

func (es eventStore) View(token, id string) (bootstrap.Config, error) {
	return es.svc.View(token, id)
}

func (es eventStore) Update(token string, cfg bootstrap.Config) error {
	if err := es.svc.Update(token, cfg); err != nil {
		return err
	}

	ev := updateConfigEvent{
		mfThing:   cfg.MFThing,
		name:      cfg.Name,
		content:   cfg.Content,
		timestamp: time.Now(),
	}

	es.add(ev)

	return nil
}

func (es eventStore) UpdateCert(token, thingKey, clientCert, clientKey, caCert string) error {
	return es.svc.UpdateCert(token, thingKey, clientCert, clientKey, caCert)
}

func (es eventStore) UpdateConnections(token, id string, connections []string) error {
	if err := es.svc.UpdateConnections(token, id, connections); err != nil {
		return err
	}

	ev := updateConnectionsEvent{
		mfThing:    id,
		mfChannels: connections,
		timestamp:  time.Now(),
	}

	es.add(ev)

	return nil
}

func (es eventStore) List(token string, filter bootstrap.Filter, offset, limit uint64) (bootstrap.ConfigsPage, error) {
	return es.svc.List(token, filter, offset, limit)
}

func (es eventStore) Remove(token, id string) error {
	if err := es.svc.Remove(token, id); err != nil {
		return err
	}

	ev := removeConfigEvent{
		mfThing:   id,
		timestamp: time.Now(),
	}

	es.add(ev)

	return nil
}

func (es eventStore) Bootstrap(externalKey, externalID string, secure bool) (bootstrap.Config, error) {
	cfg, err := es.svc.Bootstrap(externalKey, externalID, secure)

	ev := bootstrapEvent{
		externalID: externalID,
		timestamp:  time.Now(),
		success:    true,
	}

	if err != nil {
		ev.success = false
	}

	es.add(ev)

	return cfg, err
}

func (es eventStore) ChangeState(token, id string, state bootstrap.State) error {
	if err := es.svc.ChangeState(token, id, state); err != nil {
		return err
	}

	ev := changeStateEvent{
		mfThing:   id,
		state:     state,
		timestamp: time.Now(),
	}

	es.add(ev)

	return nil
}

func (es eventStore) RemoveConfigHandler(id string) error {
	return es.svc.RemoveConfigHandler(id)
}

func (es eventStore) RemoveChannelHandler(id string) error {
	return es.svc.RemoveChannelHandler(id)
}

func (es eventStore) UpdateChannelHandler(channel bootstrap.Channel) error {
	return es.UpdateChannelHandler(channel)
}

func (es eventStore) DisconnectThingHandler(channelID, thingID string) error {
	return es.svc.DisconnectThingHandler(channelID, thingID)
}

func (es eventStore) add(ev event) error {
	record := &redis.XAddArgs{
		Stream:       streamID,
		MaxLenApprox: streamLen,
		Values:       ev.encode(),
	}

	return es.client.XAdd(record).Err()
}
