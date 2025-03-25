package egts

import (
	"bufio"
	"context"
	"encoding/binary"
	"errors"
	"github.com/kuznetsovin/egts-protocol/libs/egts"
	"github.com/labstack/gommon/log"
	"io"
	"iter"
	"log/slog"
)

type Parser struct {
	reader *bufio.Reader
}

func NewParse(reader io.Reader) *Parser {
	return &Parser{
		reader: bufio.NewReader(reader),
	}
}

func (p *Parser) Points(ctx context.Context) iter.Seq2[int, Point] {
	const (
		egtsPcOk  = 0
		headerLen = 10
	)
	var (
		recvPacket []byte
		client     uint32
	)

	return func(yield func(int, Point) bool) {
		index := -1
		for {
			recvPacket = nil

			// считываем заголовок пакета
			headerBuf := make([]byte, headerLen)
			_, err := io.ReadFull(p.reader, headerBuf)
			if err != nil {
				if errors.Is(err, io.EOF) {
					return
				}
				slog.Error("read to buffer", slog.Any("error", err))
				return
			}
			// если пакет не егтс формата закрываем соединение
			if headerBuf[0] != 0x01 {
				//log.WithField("ip", conn.RemoteAddr()).Warn("Пакет не соответствует формату ЕГТС. Закрыто соединение")
				return
			}

			// вычисляем длину пакета, равную длине заголовка (HL) + длина тела (FDL) + CRC пакета 2 байта если есть FDL из приказа минтранса №285
			bodyLen := binary.LittleEndian.Uint16(headerBuf[5:7])
			pkgLen := uint16(headerBuf[3])
			if bodyLen > 0 {
				pkgLen += bodyLen + 2
			}
			// получаем концовку ЕГТС пакета
			buf := make([]byte, pkgLen-headerLen)
			if _, err := io.ReadFull(p.reader, buf); err != nil {
				slog.Error("read body package", slog.Any("error", err))
				return
			}
			// формируем полный пакет
			recvPacket = append(headerBuf, buf...)

			pkg := egts.Package{}
			resultCode, err := pkg.Decode(recvPacket)
			if resultCode != egtsPcOk {
				slog.Error("decoding packet", slog.Any("error", err))
				continue
			}
			if pkg.PacketType != egts.PtAppdataPacket {
				continue
			}
			log.Debug("received package EGTS_PT_APPDATA")

			for _, rec := range *pkg.ServicesFrameData.(*egts.ServiceDataSet) {
				point := Point{
					PacketID: uint32(pkg.PacketIdentifier),
				}

				// если в секции с данными есть oid то обновляем его
				if rec.ObjectIDFieldExists == "1" {
					client = rec.ObjectIdentifier
				}

				for _, subRec := range rec.RecordDataSet {
					srPosData, ok := subRec.SubrecordData.(*egts.SrPosData)
					if !ok {
						continue
					}

					point.Client = client
					point.Time = srPosData.NavigationTime
					point.Latitude = srPosData.Latitude
					point.Longitude = srPosData.Longitude
					point.Speed = srPosData.Speed
					point.Course = srPosData.Direction
					index++
					if !yield(index, point) {
						return
					}
					break
				}
			}
		}
	}
}
