// Package repository содержит реализацию репозиториев для работы с данными.
// Предоставляет интерфейсы и структуры для хранения и доступа к информации о маршрутах, транспорте и расписании.
package repository

import (
	"sync/atomic"
)

// SafeMapAtomic представляет потокобезопасную реализацию карты с использованием atomic.Value.
// Обеспечивает безопасный доступ к данным в многопоточной среде.
// Используется для хранения и обновления данных в репозиториях.
type SafeMapAtomic[Key comparable, Value any] struct {
	data atomic.Value // Атомарное значение для хранения карты
}

// Get возвращает значение по ключу из карты.
// Безопасно работает в многопоточной среде.
//
// Параметры:
//   - key: ключ для поиска значения
//
// Возвращает:
//   - Value: найденное значение
//   - bool: флаг наличия значения
func (s *SafeMapAtomic[Key, Value]) Get(key Key) (Value, bool) {
	m := s.data.Load().(map[Key]Value)
	value, exists := m[key]
	return value, exists
}

// Replace заменяет содержимое карты новыми данными.
// Безопасно работает в многопоточной среде.
//
// Параметры:
//   - data: новая карта для замены
func (s *SafeMapAtomic[Key, Value]) Replace(data map[Key]Value) {
	s.data.Store(data)
}

// NewSafeMapAtomic создает новый экземпляр потокобезопасной карты.
// Инициализирует пустую карту.
//
// Возвращает:
//   - SafeMapAtomic[Key, Value]: новый экземпляр карты
func NewSafeMapAtomic[Key comparable, Value any]() SafeMapAtomic[Key, Value] {
	var sm SafeMapAtomic[Key, Value]
	sm.Replace(make(map[Key]Value))
	return sm
}
