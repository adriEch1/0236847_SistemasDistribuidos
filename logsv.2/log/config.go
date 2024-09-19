package log

// estructura que nos permite asignar el tamaño máximo del store y del index
type Config struct {
	Segment struct {
		MaxStoreBytes uint64
		MaxIndexBytes uint64
		InitialOffset uint64
	}
}
