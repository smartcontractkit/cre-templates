package store

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/example/cre-confidential-authorization-template/gateway/internal/model"
)

var (
	ErrExists   = errors.New("claim already exists")
	ErrNotFound = errors.New("claim not found")
	ErrTerminal = errors.New("claim terminal state is immutable")
)

type snapshot struct {
	Version int                           `json:"version"`
	Claims  map[string]*model.ClaimRecord `json:"claims"`
}

type Store struct {
	mu     sync.Mutex
	claims map[string]*model.ClaimRecord
	path   string
}

func New() *Store { return &Store{claims: map[string]*model.ClaimRecord{}} }

func NewPersistent(path string) (*Store, error) {
	s := &Store{claims: map[string]*model.ClaimRecord{}, path: path}
	if err := s.load(); err != nil {
		return nil, err
	}
	return s, nil
}

func (s *Store) load() error {
	if s.path == "" {
		return nil
	}
	b, err := os.ReadFile(s.path)
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	if err != nil {
		return err
	}
	var snap snapshot
	if err := json.Unmarshal(b, &snap); err != nil {
		return err
	}
	if snap.Claims != nil {
		s.claims = snap.Claims
	}
	return nil
}

func (s *Store) persistLocked() error {
	if s.path == "" {
		return nil
	}
	if err := os.MkdirAll(filepath.Dir(s.path), 0o700); err != nil {
		return err
	}
	b, err := json.MarshalIndent(snapshot{Version: 1, Claims: s.claims}, "", "  ")
	if err != nil {
		return err
	}
	tmp := s.path + ".tmp"
	if err := os.WriteFile(tmp, append(b, '\n'), 0o600); err != nil {
		return err
	}
	return os.Rename(tmp, s.path)
}

func clone(r *model.ClaimRecord) *model.ClaimRecord {
	if r == nil {
		return nil
	}
	out := *r
	return &out
}

func (s *Store) AddClaim(c model.ClaimEnvelope) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.claims[c.ClaimID]; ok {
		return ErrExists
	}
	s.claims[c.ClaimID] = &model.ClaimRecord{Claim: c, Status: "PENDING", ReceivedAt: time.Now().Unix()}
	if err := s.persistLocked(); err != nil {
		delete(s.claims, c.ClaimID)
		return err
	}
	return nil
}

func (s *Store) GetClaim(id string) (*model.ClaimRecord, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	r, ok := s.claims[id]
	return clone(r), ok
}

func (s *Store) SetClaimStatus(id, status, reason string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	r, ok := s.claims[id]
	if !ok {
		return ErrNotFound
	}
	if r.Status == "AUTHORIZATION_REQUIRED" || r.Status == "REJECTED" {
		return ErrTerminal
	}
	if status != "AUTHORIZATION_REQUIRED" && status != "REJECTED" {
		return errors.New("unsupported claim status")
	}
	old := *r
	r.Status, r.Reason, r.UpdatedAt = status, reason, time.Now().Unix()
	if err := s.persistLocked(); err != nil {
		*r = old
		return err
	}
	return nil
}
