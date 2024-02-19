package benaloh

import (
	"math/big"
)

func (b *Benaloh) GenKey(test Test) error {
	var checker func(*big.Int) error

	switch test {
	case TestSolovayStrassen:
		checker = SolovayStrassen
	case TestMillerRabin:
		checker = MillerRabin
	case TestFermat:
		checker = Fermat
	default:
		checker = SolovayStrassen
	}

	for {
		b.p = GenerateRandomValueWithChecker(checker)
		if b.p == nil {
			continue
		}

		b.q = GenerateRandomValueWithChecker(checker)
		if b.q == nil {
			continue
		}

		if new(big.Int).Mod(
			new(big.Int).Mul(
				new(big.Int).Sub(b.p, big.NewInt(1)),
				new(big.Int).Sub(b.q, big.NewInt(1)),
			),
			big.NewInt(int64(b.blockLen)),
		).Cmp(big.NewInt(0)) != 0 {
			continue
		}

		tempP := new(big.Int).Div(new(big.Int).Sub(b.p, big.NewInt(1)), big.NewInt(int64(b.blockLen)))

		if new(big.Int).GCD(nil, nil, tempP, big.NewInt(int64(b.blockLen))).Cmp(big.NewInt(1)) == 0 {
			if big.NewInt(int64(b.blockLen)).Bit(0) == 0 ||
				new(big.Int).GCD(nil, nil, new(big.Int).Sub(b.q, big.NewInt(1)), big.NewInt(int64(b.blockLen))).Cmp(big.NewInt(1)) == 0 {
				break
			}
		}
	}

	n := new(big.Int).Mul(b.p, b.q)
	phi := new(big.Int).Mul(new(big.Int).Sub(b.p, big.NewInt(1)), new(big.Int).Sub(b.q, big.NewInt(1)))

	y, err := GenerateYValue(phi, n, big.NewInt(int64(b.blockLen)))
	if err != nil {
		return err
	}

	x := new(big.Int).Exp(y, new(big.Int).Div(phi, big.NewInt(int64(b.blockLen))), n)

	b.Public = []*big.Int{y, big.NewInt(int64(b.blockLen)), n}
	b.Private = []*big.Int{phi, x}

	return err
}
