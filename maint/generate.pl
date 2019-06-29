#!/usr/bin/perl

use strict;
use warnings;
use autodie;

use Encode;
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

$doc{$_} = "### `$_`\n\n`$_` $doc{$_}\n" for keys %doc;

my $body = $begin;
$body .= $doc{$_} for sort keys %doc;
$body .= $end;

my %offsets;
my $offset = length $begin;
for my $cmd (sort keys %doc) {
   my $length = length(encode('UTF-8', $doc{$cmd}, Encode::FB_CROAK));
   $offsets{$cmd} = "[$offset:" . ($offset + $length) . "]";
   $offset += $length;
}

open my $readme, '>:encoding(UTF-8)', 'README.mdwn';
print $readme $body;

close $readme;

open my $help, '>:encoding(UTF-8)', 'help_generated.go';
$body =~ s/`/` + "`" + `/g;
print $help "package main\n\n" .
   "var readme = []byte(`$body`)\n\n" .
   "var commandReadme map[string][]byte\n" .
   "func init() {\n" .
   "\tcommandReadme = map[string][]byte{\n";

print $help qq(\t\t"$_": readme$offsets{$_},\n\n) for sort keys %offsets;

print $help "\t}\n";
print $help "}\n";
