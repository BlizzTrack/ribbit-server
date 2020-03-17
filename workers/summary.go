/*
 * Copyright (c) 2020. BlizzTrack
 */

package workers

import (
	"github.com/blizztrack/ribbit-server/managers"
	"github.com/blizztrack/ribbit-server/network"
	log "github.com/sirupsen/logrus"
	"strings"
	"time"
)

type SummaryWorker struct {
	client *network.BlizzardClient
	tick   time.Duration

	running bool

	manager *managers.SummaryManager
}

var (
	summaryCommand = network.NewCommand("summary", "", "")
)

func NewSummaryWorker(client *network.BlizzardClient, tickSpeed time.Duration, manager *managers.SummaryManager) SummaryWorker {
	return SummaryWorker{
		client:  client,
		tick:    tickSpeed,
		running: false,
		manager: manager,
	}
}

func (worker SummaryWorker) Run() {
	worker.running = true
	for range time.Tick(worker.tick) {
		if !worker.running {
			return
		}

		log.Printf("checking for summary updates")

		raw, seqn, err := worker.client.Call(summaryCommand)
		if err != nil {
			log.Errorf("failed to get summary from blizzard")
			continue
		}

		if !strings.EqualFold(worker.manager.Seqn(), seqn) {
			worker.manager.Update(seqn, raw)
		}
	}
}
