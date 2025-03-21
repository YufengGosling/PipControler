use strict;
use warnings;

open(my $file_data, '<', $ARGV[0]) or die "Cannot open file $ARGV[0] \nError: $!";

while (my $line = <$file_data>) {
    if ($line =~ /(?:import\s+|from\s+)\K[\w.]+/) {
        print "$&\n";
    }
}

