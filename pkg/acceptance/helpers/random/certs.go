package random

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"math/big"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// Generate X509 returns base64 encoded certificate on a single line without the leading -----BEGIN CERTIFICATE----- and ending -----END CERTIFICATE----- markers.
func GenerateX509(t *testing.T) string {
	t.Helper()
	ca := &x509.Certificate{
		SerialNumber: big.NewInt(1658),
		Subject: pkix.Name{
			Organization: []string{"Company, INC."},
		},
		NotAfter:    time.Now().AddDate(10, 0, 0),
		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
		KeyUsage:    x509.KeyUsageDigitalSignature,
	}

	caPrivKey, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)

	caBytes, err := x509.CreateCertificate(rand.Reader, ca, ca, &caPrivKey.PublicKey, caPrivKey)
	require.NoError(t, err)

	return encode(t, "CERTIFICATE", caBytes)
}

// GenerateRSA returns an RSA public key without BEGIN and END markers, and key's hash.
func GenerateRSAPublicKey(t *testing.T) (string, string) {
	t.Helper()
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)

	pub := key.Public()
	b, err := x509.MarshalPKIXPublicKey(pub.(*rsa.PublicKey))
	require.NoError(t, err)
	return encode(t, "RSA PUBLIC KEY", b), hash(t, b)
}

func hash(t *testing.T, b []byte) string {
	t.Helper()
	hash := sha256.Sum256(b)
	return base64.StdEncoding.EncodeToString(hash[:])
}

func encode(t *testing.T, pemType string, b []byte) string {
	t.Helper()
	buffer := new(bytes.Buffer)
	err := pem.Encode(buffer,
		&pem.Block{
			Type:  pemType,
			Bytes: b,
		},
	)
	require.NoError(t, err)
	cert := strings.TrimPrefix(buffer.String(), fmt.Sprintf("-----BEGIN %s-----\n", pemType))
	cert = strings.TrimSuffix(cert, fmt.Sprintf("-----END %s-----\n", pemType))
	return cert
}
