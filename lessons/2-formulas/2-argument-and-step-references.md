# REFERENCES

A _reference_ is a formula which refers to the value of a functional
argument or another step in the function body. A reference is 
expressed as the name identifier of the argument or step being
referred to.

The following example demonstrates a reference to an argument:

```ko
Greet(name) {
  return: Print("Hello", name) // the identifier "name" is an argument reference
}
```

The following example demonstrates both types of references:

```ko
import "github.com/kocircuit/kocircuit/lib/strings"

Greet(name) {
  greeting: strings.Join(
    string: ("Hello", name) // the identifier "name" is an argument reference
    delimiter: " "
  )
  return: Print(greeting) // the identifier "greeting" is a step reference
}
```
