package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"aether-daemon"
)

type GatewaySettings struct {
	PasswordHash string `json:"password_hash"`
	DiskLimitGB  int    `json:"disk_limit_gb"`
	Theme        string `json:"theme"`
	Mode         string `json:"mode"`
}

var settingsCache *GatewaySettings

func getSettingsPath() string {
	wd, _ := os.Getwd()
	return filepath.Join(wd, ".aether", "settings.json")
}

func (a *App) GetSettings() GatewaySettings {
	if settingsCache != nil {
		return *settingsCache
	}
	path := getSettingsPath()
	bytes, err := os.ReadFile(path)
	settings := GatewaySettings{
		DiskLimitGB: 10,
		Theme:       "ctos",
		Mode:        "online",
	}
	if err == nil {
		json.Unmarshal(bytes, &settings)
	}
	settingsCache = &settings
	return settings
}

func (a *App) SaveSettings(newSettings GatewaySettings) error {
	path := getSettingsPath()
	os.MkdirAll(filepath.Dir(path), 0755)
	
	// Preserve existing password hash if not explicitly updated in SaveSettings
	current := a.GetSettings()
	if newSettings.PasswordHash == "" {
		newSettings.PasswordHash = current.PasswordHash
	}

	bytes, err := json.MarshalIndent(newSettings, "", "  ")
	if err != nil {
		return err
	}
	err = os.WriteFile(path, bytes, 0644)
	if err == nil {
		settingsCache = &newSettings
	}
	return err
}

func hashPassword(pwd string) string {
	hash := sha256.Sum256([]byte(pwd))
	return fmt.Sprintf("%x", hash)
}

func (a *App) IsLocked() bool {
	settings := a.GetSettings()
	return settings.PasswordHash != ""
}

func (a *App) SetPassword(pwd string) error {
	settings := a.GetSettings()
	settings.PasswordHash = hashPassword(pwd)
	
	// Derive 32-byte AES key and set it in daemon
	key := sha256.Sum256([]byte(pwd))
	daemon.SetMasterKey(key[:])
	
	return a.SaveSettings(settings)
}

func (a *App) VerifyPassword(pwd string) bool {
	settings := a.GetSettings()
	if settings.PasswordHash == hashPassword(pwd) {
		// Unlock daemon
		key := sha256.Sum256([]byte(pwd))
		daemon.SetMasterKey(key[:])
		return true
	}
	return false
}
