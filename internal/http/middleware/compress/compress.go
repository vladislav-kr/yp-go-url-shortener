package compress

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"slices"
	"strings"
	"sync"

)

type compressWriter struct {
	w http.ResponseWriter
	// Функция получения gzip.Writer из sync.Pool.
	// Должна вызываться в момент выбора способа записи.
	// Если клиент может обработать сжатую информацию, но Content-Typeне неподдерживаемый,
	// то создавать нельзя, т.к. в исходный writer пишеться заголовок gzip и конечные данные будут неверные
	newCompressWriter func(w io.Writer) (writer *gzip.Writer, cleanup func() (err error))
	zw                *gzip.Writer
	// Закрывает gzip.Write и возвращает в sync.Pool
	cleanup         func() (err error)
	supportCompress func(contType string) bool
}

func (cw *compressWriter) Header() http.Header {

	return cw.w.Header()
}

func (cw *compressWriter) Write(p []byte) (int, error) {

	if cw.supportCompress(cw.w.Header().Get("Content-Type")) {
		if cw.zw == nil {
			cw.zw, cw.cleanup = cw.newCompressWriter(cw.w)
		}
		return cw.zw.Write(p)
	}

	return cw.w.Write(p)
}

func (cw *compressWriter) WriteHeader(statusCode int) {

	if statusCode < http.StatusMultipleChoices &&
		cw.supportCompress(cw.w.Header().Get("Content-Type")) {
		cw.w.Header().Set("Content-Encoding", "gzip")
	}
	cw.w.WriteHeader(statusCode)
}

func (cw *compressWriter) Close() error {

	if cw.cleanup != nil {
		return cw.cleanup()
	}

	return nil
}

type compressReader struct {
	r       io.ReadCloser
	zr      *gzip.Reader
	cleanup func() (err error)
}

func (cr compressReader) Read(p []byte) (n int, err error) {

	return cr.zr.Read(p)
}

func (cr *compressReader) Close() error {

	if err := cr.cleanup(); err != nil {
		return err
	}
	return cr.r.Close()
}

type compressPool struct {
	gw           *sync.Pool
	gr           *sync.Pool
	supContTypes []string
}

func NewCompressPool(supContTypes []string) (*compressPool, error) {

	buf := bytes.NewBuffer(nil)

	gw := gzip.NewWriter(buf)
	if err := gw.Close(); err != nil {
		return nil, fmt.Errorf("component: middleware/compress, failed closing gzip.Writer: %v", err)
	}

	gr, err := gzip.NewReader(buf)
	if err != nil {
		return nil, fmt.Errorf("component: middleware/compress, failed creation gzip.Reader: %v", err)
	}
	if err := gr.Close(); err != nil {
		return nil, fmt.Errorf("component: middleware/compress, failed closing gzip.Reader: %v", err)
	}

	return &compressPool{
		gw: &sync.Pool{
			New: func() interface{} {
				return gw
			},
		},
		gr: &sync.Pool{
			New: func() interface{} {
				return gr
			},
		},
		supContTypes: supContTypes,
	}, nil
}

func (cp *compressPool) NewCompressWriter(w http.ResponseWriter) *compressWriter {

	return &compressWriter{
		w:                 w,
		newCompressWriter: cp.writer,
		supportCompress: func(contType string) bool {
			return slices.ContainsFunc(cp.supContTypes, func(sType string) bool {
				return strings.Contains(contType, sType)
			})
		},
	}
}

func (cp *compressPool) NewCompressReader(r io.ReadCloser) *compressReader {

	zr, cleanup := cp.reader(r)
	return &compressReader{
		r:       r,
		zr:      zr,
		cleanup: cleanup,
	}
}

func (cp *compressPool) reader(r io.Reader) (reader *gzip.Reader, cleanup func() (err error)) {

	reader = cp.gr.Get().(*gzip.Reader)
	reader.Reset(r)

	cleanup = func() (err error) {
		err = reader.Close()
		cp.gr.Put(reader)
		return
	}
	return
}

func (cp *compressPool) writer(w io.Writer) (writer *gzip.Writer, cleanup func() (err error)) {

	writer = cp.gw.Get().(*gzip.Writer)
	writer.Reset(w)

	cleanup = func() (err error) {
		err = writer.Close()
		cp.gw.Put(writer)
		return
	}

	return
}
