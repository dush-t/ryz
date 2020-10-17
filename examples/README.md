# Examples

1. [Installing table entries on the switch](https://github.com/dush-t/ryz/tree/master/examples/tables1) - Tables 1
2. [Installing table entries on a slightly more complex switch](https://github.com/dush-t/ryz/tree/master/examples/tables2) - Tables 2
3. [Reading counter values from the control plane](https://github.com/dush-t/ryz/tree/master/examples/counters) - Counters

## How to run
The method to run an example is the same for all examples. In these steps, I'm assuming we are trying to run the example [Tables-1](https://github.com/dush-t/ryz/tree/master/examples/tables1), which is in the directory `tables1`. Replacing this example with another, all you need to do is change the directory name (if you want to run `Tables-2`, just use the corresponding directory name `tables2`)

### 1. Move to the target example directory
```console
dushyant@p4dev:~$ cd ryz/examples/tables1
```

### 2. Compile the p4 code
```console
dushyant@p4dev:~/ryz/examples/tables1$  p4c --std p4_16 --p4runtime-files out/switch.p4info.txt -o out switch.p4
```
This will create a new folder in the `tables1` directory, named `out`. The contents of this folder look like this - 
```
ryz/examples/tables1/out
  |- switch.json         # The "binary" that will be installed on our switch
  |- switch.p4info.txt   # The p4Info file needed for using P4Runtime
  |- switch.p4i          # Ignore this for now
```

### 3. Build the controller
First move to the `controller` directory
```console
dushyant@p4dev:~/ryz/examples/tables1$ cd controller
```
Now, build the controller
```console
dushyant@p4dev:~/ryz/examples/tables1/controller$ go build
```

### 4. Start the virtual network
We'll use mininet to start a network that we can test our switch in. Open another terminal and navigate to the `examples` directory.
```console
dushyant@p4dev:~$ cd ryz/examples
```
Now run this command to start the network.
```console
dushyant@p4dev:~/ryz/examples$ sudo python run_network/run.py
```
**Note**: You do have to use sudo
This will start a network with a linear topology with one switch. Right now, this switch is "blank" as in no P4 program is installed on it. If you try to do `h1 ping h2`, you will see 100% packet loss since the switch does not have any forwarding rules yet.

### 5. Run the controller
When run, the controller will install the selected p4 binary file (`switch.json`, remember?) on the switch, and also install forwarding rules so that it can forward packets.
From within the `controller` directory, run this command to start the controller - 
```console
dushyant@p4dev:~/ryz/examples/tables1/controller$ ./controller --bin ../out/switch.json --p4Info ../out/switch.p4info.txt
```
