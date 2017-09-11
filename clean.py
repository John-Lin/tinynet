"""
Mininet 2.3.0d1 License

Modifications copyright (c) 2017 Che Wei, Lin
Copyright (c) 2013-2016 Open Networking Laboratory
Copyright (c) 2009-2012 Bob Lantz and The Board of Trustees of
The Leland Stanford Junior University

Original authors: Bob Lantz and Brandon Heller

We are making Mininet available for public use and benefit with the
expectation that others will use, modify and enhance the Software and
contribute those enhancements back to the community. However, since we
would like to make the Software available for broadest use, with as few
restrictions as possible permission is hereby granted, free of charge, to
any person obtaining a copy of this Software to deal in the Software
under the copyrights without restriction, including without limitation
the rights to use, copy, modify, merge, publish, distribute, sublicense,
and/or sell copies of the Software, and to permit persons to whom the
Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included
in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS
OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY
CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT,
TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE
SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

The name and trademarks of copyright holder(s) may NOT be used in
advertising or publicity pertaining to the Software or any derivatives
without specific, written prior permission.
"""

from subprocess import ( Popen, PIPE, check_output as co,
                         CalledProcessError )

def sh( cmd ):
    "Print a command and send it to the shell"
    print( cmd + '\n' )
    return Popen( [ '/bin/sh', '-c', cmd ], stdout=PIPE ).communicate()[ 0 ]

def main():
    print( "***  Removing OVS datapaths\n" )
    dps = sh("ovs-vsctl --timeout=1 list-br").strip().splitlines()
    if dps:
        sh( "ovs-vsctl " + " -- ".join( "--if-exists del-br " + dp
                                        for dp in dps if dp ) )
    # And in case the above didn't work...
    dps = sh( "ovs-vsctl --timeout=1 list-br" ).strip().splitlines()
    for dp in dps:
        sh( 'ovs-vsctl del-br ' + dp )

    print( "*** Removing all links of the pattern vethX\n" )
    links = sh( "ip link show | "
                "egrep -o '(veth+[[:alnum:]]+)'"
                ).splitlines()
    # Delete blocks of links
    n = 1000  # chunk size
    for i in range( 0, len( links ), n ):
        cmd = ';'.join( 'ip link del %s' % link
                            for link in links[ i : i + n ] )
        sh( '( %s ) 2> /dev/null' % cmd )

    print( "*** Removing all links of the pattern tapX\n" )
    taps = sh( "ip link show | "
               "egrep -o '(tap+[[:digit:]]+[[:alnum:]]+)'"
               ).splitlines()
    # Delete blocks of links
    n = 1000  # chunk size
    for i in range( 0, len( taps ), n ):
        cmd = ';'.join( 'ip link del %s' % tap
                            for tap in taps[ i : i + n ] )
        sh( '( %s ) 2> /dev/null' % cmd )

    print( "*** Removing all network namespaces of the pattern cni-X-X-X-X-X\n" )
    nses = sh( "ip netns | "
               "egrep -o '(cni-+[[:alnum:]]+-[[:alnum:]]+-[[:alnum:]]+-[[:alnum:]]+-[[:alnum:]]+)'"
               ).splitlines()
    # Delete blocks of links
    n = 1000  # chunk size
    for i in range( 0, len( nses ), n ):
        cmd = ';'.join( 'ip netns del %s' % ns
                            for ns in nses[ i : i + n ] )
        sh( '( %s ) 2> /dev/null' % cmd )

    print( "*** Cleanup complete.\n" )

if __name__ == "__main__":
    main()
