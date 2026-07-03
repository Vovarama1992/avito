package audit

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"
)

var mu sync.Mutex

func Logf(format string, args ...any) {
	message := fmt.Sprintf(format, args...)
	log.Println(message)
	writeDaily(message)
}

func writeDaily(message string) {
	dir := os.Getenv("AVITO_AUDIT_LOG_DIR")
	if dir == "" {
		dir = "logs"
	}

	now := time.Now()
	line := fmt.Sprintf("%s %s\n", now.Format(time.RFC3339), message)
	path := filepath.Join(dir, now.Format("2006-01-02")+"-important.log")

	mu.Lock()
	defer mu.Unlock()

	if err := os.MkdirAll(dir, 0750); err != nil {
		log.Println("AUDIT LOG ERROR: mkdir:", err)
		return
	}

	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0640)
	if err != nil {
		log.Println("AUDIT LOG ERROR: open:", err)
		return
	}
	defer f.Close()

	if _, err := f.WriteString(line); err != nil {
		log.Println("AUDIT LOG ERROR: write:", err)
	}
}
