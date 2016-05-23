# A conflict-free replicated data type.

[![Go Lang](http://kavehmz.github.io/static/gopher/gopher-front.svg)](https://golang.org/)
[![GoDoc](https://godoc.org/github.com/kavehmz/lww?status.svg)](https://godoc.org/github.com/kavehmz/lww)
[![Build Status](https://travis-ci.org/kavehmz/lww.svg?branch=master)](https://travis-ci.org/kavehmz/lww)
[![Coverage Status](https://coveralls.io/repos/kavehmz/lww/badge.svg?branch=master&service=github)](https://coveralls.io/github/kavehmz/lww?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/kavehmz/lww)](https://goreportcard.com/report/github.com/kavehmz/lww)

# What is it

In distributed computing, a __conflict-free replicated data type__ (CRDT) is a type of specially-designed data structure used to achieve strong eventual consistency (SEC) and monotonicity (absence of rollbacks).

One type of data structure used in implementing CRDT is LWW-element-set.

LWW-element-set is a set that its elements have timestamp. Add and remove will save the timestamp along with data in two different sets for each element.

Queries over LWW-set will check both add and remove timestamps to decide about state of each element is being existed to removed from the list.
