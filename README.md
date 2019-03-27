[![Build Status](https://travis-ci.org/gregoryv/golden.svg?branch=master)](https://travis-ci.org/gregoryv/golden)
[![codecov](https://codecov.io/gh/gregoryv/golden/branch/master/graph/badge.svg)](https://codecov.io/gh/gregoryv/golden)
[![Maintainability](https://api.codeclimate.com/v1/badges/df2736e1ac63580b49d7/maintainability)](https://codeclimate.com/github/gregoryv/golden/maintainability)

[golden](https://godoc.org/github.com/gregoryv/golden) - package defines test helper for working with golden files

Golden files contain expected values within your tests. They are useful
when you want to check more complex outputs. This package makes it easy
to Save and Load such files within the testdata directory.


    func TestMe(t *testing.T) {
        got := doSomething()
        exp := golden.LoadString()
        if got != exp {
            t.Fail()
        }
        golden.SaveString(t, got)
    }

Golden file is saved in testdata/package.TestMe and an entry is added to
testdata/golden.files

To update the golden files use

    go test -args -update-golden
