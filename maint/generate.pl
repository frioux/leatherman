#!/usr/bin/perl

use strict;
use warnings;
use autodie;

use File::Basename 'fileparse';
use File::Spec ();
use Encode;
use JSON::PP;

no warnings 'uninitialized';

my %doc;

while (<STDIN>) {
   my $c = decode_json($_);

   my $doc_path = ($c->{path} =~ s/\.go$/.md/r);
   my $d = do { open my $fh, '<:encoding(UTF-8)', $doc_path; local $/; <$fh> };

   my ($tool, $dir) = fileparse($c->{path}, '.go');
   
   my ($cat) = ($dir =~ m{/internal/tool/([^/]+)/[^/]+});

   if (!$cat) {
      warn "no category for $tool, skipping...\n";
      next;
   }

   $doc{$cat}{$tool} = { path => File::Spec->abs2rel($doc_path), doc => $d };
}

open my $fh, '<:encoding(UTF-8)', 'maint/README_begin.md';
my $begin = "<!-- Code generated by maint/generate-README. DO NOT EDIT. -->\n" .
            do { local $/; <$fh> };
close $fh;

open $fh, '<:encoding(UTF-8)', 'maint/README_end.md';
my $end = do { local $/; <$fh> };
close $fh;

for my $category (keys %doc) {
   for my $tool (keys %{$doc{$category}}) {
      $doc{$category}{$tool}{doc} = "#### `$tool`\n\n$doc{$category}{$tool}{doc}\n"
   }
}

my $body = $begin;

for my $category (sort keys %doc) {
   $body .= "### $category\n\n";

   for my $tool (sort keys %{$doc{$category}}) {
      $body .= $doc{$category}{$tool}{doc};
   }
}

$body .= $end;

open my $readme, '>:encoding(UTF-8)', 'README.mdwn';
print $readme $body;

close $readme;

open my $help, '>:encoding(UTF-8)', 'help_generated.go';
print $help "package main\n\n";
print $help qq(import "embed"\n\n);

for my $category (sort keys %doc) {
   for my $tool (sort keys %{$doc{$category}}) {
      print $help "//go:embed $doc{$category}{$tool}{path}\n";
   }
}

print $help "var helpFS embed.FS\n\n";

print $help "var helpPaths = map[string]string{\n";

for my $category (sort keys %doc) {
   for my $tool (sort keys %{$doc{$category}}) {
      print $help qq(\t"$tool": "$doc{$category}{$tool}{path}",\n\n);
   }
}

print $help "}\n";

close $help;

system 'go', 'fmt';
