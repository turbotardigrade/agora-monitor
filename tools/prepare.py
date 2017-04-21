from subprocess import Popen, PIPE, STDOUT
import libtmux

import conf

nodes = conf.nodes
peers = [None for _ in range(len(nodes))]

for i, n in enumerate(nodes):
    IP = n[0]
    NUM = n[1]
    
    p = Popen(['ssh', '-oStrictHostKeyChecking=no', 'ubuntu@'+IP], stdin=PIPE, stdout=PIPE, stderr=STDOUT)
    p.stdin.write(b'cd feed; rm -rf feeder*; ./create_feeders.sh '+NUM+' abc\n')

    for line in iter(p.stdout.readline, b''):
        print line,
        if 'peer identity: ' in line:
            peer = line.split(' ')[-1][:-1]
            print peer
            peers[i] = peer
            
            p.stdin.close()
            p.stdout.close()
            p.wait()
            break

for p in peers:
    print p

print '---------------- for monitor\'s main.go:'
for i, p in enumerate(peers):
    is_spammer = 'false'
    if int(nodes[i][1]) > 18:
            is_spammer = 'true'

    print '{"%s", %s},' % (p, is_spammer)
