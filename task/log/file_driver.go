package log

import (
	"bufio"
	"bytes"
	"errors"
	"github.com/apex/log"
	"github.com/crawlab-team/crawlab-core/utils"
	"github.com/crawlab-team/go-trace"
	"github.com/spf13/viper"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

type FileLogDriver struct {
	// settings
	opts *FileLogDriverOptions // options

	// internals
	mu          sync.Mutex
	logFileName string
}

type FileLogDriverOptions struct {
	BaseDir string
	Ttl     time.Duration
}

func (d *FileLogDriver) Init() (err error) {
	go d.cleanup()

	return nil
}

func (d *FileLogDriver) Close() (err error) {
	return nil
}

func (d *FileLogDriver) WriteLine(id string, line string) (err error) {
	d.initDir(id)

	d.mu.Lock()
	defer d.mu.Unlock()
	filePath := d.getLogFilePath(id, d.logFileName)

	f, err := os.OpenFile(filePath, os.O_WRONLY|os.O_APPEND|os.O_CREATE, os.FileMode(0760))
	if err != nil {
		return trace.TraceError(err)
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			log.Errorf("close file error: %s", err.Error())
		}
	}(f)

	_, err = f.WriteString(line + "\n")
	if err != nil {
		return trace.TraceError(err)
	}

	return nil
}

func (d *FileLogDriver) WriteLines(id string, lines []string) (err error) {
	linesString := strings.Join(lines, "\n")
	if err := d.WriteLine(id, linesString); err != nil {
		return err
	}
	return nil
}

func (d *FileLogDriver) Find(id string, pattern string, skip int, limit int) (lines []string, err error) {
	if pattern != "" {
		return lines, errors.New("not implemented")
	}
	if !utils.Exists(d.getLogFilePath(id, d.logFileName)) {
		return nil, nil
	}

	f, err := os.Open(d.getLogFilePath(id, d.logFileName))
	if err != nil {
		return nil, trace.TraceError(err)
	}
	defer f.Close()

	sc := bufio.NewReaderSize(f, 1024*1024*10)

	i := -1
	for {
		line, err := sc.ReadString(byte('\n'))
		if err != nil {
			break
		}
		line = strings.TrimSuffix(line, "\n")

		i++

		if i < skip {
			continue
		}

		if i >= skip+limit {
			break
		}

		lines = append(lines, line)
	}

	return lines, nil
}

func (d *FileLogDriver) Count(id string, pattern string) (n int, err error) {
	if pattern != "" {
		return n, errors.New("not implemented")
	}
	if !utils.Exists(d.getLogFilePath(id, d.logFileName)) {
		return 0, nil
	}

	f, err := os.Open(d.getLogFilePath(id, d.logFileName))
	if err != nil {
		return n, trace.TraceError(err)
	}
	return d.lineCounter(f)
}

func (d *FileLogDriver) Flush() (err error) {
	return nil
}

func (d *FileLogDriver) getBasePath(id string) (filePath string) {
	return filepath.Join(d.opts.BaseDir, id)
}

func (d *FileLogDriver) getMetadataPath(id string) (filePath string) {
	return filepath.Join(d.opts.BaseDir, id, MetadataName)
}

func (d *FileLogDriver) getLogFilePath(id, fileName string) (filePath string) {
	return filepath.Join(d.opts.BaseDir, id, fileName)
}

func (d *FileLogDriver) initDir(id string) {
	if !utils.Exists(d.getBasePath(id)) {
		if err := os.MkdirAll(d.getBasePath(id), os.FileMode(0770)); err != nil {
			trace.PrintError(err)
		}
	}
}

func (d *FileLogDriver) lineCounter(r io.Reader) (n int, err error) {
	buf := make([]byte, 32*1024)
	count := 0
	lineSep := []byte{'\n'}

	for {
		c, err := r.Read(buf)
		count += bytes.Count(buf[:c], lineSep)

		switch {
		case err == io.EOF:
			return count, nil

		case err != nil:
			return count, err
		}
	}
}

func (d *FileLogDriver) cleanup() {
	for {
		// 增加对目录不存在的判断
		dirs, err := utils.ListDir(d.opts.BaseDir)
		if err != nil {
			trace.PrintError(err)
			time.Sleep(10 * time.Minute)
			continue
		}
		for _, dir := range dirs {
			info, err := dir.Info()
			if err != nil {
				trace.PrintError(err)
				continue
			}
			if time.Now().After(info.ModTime().Add(d.opts.Ttl)) {
				if err := os.RemoveAll(d.getBasePath(dir.Name())); err != nil {
					trace.PrintError(err)
					continue
				}
				log.Infof("removed outdated log directory: %s", d.getBasePath(dir.Name()))
			}
		}

		time.Sleep(10 * time.Minute)
	}
}

var logDriver Driver

func newFileLogDriver(options *FileLogDriverOptions) (driver Driver, err error) {
	if options == nil {
		options = &FileLogDriverOptions{}
	}

	// normalize BaseDir
	baseDir := options.BaseDir
	if baseDir == "" {
		baseDir = "/var/log/crawlab"
	}
	options.BaseDir = baseDir

	// normalize Ttl
	ttl := options.Ttl
	if ttl == 0 {
		ttlSeconds := viper.GetInt("log.ttl")
		if ttlSeconds == 0 {
			ttl = 30 * 24 * time.Hour
		} else {
			ttl = time.Second * time.Duration(ttlSeconds)
		}
	}
	options.Ttl = ttl

	// driver
	driver = &FileLogDriver{
		opts:        options,
		logFileName: "log.txt",
		mu:          sync.Mutex{},
	}

	// init
	if err := driver.Init(); err != nil {
		return nil, err
	}

	return driver, nil
}

func GetFileLogDriver(options *FileLogDriverOptions) (driver Driver, err error) {
	if logDriver != nil {
		return logDriver, nil
	}
	logDriver, err = newFileLogDriver(options)
	if err != nil {
		return nil, err
	}
	return logDriver, nil
}
