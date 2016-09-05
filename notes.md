- so, instead of wrapping a generic network interface
- let's provide a generic RPC interface
- where it is sync RPC interface by default
- and in the future, it can be just an async RPC interface of some sort.
- an RPC interface
- flag, reserve, service.method. payload addr, payload length
  return addr, return size

0: flag
1: error code
4-6: service id
6-8: method id
8-12: request addr
12-16: request length
16-20: response addr
20-24: response length
24-28: response length (used)
28-32: reserved