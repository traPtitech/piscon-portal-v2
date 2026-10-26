package runner

import (
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// ログはベンチマーカーの一時ディレクトリ外に保存し、runnerの再起動後も残す。
func openLogArchive(jobID string) (stdout, stderr io.WriteCloser, err error) {
	if !filepath.IsLocal(jobID) || filepath.Base(jobID) != jobID {
		return nil, nil, fmt.Errorf("invalid job ID for log archive")
	}
	root := os.Getenv("PISCON_LOG_DIR")
	if root == "" {
		root = "logs"
	}
	dir := filepath.Join(root, jobID)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return nil, nil, fmt.Errorf("create log archive directory: %w", err)
	}
	open := func(name string) (io.WriteCloser, error) {
		f, err := os.OpenFile(filepath.Join(dir, name+".gz"), os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0600)
		if err != nil {
			return nil, fmt.Errorf("create %s archive: %w", name, err)
		}
		z, err := gzip.NewWriterLevel(f, gzip.BestSpeed)
		if err != nil {
			f.Close()
			return nil, fmt.Errorf("create gzip writer: %w", err)
		}
		return &logArchive{Writer: z, file: f}, nil
	}
	stdout, err = open("stdout")
	if err != nil {
		return nil, nil, err
	}
	stderr, err = open("stderr")
	if err != nil {
		return nil, nil, errors.Join(err, stdout.Close())
	}
	return stdout, stderr, nil
}

type logArchive struct {
	*gzip.Writer
	file *os.File
}

func (a *logArchive) Close() error { return errors.Join(a.Writer.Close(), a.file.Close()) }
