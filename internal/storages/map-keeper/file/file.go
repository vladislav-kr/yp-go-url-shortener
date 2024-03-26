// file отвечает за чтение и запись в файл из in-memory хранилища
package file

import (
	"encoding/json"
	"os"

	"github.com/vladislav-kr/yp-go-url-shortener/internal/domain/models"
)

// Producer хранит параметры для записи в файл.
type Producer struct {
	file    *os.File
	encoder *json.Encoder
}

// NewProducer конструктор Producer.
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

// Write запись в файл
func (p *Producer) Write(shortURL *models.FileURL) error {
	return p.encoder.Encode(shortURL)
}

// Close закрывает файл
func (p *Producer) Close() error {
	if p.file != nil {
		return p.file.Close()
	}
	return nil
}

// Consumer хранит параметры для чтения из файла.
type Consumer struct {
	file    *os.File
	decoder *json.Decoder
}

// NewConsumer конструктор Consumer.
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

// More сообщает, есть ли еще один элемент в
// текущем массиве или анализируемом объекте.
func (c *Consumer) More() bool {
	return c.decoder.More()
}

// Decode читает из файла
func (c *Consumer) Decode(url *models.FileURL) error {
	return c.decoder.Decode(url)
}

// Close закрывает файл.
func (c *Consumer) Close() error {
	return c.file.Close()
}
