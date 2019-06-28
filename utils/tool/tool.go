package tool

import (
	"ark-common/constants"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/ssh"
)

var iv = []byte{
	0x00, 0x01, 0x02, 0x03,
	0x04, 0x05, 0x06, 0x07,
	0x08, 0x09, 0x0a, 0x0b,
	0x0c, 0x0d, 0x0e, 0x0f,
}

var mk = []string{
	"T", "M", "0", "N", "A", "j", "y", "B", "C", "K",
	"4", "t", "b", "c", "I", "L", "e", "u", "5", "t",
	"i", "0", "l", "Q", "A", "i", "t", "C", "w", "3",
	"F", "t", "b", "J", "8", "Y", "9", "K", "Y", "l",
	"R", "M", "a", "E", "x", "T", "o", "p", "q", "8",
	"y", "d", "o", "2", "x", "h", "B", "4", "K", "s",
	"Q", "0", "R", "q", "e", "F", "7", "W", "Q", "R",
	"l", "C", "t", "v", "o", "v", "S", "D", "V", "F",
	"3", "h", "X", "p", "s", "W", "R", "s", "l", "v",
	"P", "a", "/", "l", "V", "X", "X", "l", "D", "G",
	"u", "9", "n", "9", "S", "W", "J", "S", "0", "h",
	"O", "I", "+", "/", "+", "v", "w", "2", "T", "Z",
	"2", "W", "1", "8", "R", "b", "T", "S"}[23:55]

var mi = []string{
	"K", "1", "H", "4", "9", "p", "4", "/", "9", "1",
	"m", "o", "0", "N", "0", "=", "c", "e", "3", "k",
	"o", "o", "g", "3", "c", "s", "u", "W", "Z", "9",
	"B", "3", "V", "R", "o", "b", "1", "+", "6", "o",
	"1", "m", "r", "5", "g", "b", "O", "f", "u", "1",
	"W", "7", "t", "H", "A", "z", "w", "o", "H", "P",
	"d", "X", "7", "Y", "U", "o", "L", "7", "M", "M",
	"S", "S", "P", "8", "f", "Z", "w", "2", "u", "G",
}[21:37]

// UUID 返回一个UUID
func UUID() string {
	return primitive.NewObjectID().Hex()
}

// TimeForISO8601 转换为time.Time类型
func TimeForISO8601(t string) (rt time.Time) {
	rt, _ = time.Parse(constants.ISO8601, t)
	return rt
}

// ECP 字符串加密
func ECP(data string) string {
	key := []byte(strings.Join(mk, ""))
	iv := []byte(strings.Join(mi, ""))
	block, err := aes.NewCipher(key)
	if err != nil {
		return ""
	}
	encrypter := cipher.NewCFBEncrypter(block, iv)
	encrypted := make([]byte, len(data))
	encrypter.XORKeyStream(encrypted, []byte(data))
	return hex.EncodeToString(encrypted)
}

// DCP 字符串解密
func DCP(data string) string {
	var err error
	defer func() {
		if e := recover(); e != nil {
			err = e.(error)
		}
	}()
	src, err := hex.DecodeString(data)
	if err != nil {
		return ""
	}

	key := []byte(strings.Join(mk, ""))
	iv := []byte(strings.Join(mi, ""))

	decrypted := make([]byte, len(src))
	block, err := aes.NewCipher(key)
	if err != nil {
		return ""
	}
	decrypter := cipher.NewCFBDecrypter(block, iv)
	decrypter.XORKeyStream(decrypted, src)
	return string(decrypted)
}

// NewRSAKeyPair 生成一对RAS密钥对
func NewRSAKeyPair() (privateKey, publicKey []byte) {
	// 生成私钥
	pk, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return
	}
	pb := pem.Block{
		Type:    "RSA PRIVATE KEY",
		Headers: nil,
		Bytes:   x509.MarshalPKCS1PrivateKey(pk),
	}
	privateKey = pem.EncodeToMemory(&pb)
	// 这里必须传入指针
	pubKey, err := ssh.NewPublicKey(&pk.PublicKey)
	if err != nil {
		return nil, nil
	}
	publicKey = ssh.MarshalAuthorizedKey(pubKey)
	return
}
