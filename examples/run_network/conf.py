import os

NUM_HOSTS = 3
BEHAVIORAL_EXE = 'simple_switch_grpc'
LOG_FILE = os.path.join('/var/log', 'monitor.p4.log')
THRIFT_PORT = 9090
PCAP_DUMP = False

SWITCH_CONFIG = "switch/congestion_monitor.config"
SWITCH_PROGRAM_PATH = 'switch/monitor.p4'

NOTIFICATIONS_ADDR = 'ipc:///tmp/bmv2-0-notifications.ipc'