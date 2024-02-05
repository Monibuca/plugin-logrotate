package logrotate

import (
	"fmt"
	"net/http"
	"os/exec"
)

func (l *LogRotateConfig) API_find(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("query")
	if query == "" {
		w.Write([]byte("query is empty"))
		return
	}
	cmd := exec.Command("grep", "-r",  query, l.Path)
	cmd.Stdout = w
	cmd.Stderr = w
	cmd.Run()
}
