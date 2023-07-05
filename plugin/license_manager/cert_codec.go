package licensemanager

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"reflect"
	"unsafe"
)

func decodeCert(data []byte) (*Certificate, error) {
	// 解析证书
	block, _ := pem.Decode(data)
	if block == nil {
		return nil, fmt.Errorf("illegal certificate")
	}
	xcert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, err
	}

	key, ok := xcert.PublicKey.(*rsa.PublicKey)
	if !ok {
		return nil, ErrUnsupportedPublicKey
	}

	var cert = &Certificate{
		Raw:       b2s(data),
		Issuer:    xcert.Issuer,
		Subject:   xcert.Subject,
		NotBefore: xcert.NotBefore,
		NotAfter:  xcert.NotAfter,
		PublicKey: key,
	}
	return cert, nil
}

func b2s(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

func s2b(s string) (b []byte) {
	var sh = (*reflect.StringHeader)(unsafe.Pointer(&s))
	var bh = (*reflect.SliceHeader)(unsafe.Pointer(&b))
	bh.Len = sh.Len
	bh.Cap = sh.Len
	bh.Data = sh.Data
	return
}
