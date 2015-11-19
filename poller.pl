#!/usr/bin/env perl

use 5.22.0;
use warnings;

use List::Util qw(max min);

my %colors = (
   r => 0,
   g => 0,
   b => 0,
);

while (1) {
   while (<STDIN>) {
      if (my ($color, $direction, $amount) = m/^([rgb])([+-])?(\d+)/) {
         $direction ||= '';
         print "c: $color, a: $amount, d: $direction\n";

         if ($direction eq '-') {
            $colors{$color} -= $amount;
         } elsif ($direction eq '+') {
            $colors{$color} += $amount;
         } else {
            $colors{$color}  = $amount;
         }

         # verify range
         $colors{$color} = max(0,   $colors{$color});
         $colors{$color} = min(255, $colors{$color});

         system 'blink1-tool', sprintf '--rgb=%d,%d,%d', @colors{qw(r g b)}
      } else {
         warn "wtf? $_\n"
      }
   }
}
