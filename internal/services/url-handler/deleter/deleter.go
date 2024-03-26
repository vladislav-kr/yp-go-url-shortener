package deleter

import (
	"context"
	"sync"
	"time"

	"github.com/vladislav-kr/yp-go-url-shortener/internal/domain/models"
)

type Deleter struct {
	context    context.Context
	jobs       chan models.MassDeleteURL
	result     chan models.DeleteURL
	jobsClosed bool
	callback   func(urls []models.DeleteURL)
}

func NewDeleter(ctx context.Context,
	bufLen int,
	callback func(urls []models.DeleteURL),
) *Deleter {
	d := &Deleter{
		context:  ctx,
		jobs:     make(chan models.MassDeleteURL, bufLen),
		callback: callback,
	}

	go func() {
		<-d.context.Done()
		d.jobsClosed = true
		close(d.jobs)
	}()

	workers := d.fanOut()

	d.result = d.fanIn(workers...)

	d.deleter()

	return d
}

func (d *Deleter) AddMessages(shortURLS []string, userID string) {

	if d.jobsClosed {
		return
	}

	go func() {
		select {
		case d.jobs <- models.MassDeleteURL{ShortURLS: shortURLS, UserID: userID}:
		case <-d.context.Done():
			return
		}
	}()
}

func (d *Deleter) deleter() {
	deleteURLS := []models.DeleteURL{}
	go func() {
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case url, ok := <-d.result:
				if !ok {
					d.callCallback(deleteURLS)
					deleteURLS = deleteURLS[:0]
					return
				}

				deleteURLS = append(deleteURLS, url)

				if len(deleteURLS) > 9 {
					d.callCallback(deleteURLS)
					deleteURLS = deleteURLS[:0]
				}

			case <-d.context.Done():
				d.callCallback(deleteURLS)
				deleteURLS = deleteURLS[:0]
				return
			case <-ticker.C:
				d.callCallback(deleteURLS)
				deleteURLS = deleteURLS[:0]
			}
		}
	}()
}

func (d *Deleter) callCallback(urls []models.DeleteURL) {
	if len(urls) > 0 {
		d.callback(urls)
	}
}
func (d *Deleter) convertor() chan models.DeleteURL {
	result := make(chan models.DeleteURL)

	go func() {
		defer close(result)
		for job := range d.jobs {
			for _, url := range job.ShortURLS {
				select {
				case result <- models.DeleteURL{ShortURL: url, UserID: job.UserID}:
				case <-d.context.Done():
					return
				}
			}

		}
	}()

	return result
}

func (d *Deleter) fanOut() []chan models.DeleteURL {

	workers := 3
	channels := make([]chan models.DeleteURL, workers)

	for i := 0; i < workers; i++ {
		channels[i] = d.convertor()
	}

	return channels
}

func (d *Deleter) fanIn(results ...chan models.DeleteURL) chan models.DeleteURL {
	final := make(chan models.DeleteURL)
	var wg sync.WaitGroup

	for _, ch := range results {
		chClosure := ch

		wg.Add(1)

		go func() {
			defer wg.Done()
			for url := range chClosure {
				select {
				case <-d.context.Done():
					return
				case final <- url:
				}
			}
		}()
	}

	go func() {
		wg.Wait()
		close(final)
	}()

	return final
}
