package transport_type

//go:generate go tool go-enum --lower --marshal --names --values transport_type

// ENUM(
//
//	BUS,
//	TROLLEYBUS,
//	TRAMWAY,
//	MINIBUS,
//
// )
type Type int
