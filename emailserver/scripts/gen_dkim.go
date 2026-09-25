package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"flag"
	"fmt"
	"os"
)

func main() {
	domain := flag.String("domain", "atonixcorp.com", "Domain for DKIM")
	selector := flag.String("selector", "atonix", "DKIM selector")
	out := flag.String("out", "dkim_private.pem", "Output private key file")
	flag.Parse()

	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		fmt.Println("Error generating key:", err)
		os.Exit(1)
	}

	// write private key
	f, err := os.Create(*out)
	if err != nil {
		fmt.Println("Error creating file:", err)
		os.Exit(1)
	}
	defer f.Close()

	block := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(priv),
	}
	if err := pem.Encode(f, block); err != nil {
		fmt.Println("Error writing key:", err)
		os.Exit(1)
	}

	// build public key for DNS
	pubBytes, err := x509.MarshalPKIXPublicKey(&priv.PublicKey)
	if err != nil {
		fmt.Println("Error marshaling public key:", err)
		os.Exit(1)
	}

	fmt.Println("DKIM private key written to:", *out)
	fmt.Println()
	fmt.Println("Add this DNS TXT record:")
	fmt.Printf("%s._domainkey.%s. IN TXT \"v=DKIM1; k=rsa; p=%s\"\n",
		*selector, *domain, base64NoWrap(pubBytes))
}

func base64NoWrap(b []byte) string {
	const table = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/"
	var out []byte
	var val uint32
	var valb int
	for _, c := range b {
		val = (val << 8) | uint32(c)
		valb += 8
		for valb >= 6 {
			out = append(out, table[(val>>(valb-6))&0x3F])
			valb -= 6
		}
	}
	if valb > 0 {
		out = append(out, table[(val<<(6-valb))&0x3F])
	}
	for len(out)%4 != 0 {
		out = append(out, '=')
	}
	return string(out)
}
