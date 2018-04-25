#!/usr/bin/perl

use 5.22.0;
use warnings;

$|++;
my $off_count = 0;

while (1) {
   sleep 1;

   if (playing_sounds()) {
      print "r255\n";
      $off_count = 0;
   } elsif ($off_count++ >= 2) {
      print "r0\n";
   }
}

sub _parse_plsi {
   my @lines = `pacmd list-sink-inputs`;

   my @data;
   my $current;

   for my $line (@lines) {
      chomp $line;

      if ($line =~ m/index:\s/) {
         push @data, $current if $current;
         $current = {};
      }
      my $re = qr/\s*[:=]\s*/;
      if ($line =~ $re) {
         my ($l, $r) = split $re, $line, 2;
         $l =~ s/^\s*//;
         $r =~ s/^"//;
         $r =~ s/"$//;
         $current->{$l} = $r
      }
   }
   push @data, $current;

   return @data;
}

sub playing_sounds {
   my @sinks = _parse_plsi();

   my $chrome = 0;
   for (grep { ($_->{state}||'') eq 'RUNNING' } @sinks) {
      return 1 if $_->{'application.name'} !~ m/chrome/i;
      $chrome ++;
   }

   return 1 if $chrome > 0;

   return 0
}
