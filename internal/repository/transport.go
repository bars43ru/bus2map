package repository

import (
	"bufio"
	"context"
	"fmt"
	"log/slog"
	"os"
	"regexp"
	"strings"

	"github.com/bars43ru/bus2map/internal/model"
	"github.com/bars43ru/bus2map/internal/model/transport_type"

	"github.com/yaacov/observer/observer"
)

// Регулярное выражение для парсинга строк: uid;state;type
const patternTransport = `(?P<uid>[^;]*);(?P<state>[^;]*);(?P<type>[^;]*)`

type Transport struct {
	file  string
	regex *regexp.Regexp
	data  SafeMapAtomic[string, model.Transport]
}

type transportGUID struct {
	GUID string
}

func NewTransport(file string) *Transport {
	return &Transport{
		file:  file,
		regex: regexp.MustCompile(patternTransport),
		data:  NewSafeMapAtomic[string, model.Transport](),
	}
}

func (s *Transport) Get(uuid string) (model.Transport, error) {
	t, ok := s.data.Get(uuid)
	if !ok {
		return t, ErrNotFound
	}
	return t, nil
}

func (s *Transport) Run(ctx context.Context) error {
	o := observer.Observer{}
	err := o.Watch([]string{s.file})
	if err != nil {
		return fmt.Errorf("subscribe watch %s: %w", s.file, err)
	}
	defer func(o *observer.Observer) {
		err := o.Close()
		if err != nil {
			slog.ErrorContext(ctx, "close file change watch", slog.Any("error", err))
		}
	}(&o)

	replaceDatasource := func() {
		transports, err := s.readFromFile(ctx)
		if err != nil {
			slog.ErrorContext(ctx, "load datasource transport", slog.Any("error", err))
			return
		}
		s.replace(transports)
	}

	o.AddListener(func(e interface{}) {
		slog.InfoContext(ctx, fmt.Sprintf("file modified: %v", e))
		replaceDatasource()
	})
	replaceDatasource()
	<-ctx.Done()
	return nil
}

func (s *Transport) replace(transports []model.Transport) {
	data := make(map[string]model.Transport, len(transports))
	for _, t := range transports {
		data[t.GUID] = t
	}
	s.data.Replace(data)
}

func (s *Transport) readFromFile(ctx context.Context) ([]model.Transport, error) {
	file, err := os.Open(s.file)
	if err != nil {
		return nil, err
	}
	defer func(*os.File) {
		if err := file.Close(); err != nil {
			slog.ErrorContext(ctx, "close file",
				slog.String("file", s.file),
				slog.Any("error", err),
			)
		}
	}(file)

	var transports []model.Transport

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		record := scanner.Text()
		if strings.TrimSpace(record) == "" {
			continue
		}

		transport, err := s.parseRawRecord(record)
		if err != nil {
			return nil, fmt.Errorf("parse raw record `%s`: %w", record, err)
		}

		transports = append(transports, transport)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return transports, nil
}

// Parse парсит строку в структуру Transport
func (s *Transport) parseRawRecord(record string) (model.Transport, error) {
	match := s.regex.FindStringSubmatch(record)
	if match == nil {
		return model.Transport{}, fmt.Errorf("raw record `%s` doesn't match the format `%s`", record, patternTransport)
	}

	groupNames := s.regex.SubexpNames()
	result := model.Transport{}

	for i, name := range groupNames {
		if i != 0 {
			switch name {
			case "uid":
				result.GUID = match[i]
			case "state":
				result.StateNumber = model.StateNumber(match[i])
			case "type":
				_type, err := transport_type.ParseType(match[i])
				if err != nil {
					return model.Transport{}, fmt.Errorf("unexpected value `%s` for `transport_type`: %w", match[i], err)
				}
				result.Type = _type
			}
		}
	}
	return result, nil
}
