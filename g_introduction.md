# G Language Introduction

## If you are a Go programmer

G is very similar to Go, almost a almost a strict subset. Thus G is not a new languge for Go programmer.

Difference between G and Go:
* Delaration of struct and interface
In G, interface and struct is declared as, struct + *name*, but not type *name* struct. e.g.:

    struct circle {
        center int
        radium int
    }

    or

    interface symbol {
        String() string
    }

* if/else block
In G, if there is only one *break*, *continue*, or *return* statement after *if* or *else*
the the brackets for the block can be omitted. e.g. in G, your can write:
    
    func isEven (i int) bool {
        if i%2==0  return true
        return false
    }

    or:

    for ... {
        ...
        if flag==true break;
        ...
    }

Since we are still developing our tool chain, so there are many features of Golang current not supported by G.

## If you are a programmer of other languages like C/C++, Java

