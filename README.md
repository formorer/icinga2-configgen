icinga2 config generator
------------------------

Sometimes you need testdata, sometimes a lot.
Therefore I wrote this simple Icinga2 configgenerator, 
it will create as many hosts, with as many services as you wish. 

It takes the host and service template from templates/, feel free to adapt or
change if you want.

Usage:

icinga2-configgen --hosts=50 --services=100 --confdir=/etc/icinga2/conf.d/

This will generate 50 hosts with 100 services each and writes them, one file per host,
into /etc/icinga2/conf.d/.

The template directory defaults to /etc/icinga2-configgen, but falls back to
templates/.
