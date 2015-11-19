#!/usr/bin/env perl

use 5.22.0;
use warnings;

$|++;
use LWP::UserAgent;
use JSON::XS 'decode_json';

use experimental 'postderef', 'signatures';

my $ua = LWP::UserAgent->new;
my $token = $ENV{SLACK_TOKEN} or die "you must set SLACK_TOKEN\n";
my $pal   = shift;

my $user_id = get_userid($pal);

warn "User ID for $pal is $user_id\n";

$SIG{INT} = sub {
   warn "Resetting!\n";
   say 'g0';
   exit 1;
};

while (1) {
   sleep 2;

   if (is_online($user_id)) {
      say 'g255';
   } else {
      say 'g0';
   }

}

sub get_userid ($pal) {
   my $req = $ua->get(
      "https://slack.com/api/users.list?token=$token"
   );

   die "could not get user id!\n" unless $req->is_success;

   my ($ret) =
      map $_->{id},
      grep $_->{name} =~ m/$pal/,
      decode_json($req->decoded_content)->{members}->@*;

   $ret;
}

sub is_online ($user_id) {
   my $req = $ua->get(
      "https://slack.com/api/users.getPresence?token=$token&user=$user_id"
   );

   return unless $req->is_success;

   decode_json($req->decoded_content)->{presence} eq 'active'
}
