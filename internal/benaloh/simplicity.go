package benaloh

import (
	"math/big"
	"math/rand"
)

type Test int

const (
	TestSolovayStrassen Test = iota
	TestMillerRabin
	TestFermat
)

func SolovayStrassen(number *big.Int) error {
	min := big.NewInt(2)
	max := new(big.Int).Sub(number, big.NewInt(2))

	for i := 0; i < 100; i++ {
		random_number, err := RandomIntInRange(min, max)
		if err != nil {
			return err
		}

		rightDiv := ExpMod(random_number, new(big.Int).Sub(number, big.NewInt(1)), number)
		leftExp := new(big.Int).Exp(random_number, new(big.Int).Div(
			new(big.Int).Sub(number, big.NewInt(1)),
			big.NewInt(2),
		), number)

		if EuclidAlgorithm(random_number, number).Cmp(big.NewInt(1)) > 0 {
			return ErrBadNumber
		}

		if leftExp.Cmp(rightDiv) != 0 {
			return ErrBadNumber
		}
	}

	return nil
}

func MillerRabin(number *big.Int) error {
	if number.Bit(0) == 0 {
		return ErrBadNumber
	}

	t := new(big.Int).Sub(number, big.NewInt(1))
	exp := big.NewInt(0)

	for t.Bit(0) == 0 {
		t.Rsh(t, 1)
		exp.Add(exp, big.NewInt(1))
	}

	for i := 0; i < 100; i++ {
		temp := new(big.Int).Sub(number, big.NewInt(3))

		a := new(big.Int).Add(big.NewInt(2), new(big.Int).SetUint64(uint64(rand.Intn(int(temp.Int64())))))

		x := new(big.Int).Exp(a, t, number)

		if x.Cmp(big.NewInt(1)) != 0 && x.Cmp(new(big.Int).Sub(number, big.NewInt(1))) != 0 {
			for j := big.NewInt(0); j.Cmp(new(big.Int).Sub(exp, big.NewInt(1))) != 0; j.Add(j, big.NewInt(1)) {
				x = x.Exp(x, big.NewInt(2), number)

				if x.Cmp(new(big.Int).Sub(number, big.NewInt(1))) == 0 {
					break
				}

				if x.Cmp(big.NewInt(1)) == 0 {
					return ErrBadNumber
				}
			}

			if x.Cmp(new(big.Int).Sub(number, big.NewInt(1))) != 0 {
				return ErrBadNumber
			}
		}
	}

	return nil
}

func Fermat(number *big.Int) error {
	if number.Bit(0) == 0 {
		return ErrBadNumber
	}

	for i := 0; i < 100; i++ {
		temp := new(big.Int).Sub(number, big.NewInt(3))

		a := new(big.Int).Add(big.NewInt(2), new(big.Int).SetUint64(uint64(rand.Intn(int(temp.Int64())))))

		if a.Exp(a, new(big.Int).Sub(number, big.NewInt(1)), number).Cmp(big.NewInt(1)) != 0 {
			return ErrBadNumber
		}
	}

	return nil
}
