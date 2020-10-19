# DirectCounters

This example illustrates how you can use ryz to read values of DirectCounters. You can use DirectCounters for counting how many packets/bytes have matched with a specific entry in a specific table. If that's what you want to do with ryz, give this example a read. To run the example for yourself, follow [these steps](https://github.com/dush-t/ryz/tree/master/examples) and then come back here.

### On running the controller
The controller will install the p4 binary on the switch, and will install entries in the `ipv4_lpm` table (defined in `switch.p4`). After these entries are installed, the switch should be able to forward packets correctly i.e. you should be able to ping host `h2` from host `h1` (try it out, run `h1 ping h2` on the mininet prompt)

The controller will then read the value of a `DirectCounter` against the entry `10.0.0.10` in the table `ipv4_lpm` every two seconds and display the results. In the example topology we're using, there are three hosts (`h1`, `h2` and `h3` connected to switch ports 1, 2 and 3 respectively). 

If you run `h1 ping h2` now, `h1` will send packets to `h2` and `h2` will reply by sending packets to `h1`. The packets sent by `h2` will match with the table entry `10.0.0.10` in the table `ipv4_lpm` (since `10.0.0.10` is `h1`'s IP address), and the value of the counter will be incremented. Subsequently, you will see that the values read by the controller will increase, showing that the controller is working as intended.

Go ahead, try it yourself. It's not hard to set up.
