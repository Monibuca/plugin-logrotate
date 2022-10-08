package logrotate

import "net/http"

func (l *LogRotateConfig) API_find(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("not support windows"))
}
