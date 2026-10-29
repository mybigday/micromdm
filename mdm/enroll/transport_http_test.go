package enroll

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"errors"
	"math/big"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/fullsailor/pkcs7"
)

func TestDecodeMDMEnrollRequestGET(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/mdm/enroll", nil)
	result, err := decodeMDMEnrollRequest(context.Background(), req)
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := result.(mdmEnrollRequest); !ok {
		t.Fatalf("unexpected request type %T", result)
	}
}

func TestDecodeMDMEnrollRequestRejectsInvalidCMS(t *testing.T) {
	for _, body := range []string{"", "not CMS", "<plist><dict/></plist>"} {
		req := httptest.NewRequest(http.MethodPost, "/mdm/enroll", strings.NewReader(body))
		if _, err := decodeMDMEnrollRequest(context.Background(), req); err == nil {
			t.Errorf("accepted invalid CMS: %q", body)
		}
	}
}

func TestDecodeMDMEnrollRequestRejectsNonAppleSigner(t *testing.T) {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal(err)
	}
	template := &x509.Certificate{
		SerialNumber:       big.NewInt(1),
		Subject:            pkix.Name{CommonName: "untrusted enrollment test"},
		NotBefore:          time.Now().Add(-time.Hour),
		NotAfter:           time.Now().Add(time.Hour),
		KeyUsage:           x509.KeyUsageDigitalSignature,
		SignatureAlgorithm: x509.SHA256WithRSA,
	}
	der, err := x509.CreateCertificate(rand.Reader, template, template, &key.PublicKey, key)
	if err != nil {
		t.Fatal(err)
	}
	cert, err := x509.ParseCertificate(der)
	if err != nil {
		t.Fatal(err)
	}
	signed, err := pkcs7.NewSignedData([]byte(`<?xml version="1.0"?><plist version="1.0"><dict><key>SERIAL</key><string>TEST-SERIAL</string></dict></plist>`))
	if err != nil {
		t.Fatal(err)
	}
	if err := signed.AddSigner(cert, key, pkcs7.SignerInfoConfig{}); err != nil {
		t.Fatal(err)
	}
	body, err := signed.Finish()
	if err != nil {
		t.Fatal(err)
	}
	// The CMS is valid: it must be rejected specifically because its signer is
	// not signed by the pinned Apple Device CA, even if it claims a serial.
	parsed, err := pkcs7.Parse(body)
	if err != nil {
		t.Fatal(err)
	}
	if err := parsed.Verify(); err != nil {
		t.Fatalf("test CMS must have a valid signature: %v", err)
	}
	req := httptest.NewRequest(http.MethodPost, "/mdm/enroll", bytes.NewReader(body))
	_, err = decodeMDMEnrollRequest(context.Background(), req)
	if err == nil {
		t.Fatal("accepted enrollment signed by an untrusted key")
	}
	for _, want := range []string{"not signed by Apple Device CA", "signature_algorithm=SHA256-RSA"} {
		if !strings.Contains(err.Error(), want) {
			t.Errorf("missing diagnostic %q: %v", want, err)
		}
	}
	if !errors.Is(err, rsa.ErrVerification) {
		t.Errorf("underlying verification error was lost: %v", err)
	}
}
