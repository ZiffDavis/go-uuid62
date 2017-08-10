package uuid62

import (
	"math/big"
	"errors"
	"strings"
	"github.com/google/uuid"
)

var alphabet = string("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func Uuuid2Base62String(id uuid.UUID) (string, error) {
	idBytes, err := id.MarshalBinary()
	if err != nil { return "", err}
	integer := big.NewInt(int64(0))
	integer.SetBytes(idBytes)
	return BigInt2String(integer, 62)
}

func Base62String2Uuid(s string) (uuid.UUID, error) {
	integer, err := String2BigInt(s, 62)
	if err != nil { return uuid.UUID{}, nil}

	integerBytes := integer.Bytes()
	return uuid.ParseBytes(integerBytes)
}

func BigInt2String(n *big.Int, radix int) (string, error) {
	if(radix < 2 || radix > 62) {
		return "", errors.New("Radix must be between 2-62 inclusive")
	}

	var resultBytes []byte
	zero := big.NewInt(int64(0))
	base := big.NewInt(int64(radix))
	r := big.NewInt(int64(0))

	// create new i from n to avoid mutating n
	i := big.NewInt(int64(0))
	i.Add(i, n)

	for ; 1 == i.Cmp(zero);  {
		i.QuoRem(i, base, r)
		resultBytes = append([]byte{alphabet[r.Int64()]}, resultBytes...)
	}
	return string(resultBytes[:]), nil
}

func String2BigInt(s string, radix int) (*big.Int, error) {
	zero := big.NewInt(int64(0))
	if(radix < 2 || radix > 62) {
		return zero, errors.New("Radix must be between 2-62 inclusive")
	}

	acc := big.NewInt(int64(0))
	length := len(s)
	for i := 0; i < length; i++ {
		place := length - i - 1
		// acc += radix ** place
		base := big.NewInt(int64(radix))
		power := big.NewInt(int64(place))
		digitChar := s[i]
		digit := strings.IndexByte(alphabet, digitChar)

		if(digit < 0 || digit > radix - 1) {
			return zero, errors.New("Digit in string is outside specified radix's range")
		}

		value := big.NewInt(int64(0))
		value.Exp(base, power, nil)
		value.Mul(value, big.NewInt(int64(digit)))
		acc.Add(acc, value)
	}
	return acc, nil
}

