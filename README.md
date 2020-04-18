[![Build Status](https://travis-ci.org/gregoryv/golden.svg?branch=master)](https://travis-ci.org/gregoryv/golden)
[![codecov](https://codecov.io/gh/gregoryv/golden/branch/master/graph/badge.svg)](https://codecov.io/gh/gregoryv/golden)
[![Maintainability](https://api.codeclimate.com/v1/badges/df2736e1ac63580b49d7/maintainability)](https://codeclimate.com/github/gregoryv/golden/maintainability)

[golden](https://godoc.org/github.com/gregoryv/golden) - package defines test helper for working with golden files

Golden files contain expected values within your tests. They are useful
when you want to check more complex outputs. This package makes it easy
to Save and Load such files within the testdata directory.

Simplest example

    func TestMe(t *testing.T) {
        got := doSomething()
        golden.Assert(t, got)
    }

The golden file for above test is saved in `testdata/package.TestMe`
and an entry is added to `testdata/golden.files` which keeps track of
used files. If you eg. rename a test the golden file will be saved
under a new name.

Keep `golden.files` under revision control to quickly spot which files
are no longer used.

To update the golden files use

    go test -args -update-golden


Article at [7de.se/golden](https://www.7de.se/golden/)
