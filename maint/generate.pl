#!/usr/bin/perl -CO

use strict;
use warnings;

use JSON::PP;

no warnings 'uninitialized';

my %doc;

while (<STDIN>) {
   my $c = decode_json($_);

   die "Command should have exactly one comment\n" if @$c != 1;

   my $d = $c->[0];

   $d =~ s/^ \/\*\s+  //x;
   $d =~ s/  \s+\*\/ $//x;

   my ($body, $cmd) = ($d =~ m/^(?:\S+\s+)(.+)\s+Command:\s+(.+)$/s);

   $doc{$cmd} = $body;
}

print "### `$_`\n\n`$_` $doc{$_}\n" for sort keys %doc;
