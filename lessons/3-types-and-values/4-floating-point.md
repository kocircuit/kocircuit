# FLOATING-POINT VALUES

Supported floating-point types are: `Float32` and `Float64`.

## FLOATING-POINT CONVERSIONS

Floating-point values can be converted to other floating-point types using the builtin
conversion functions: `Float32` and `Float64`.

Both floating-point conversion functions expect to be called with a single unnamed
floating-point argument.

For instance:

```ko
ConvertToFloat32(x) {
  return: Float32(x) // if the value of x is not Float32 or Float64, a type panic will be issued
}
```
