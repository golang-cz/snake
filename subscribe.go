package main

import "github.com/golang-cz/snake/proto"

func (s *Server) subscribe(c chan *proto.Update) (*proto.Update, uint64) {
	s.mu.Lock()
	defer s.mu.Unlock()

	id := s.lastSubId
	s.subs[id] = c
	s.lastSubId++

	u := &proto.Update{
		State: s.state,
	}

	return u, id
}

func (s *Server) unsubscribe(subscriptionId uint64) {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.subs, subscriptionId)
}
