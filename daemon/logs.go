package main

import (
	"encoding/json"
	"net/http"
)

func (a *App) logsHandler(w http.ResponseWriter, r *http.Request) {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	json.NewEncoder(w).Encode(a.timerLogs)
}
