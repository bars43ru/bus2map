// Package yandex содержит клиент для отправки данных о местоположении транспортных средств в Яндекс.Транспорт.
// Реализует протокол обмена данными с сервисом Яндекс.Транспорт для отображения общественного транспорта на карте.
package yandex

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// Client определяет интерфейс для отправки данных о местоположении транспортных средств в Яндекс.Транспорт.
// Используется для передачи информации о маршрутах и текущем положении транспорта.
type Client interface {
	// Send отправляет данные о местоположении транспортных средств в Яндекс.Транспорт.
	// Принимает контекст и массив треков с информацией о местоположении.
	Send(ctx context.Context, t []Track) error
}

// HttpClient реализует интерфейс Client для отправки данных через HTTP.
// Содержит идентификатор клиента (clid) и URL сервера Яндекс.Транспорт.
type HttpClient struct {
	clid     string // Идентификатор клиента в системе Яндекс.Транспорт
	url      string // URL сервера Яндекс.Транспорт
	compress bool
}

func New(clid string, url string, compress bool) *HttpClient {
	return &HttpClient{
		clid:     clid,
		url:      url,
		compress: compress,
	}
}

// Send отправляет данные о местоположении транспортных средств в Яндекс.Транспорт.
// Преобразует данные в XML-формат и отправляет их на сервер через HTTP POST запрос.
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

// sendRequest выполняет HTTP POST запрос к серверу Яндекс.Транспорт.
// Отправляет XML-данные в формате application/x-www-form-urlencoded.
// Возвращает ответ сервера или ошибку в случае неудачи.
func (c *HttpClient) sendRequest(ctx context.Context, xml []byte) ([]byte, error) {
	var (
		request *http.Request
		err     error
	)
	if c.compress {
		request, err = c.makeRequestWithCompress(ctx, xml)
		if err != nil {
			return nil, fmt.Errorf("make request with compress: %w", err)
		}
	} else {
		request, err = c.makeRequestWithoutCompress(ctx, xml)
		if err != nil {
			return nil, fmt.Errorf("make request without compress: %w", err)
		}
	}

	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, fmt.Errorf("send: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("code status response %d", resp.StatusCode)
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response: %w", err)
	}
	return b, nil
}

func (c *HttpClient) makeRequestWithoutCompress(ctx context.Context, xml []byte) (*http.Request, error) {
	data := url.Values{}
	data.Set("compressed", "0")
	data.Set("data", string(xml))
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.url, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, fmt.Errorf("prepare request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	return req, nil
}

func (c *HttpClient) makeRequestWithCompress(ctx context.Context, xml []byte) (*http.Request, error) {
	data := url.Values{}
	data.Set("compressed", "1")
	b, err := c.compressBytes(xml)
	if err != nil {
		return nil, fmt.Errorf("compress xml: %w", err)
	}
	data.Set("data", string(b))
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.url, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, fmt.Errorf("prepare request: %w", err)
	}
	req.Header.Set("Content-Type", "multipart/form-data")
	return req, nil
}

func (c *HttpClient) compressBytes(data []byte) ([]byte, error) {
	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)
	_, err := gz.Write(data)
	if err != nil {
		return nil, err
	}
	// Закрываем gzip, чтобы flush-нуть все данные в буфер
	if err := gz.Close(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
