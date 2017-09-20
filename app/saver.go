package app

import (
	"log"
	"os"
)

type (
	Saver struct {
		file *os.File
	}
)

func NewSaver(path string) *Saver {
	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_SYNC, 0666)

	if err != nil {
		log.Fatal(err)
	}

	return &Saver{
		file: file,
	}
}

func (s *Saver) Read() ([]byte, error) {
	// 4K buffer
	data := make([]byte, 4096)
	n, err := s.file.Read(data)

	if err != nil {
		return nil, err
	}

	// Trimming the data.
	return data[:n], nil
}

func (s *Saver) Save(data []byte) error {
	s.file.Seek(0, 0)
	s.file.Truncate(0)

	_, err := s.file.WriteAt(data, 0)

	if err != nil {
		return err
	}

	return nil
}

func (s *Saver) Close() error {
	return s.file.Close()
}
