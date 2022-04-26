# Examples

The following CLI arguments can be used to verify the behaviour of the package:

```bash
TEST_STRING_FIELD="i like trains" go run example.go -test-string-field "ananas on pizza is good" -test-int-field 69
TEST_INT_FIELD=420 go run example.go -test-int-field 69
```