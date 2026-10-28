package crypto

import (
	"crypto/x509"
	"encoding/pem"
	"testing"
)

// These public certificate fixtures come from upstream MicroMDM. The SHA-256
// fixture was supplied by an Apple TV during ABM enrollment (upstream #989).
func TestVerifyFromAppleDeviceCA(t *testing.T) {
	for _, name := range []string{"appleca_signed.pem", "appleca_signed_newer.pem"} {
		t.Run(name, func(t *testing.T) {
			cert, err := ReadPEMCertificateFile("testdata/" + name)
			if err != nil {
				t.Fatal(err)
			}
			if err := VerifyFromAppleDeviceCA(cert); err != nil {
				t.Fatalf("valid Apple signature (%s): %v", cert.SignatureAlgorithm, err)
			}
			t.Run("tampered signature", func(t *testing.T) {
				modified := *cert
				modified.Signature = append([]byte(nil), cert.Signature...)
				modified.Signature[0] ^= 1
				if VerifyFromAppleDeviceCA(&modified) == nil {
					t.Fatal("accepted a tampered signature")
				}
			})
			t.Run("tampered content", func(t *testing.T) {
				modified := *cert
				modified.RawTBSCertificate = append([]byte(nil), cert.RawTBSCertificate...)
				modified.RawTBSCertificate[len(modified.RawTBSCertificate)-1] ^= 1
				if VerifyFromAppleDeviceCA(&modified) == nil {
					t.Fatal("accepted modified certificate content")
				}
			})
		})
	}
}

func TestVerifyFromAppleDeviceCARejectsWrongSigner(t *testing.T) {
	block, _ := pem.Decode([]byte(appleiPhoneDeviceCAPEM))
	if block == nil {
		t.Fatal("invalid CA fixture")
	}
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		t.Fatal(err)
	}
	// The Device CA is signed by a different CA, not by its own public key.
	if VerifyFromAppleDeviceCA(cert) == nil {
		t.Fatal("accepted a certificate signed by a different key")
	}
}

func TestVerifyFromAppleDeviceCARejectsUnsupportedAlgorithm(t *testing.T) {
	cert, err := ReadPEMCertificateFile("testdata/appleca_signed.pem")
	if err != nil {
		t.Fatal(err)
	}
	cert.SignatureAlgorithm = x509.MD5WithRSA
	if VerifyFromAppleDeviceCA(cert) == nil {
		t.Fatal("accepted an unsupported signature algorithm")
	}
}
