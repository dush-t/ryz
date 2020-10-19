# counters

This example illustrates how you can use ryz to read values of counters. You can use counters for stuff like monitoring traffic, counting certain events, etc. If that's what you want to do with ryz, give this example a read. To run the example for yourself, follow [these steps](https://github.com/dush-t/ryz/tree/master/examples) and then come back here.

### On running the controller
The controller will install the p4 binary on the switch, and will install entries in the `ipv4_lpm` table (defined in `switch.p4`). After these entries are installed, the switch should be able to forward packets correctly i.e. you should be able to ping host `h2` from host `h1` (try it out, run `h1 ping h2` on the mininet prompt)

The controller will then read the counter `port_counter` on the switch every two seconds and display the results. In the example topology we're using, there are three hosts (`h1`, `h2` and `h3` connected to switch ports 1, 2 and 3 respectively). 

If you run `h1 ping h2` now, the counter will count packets on ports 1 and 2 (since `h1` is sending packets to `h2`, and `h2` sends packets to `h1` in response) and the values read by the controller will increase accordingly. Once you stop the `ping` operation, you'll notice that the counter values will not change anymore.

Go ahead, try it yourself. It's not hard to set up.
