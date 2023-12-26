package file

import (
	"encoding/json"
	"os"

	"github.com/vladislav-kr/yp-go-url-shortener/internal/domain/models"
)

type Producer struct {
	file    *os.File
	encoder *json.Encoder
}

func NewProducer(path string) (*Producer, error) {

	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return nil, err
	}

	return &Producer{
		file:    file,
		encoder: json.NewEncoder(file),
	}, nil
}

func (p *Producer) Write(shortURL *models.FileURL) error {
	return p.encoder.Encode(shortURL)
}

func (p *Producer) Close() error {
	if p.file != nil {
		return p.file.Close()
	}
	return nil
}

type Consumer struct {
	file    *os.File
	decoder *json.Decoder
}

func NewConsumer(path string) (*Consumer, error) {
	file, err := os.OpenFile(path, os.O_RDONLY, 0644)
	if err != nil {
		return nil, err
	}

	return &Consumer{
		file:    file,
		decoder: json.NewDecoder(file),
	}, nil
}

func (c *Consumer) More() bool {
	return c.decoder.More()
}

func (c *Consumer) Decode(url *models.FileURL) error {
	return c.decoder.Decode(url)
}

func (c *Consumer) Close() error {
	return c.file.Close()
}
