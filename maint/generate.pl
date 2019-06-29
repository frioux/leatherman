#!/usr/bin/perl -CO

use strict;
use warnings;
use autodie;

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

open my $readme, '>:encoding(UTF-8)', 'README.mdwn';

open my $fh, '<:encoding(UTF-8)', 'maint/README_begin.md';
print $readme do { local $/; <$fh> };
close $fh;

print $readme "### `$_`\n\n`$_` $doc{$_}\n" for sort keys %doc;

open $fh, '<:encoding(UTF-8)', 'maint/README_end.md';
print $readme do { local $/; <$fh> };
close $fh;
close $readme;
