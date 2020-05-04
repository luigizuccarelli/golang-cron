package main

import (
	"crypto/tls"
	"fmt"
	"github.com/microlib/simple"
	"gopkg.in/robfig/cron.v2"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

var (
	logger *simple.Logger
)

// This microservice will be setup for the various q's
// ORDER_Q - get quotes (batch up to 100) for each order in this q (list)
// PENDING_Q - get quotes (batch up to 100) for each item in this q

func main() {

	logger = &simple.Logger{Level: os.Getenv("LOG_LEVEL")}

	err := ValidateEnvars(logger)
	if err != nil {
		os.Exit(1)
	}

	cr := cron.New()
	cr.AddFunc(os.Getenv("CRON"),
		func() {
			updatePrices(logger)
		})
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)

	go func() {
		<-c
		cleanup(cr)
		os.Exit(1)
	}()

	cr.Start()

	for {
		logger.Debug(fmt.Sprintf("NOP sleeping for %s seconds\n", os.Getenv("SLEEP")))
		s, _ := strconv.Atoi(os.Getenv("SLEEP"))
		time.Sleep(time.Duration(s) * time.Second)
	}
}

func getData(logger *simple.Logger) error {
	// set up http object
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	httpClient := &http.Client{Transport: tr}

	// get the list from redis
	list := redis.Get(os.Getenv("OBJECT_Q"))

	// replace the {symbols} object with list
	// replace the {token} object with os.Getenv("PROVIDER_TOKEN")

	req, _ := http.NewRequest("GET", os.Getenv("URL"), nil)
	//req.Header.Set("X-Api-Key", os.Getenv("APIKEY"))
	req.Header.Set("Content-Type", "application/json")
	resp, err := httpClient.Do(req)
	if err != nil {
		logger.Error(fmt.Sprintf("Http request %v", err))
		return err
	}

	defer resp.Body.Close()
	body, e := ioutil.ReadAll(resp.Body)
	if e != nil {
		logger.Error(fmt.Sprintf("Cron service updatePrices %v", e))
		return e
	}
	logger.Debug(fmt.Sprintf("Response from server %s", string(body)))

	// now for each stock save the resultant data redis
	// redis.Set(stock,result,0)

	return nil
}

func cleanup(c *cron.Cron) {
	logger.Warn("Cleanup resources")
	logger.Info("Terminating")
	c.Stop()
}
