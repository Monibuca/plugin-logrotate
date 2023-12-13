package logrotate

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	. "m7s.live/engine/v4"
	"m7s.live/engine/v4/log"
	"m7s.live/engine/v4/util"
	. "m7s.live/engine/v4/util"
)

type FileInfo struct {
	Name string
	Size int64
}
type LogRotateConfig struct {
	Path        string `default:"./logs" desc:"日志文件存放目录"`
	Size        int64  `desc:"日志文件大小，单位：字节"`
	Days        int    `default:"1" desc:"日志文件保留天数"`
	Formatter   string `default:"2006-01-02T15" desc:"日志文件名格式"`
	file        *os.File
	currentSize int64
	createTime  time.Time
	hours       float64
	splitFunc   func() bool
}

var _ = InstallPlugin(&LogRotateConfig{})

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

func (l *LogRotateConfig) API_list(w http.ResponseWriter, r *http.Request) {
	dir, err := os.Open(l.Path)
	if err == nil {
		var files []os.FileInfo
		if files, err = dir.Readdir(0); err == nil {
			var fileInfos []*FileInfo
			for _, info := range files {
				fileInfos = append(fileInfos, &FileInfo{
					info.Name(), info.Size(),
				})
			}
			err = json.NewEncoder(w).Encode(fileInfos)
		}
	}
	if err != nil {
		ReturnError(APIErrorOpen, err.Error(), w, r)
	} else {
		ReturnOK(w, r)
	}
}

func (l *LogRotateConfig) API_download(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Disposition", "attachment; filename="+r.URL.Query().Get("file"))
	l.API_open(w, r)
}

func (l *LogRotateConfig) API_open(w http.ResponseWriter, r *http.Request) {
	path := filepath.Join(l.Path, r.URL.Query().Get("file"))
	if !util.IsSubdir(l.Path, path) {
		http.Error(w, "invalid file", http.StatusBadRequest)
		return
	}
	file, err := os.Open(path)
	if err == nil {
		defer file.Close()
		_, err = io.Copy(w, file)
	}
	if err != nil {
		ReturnError(APIErrorOpen, err.Error(), w, r)
	} else {
		ReturnOK(w, r)
	}
}
