# HANGING

The builtin function `Hang` hangs forever, never returning a value.

For instance, this function will block forever:

```ko
BlockForever() {
  return: Hang()
}
```
