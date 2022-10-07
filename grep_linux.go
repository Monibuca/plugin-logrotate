package logrotate

import (
	"fmt"
	"net/http"
	"os/exec"
)

func (l *LogRotateConfig) API_find(w http.ResponseWriter, r *http.Request) {
	cmd := exec.Command("grep", fmt.Sprintf("\"%s\"", r.URL.Query().Get("query")), l.Path)
	cmd.Stdout = w
	cmd.Stderr = w
	cmd.Run()
}
