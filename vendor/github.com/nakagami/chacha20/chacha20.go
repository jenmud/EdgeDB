package chacha20

import (
	"crypto/cipher"
	"encoding/binary"
	"errors"
	"math/bits"
)

type Cipher struct {
	constant  [4]uint32
	key       [8]uint32
	counter   uint64
	nonce     []uint32
	block     [64]byte
	block_pos int
}

var _ cipher.Stream = (*Cipher)(nil)

func NewCipher(key []byte, nonce []byte, count uint64) (*Cipher, error) {
	if len(key) != 32 {
		err := errors.New("key must be 32 bytes")
		return nil, err
	}
	if len(nonce) != 12 && len(nonce) != 8 {
		err := errors.New("nonce must be 12 bytes or 8 bytes")
		return nil, err
	}

	c := new(Cipher)
	c.constant = [4]uint32{0x61707865, 0x3320646e, 0x79622d32, 0x6b206574}
	c.key = [8]uint32{
		binary.LittleEndian.Uint32(key[0:4]),
		binary.LittleEndian.Uint32(key[4:8]),
		binary.LittleEndian.Uint32(key[8:12]),
		binary.LittleEndian.Uint32(key[12:16]),
		binary.LittleEndian.Uint32(key[16:20]),
		binary.LittleEndian.Uint32(key[20:24]),
		binary.LittleEndian.Uint32(key[24:28]),
		binary.LittleEndian.Uint32(key[28:32]),
	}
	c.counter = count
	c.nonce = make([]uint32, len(nonce)/4)
	c.nonce[0] = binary.LittleEndian.Uint32(nonce[0:4])
	c.nonce[1] = binary.LittleEndian.Uint32(nonce[4:8])

	if len(nonce) == 12 {
		c.nonce[2] = binary.LittleEndian.Uint32(nonce[8:12])
	}

	c.setChaCha20RoundBlock()

	return c, nil
}

func (c *Cipher) toState() [16]uint32 {
	var cn1, cn2, cn3, cn4 uint32
	if len(c.nonce) == 3 {
		cn1 = uint32(c.counter)
		cn2 = c.nonce[0]
		cn3 = c.nonce[1]
		cn4 = c.nonce[2]
	} else {
		cn1 = uint32(c.counter)
		cn2 = uint32(c.counter >> 32)
		cn3 = c.nonce[0]
		cn4 = c.nonce[1]
	}

	return [16]uint32{
		c.constant[0], c.constant[1], c.constant[2], c.constant[3],
		c.key[0], c.key[1], c.key[2], c.key[3],
		c.key[4], c.key[5], c.key[6], c.key[7],
		cn1, cn2, cn3, cn4,
	}
}

func (c *Cipher) XORKeyStream(dst, src []byte) {
	for i := range len(src) {
		dst[i] = src[i] ^ c.block[c.block_pos]
		c.block_pos++
		if len(c.block) == c.block_pos {
			c.counter++
			c.setChaCha20RoundBlock()
		}
	}
}

func (c *Cipher) setChaCha20RoundBlock() {
	c.block = c.keyStream()
	c.block_pos = 0
}

func (c *Cipher) keyStream() [64]byte {
	x := c.toState()
	for i := 0; i < 10; i++ {
		// column round
		x[0], x[4], x[8], x[12] = qr(x[0], x[4], x[8], x[12])
		x[1], x[5], x[9], x[13] = qr(x[1], x[5], x[9], x[13])
		x[2], x[6], x[10], x[14] = qr(x[2], x[6], x[10], x[14])
		x[3], x[7], x[11], x[15] = qr(x[3], x[7], x[11], x[15])
		// diagonal round
		x[0], x[5], x[10], x[15] = qr(x[0], x[5], x[10], x[15])
		x[1], x[6], x[11], x[12] = qr(x[1], x[6], x[11], x[12])
		x[2], x[7], x[8], x[13] = qr(x[2], x[7], x[8], x[13])
		x[3], x[4], x[9], x[14] = qr(x[3], x[4], x[9], x[14])
	}
	initial := c.toState()
	for i := range x {
		x[i] += initial[i]
	}
	var stream [64]byte
	for i, v := range x {
		stream[i*4] = byte(v)
		stream[i*4+1] = byte(v >> 8)
		stream[i*4+2] = byte(v >> 16)
		stream[i*4+3] = byte(v >> 24)
	}
	return stream
}

func qr(a, b, c, d uint32) (uint32, uint32, uint32, uint32) {
	a += b
	d ^= a
	d = bits.RotateLeft32(d, 16)
	c += d
	b ^= c
	b = bits.RotateLeft32(b, 12)
	a += b
	d ^= a
	d = bits.RotateLeft32(d, 8)
	c += d
	b ^= c
	b = bits.RotateLeft32(b, 7)
	return a, b, c, d
}
