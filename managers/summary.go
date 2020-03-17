/*
 * Copyright (c) 2020. BlizzTrack
 */

package managers

import (
	"github.com/blizztrack/ribbit-server/storage/mongo"
	log "github.com/sirupsen/logrus"
	"strings"

	"github.com/blizztrack/ribbit-go"
)

type SummaryManager struct {
	current string
	items   []ribbit.SummaryItem
	raw     string
}

func NewSummaryManager() *SummaryManager {
	sm := new(SummaryManager)
	sm.loadLatest()
	return sm
}

func (sm *SummaryManager) Seqn() string {
	return sm.current
}

func (sm *SummaryManager) Raw() string {
	return sm.raw
}

func (sm *SummaryManager) Get(code, file string) ribbit.SummaryItem {
	if file == "versions" {
		file = ""
	}

	// this is cause summary cdn = cdns in the file look up table
	if file == "cdns" {
		file = "cdn"
	}

	for _, item := range sm.items {
		if strings.EqualFold(item.Product, code) && strings.EqualFold(item.Flags, file) {
			return item
		}
	}

	return ribbit.SummaryItem{}
}

func (sm *SummaryManager) Update(seqn, raw string) error {
	oldSeqn := sm.current

	sm.current = seqn
	sm.items = ToSummaryList(raw)
	sm.raw = raw

	err := mongo.SetRibbitItem(mongo.CacheItem{
		Code: "summary",
		File: "summary",
		Seqn: seqn,
		Raw:  raw,
	})

	if err != nil {
		return err
	}

	log.Printf("[Local] Updated Summary file to %s from %s", seqn, oldSeqn)

	return nil
}

func (sm *SummaryManager) loadLatest() {
	data, err := mongo.LatestSummary()
	if err != nil {
		return
	}

	sm.current = data.Seqn
	sm.items = ToSummaryList(data.Raw)
	sm.raw = data.Raw

	log.Infof("Loaded latest cached summary file: %s", sm.current)
}
