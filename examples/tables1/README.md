# tables-1

This example illustrates how you can use ryz to control a very simple ipv4-forwarding switch. If managing table entries is what you want to use ryz for, give this example a read. To run the example for yourself, follow [these steps](https://github.com/dush-t/ryz/tree/master/examples) and then come back here.

### On running the controller
The controller will install the p4 binary on the switch, and will install entries in the `ipv4_lpm` table (defined in `switch.p4`). After these entries are installed, the switch should be able to forward packets correctly i.e. you should be able to ping host `h2` from host `h1` (try it out, run `h1 ping h2` on the mininet prompt)
