package pcache

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test_If_NewPredicate creates a new Predicate.
func Test_If_NewPredicate(t *testing.T) {
	called := false
	p := NewPredicate(func() bool {
		called = true
		return true
	})
	c := New()
	assert.False(t, called, "f was called")
	assert.True(t, p.f(c), "f returned false")
	assert.True(t, called, "f was not called")
}

// Test_If_Eval_Calls_F_When_Not_In_Cache tests that Eval calls the function
// when the predicate is not in the cache. It also tests that the result is
// cached.
func Test_If_Eval_Calls_F_When_Not_In_Cache(t *testing.T) {
	c := New()
	called := false
	p := &Predicate{
		f: func(*Cache) bool {
			called = true
			return true
		},
	}
	assert.True(t, p.Eval(c), "Eval returned false")
	assert.True(t, called, "f was not called")
	_, ok := c.isInCache(p)
	assert.True(t, ok, "predicate was not cached")
}

// Test_If_Eval_Does_Not_Call_F_When_In_Cache tests that Eval does not call the
// function when the predicate is in the cache.
func Test_If_Eval_Does_Not_Call_F_When_In_Cache(t *testing.T) {
	c := New()
	called := false
	p := &Predicate{
		f: func(*Cache) bool {
			called = true
			return true
		},
	}
	c.addToCache(p, true)
	assert.True(t, p.Eval(c), "Eval returned false")
	assert.False(t, called, "f was called")
}

// Test_If_Boolean_Works tests that And/Or returns a new Predicate that is the
// logical AND/OR of the given Predicates.
func Test_If_Boolean_Works(t *testing.T) {
	type testcase struct {
		input             [3]bool
		expectedOutputAnd bool
		expectedCalledAnd [3]bool
		expectedOutputOr  bool
		expectedCalledOr  [3]bool
	}
	testcases := []testcase{
		{
			input:             [3]bool{true, true, true},
			expectedOutputAnd: true,
			expectedCalledAnd: [3]bool{true, true, true},
			expectedOutputOr:  true,
			expectedCalledOr:  [3]bool{true, false, false},
		},
		{
			input:             [3]bool{true, true, false},
			expectedOutputAnd: false,
			expectedCalledAnd: [3]bool{true, true, true},
			expectedOutputOr:  true,
			expectedCalledOr:  [3]bool{true, false, false},
		},
		{
			input:             [3]bool{true, false, true},
			expectedOutputAnd: false,
			expectedCalledAnd: [3]bool{true, true, false},
			expectedOutputOr:  true,
			expectedCalledOr:  [3]bool{true, false, false},
		},
		{
			input:             [3]bool{true, false, false},
			expectedOutputAnd: false,
			expectedCalledAnd: [3]bool{true, true, false},
			expectedOutputOr:  true,
			expectedCalledOr:  [3]bool{true, false, false},
		},
		{
			input:             [3]bool{false, true, true},
			expectedOutputAnd: false,
			expectedCalledAnd: [3]bool{true, false, false},
			expectedOutputOr:  true,
			expectedCalledOr:  [3]bool{true, true, false},
		},
		{
			input:             [3]bool{false, true, false},
			expectedOutputAnd: false,
			expectedCalledAnd: [3]bool{true, false, false},
			expectedOutputOr:  true,
			expectedCalledOr:  [3]bool{true, true, false},
		},
		{
			input:             [3]bool{false, false, true},
			expectedOutputAnd: false,
			expectedCalledAnd: [3]bool{true, false, false},
			expectedOutputOr:  true,
			expectedCalledOr:  [3]bool{true, true, true},
		},
		{
			input:             [3]bool{false, false, false},
			expectedOutputAnd: false,
			expectedCalledAnd: [3]bool{true, false, false},
			expectedOutputOr:  false,
			expectedCalledOr:  [3]bool{true, true, true},
		},
	}
	for _, tc := range testcases {
		tc := tc
		t.Run(fmt.Sprintf("AND %v-%v-%v", tc.input[0], tc.input[1], tc.input[2]), func(t *testing.T) {
			c := New()
			called := [3]bool{}
			p := [3]*Predicate{}
			for i := range p {
				i := i
				p[i] = &Predicate{
					f: func(*Cache) bool {
						called[i] = true
						return tc.input[i]
					},
				}
			}
			and := And(p[0], p[1], p[2])
			assert.Equal(t, tc.expectedOutputAnd, and.Eval(c), "Eval returned wrong value")
			assert.Equal(t, tc.expectedCalledAnd, called, "f was not called correctly")
		})
		t.Run(fmt.Sprintf("OR %v-%v-%v", tc.input[0], tc.input[1], tc.input[2]), func(t *testing.T) {
			c := New()
			called := [3]bool{}
			p := [3]*Predicate{}
			for i := range p {
				i := i
				p[i] = &Predicate{
					f: func(*Cache) bool {
						called[i] = true
						return tc.input[i]
					},
				}
			}
			or := Or(p[0], p[1], p[2])
			assert.Equal(t, tc.expectedOutputOr, or.Eval(c), "Eval returned wrong value")
			assert.Equal(t, tc.expectedCalledOr, called, "f was not called correctly")
		})
	}
}

// Test_If_Not_Returns_The_Negation_Of_The_Given_Predicate tests that Not returns
// the negation of the given predicate.
func Test_If_Not_Returns_The_Negation_Of_The_Given_Predicate(t *testing.T) {
	c := New()
	assert.False(t, Not(True()).Eval(c), "Not(True()) returned true")
	assert.True(t, Not(False()).Eval(c), "Not(False()) returned false")
}
