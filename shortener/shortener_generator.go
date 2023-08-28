package shortener

import (
	"crypto/sha256"
	"fmt"
	"math/big"

	"github.com/itchyny/base58-go"
)

func Sha256Of(originalUrl string) []byte {
	algorithm := sha256.New()
	algorithm.Write([]byte(originalUrl))
	hashedUrl := algorithm.Sum(nil)
	return hashedUrl
}

func Base58Encoder(hashedUrl []byte) string {
	encoding := base58.BitcoinEncoding
	encodedHash, err := encoding.Encode(hashedUrl)
	if err != nil {
		fmt.Printf("Something went wrong while encoding the hashedUrl to base58")
		return ""
	}

	return string(encodedHash)
}

func GenerateShortLink(originalUrl string) string {
	hashedUrl := Sha256Of(originalUrl)
	generatedNumber := new(big.Int).SetBytes(hashedUrl).Uint64()
	finalString := Base58Encoder([]byte(fmt.Sprintf("%d", generatedNumber)))
	return finalString[:8]
}
