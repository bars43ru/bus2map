package yandex

import (
	"context"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type Client interface {
	Send(ctx context.Context, t []Track) error
}

type HttpClient struct {
	clid string
	url  string
}

func New(clid string, url string) *HttpClient {
	return &HttpClient{
		clid: clid,
		url:  url,
	}
}

func (c *HttpClient) Send(ctx context.Context, t []Track) error {
	v := tracks{Clid: c.clid, Tracks: t}
	xmlReq, err := xml.Marshal(v)
	if err != nil {
		return fmt.Errorf("marshal to xml: %w", err)
	}
	_, err = c.sendRequest(ctx, xmlReq)
	if err != nil {
		return fmt.Errorf("sending yandex: %w", err)
	}
	return nil
}

func (c *HttpClient) sendRequest(ctx context.Context, xml []byte) ([]byte, error) {
	client := &http.Client{}
	data := url.Values{}
	data.Set("compressed", "0")
	data.Set("data", string(xml))

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.url, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, fmt.Errorf("prepare request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("send: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		return nil, errors.New("status response")
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response: %w", err)
	}
	return b, nil
}
