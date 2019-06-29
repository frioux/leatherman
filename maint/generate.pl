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

open my $fh, '<:encoding(UTF-8)', 'maint/README_begin.md';
my $begin = do { local $/; <$fh> };
close $fh;

open $fh, '<:encoding(UTF-8)', 'maint/README_end.md';
my $end = do { local $/; <$fh> };
close $fh;

my $body = $begin;
$body .= "### `$_`\n\n`$_` $doc{$_}\n" for sort keys %doc;
$body .= $end;

open my $readme, '>:encoding(UTF-8)', 'README.mdwn';
print $readme $body;

close $readme;

open my $help, '>:encoding(UTF-8)', 'cmd/leatherman/help_generated.go';
$body =~ s/`/` + "`" + `/g;
print $help "package main\n\n" .
   "var readme = []byte(`$body`)"
