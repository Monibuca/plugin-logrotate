package logrotate

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	. "m7s.live/engine/v4"
	"m7s.live/engine/v4/log"
	. "m7s.live/engine/v4/util"
)

type FileInfo struct {
	Name string
	Size int64
}
type LogRotateConfig struct {
	Path        string
	Size        int64
	Days        int
	Formatter   string
	file        *os.File
	currentSize int64
	createTime  time.Time
	hours       float64
	splitFunc   func() bool
}

var plugin = InstallPlugin(&LogRotateConfig{
	Path:      "./logs",
	Days:      1,
	Formatter: "2006-01-02T15",
})

func (config *LogRotateConfig) OnEvent(event any) {
	switch event.(type) {
	case FirstConfig:
		if config.Size > 0 {
			config.splitFunc = config.splitBySize
		} else {
			if config.Days == 0 {
				config.Days = 1
			}
			config.hours = float64(config.Days) * 24
			config.splitFunc = config.splitByTime
		}
		config.createTime = time.Now()
		if config.Formatter == "" {
			config.Formatter = "2006-01-02T15"
		}
		err := os.MkdirAll(config.Path, 0766)
		config.file, err = os.OpenFile(filepath.Join(config.Path, fmt.Sprintf("%s.log", config.createTime.Format(config.Formatter))), os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0666)
		if err == nil {
			stat, _ := config.file.Stat()
			config.currentSize = stat.Size()
			log.AddWriter(config)
		} else {
			log.Error(err)
		}
	}
}

func (l *LogRotateConfig) splitBySize() bool {
	return l.currentSize >= l.Size
}
func (l *LogRotateConfig) splitByTime() bool {
	return time.Since(l.createTime).Hours() > l.hours
}
func (l *LogRotateConfig) Write(data []byte) (n int, err error) {
	n, err = l.file.Write(data)
	l.currentSize += int64(n)
	if err == nil {
		if l.splitFunc() {
			l.createTime = time.Now()
			if file, err := os.OpenFile(filepath.Join(l.Path, fmt.Sprintf("%s.log", l.createTime.Format(l.Formatter))), os.O_TRUNC|os.O_WRONLY|os.O_CREATE, 0666); err == nil {
				l.file = file
				l.currentSize = 0
			}
		}
	}
	return
}
func (l *LogRotateConfig) API_tail(w http.ResponseWriter, r *http.Request) {
	writer := NewSSE(w, r.Context())
	log.AddWriter(writer)
	<-r.Context().Done()
	log.DeleteWriter(writer)
}
func (l *LogRotateConfig) API_find(w http.ResponseWriter, r *http.Request) {
	cmd := exec.Command("grep", fmt.Sprintf("\"%s\"", r.URL.Query().Get("query")), l.Path)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Run()
	w.Write([]byte(out.String()))
}
func (l *LogRotateConfig) API_list(w http.ResponseWriter, r *http.Request) {
	dir, err := os.Open(l.Path)
	defer func() {
		if err != nil {
			w.Write([]byte(err.Error()))
		}
	}()
	if err != nil {
		return
	}
	var files []os.FileInfo
	files, err = dir.Readdir(0)
	if err != nil {
		return
	}
	var fileInfos []*FileInfo
	for _, info := range files {
		fileInfos = append(fileInfos, &FileInfo{
			info.Name(), info.Size(),
		})
	}
	var bytes []byte
	bytes, err = json.Marshal(fileInfos)
	if err != nil {
		return
	}
	w.Write(bytes)
}
func (l *LogRotateConfig) API_download(w http.ResponseWriter, r *http.Request) {
	filename := r.URL.Query().Get("file")
	file, err := os.Open(filepath.Join(l.Path, filename))
	defer func() {
		if err != nil {
			w.Write([]byte(err.Error()))
		}
	}()
	w.Header().Add("Content-Disposition", "attachment; filename="+filename)
	if err != nil {
		return
	}
	_, err = io.Copy(w, file)
	if err != nil {
		return
	}
}
func (l *LogRotateConfig) API_open(w http.ResponseWriter, r *http.Request) {
	filename := r.URL.Query().Get("file")
	file, err := os.Open(filepath.Join(l.Path, filename))
	defer func() {
		if err != nil {
			w.Write([]byte(err.Error()))
		}
	}()
	if err != nil {
		return
	}
	_, err = io.Copy(w, file)
	if err != nil {
		return
	}
}
