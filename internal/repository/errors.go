package repository

import "errors"

var ErrNotFound = errors.New("not found")

const (
	FileDatasourceRoute     = "./datasource/route.txt"
	FileDatasourceSchedule  = "./datasource/schedule.txt"
	FileDatasourceTransport = "./datasource/transport.txt"
)
