package ws

import (
	"bufio"
	"context"
	"errors"
	"io"
	"log"
	"os"
	"time"
)

func FileMonitor(ctx context.Context, filename string, group string, hook func(string, []byte)) {
	// TODO: 文件不存在时处理
	f, err := os.Open(filename)
	if err != nil {
		log.Fatalln(err)
	}
	defer f.Close()
	f.Seek(0, io.SeekEnd) // nolint errcheck
	reader := bufio.NewReader(f)
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		line, _, err := reader.ReadLine()
		if err != nil {
			if errors.Is(err, io.EOF) {
				time.Sleep(500 * time.Millisecond)
				continue
			}
			return
		}
		go hook(group, line)
	}
}
