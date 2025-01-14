package service

import (
	"github.com/google/uuid"
	"sync"
	"time"
)

type UploadSession struct {
	ID             string
	ExpiresAt      time.Time
	TotalChunks    int
	ReceivedChunks map[int]bool
	FileName       string
}

type UploadService struct {
	sessions map[string]UploadSession
	mutex    sync.Mutex
	timeout  time.Duration
}

func NewUploadService(timeout time.Duration) *UploadService {
	return &UploadService{
		sessions: make(map[string]UploadSession),
		timeout:  timeout,
	}
}

func (s *UploadService) StartSession(fileName string, totalChunks int) string {
	id := uuid.New().String()
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.sessions[id] = UploadSession{
		ID:             id,
		ExpiresAt:      time.Now().Add(s.timeout),
		FileName:       fileName,
		TotalChunks:    totalChunks,
		ReceivedChunks: make(map[int]bool),
	}
	return id
}

func (s *UploadService) AddChunk(sessionID string, chunkIndex int) bool {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	session, exists := s.sessions[sessionID]
	if !exists || session.ExpiresAt.Before(time.Now()) {
		return false
	}
	session.ReceivedChunks[chunkIndex] = true
	s.sessions[sessionID] = session
	return true
}

func (s *UploadService) IsUploadComplete(sessionID string) bool {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	session, exists := s.sessions[sessionID]
	return exists && len(session.ReceivedChunks) == session.TotalChunks
}

func (s *UploadService) CleanExpiredSessions() {
	for {
		time.Sleep(1 * time.Minute)
		s.mutex.Lock()
		for id, session := range s.sessions {
			if session.ExpiresAt.Before(time.Now()) {
				delete(s.sessions, id)
			}
		}
		s.mutex.Unlock()
	}
}
