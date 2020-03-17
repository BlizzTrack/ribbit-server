/*
 * Copyright (c) 2020. BlizzTrack
 */

package mongo

import (
	"github.com/globalsign/mgo/bson"
	"time"
)

type CacheItem struct {
	ID      bson.ObjectId `bson:"_id,omitempty" json:"-"`
	Code    string        `bson:"code"`
	Seqn    string        `bson:"seqn"`
	File    string        `bson:"file"`
	Indexed time.Time     `bson:"indexed"`
	Raw     string        `bson:"raw"`
}

const (
	ribbitCol = "ribbit"
)

func GetRibbitItem(code, seqn, file string) (CacheItem, error) {
	var item CacheItem
	if err := One(ribbitCol, bson.M{"code": code, "seqn": seqn, "file": file}, &item); err != nil {
		return item, err
	}
	return item, nil
}

func SetRibbitItem(item CacheItem) error {
	item.Indexed = time.Now()
	return Insert(ribbitCol, item)
}

func LatestSummary() (CacheItem, error) {
	session, c := Collection(ribbitCol)
	defer session.Close()

	var item CacheItem

	if err := c.Find(bson.M{"code": "summary", "file": "summary"}).Sort("-seqn").One(&item); err != nil {
		return item, err
	}

	return item, nil
}
