package wialonips

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"iter"
	"log/slog"
	"regexp"
	"strconv"
	"time"

	"github.com/bars43ru/bus2map/pkg/xslog"
)

type Parser struct {
	reader  *bufio.Reader
	msgLExp *regexp.Regexp
	msgDExp *regexp.Regexp
	msgL    messageL
}

func NewParse(reader io.Reader) (*Parser, error) {
	parse := &Parser{
		reader:  bufio.NewReader(reader),
		msgLExp: regexp.MustCompile(patternL),
		msgDExp: regexp.MustCompile(patternD),
	}
	if err := parse.readHeader(); err != nil {
		return nil, err
	}
	return parse, nil
}

func (p *Parser) uid() string {
	return p.msgL.UID
}

func (p *Parser) nextMessage() (messageD, error) {
	for {
		s, err := p.reader.ReadString(delimiter)
		if s == "" && err != nil {
			return messageD{}, fmt.Errorf("read data: %w", err)
		}

		msgD, err := p.parseD(s)
		if err != nil && errors.Is(err, ErrFormat) {
			continue
		}
		if err != nil {
			return messageD{}, fmt.Errorf("parse data message: %w", err)
		}

		// Причина появления этого условия см. https://github.com/bars43ru/gps2Yandex/issues/11
		if int(msgD.Latitude) == 90 && int(msgD.Longitude) == 0 {
			continue
		}
		return msgD, nil
	}
}

func (p *Parser) readHeader() error {
	s, err := p.reader.ReadString(delimiter)
	if s == "" && err != nil {
		return fmt.Errorf("read message L: %w", err)
	}
	p.msgL, err = p.parseL(s)
	if err != nil {
		return fmt.Errorf("parse data header: %w", err)
	}
	return nil
}

func (p *Parser) parseL(s string) (messageL, error) {
	match := p.msgLExp.FindStringSubmatch(s)
	if match == nil {
		return messageL{}, fmt.Errorf("parse message L: %w", ErrFormat)
	}
	return messageL{UID: match[p.msgLExp.SubexpIndex(uidField)]}, nil
}

func (p *Parser) parseD(s string) (messageD, error) {
	match := p.msgDExp.FindStringSubmatch(s)
	if match == nil {
		return messageD{}, fmt.Errorf("parse message D: %w", ErrFormat)
	}

	// Не проверяем на ошибки, т.к. patternD и patternL не допускает не корректный формат входящих данных
	dt, _ := time.Parse(layoutTime, match[p.msgDExp.SubexpIndex(dateField)]+match[p.msgDExp.SubexpIndex(timeField)])
	lat1, _ := strconv.ParseFloat(match[p.msgDExp.SubexpIndex(lat1Field)], 64)
	lon1, _ := strconv.ParseFloat(match[p.msgDExp.SubexpIndex(lon1Field)], 64)
	speed, _ := strconv.ParseUint(match[p.msgDExp.SubexpIndex(speedField)], 10, 32)
	course, _ := strconv.ParseUint(match[p.msgDExp.SubexpIndex(courseField)], 10, 8)

	return messageD{
		Time:      dt,
		Latitude:  Coordinate(lat1),
		Longitude: Coordinate(lon1),
		Speed:     uint(speed),
		Course:    uint8(course),
	}, nil
}

func (p *Parser) Points(ctx context.Context) iter.Seq2[int, Point] {
	return func(yield func(int, Point) bool) {
		index := -1
		for {
			s, err := p.reader.ReadString(delimiter)
			if err != nil {
				if errors.Is(err, io.EOF) {
					return
				}
				slog.ErrorContext(ctx, "data read",
					xslog.Error(err),
					slog.Any("uid", p.uid()),
					slog.String("data", s),
				)
				return
			}
			if s == "" {
				slog.ErrorContext(ctx, "read empty data", slog.Any("uid", p.uid()))
				return
			}

			msgD, err := p.parseD(s)
			if err != nil && errors.Is(err, ErrFormat) {
				continue
			}
			if err != nil {
				slog.ErrorContext(ctx, "parsing of data stream stopped due to unexpected incoming data",
					xslog.Error(err),
					slog.Any("uid", p.uid()),
					slog.String("data", s),
				)
				return
			}

			// Причина появления этого условия см. https://github.com/bars43ru/gps2Yandex/issues/11
			if int(msgD.Latitude) == 90 && int(msgD.Longitude) == 0 {
				continue
			}

			point := Point{
				messageL: p.msgL,
				messageD: msgD,
			}
			index++
			if !yield(index, point) {
				return
			}
		}
	}
}
