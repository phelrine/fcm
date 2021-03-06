package fcm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const (
	PriorityNormal = "normal"
	PriorityHigh   = "high"

	fcnSendEndpoint = "https://fcm.googleapis.com/fcm/send"
)

type Message struct {
	To              string            `json:"to,omitempty"`
	RegistrationIds []string          `json:"registration_ids,omitempty"`
	Priority        string            `json:"priority,omitempty"`
	TimeToLive      int64             `json:"time_to_live,omitempty"`
	DryRun          bool              `json:"dry_run,omitempty"`
	CollapseKey     string            `json:"collapse_key,omitempty"`
	Notification    Notification      `json:"notification,omitempty"`
	Data            map[string]string `json:"data,omitempty"`
}

type Notification struct {
	Title       string `json:"title,omitempty"`
	Body        string `json:"body,omitempty"`
	Sound       string `json:"sound,omitempty"`
	Icon        string `json:"icon,omitempty"`
	Tag         string `json:"tag,omitempty"`
	Color       string `json:"color,omitempty"`
	ClickAction string `json:"click_action,omitempty"`
}

type Result struct {
	MessageId      string `json:"message_id,omitempty"`
	RegistrationId string `json:"registration_id,omitempty"`
	Error          string `json:"error,omitempty"`
}

type Response struct {
	MulticastId  int64    `json:"multicast_id,omitempty"`
	Success      int64    `json:"success,omitempty"`
	Failure      int64    `json:"failure,omitempty"`
	CanonicalIds int64    `json:"canonical_ids,omitempty"`
	Results      []Result `json:"results,omitempty"`
}

type Client struct {
	client httpClient
	auth   string
}

func New(apiKey string) *Client {
	return &Client{
		client: &http.Client{},
		auth:   "key=" + apiKey,
	}
}

func (c *Client) Send(msg *Message) (*Response, error) {
	body, err := json.Marshal(&msg)
	if err != nil {
		return nil, err
	}
	r := bytes.NewReader(body)

	req, err := http.NewRequest("POST", fcnSendEndpoint, r)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", c.auth)

	res, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("%s %s", res.Status, string(resBody))
	}

	if len(resBody) == 0 {
		return nil, fmt.Errorf("Response body is empty")
	}

	var out Response
	if err := json.Unmarshal(resBody, &out); err != nil {
		return nil, fmt.Errorf("%s %s", err.Error(), string(resBody))
	}
	return &out, nil
}

type httpClient interface {
	Do(*http.Request) (*http.Response, error)
}
