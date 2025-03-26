package repository

import (
	"bufio"
	"context"
	"fmt"
	"log/slog"
	"os"
	"regexp"
	"slices"
	"strings"
	"time"

	"github.com/yaacov/observer/observer"

	"github.com/bars43ru/bus2map/internal/model"
)

const patternSchedule = `(?P<route>[^;]*);(?P<transport>[^;]*);(?P<begin>[^;]+);(?P<end>[^;]+)`

type Schedule struct {
	file  string
	regex *regexp.Regexp
	data  SafeMapAtomic[model.StateNumber, []model.Schedule]
}

func NewSchedule(file string) *Schedule {
	return &Schedule{
		file:  file,
		regex: regexp.MustCompile(patternSchedule),
		data:  NewSafeMapAtomic[model.StateNumber, []model.Schedule](),
	}
}

func (s *Schedule) GetCurrent(stateNumber model.StateNumber, currentTime time.Time) (model.Schedule, error) {
	items, ok := s.data.Get(stateNumber)
	if !ok {
		return model.Schedule{}, ErrNotFound
	}
	for _, item := range items {
		if item.From.Compare(currentTime) <= 0 &&
			item.To.Compare(currentTime) >= 0 {
			return item, nil
		}
	}
	return model.Schedule{}, ErrNotFound
}

func (s *Schedule) Replace(schedules []model.Schedule) {
	data := make(map[model.StateNumber][]model.Schedule, len(schedules))
	for _, schedule := range schedules {
		v := data[schedule.StateNumber]
		v = append(v, schedule)
		data[schedule.StateNumber] = v
	}
	for _, v := range data {
		slices.SortStableFunc(v,
			func(a, b model.Schedule) int {
				return a.From.Compare(b.From)
			},
		)
	}
	s.data.Replace(data)
}

func (s *Schedule) Run(ctx context.Context) error {
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
		schedules, err := s.readFromFile(ctx)
		if err != nil {
			slog.ErrorContext(ctx, "load datasource transport", slog.Any("error", err))
			return
		}
		s.Replace(schedules)
	}

	o.AddListener(func(e interface{}) {
		slog.InfoContext(ctx, fmt.Sprintf("file modified: %v", e))
		replaceDatasource()
	})
	replaceDatasource()
	<-ctx.Done()
	return nil
}

func (s *Schedule) readFromFile(ctx context.Context) ([]model.Schedule, error) {
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

	var schedules []model.Schedule

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		record := scanner.Text()
		if strings.TrimSpace(record) == "" {
			continue
		}

		schedule, err := s.parseRawRecord(record)
		if err != nil {
			return nil, fmt.Errorf("parse raw record `%s`: %w", record, err)
		}

		schedules = append(schedules, schedule)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return schedules, nil
}

// Parse парсит строку в структуру Schedule
func (s *Schedule) parseRawRecord(record string) (model.Schedule, error) {
	match := s.regex.FindStringSubmatch(record)
	if match == nil {
		return model.Schedule{}, fmt.Errorf("raw record `%s` doesn't match the format `%s`", record, patternSchedule)
	}

	groupNames := s.regex.SubexpNames()
	schedule := model.Schedule{}

	for i, name := range groupNames {
		if i != 0 {
			switch name {
			case "route":
				schedule.Number = model.RouteNumber(match[i])
			case "transport":
				schedule.StateNumber = model.StateNumber(match[i])
			case "begin":
				t, err := s.ParseDateTime(match[i])
				if err != nil {
					return model.Schedule{}, fmt.Errorf("invalid begin datetime: %w", err)
				}
				schedule.From = t
			case "end":
				t, err := s.ParseDateTime(match[i])
				if err != nil {
					return model.Schedule{}, fmt.Errorf("invalid end datetime: %w", err)
				}
				schedule.To = t
			}
		}
	}
	return schedule, nil
}

func (s *Schedule) ParseDateTime(value string) (time.Time, error) {
	const dateTimeFormat = "02/01/2006T15:04:05ZZ07:00"
	t, err := time.Parse(dateTimeFormat, value)
	if err != nil {
		return time.Time{}, fmt.Errorf("value `%s` is not in proper format `%s`", value, dateTimeFormat)
	}
	return t, nil
}
