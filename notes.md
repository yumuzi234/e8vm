- for handling keyboard inputs, the vm needs to be able to enter a suspending
state
- when an input comes, the vm continues running
- does this thing has anything to do with interrupt?
- no, it does not have to...
- wake-up events can be: keyboard input, mouse input, network input.