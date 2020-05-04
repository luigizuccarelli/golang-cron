package main

import (
	"gopkg.in/mgo.v2/bson"
	"strings"
	"time"
)

type CustomTime struct {
	time.Time
}

func (ct *CustomTime) UnmarshalJSON(b []byte) (err error) {
	s := strings.Trim(string(b), "\"")
	if s == "null" {
		ct.Time = time.Time{}
		return
	}
	ct.Time, err = time.Parse("2006-01-02 15:04", s)
	return
}

// ShcemaInterface - acts as an interface wrapper for our profile schema
// All the go microservices will using this schema
type SchemaInterface struct {
	ID         bson.ObjectId `json:"_id" bson:"_id,omitempty"`
	LastUpdate int64         `json:"lastupdate,omitempty"`
	MetaInfo   string        `json:"metainfo,omitempty"`
}

// Response schema
type Response struct {
	StatusCode string `json:"statuscode"`
	Status     string `json:"status"`
	Message    string `json:"message"`
}
