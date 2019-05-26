package utils

type IdSequenceInterface interface {
	NextId() uint64
}

func NewIdSequence() *IdSequence {
	return &IdSequence{}
}

type IdSequence struct {
	nextId uint64
}

func (r *IdSequence) NextId() uint64 {
	r.nextId++
	return r.nextId
}
