/*
 * Copyright (c) 2020. BlizzTrack
 */

package main

import (
	"github.com/blizztrack/ribbit-server/managers"
	"github.com/blizztrack/ribbit-server/network"
	"github.com/blizztrack/ribbit-server/stats"
	"github.com/blizztrack/ribbit-server/storage/mongo"
	"github.com/blizztrack/ribbit-server/workers"
	log "github.com/sirupsen/logrus"

	"gopkg.in/alecthomas/kingpin.v2"
	"os"
	"time"
)

var (
	/* Mongo DB */
	mongoHost = kingpin.Flag("mongo_host", "Mongo Host").Envar("MONGO_HOST").Default("localhost:27017").String()
	mongoUser = kingpin.Flag("mongo_user", "Mongo User").Envar("MONGO_USER").Default("").String()
	mongoPass = kingpin.Flag("mongo_pass", "Mongo Pass").Envar("MONGO_PASS").Default("").String()
	mongoDB   = kingpin.Flag("mongo_db", "Mongo Database").Envar("MONGO_DB").Default("blizztrack").String()

	listen = kingpin.Flag("listen", "Listening port").Short('l').Envar("LISTEN").Default("127.0.0.1:1119").String()

	client          *network.BlizzardClient
	versionsManager *managers.VersionsManager
	summaryManager  *managers.SummaryManager
)

func init() {
	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&log.TextFormatter{
		ForceColors:               true,
		EnvironmentOverrideColors: true,
	})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)

	client = network.NewBlizzardClient("us", "")
}

func main() {
	kingpin.Parse()

	mongo.New(mongo.MongoSettings{
		Host:     *mongoHost,
		Username: *mongoUser,
		Password: *mongoPass,
		Database: *mongoDB,
	})

	summaryManager = managers.NewSummaryManager()
	versionsManager = managers.NewVersionsWorker(summaryManager, client)

	summaryWorker := workers.NewSummaryWorker(client, 1*time.Minute, summaryManager)
	go summaryWorker.Run()

	server := network.NewServer(*listen, client, onData)
	log.Error(server.Run())
}

func onData(command network.Command) (bool, []byte, string, error) {
	if command.Method == "stats" {
		s := stats.Get()
		return false, []byte(s.String()), "stats", nil
	}

	if command.Method == "summary" {
		stats.HitCommand(command.String())
		return false, []byte(summaryManager.Raw()), summaryManager.Seqn(), nil
	}

	remote, raw, seqn, err := versionsManager.Get(command.Product, command.File)
	if err != nil {
		stats.Miss()
		stats.MissCommand(command.String())
		return remote, nil, "", err
	}

	stats.Hit()
	stats.HitCommand(command.String())
	return remote, []byte(raw), seqn, nil
}
