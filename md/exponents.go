package md

import "strconv"

type Exponents struct {
	storage map[string]map[string]int
}

func (e *Exponents) Get(provider, pair string) (int, bool) {
	providerPairsExponents, ok := e.storage[provider]
	if !ok {
		return -8, ok
	}

	pairExponents, ok := providerPairsExponents[pair]
	return pairExponents, ok
}

func (e *Exponents) Set(provider string, pair string, exp int) {
	pairsExponents, ok := e.storage[provider]
	if !ok || pairsExponents == nil {
		pairsExponents = make(map[string]int, 100)
		pairsExponents[pair] = exp
		e.storage[provider] = pairsExponents
	} else {
		e.storage[provider][pair] = exp
	}
}

func (e *Exponents) SetString(provider, pair, exp string) error {
	intExp, err := strconv.Atoi(exp[2:])
	if err != nil {
		return err
	}

	e.Set(provider, pair, intExp)

	return nil
}

func NewExponents() *Exponents {
	return &Exponents{
		storage: make(map[string]map[string]int, 5),
	}
}


