package test

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

)

func StartS(
	sService string,
	sPollPeriod int,
	db *m.DB) {

	if !checkInit() {
		return
	}

	logger.Infof("HMS url=%s period=%d", sService, sPollPeriod)

	t := os.Getenv("S_T")
	if len(t) == 0 {
		logger.Panicf("S_T is not set")
		return
	}

	EH(db, sService, t, mgr)

	go sSubscribe(db, sService, t)

	// poll periodically in case some event get lost
	if sPollPeriod < SMinPollPeriod {
		logger.Warnf("Setting S Poll Period to minimum value of %d", SMinPollPeriod)
		sPollPeriod = SMinPollPeriod
	}

	for range time.Tick(time.Duration(sPollPeriod) * time.Second) {
		logger.Infof("polling HMS ------------")
		EH(db, sService, t, mgr)
	}
}

