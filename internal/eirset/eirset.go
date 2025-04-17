package util

import (
	"container/ring"
	"errors"

	"github.com/free5gc/amf/internal/logger"
)

type EirSet struct {
	head   *ring.Ring
	lookup map[string]*ring.Ring
	size   int
}

func New() *EirSet {
	return &EirSet{
		lookup: make(map[string]*ring.Ring),
	}
}

func (s *EirSet) Add(v string) {
	if _, exists := s.lookup[v]; exists {
		logger.UtilLog.Debugln("value already present")
	} else {
		n := ring.New(1)
		n.Value = v

		if s.head == nil {
			s.head = n
		} else {
			s.head.Link(n)
		}

		s.lookup[v] = n
		s.size++
	}
}

func (s *EirSet) Remove(v string) {
	n, exists := s.lookup[v]
	if !exists {
		logger.UtilLog.Debugln("value not found")
	}

	if s.size == 1 {
		s.head = nil
	} else {
		prev := n.Prev()
		prev.Unlink(1)
		if n == s.head {
			s.head = prev.Next()
		}
	}
	delete(s.lookup, v)
	s.size--
}

func (s *EirSet) Next() (string, error) {
	if s.head == nil {
		return "", errors.New("set is empty")
	}
	s.head = s.head.Next()
	return s.head.Value.(string), nil
}
