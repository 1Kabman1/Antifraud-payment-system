package hashStorage

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewStorage(t *testing.T) {
	s := NewStorage()
	assert.Equal(t, len(s.rules), 0)
	assert.Equal(t, len(s.counter), 0)
	assert.Equal(t, len(s.archivist), 0)
}

func TestStorage_Counter(t *testing.T) {
	s := NewStorage()
	err, _ := s.Counter([16]byte{})
	if err == nil {
		panic(err)
	}
}
