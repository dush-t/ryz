import os
from subprocess import PIPE, Popen
from time import sleep

from mininet.net import Mininet
from mininet.log import setLogLevel, info
from mininet.cli import CLI

from p4_mininet import P4Switch, P4Host
from p4runtime_switch import P4RuntimeSwitch
import conf
from topology import Topology
from compile_switch import run_compile_bmv2


def configure_switch(queue_rate):
    info(' ')

    info("Configuring switch")
    
    proc = Popen(["simple_switch_CLI"], stdin=PIPE)
    proc.communicate(input="set_queue_rate " + str(queue_rate))

    info("Configuration complete")
    info(' ')


def run():
    # output_file = run_compile_bmv2(conf.SWITCH_PROGRAM_PATH)

    num_hosts = conf.NUM_HOSTS

    topo = Topology(conf.BEHAVIORAL_EXE,
                    None,
                    conf.LOG_FILE,
                    conf.THRIFT_PORT,
                    conf.PCAP_DUMP,
                    num_hosts,
                    conf.NOTIFICATIONS_ADDR)

    info('Topology generated\n')

    net = Mininet(topo = topo,
                  host = P4Host,
                  switch = P4RuntimeSwitch,
                  controller = None)

    info('Network configuration generated\n')
    info('Starting network')
    net.start()

    info('Network started\n')

    sw_mac = ["00:aa:bb:00:00:%02x" % n for n in xrange(num_hosts)]
    sw_addr = ["10.0.%d.1" % n for n in xrange(num_hosts)]

    for n in xrange(num_hosts):
        h = net.get('h%d' % (n + 1))
        h.setARP(sw_addr[n], sw_mac[n])
        h.setDefaultRoute("dev %s via %s" % (h.defaultIntf().name, sw_addr[n]))

        h.describe(sw_addr[n], sw_mac[n])

    sleep(1)

    configure_switch(10)

    CLI( net )
    print('hello')
    net.stop()


setLogLevel('info')
run()
    
