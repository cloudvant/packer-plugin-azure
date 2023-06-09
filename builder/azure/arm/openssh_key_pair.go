// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package arm

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"time"

	"golang.org/x/crypto/ssh"
)

const (
	KeySize = 2048
)

type OpenSshKeyPair struct {
	privateKey *rsa.PrivateKey
	publicKey  ssh.PublicKey
}

func NewOpenSshKeyPair() (*OpenSshKeyPair, error) {
	return NewOpenSshKeyPairWithSize(KeySize)
}

func NewOpenSshKeyPairWithSize(keySize int) (*OpenSshKeyPair, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, keySize)
	if err != nil {
		return nil, err
	}

	publicKey, err := ssh.NewPublicKey(&privateKey.PublicKey)
	if err != nil {
		return nil, err
	}

	return &OpenSshKeyPair{
		privateKey: privateKey,
		publicKey:  publicKey,
	}, nil
}

func (s *OpenSshKeyPair) AuthorizedKey() string {
	return fmt.Sprintf("%s %s packer Azure Deployment%s",
		s.publicKey.Type(),
		base64.StdEncoding.EncodeToString(s.publicKey.Marshal()),
		time.Now().Format(time.RFC3339))
}

func (s *OpenSshKeyPair) PrivateKey() []byte {
	privateKey := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(s.privateKey),
	})

	return privateKey
}
