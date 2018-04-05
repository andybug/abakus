package cmd

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"os"
	"path"

	"github.com/BurntSushi/toml"
	"golang.org/x/crypto/argon2"
	log "github.com/sirupsen/logrus"
)

type Config struct {
	Aws *AWSConfig `toml:"aws"`
	Gcs *GCSConfig `toml:"gcs"`
}

type AWSConfig struct {
	Bucket string `toml:"bucket"`
	Prefix string `toml:"prefix"`
	AccessKey string `toml:"-"`
	SecretKey string `toml:"-"`
	EncAccessKey string `toml:"accessKey"`
	EncSecretKey string `toml:"secretKey"`
}

type GCSConfig struct {
	Bucket string `toml:"bucket"`
	Prefix string `toml:"prefix"`
}

func (config *Config) write(path string, password string) {
	config.encrypt(password)
	
	f, err := os.OpenFile(path, os.O_WRONLY | os.O_CREATE | os.O_TRUNC, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	
	encoder := toml.NewEncoder(f)
	err = encoder.Encode(config)
	if err != nil {
		log.Fatal(err)
	}

	log.WithFields(log.Fields{
		"config": path,
	}).Debug("Wrote abakus configuration")
}

func (config *Config) encrypt(password string) {
	// encrypt aws configuration
	if config.Aws != nil {
		config.Aws.encrypt(password)
	}

	// encrypt gcs configuration
	if config.Gcs != nil {
		//config.Gcs.encrypt(password)
		log.Fatal("Unimplemented")
	}
}

func (aws *AWSConfig) encrypt(password string) {
	var salt = path.Join(aws.Bucket, aws.Prefix, "aws.accessKey")
	aws.EncAccessKey = encryptConfigString(password, salt, aws.AccessKey)

	salt = path.Join(aws.Bucket, aws.Prefix, "aws.secretKey")
	aws.EncSecretKey = encryptConfigString(password, salt, aws.SecretKey)
}

func encryptConfigString(password string, salt string, plaintext string) string {
	key := argon2.IDKey([]byte(password), []byte(salt), 1, 64*1024, 4, 32)

	c, err := aes.NewCipher(key)
	if err != nil {
		log.Fatal(err)
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		log.Fatal(err)
	}

	// since we don't reuse key, nonce doesn't really matter
	nonce := make([]byte, gcm.NonceSize())
	for i, _ := range nonce {
		nonce[i] = 0x56
	}

	encrypted := gcm.Seal(nil, nonce, []byte(plaintext), nil)

	//test, err := gcm.Open(nil, nonce, encrypted, nil)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//fmt.Printf("plz: %s\n", test)
	return base64.StdEncoding.EncodeToString(encrypted)
}
