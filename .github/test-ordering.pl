#!/usr/bin/perl
#
# Trivial script to find *.go files, and ensure their
# functions are all defined in alphabetical order.
#
# Ignore "init" and any "BenchmarkXXX" functions.
#

use strict;
use warnings;

use File::Find;

# Failure count
my $failed = 0;

# Find all files beneath the current directory
find( { wanted => \&process, no_chdir => 1, follow => 0 }, '.' );

# Return the result as an exit code
exit $failed;


# Process a file
sub process
{
    # Get the filename, and make sure it is a file.
    my $file = $File::Find::name;
    return unless ( $file =~ /.go$/ );

    print "$file\n";

    open( my $handle, "<", $file ) or
      die "Failed to read $file: $!";

    my @subs;

    foreach my $line (<$handle>)
    {
        if ( $line =~ /^func\s+([^(]+)\(/ )
        {
            my $func = $1;

            # Skip init
            next if $func eq "init";

            # Skip BenchmarkXXX
            next if $func =~ /^Benchmark/;

            # Record the function now.
            push( @subs, $func );
        }
    }
    close $handle;

    # Is the list of functions sorted?
    my @sorted = sort @subs;
    my $len    = $#sorted;

    my $i = 0;
    while ( $i < $len )
    {
        if ( $sorted[$i] ne $subs[$i] )
        {
            print "$sorted[$i] ne $subs[$i]\n";
            $failed++;
        }
        $i++;

    }
}
