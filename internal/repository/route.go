package repository

import (
	"bufio"
	"context"
	"fmt"
	"log/slog"
	"os"
	"regexp"
	"strings"

	"github.com/yaacov/observer/observer"

	"github.com/bars43ru/bus2map/internal/model"
)

const patternRoute = `(?P<internal>[^;]*);(?P<yandex>[^;]*);(?P<2gis>[^;]*)`

type Route struct {
	file  string
	regex *regexp.Regexp
	data  SafeMapAtomic[model.RouteNumber, model.Route]
}

func NewRoute(file string) *Route {
	return &Route{
		file:  file,
		regex: regexp.MustCompile(patternRoute),
		data:  NewSafeMapAtomic[model.RouteNumber, model.Route](),
	}
}

func (s *Route) GetRoute(number model.RouteNumber) (model.Route, error) {
	r, ok := s.data.Get(number)
	if !ok {
		return r, ErrNotFound
	}
	return r, nil
}

func (s *Route) Run(ctx context.Context) error {
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
		routes, err := s.readFromFile(ctx)
		if err != nil {
			slog.ErrorContext(ctx, "load datasource route", slog.Any("error", err))
			return
		}
		s.replace(routes)
	}

	o.AddListener(func(e interface{}) {
		slog.InfoContext(ctx, fmt.Sprintf("file modified: %v", e))
		replaceDatasource()
	})
	replaceDatasource()
	<-ctx.Done()
	return nil
}

func (s *Route) replace(routes []model.Route) {
	data := make(map[model.RouteNumber]model.Route, len(routes))
	for _, route := range routes {
		data[route.Number] = route
	}
	s.data.Replace(data)
}

func (s *Route) readFromFile(ctx context.Context) ([]model.Route, error) {
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

	var routes []model.Route

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		record := scanner.Text()
		if strings.TrimSpace(record) == "" {
			continue
		}

		route, err := s.parseRawRecord(record)
		if err != nil {
			return nil, fmt.Errorf("parse raw record `%s`: %w", record, err)
		}

		routes = append(routes, route)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return routes, nil
}

// Parse парсит строку в структуру Route
func (s *Route) parseRawRecord(record string) (model.Route, error) {
	match := s.regex.FindStringSubmatch(record)
	if match == nil {
		return model.Route{}, fmt.Errorf("raw record `%s` doesn't match the format `%s`", record, patternRoute)
	}

	groupNames := s.regex.SubexpNames()
	result := model.Route{}

	for i, name := range groupNames {
		if i != 0 {
			switch name {
			case "internal":
				result.Number = model.RouteNumber(match[i])
			case "yandex":
				result.YandexNumber = match[i]
			case "2gis":
				result.TwoGISNumber = match[i]
			}
		}
	}
	return result, nil
}
