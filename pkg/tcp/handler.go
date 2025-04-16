// Package tcp предоставляет компоненты для работы с TCP-соединениями.
// Содержит сервер и обработчики для приема и обработки TCP-соединений.
package tcp

import (
	"context"
	"io"
)

// ConnectionHandler определяет интерфейс для обработки TCP-соединений.
// Используется для обработки входящих данных от клиентов.
type ConnectionHandler interface {
	// Accept обрабатывает входящее TCP-соединение.
	// Принимает контекст и читатель для получения данных от клиента.
	Accept(ctx context.Context, r io.Reader) error
}

// ConnectionHandlerFunc представляет функцию-обработчик TCP-соединения.
// Реализует интерфейс ConnectionHandler.
type ConnectionHandlerFunc func(ctx context.Context, r io.Reader) error

// Accept реализует метод интерфейса ConnectionHandler для ConnectionHandlerFunc.
// Вызывает функцию-обработчик с переданными параметрами.
func (h ConnectionHandlerFunc) Accept(ctx context.Context, r io.Reader) error {
	return h(ctx, r)
}
