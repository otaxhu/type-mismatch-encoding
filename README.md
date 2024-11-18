# Type mismatched encoding (Golang)

This repository is a fork of Go's 1.23.3 `encoding/xml` and `encoding/json` packages that adds the functionality to allow members of a struct with a type that does not match the type of the JSON/XML member it is mapped to, to be deserialized to their zero-values with no returned errors.

This "fork" is not a conventional Github fork, it's only copies of those packages.

## Purpose

The purpose of this library, mainly, is to solve [another issue](https://github.com/otaxhu/problem/issues/14) of another repository (https://github.com/otaxhu/problem), I thought it could be more useful if it is decoupled of that repository.

## Audit

In order to be the most transparent, you can check the differences between this fork and Go's 1.23.3 packages by downloading a copy of Go 1.23.3 and cloning this repository, then executing the `diff` command:

```sh
# diff on encoding/json
diff -ur path/to/go/src/encoding/json path/to/type-mismatch-encoding/encoding/json

# diff on encoding/xml
diff -ur path/to/go/src/encoding/xml path/to/type-mismatch-encoding/encoding/xml
```
