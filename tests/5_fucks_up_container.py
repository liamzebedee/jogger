# Tested this inside of a container

from subprocess import call
call("rm -rf home/".split(' '))
call("ls /".split(' '))