package storage

import (
	"bytes"
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
)

type Backend interface {
	Save(ctx context.Context, name string, r io.Reader) (string, error)
}

type LocalStorage struct {
	BaseDir string
}

func (l *LocalStorage) Save(_ context.Context, name string, r io.Reader) (string, error) {
	if err := os.MkdirAll(l.BaseDir, 0o755); err != nil {
		return "", err
	}
	path := filepath.Join(l.BaseDir, filepath.Base(name))
	f, err := os.Create(path)
	if err != nil {
		return "", err
	}
	defer f.Close()
	_, err = io.Copy(f, r)
	return path, err
}

type Service struct {
	Backend Backend
	MaxSize int64
	AESKey  []byte
}

func (s *Service) SaveMultipart(ctx context.Context, fh *multipart.FileHeader) (string, error) {
	if fh.Size > s.MaxSize {
		return "", fmt.Errorf("file too large: %d > %d", fh.Size, s.MaxSize)
	}
	f, err := fh.Open()
	if err != nil {
		return "", err
	}
	defer f.Close()
	if len(s.AESKey) == 0 {
		return s.Backend.Save(ctx, fh.Filename, f)
	}
	enc, err := encryptStream(f, s.AESKey)
	if err != nil {
		return "", err
	}
	return s.Backend.Save(ctx, fh.Filename+".enc", enc)
}

func encryptStream(r io.Reader, key []byte) (io.Reader, error) {
	if len(key) != 16 && len(key) != 24 && len(key) != 32 {
		return nil, errors.New("invalid AES key length")
	}
	plain, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return nil, err
	}
	sealed := gcm.Seal(nonce, nonce, plain, nil)
	return bytes.NewReader(sealed), nil
}
