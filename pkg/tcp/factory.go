// Package tcp предоставляет компоненты для работы с TCP-соединениями.
// Содержит сервер и обработчики для приема и обработки TCP-соединений.
package tcp

import (
	"fmt"
)

// New создает новый экземпляр TCP-сервера.
// Инициализирует сервер с указанным адресом и обработчиком соединений.
//
// Параметры:
//   - address: адрес для прослушивания в формате "host:port"
//   - handler: обработчик входящих соединений
//
// Возвращает:
//   - *Server: указатель на созданный сервер
//   - error: ошибка в случае некорректных параметров
//
// Ошибки:
//   - Возвращает ошибку при пустом адресе
//   - Возвращает ошибку при отсутствии обработчика
func New(
	address string,
	handler ConnectionHandler,
) (*Server, error) {
	if address == "" {
		return nil, fmt.Errorf("param `address`empty")
	}
	if handler == nil {
		return nil, fmt.Errorf("param `handler` empty")
	}
	return &Server{
		address: address,
		handler: handler,
	}, nil
}
