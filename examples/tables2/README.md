# tables-2
This example illustrates how you can use ryz to control a very simple ipv4-forwarding switch. If managing table entries is what you want to use ryz for, give this example a read. This is basically the same thing as [tables-1](https://github.com/dush-t/ryz/tree/master/examples/tables1), except the table structure in the switch is a bit more complicated than the former.

**Disclaimer**: The only reason this example is here is because I had some P4 code from one of my earlier projects, which I reused to create this example. I later realized that this example is unnecessarily complex for beginners. To fix this, I wrote [tables-1](https://github.com/dush-t/ryz/tree/master/examples/tables1). So, I recommend that you read the `tables-1` example before (or even instead of) this one.


### Running the example
To run the example for yourself, follow [these steps](https://github.com/dush-t/ryz/tree/master/examples) and then come back here.

### What the controller does
The controller will install the p4 binary on the switch, and will install entries in the `ipv4_lpm` table (defined in `switch.p4`). After these entries are installed, the switch should be able to forward packets correctly i.e. you should be able to ping host `h2` from host `h1` (try it out, run `h1 ping h2` on the mininet prompt)
