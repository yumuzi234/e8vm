## Why write this project?

Why another language? Why another operating system?

All code dies. Some dies in hours, some dies in years. Seldom lives a
long life.

As a result, programmers keep writing essentially the same code again
and again, in different times, maybe under different stories.

Some code dies because it is no longer needed; that's fine. But
some code dies because it is too complex to maintain or modify. Often
times, it is not because the algorithm is too complex to comprehend,
but the architecture losses its structure due to unmanaged code debt.
No one understands the code any more.

To avoid replaying this history, Small VM proposes an entire set of
language tool-chain that puts code readability, or more precisely,
code *comprehensibility* as the one and only first priority.

Note that, even today, many programming languages and systems put coding
efficiency, performance and safety as first considerations, but
readability and comprehensibility as the last 
([for example](http://andrewkelley.me/post/intro-to-zig.html)).

I disagree with this ordering. At the scale of today's software engineering,
good code quality fundamentally comes
from relentless iterations, and iterations require good code
understandings by *human* programmers (before AI's can read and write 
programs). In the long run, code comprehensibility dominates all.

To achieve high comprehensibility, Small VM project:

- Uses simple programming languages with only a small set of
  language features.
- Limits file sizes (80 chars per line, 300 lines maximum), and
  forbids circular dependencies among files with compiler-enforced
  checking.
- Minimizes tool-chain setup frictions; compiles and runs
  right inside Web browsers.

These principles do not immediately lead to the best possible code
comprehensibility, but they at least encourage code reading, and with
some proper execution, they will hopefully close the feedback loop on
code comprehensibility so that comprehensibility can improve over iterations.

Long live our code.
