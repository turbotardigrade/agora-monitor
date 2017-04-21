#!/usr/bin/env python2.7

import libtmux
import conf

nodes = conf.nodes

server = libtmux.Server()
session = server.find_where({ "session_name": "agora" })

w = session.new_window(attach=False, window_name="network control")

panes = [None for _ in range(8)]

panes[0] = w.list_panes()[0]

for i in range(1, len(nodes)):
    panes[i] = w.split_window()
    w.select_layout(layout='tiled')

for i, n in enumerate(nodes):
    IP = n[0]
    panes[i].send_keys('ssh ubuntu@%s' % IP)
    panes[i].send_keys('cd feed/feeder\t')
    panes[i].send_keys('echo %s' %n[1])
