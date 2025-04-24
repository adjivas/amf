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

func (s *EirSet) Add(v string) error {
	if _, exists := s.lookup[v]; exists {
		logger.UtilLog.Debugln("EIR value already exists")
		return errors.New("EIR value already exists")
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
		return nil
	}
}

func (s *EirSet) Remove(v string) error {
	n, exists := s.lookup[v]
	if !exists {
		logger.UtilLog.Debugln("EIR missing value")
		return errors.New("EIR missing value")
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
	return nil
}

func (s *EirSet) Next() (string, error) {
	if s.head == nil {
		logger.UtilLog.Debugln("EIR set is empty")
		return "", errors.New("EIR set is empty")
	} else {
		s.head = s.head.Next()
		return s.head.Value.(string), nil
	}
}

func (s *EirSet) PrintAll() {
	eirList := []string { }
	for index := 0; index < s.size; index++ {
		eir, _ := s.Next()
		eirList = append(eirList, eir)
	}
	logger.UtilLog.Infof("EirSet List: (%s)", eirList)
}
