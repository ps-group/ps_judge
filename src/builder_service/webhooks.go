package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/sirupsen/logrus"
)

type WebhookData struct {
	UUID    string `json:"uuid"`
	Succeed bool   `json:"succeed"`
}

type WebhookRequest struct {
	Data WebhookData
	URL  string
}

type WebhookService interface {
	postBuildFinished(url string, uuid string, succeed bool)
}

type webhookServiceImpl struct {
	client  *http.Client
	input   chan WebhookRequest
	closing chan bool
}

func newWebhookService() *webhookServiceImpl {
	service := new(webhookServiceImpl)
	service.client = new(http.Client)
	service.input = make(chan WebhookRequest)
	service.closing = make(chan bool)
	go service.listenWebhooks()
	return service
}

func (ws *webhookServiceImpl) listenWebhooks() {
	for {
		select {
		case request := <-ws.input:
			err := ws.postBuildFinishedImpl(request)
			if err != nil {
				logrus.Errorf("webhook error: %v", err)
			}
		case <-ws.closing:
			ws.closing <- true
			return
		}
	}
}

func (ws *webhookServiceImpl) close() {
	ws.closing <- true
	<-ws.closing
}

func (ws *webhookServiceImpl) postBuildFinished(url string, uuid string, succeed bool) {
	webhookReq := WebhookRequest{
		URL: url,
		Data: WebhookData{
			Succeed: succeed,
			UUID:    uuid,
		},
	}
	ws.input <- webhookReq
}

func (ws *webhookServiceImpl) postBuildFinishedImpl(webhookReq WebhookRequest) error {
	jsonBytes, err := json.Marshal(webhookReq.Data)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", webhookReq.URL, bytes.NewReader(jsonBytes))
	if err != nil {
		return err
	}
	res, err := ws.client.Do(req)
	if err != nil {
		return err
	}
	if res.StatusCode != http.StatusOK {
		return errors.New("request failed: " + res.Status)
	}
	return nil
}
