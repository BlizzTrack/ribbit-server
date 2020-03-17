/*
 * Copyright (c) 2020. BlizzTrack
 */

package managers

import (
	"errors"
	"fmt"
	"github.com/blizztrack/ribbit-server/network"
	"github.com/blizztrack/ribbit-server/storage/mongo"
	"time"
)

type VersionsManager struct {
	manager *SummaryManager
	client  *network.BlizzardClient
}

func NewVersionsWorker(manager *SummaryManager, client *network.BlizzardClient) *VersionsManager {
	return &VersionsManager{
		manager: manager,
		client:  client,
	}
}

func (vm *VersionsManager) Get(code, file string) (bool, string, string, error) {
	 item := vm.manager.Get(code, file)
	 if item.Product == "" || item.Product != code {
	 	return false, "", "", errors.New("failed to look up product")
	 }

	 current, _ := mongo.GetRibbitItem(code, item.Seqn, file)
	 if current.Seqn != item.Seqn {
	 	raw, seqn, err := vm.client.Call(network.NewCommand("products", code, file))
	 	if err != nil {
	 		return false, "", "", err
		}

		if seqn != "" {
			err = mongo.SetRibbitItem(mongo.CacheItem{
				Code:    code,
				Seqn:    seqn,
				File:    file,
				Indexed: time.Time{},
				Raw:     raw,
			})
		} else {
			return true, "", "", fmt.Errorf("%s %s was empty", code, file)
		}
		return true, raw, seqn, nil
	 }

	 return false, current.Raw, current.Seqn, nil
}
