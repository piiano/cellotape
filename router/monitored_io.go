package router

import (
	"io"
	"net/http"
	"time"

	"github.com/piiano/cellotape/router/utils"
)

type HTTPDurations interface {
	ReadDuration() time.Duration
	WriteDuration() time.Duration
}

type monitoredHTTP struct {
	reader        io.ReadCloser
	writer        http.ResponseWriter
	readStart     *time.Time
	writeStart    *time.Time
	readDuration  time.Duration
	writeDuration time.Duration
}

type MonitoredHTTP interface {
	io.ReadCloser
	http.ResponseWriter
	HTTPDurations
}

func NewMonitoredHTTP(writer http.ResponseWriter, reader io.ReadCloser) MonitoredHTTP {
	return &monitoredHTTP{
		reader: reader,
		writer: writer,
	}
}

func (mr *monitoredHTTP) Read(p []byte) (int, error) {
	if mr.readStart == nil {
		mr.readStart = utils.Ptr(time.Now())
	}

	n, err := mr.reader.Read(p)
	mr.readDuration = time.Since(*mr.readStart)
	return n, err
}

func (mr *monitoredHTTP) Close() error {
	return mr.reader.Close()
}

func (mr *monitoredHTTP) initWriteStart() {
	if mr.writeStart == nil {
		mr.writeStart = utils.Ptr(time.Now())
	}
}

func (mr *monitoredHTTP) Write(p []byte) (int, error) {
	mr.initWriteStart()

	n, err := mr.writer.Write(p)
	mr.writeDuration += time.Since(*mr.writeStart)
	return n, err
}

func (mr *monitoredHTTP) Header() http.Header {
	return mr.writer.Header()
}

func (mr *monitoredHTTP) WriteHeader(statusCode int) {
	mr.initWriteStart()

	mr.writer.WriteHeader(statusCode)
}

func (mr *monitoredHTTP) ReadDuration() time.Duration {
	return mr.readDuration
}

func (mr *monitoredHTTP) WriteDuration() time.Duration {
	return mr.writeDuration
}
