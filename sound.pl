#!/usr/bin/perl

use 5.22.0;
use warnings;

$|++;
my $off_count = 0;

while (1) {
   sleep 1;

   if (playing_sounds()) {
      warn "red\n";
      print "r255\n";
      $off_count = 0;
   } elsif ($off_count++ >= 2) {
      warn "black\n";
      print "r0\n";
   }
}

sub playing_sounds {
   my @lines =
      grep m/RUNNING/,
      split /\n/,
      `pacmd list-sink-inputs`;

   warn "sound is playing\n" if @lines;
   warn "silence\n" if !@lines;

   scalar @lines
}
