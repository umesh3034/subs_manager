TL project : Disaggregated BNG demonstration:
---------------------------------------------

Aim:
* To showcase a disaggregated BNG gateway on cloud environment(see illustration).
* Propose standardized APIs/Messaging between access protocol entities and Subscriber manager.
* vBNG/IPOS components are _not_ getting reused here. This has been extensively debated in the TL meetings.
* Showcase NATs messaging wrt to protocol entities and Subscriber manager.

see basic design.doc attached

======================================================

Prerequisites:
--------------
1)basic DHCP knowledge
2)nats messaging : (http://nats.io/)
3)understanding of radius authentication
4)basic understanding of scapy tool- for testing dhcp client : (https://www.howtoinstall.co/en/ubuntu/trusty/python-scapy)
5)use a Linux machine
6)install freeradius server : (http://blog.moatazthenervous.com/installing-radius-on-ubuntu-14-04/)
7)install gnats server : (https://github.com/nats-io/gnatsd)

=======================================================

https://github.com/umesh3034/Mirror-of-ISC-DHCP
https://github.com/umesh3034/cnats
https://github.com/umesh3034/scapy
https://github.com/umesh3034/subs_manager.git


How to build:
-------------

1)/* build cnats client only if any changes are required, else the libnats_static.a is already built and copied to dhcp-4.3.5/cnats . you can directly go to #2 below in that case. */
Build the cnats client static library. (https://github.com/nats-io/cnats)
First build the cnats-master - i.e the cnats client to be used with dhcp server

cd cmake-master
mkdir build
cd build
cmake ..

mkdir ../../dhcp-4.3.5/cnats
cp -r ../install/* ../../dhcp-4.3.5/cnats
rm ../../dhcp-4.3.5/cnats/lib/libnats.so  /* not using .so, libnats_static.a is used */

2)Build the dhcp server application. (https://www.isc.org/downloads/)

#1 above copied the nats client static binary to dhcp folder. Makefile changes and application changes are already done to build nats client into dhcp server and calling required APIs to send and receive nats messages to nats server(see the design.doc)

cd dhcp-4.3.5
./configure --enable-debug 
make


3)To Build GO application for subscriber manager, set up the environment as per : (https://golang.org/doc/code.html)

=======================================================

How to Run and test the full flow:
-------------------------------


1)Run radius server
  cd /etc/freeradius
  vim clients.conf
  add 
	client 0.0.0.0/0 {
	  secret = "redback"
	  shortname = name
	}
  sudo vim users - add test line to end of this file : umesh Auth-Type := Accept, User-Password == "redback"
  sudo freeradius -fxxX
  
2)Run installed nats server

  eumesbs@elxa4r7r022:~/GO/bin$ ./gnatsd -D -V


3)run dhcp server

  cd dhcp-4.3.5

  mkdir -p /var/db
  touch /var/db/dhcpd.leases
  sudo vim /etc/dhcpd.conf /* add the subnet range */
  sudo ./server/dhcpd -cf /etc/dhcpd.conf -lf /var/db/dhcpd.leases

4)start the subscriber manager application

eumesbs@elxa4r7r022$ go run subs_manager.go AUTH


5)start the client
make sure the scapy python packages are installed and working fine

 cd scapy-master/umesh/DHCPig-master
 sudo ./DHCP.py -c -v99 wlan0 -> wlan0 is the interface for which the subnet is defined for the dhcp server, change as per needed

========================================================

