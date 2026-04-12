package unitrade

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"time"
)

const (
	defaultProfileName = "default"
	defaultHost        = "http://127.0.0.1:8888"
)

type Profile struct {
	Name        string    `json:"name"`
	Host        string    `json:"host"`
	Token       string    `json:"token"`
	UserID      uint      `json:"userId"`
	AuthorityID uint      `json:"authorityId"`
	Username    string    `json:"username"`
	NickName    string    `json:"nickName"`
	ExpiresAt   time.Time `json:"expiresAt"`
}

type ProfileStore struct {
	Current  string             `json:"current"`
	Profiles map[string]Profile `json:"profiles"`
}

func defaultStore() *ProfileStore {
	return &ProfileStore{
		Current:  defaultProfileName,
		Profiles: map[string]Profile{},
	}
}

func profileStorePath() (string, error) {
	base, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(base, "unitrade", "profiles.json"), nil
}

func loadProfileStore() (*ProfileStore, string, error) {
	path, err := profileStorePath()
	if err != nil {
		return nil, "", err
	}
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return defaultStore(), path, nil
	}
	if err != nil {
		return nil, "", err
	}

	store := defaultStore()
	if err := json.Unmarshal(data, store); err != nil {
		return nil, "", err
	}
	if store.Profiles == nil {
		store.Profiles = map[string]Profile{}
	}
	if store.Current == "" {
		store.Current = defaultProfileName
	}
	return store, path, nil
}

func saveProfileStore(path string, store *ProfileStore) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(store, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o600)
}

func (s *ProfileStore) activeProfileName(preferred string) string {
	if preferred != "" {
		return preferred
	}
	if s.Current != "" {
		return s.Current
	}
	return defaultProfileName
}

func (s *ProfileStore) getProfile(name string) (Profile, bool) {
	profile, ok := s.Profiles[name]
	return profile, ok
}

func (s *ProfileStore) setProfile(profile Profile) {
	if s.Profiles == nil {
		s.Profiles = map[string]Profile{}
	}
	if profile.Name == "" {
		profile.Name = defaultProfileName
	}
	s.Profiles[profile.Name] = profile
	s.Current = profile.Name
}

func (s *ProfileStore) deleteProfile(name string) {
	delete(s.Profiles, name)
	if s.Current == name {
		s.Current = defaultProfileName
	}
}
