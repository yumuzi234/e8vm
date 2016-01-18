how do we organize user level code and kernel level code?

- what we really want now is a consistent way to save the source code.
- we need a user account that belongs to the shanhu system

e8vm.io / os8 /

- the standard library will be very commonly imported.
- one solution is that we have special user accounts, these are system users accounts
- like shanhu.io/~std/
- for now, it is too early to start the ~std, as we do not have a working os yet
- so instead, we will start two repositories
- shanhu.io/~os8/ and shanhu.io/~toybox
- the ~std is a special repository for standard library that can be imported directly
- okay, so this one is solved
- another issue is how to deal with kernel bare-metal binaries and user level binaries
- the current difference is that kernel binaries and user binaries start at
  different program counters
