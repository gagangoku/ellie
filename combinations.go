package main

import "gonum.org/v1/gonum/stat/combin"

// Combination generator wrapper over combin, but allows the first empty combination.
// Useful so you dont have to write complex code to handle the initial case where there are no valid combinations
// Usage pattern:
// c, comb := NewCombinator(lens, true)
//
//	for c.Next() {
//	    c.Product(comb)
//	    ...
//	}
//
// NOTE: Call Product() only if Next() allows.

type Combinator struct {
	cgen           *combin.CartesianGenerator
	comb           []int
	iter           int
	numCombs       int
	forceFirstComb bool
}

func NewCombinator(lens []int, forceFirstComb bool) (*Combinator, []int) {
	var c Combinator
	cgen, comb := _initCombinationGenerator(lens)
	c.cgen = cgen
	c.comb = comb
	c.iter = -1
	c.forceFirstComb = forceFirstComb

	numCombs := 1
	for _, l := range lens {
		numCombs *= l
	}
	c.numCombs = numCombs
	if len(lens) == 0 {
		c.numCombs = 0
	}
	return &c, comb
}

func (c *Combinator) Next() bool {
	if c.forceFirstComb && c.iter == -1 && c.numCombs == 0 {
		return true
	}
	return c.cgen.Next()
}

func (c *Combinator) Product(dst []int) []int {
	c.iter++
	if c.forceFirstComb && c.iter == 0 && c.numCombs == 0 {
		return dst
	}
	return c.cgen.Product(dst)
}

func _initCombinationGenerator(lens []int) (*combin.CartesianGenerator, []int) {
	comb := make([]int, len(lens))
	cgen := combin.NewCartesianGenerator(lens)
	return cgen, comb
}
