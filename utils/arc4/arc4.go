package arc4

import (
	"crypto/md5"
	"crypto/rc4"
	"encoding/hex"
)

func XORKeyStream(eamuseKey string, data []byte) error {
	eamuseKeyBytes, err := hex.DecodeString(eamuseKey)
	if err != nil {
		return err
	}

	key := append(eamuseKeyBytes[:6], []byte{0x69, 0xd7, 0x46, 0x27, 0xd9, 0x85, 0xee, 0x21, 0x87, 0x16, 0x15, 0x70, 0xd0, 0x8d, 0x93, 0xb1, 0x24, 0x55, 0x03, 0x5b, 0x6d, 0xf0, 0xd8, 0x20, 0x5d, 0xf5}...)
	h := md5.Sum(key)

	cipher, err := rc4.NewCipher(h[:])
	if err != nil {
		return err
	}

	cipher.XORKeyStream(data, data)

	return nil
}
