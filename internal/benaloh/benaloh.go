package benaloh

import (
	"crypto/rand"
	"errors"
	"math/big"
)

var (
	ErrBadNumber = errors.New("bad number")
)

type Benaloh struct {
	Public, Private []*big.Int
	p, q            *big.Int
	blockLen        int
}

func NewBenaloh() *Benaloh {
	return &Benaloh{
		Public:   make([]*big.Int, 3),
		Private:  make([]*big.Int, 2),
		blockLen: 257,
	}
}

func (b *Benaloh) Encrypt(m *big.Int) (*big.Int, error) {
	y, r, n := b.Public[0], b.Public[1], b.Public[2]

	u, err := rand.Int(rand.Reader, new(big.Int).Sub(n, big.NewInt(2)))
	if err != nil {
		return nil, err
	}

	u.Add(u, big.NewInt(2))

	c := new(big.Int).Mul(
		new(big.Int).Exp(y, m, nil),
		new(big.Int).Exp(u, r, nil),
	)

	c.Mod(c, n)

	return c, nil
}

func (b *Benaloh) Decrypt(c *big.Int) (*big.Int, error) {
	phi, x := b.Private[0], b.Private[1]
	r, n := b.Public[1], b.Public[2]

	a := new(big.Int).Exp(c, new(big.Int).Div(phi, r), n)
	m := Logarithm(a, x, r, n)

	return m, nil
}
