package uuid62

import "testing"
import (
	. "gopkg.in/check.v1"
	"math/big"
	"math/rand"
	"github.com/google/uuid"
)


func Test(t *testing.T) { TestingT(t) }

type Uuid62Suite struct{}
var _ = Suite(&Uuid62Suite{})

var seed = int64(62)

type intTestVector struct {
	i *big.Int
	s string
	radix int
}

type uuidTestVector struct {
	id string
	i *big.Int
}

// Creates a new big.Int from the string s.
// Returns a pointer to the new big.Int and a bool indicating success
func newBigInt(s string) (*big.Int, bool) {
	i := big.Int{}
	return i.SetString(s, 10)
}

func mustNewBigInt(s string) *big.Int {
	i, _ := newBigInt(s)
	return i
}

var intTestVectors = []intTestVector{
	{i: mustNewBigInt("13"), s: "1101", radix: 2},
	{i: mustNewBigInt("59774123759"), s: "dead0beef", radix: 16},
	{i: mustNewBigInt("1"), s: "1", radix: 62},
	{i: mustNewBigInt("63"), s: "11", radix: 62},
	{i: mustNewBigInt("61"), s: "Z", radix: 62},
}

var uuidTestVectors = []uuidTestVector{
	{id: "00000000-0000-0000-0000-000000000000", i: mustNewBigInt("0")},
}

func (s *Uuid62Suite) TestUuid2Base62String(c *C) {
	for _, tv := range uuidTestVectors {
		id, _ := uuid.Parse(tv.id)
		base62, err := Uuuid2Base62String(id)
		c.Assert(err, IsNil)
		c.Assert(base62, Equals, tv.i.String())
	}
}

func (s *Uuid62Suite) TestBigInt2String(c *C) {
	for _, tv := range intTestVectors {
		base62, err := BigInt2String(tv.i, tv.radix)
		c.Assert(err, Equals, nil)
		c.Assert(base62, Equals, tv.s)
	}
}

func (s *Uuid62Suite) TestEdgeCaseBigInt2String(c *C) {
	// empty string
	result, err := BigInt2String(big.NewInt(int64(0)), 62)
	c.Assert(err, IsNil)
	c.Assert(result, Equals, "0")
}

func (s *Uuid62Suite) TestString2BigInt(c *C) {
	for _, tv := range intTestVectors {
		i, err := String2BigInt(tv.s, tv.radix)
		c.Assert(err, IsNil)
		c.Assert(i.Cmp(tv.i), Equals, 0)
	}
}

func (s *Uuid62Suite) TestEdgeCaseString2BigInt(c *C) {
	// empty string
	result, err := String2BigInt("", 62)
	c.Assert(err, IsNil)
	c.Assert(result.Int64(), Equals, int64(0))
}

func (s *Uuid62Suite) TestInvalidString2BigInt(c *C) {
	// empty string
	_, err1 := String2BigInt("a", 63)
	c.Assert(err1, ErrorMatches, "Radix must be.*")

	_, err2 := String2BigInt("a", 1)
	c.Assert(err2, ErrorMatches, "Radix must be.*")

	_, err3 := String2BigInt("-", 62)
	c.Assert(err3, ErrorMatches, "Digit.*")
}

// Do 1000 runs with random big ints, to test
// that identity holds when composing bi2str and str2bi
func (s *Uuid62Suite) TestRandomIdentity(c *C) {
	prng := rand.New(rand.NewSource(seed))
	two := big.NewInt(int64(2))
	max := two.Exp(two, big.NewInt(int64(129)), nil)
	runs := 1000
	bi := big.NewInt(int64(0))
	for i := 0; i < runs; i++ {
		bi.Rand(prng, max)
		s, err := BigInt2String(bi, 62)
		c.Assert(err, IsNil)
		result, err := String2BigInt(s, 62)
		c.Assert(err, IsNil)
		c.Assert(result.String(), Equals, bi.String())
	}
}