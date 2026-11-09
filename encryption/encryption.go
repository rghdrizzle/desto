package encryption

import (
    "crypto/rand"

    "golang.org/x/crypto/chacha20poly1305"
	"context"
)

// Implementation of envelope encryption

// Interface to wrap any kind of KMS, including local. For now we will go with local
type KeyEncryptionWrapper interface {
    EncryptKey(ctx context.Context, keyPlain []byte) ([]byte, error)
    DecryptKey(ctx context.Context, keyCipher []byte) ([]byte, error)
}

type Enveloped struct {
    kms KeyEncryptionWrapper
}

func NewEnveloped(kms KeyEncryptionWrapper) *Enveloped {
    return &Enveloped{
        kms: kms,
    }
}

func GenerateNewKey(){

}

func DecryptDataWithDEK(){

}

func EncryptDataWithDEK(){

}

func RetreiveEncryptedDEK(){

}

func EncryptKey(ctx context.Context, keyPlain []byte) ([]byte, error){

}

func DecryptKey(ctx context.Context, keyCipher []byte) ([]byte, error){

}