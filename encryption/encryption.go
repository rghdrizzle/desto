package encryption

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
	"strings"

	"golang.org/x/crypto/chacha20poly1305"
)

// Implementation of envelope encryption

// Interface to wrap any kind of KMS, including local. For now we will go with local only, in the future we might consider giving options to store in user's own KMS for more control
type KeyEncryptionWrapper interface {
    EncryptKey(ctx context.Context, keyPlain []byte) ([]byte, error)
    DecryptKey(ctx context.Context, keyCipher []byte) ([]byte, error)
}

type Enveloped struct {
    kms KeyEncryptionWrapper
}

type LocalKMS struct{
    Key []byte
}

func NewEnveloped(kms KeyEncryptionWrapper) *Enveloped {
    return &Enveloped{
        kms: kms,
    }
}


const KeySize = chacha20poly1305.KeySize


func GenerateNewKey()[]byte{
    b := make([]byte, KeySize)
    _, err := rand.Read(b)
    if err != nil {
        panic("unexpected failure: could not generate random data")
    }
    return b    
}

func DecryptDataWithDEK(key []byte,encryptedData []byte)([]byte,error){
    aead, err := chacha20poly1305.NewX(key)
    if err != nil {
        return nil, fmt.Errorf("failed to create encryption AEAD: %w", err)
    }
    if len(encryptedData) < aead.NonceSize() {
        return nil, fmt.Errorf("ciphertext is too short: %d < %d", len(encryptedData), aead.NonceSize())
    }

    // Split nonce and ciphertext.
    nonce, ciphertext := encryptedData[:aead.NonceSize()], encryptedData[aead.NonceSize():]

    // Decrypt the message and check it wasn't tampered with.
    plainData, err := aead.Open(nil, nonce, ciphertext, nil)
    if err != nil {
        return nil, fmt.Errorf("failed to decrypt: %w", err)
    }

    return plainData, nil
}

func EncryptDataWithDEK(key []byte, plainData []byte) ([]byte, error){
    aead, err := chacha20poly1305.NewX(key)
    if err != nil {
        return nil, fmt.Errorf("failed to create encryption AEAD: %w", err)
    }

    // Select a random nonce, and leave capacity for the ciphertext.
    nonce := make([]byte, aead.NonceSize(), aead.NonceSize()+len(plainData)+aead.Overhead())
    if _, err := rand.Read(nonce); err != nil {
        return nil, fmt.Errorf("failed to encrypt: %w", err)
    }
    // Encrypt the message and append the ciphertext to the nonce.
    encryptedMsg := aead.Seal(nonce, nonce, plainData, nil)

    return encryptedMsg, nil
}

func RetreiveEncryptedDEK(){

}

func (e *Enveloped) Encrypt(ctx context.Context, data []byte) (string, error){
    var err error
    dek := GenerateNewKey() // create dek key
    encryptedDek , err := e.kms.EncryptKey(ctx,dek) // encrypt the dek using KEK
    if err !=nil{
        return "", fmt.Errorf("could not encrypt dek: %w", err)
    }
    encryptedData , err:= EncryptDataWithDEK(dek,data) // encrypt the data with dek
    if err != nil {
        return "", fmt.Errorf("could not encrypt data: %w", err)
    }
    return fmt.Sprintf(
        "%s.%s",
        base64.StdEncoding.EncodeToString(encryptedDek), // only store encrypted dek with the data
        base64.StdEncoding.EncodeToString(encryptedData),
    ), nil



}

func (e *Enveloped) Decrypt(ctx context.Context, cipherData string) ([]byte, error){
    var err error
    partsDot := strings.SplitN(cipherData, ".", 3) // this is based on the format we created during Encrypt which splits the string into 2 parts with dot being the delimitter
    if len(partsDot) != 2 {
        return nil, fmt.Errorf("invalid cipher data blob")
    }
    encryptedDek, err:= base64.StdEncoding.DecodeString(partsDot[0])
    if err!=nil{
        log.Fatal("error in decoding dek")
    }
    encryptedData,err:= base64.StdEncoding.DecodeString(partsDot[1])
    if err!=nil{
        log.Fatal("error in decoding data")
    }
    if encryptedDek == nil || encryptedData == nil {
        return nil, fmt.Errorf("invalid cipher data blob")
    }

    dekKey, err := e.kms.DecryptKey(ctx, encryptedDek) // decrypt the key using KEK
    if err != nil {
        return nil, fmt.Errorf("could not decrypt DEK: %w", err)
    }
    plainData, err := DecryptDataWithDEK(
        dekKey,
        encryptedData,
    )
    if err != nil {
        return nil, fmt.Errorf("could not decrypt data with dek, data could be corrupted: %w", err)
    }
    return plainData, nil
}

func (l *LocalKMS) EncryptKey(ctx context.Context, keyPlain []byte) ([]byte, error){
    // Encrypt DEK with KEK

    aead, err := chacha20poly1305.NewX(l.Key)
    if err != nil {
        return nil, fmt.Errorf("failed to create encryption AEAD: %w", err)
    }

    // Select a random nonce, and leave capacity for the ciphertext.
    nonce := make([]byte, aead.NonceSize(), aead.NonceSize()+len(keyPlain)+aead.Overhead())
    if _, err := rand.Read(nonce); err != nil {
        return nil, fmt.Errorf("failed to encrypt: %w", err)
    }
    // Encrypt the key and append the ciphertext to the nonce.
    encryptedKey := aead.Seal(nonce, nonce, keyPlain, nil)

    return encryptedKey, nil

}

func (l *LocalKMS) DecryptKey(ctx context.Context, keyCipher []byte) ([]byte, error){
    aead, err := chacha20poly1305.NewX(l.Key)
    if err != nil {
        return nil, fmt.Errorf("failed to create encryption AEAD: %w", err)
    }
    if len(keyCipher) < aead.NonceSize() {
        return nil, fmt.Errorf("ciphertext is too short: %d < %d", len(keyCipher), aead.NonceSize())
    }

    // Split nonce and ciphertext.
    nonce, ciphertext := keyCipher[:aead.NonceSize()], keyCipher[aead.NonceSize():]

    // Decrypt the key and check it wasn't tampered with.
    dek, err := aead.Open(nil, nonce, ciphertext, nil)
    if err != nil {
        return nil, fmt.Errorf("failed to decrypt: %w", err)
    }
    return dek, nil
}

