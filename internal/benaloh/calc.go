package benaloh

import (
	"crypto/rand"
	"math/big"
)

func RandomIntInRange(min, max *big.Int) (*big.Int, error) {
	diff := new(big.Int).Sub(max, min)
	diff.Add(diff, big.NewInt(1))

	result, err := rand.Int(rand.Reader, diff)
	if err != nil {
		return nil, err
	}

	return result.Add(result, min), nil
}

func GenerateRandomValueWithChecker(checker func(*big.Int) error) *big.Int {
	min, _ := new(big.Int).SetString("10000000000000000001000000000000000000", 10)
	max, _ := new(big.Int).SetString("92233720368547758079223372036854775807", 10)

	for {
		val, err := RandomIntInRange(min, max)
		if err != nil {
			continue
		}

		if checker(val) != nil {
			continue
		}

		return val
	}
}

func Factorize(n *big.Int) []*big.Int {
	var factors []*big.Int

	for new(big.Int).Mod(n, big.NewInt(2)).Cmp(big.NewInt(0)) == 0 {
		factors = append(factors, big.NewInt(2))
		n.Div(n, big.NewInt(2))
	}

	t := big.NewInt(3)

	for new(big.Int).Sqrt(n).Cmp(t) >= 0 {
		for new(big.Int).Mod(n, big.NewInt(2)).Cmp(big.NewInt(0)) == 0 {
			factors = append(factors, new(big.Int).Set(t))
			n.Div(n, t)
		}

		t.Add(t, big.NewInt(2))
	}

	if n.Cmp(big.NewInt(2)) > 0 {
		factors = append(factors, new(big.Int).Set(n))
	}

	return factors
}

func GenerateYValue(phi, n *big.Int, r *big.Int) (*big.Int, error) {
	for {
		y, err := rand.Int(rand.Reader, new(big.Int).Sub(n, big.NewInt(2)))
		if err != nil {
			return nil, err
		}

		y.Add(y, big.NewInt(2))

		if MillerRabin(r) != nil {
			var tmp []*big.Int
			factor := Factorize(y)

			for _, i := range factor {
				exp := new(big.Int).Div(phi, r)
				x := new(big.Int).Exp(i, exp, n)

				if x.Cmp(big.NewInt(1)) != 0 {
					tmp = append(tmp, i)
				} else {
					break
				}

				if len(tmp) == len(factor) {
					return y, nil
				}
			}
		} else {
			exp := new(big.Int).Div(phi, r)
			x := new(big.Int).Exp(y, exp, n)

			if x.Cmp(big.NewInt(1)) != 0 {
				return y, nil
			}
		}
	}
}

func EuclidAlgorithm(a, b *big.Int) *big.Int {
	zero := big.NewInt(0)

	for b.Cmp(zero) != 0 {
		a, b = b, new(big.Int).Mod(a, b)
	}

	return a
}

func Logarithm(a, x, r, n *big.Int) *big.Int {
	for m := big.NewInt(0); m.Cmp(r) < 0; m.Add(m, big.NewInt(1)) {
		if new(big.Int).Exp(x, m, n).Cmp(a) == 0 {
			return m
		}
	}

	return nil
}

func ExpMod(b, n, m *big.Int) *big.Int {
	result := big.NewInt(1)

	for n.Cmp(big.NewInt(0)) > 0 {
		if n.Bit(0) == 1 {
			result.Mul(result, b)
			result.Mod(result, m)
		}

		b.Mul(b, b).Mod(b, m)
		n.Rsh(n, 1)
	}

	return result
}
