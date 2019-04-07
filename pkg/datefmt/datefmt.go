package datefmt

import "strings"

type replacement struct {
	from, to string
}

// Mon Jan 2 15:04:05 -0700 MST 2006
var replacements = []replacement{
	{"%a", "Mon"},
	{"%A", "Monday"},
	{"%b", "Jan"},
	{"%B", "January"},
	// %c     The preferred date and time representation for the current locale.
	{"%C", "06"},
	{"%d", "02"},
	{"%D", "01/02/06"},
	// %e     Like %d, the day of the month as a decimal number, but a leading zero is replaced by a space. (SU) (Calculated from tm_mday.)
	// %E     Modifier: use alternative format, see below. (SU)
	{"%F", "2006-01-02"},
	// %G     The  ISO 8601  week-based  year (see NOTES) with century as a decimal number.  The 4-digit year corresponding to the ISO week number (see %V).  This has the same format and
	//        value as %Y, except that if the ISO week number belongs to the previous or next year, that year is used instead. (TZ) (Calculated from tm_year, tm_yday, and tm_wday.)
	// %g     Like %G, but without century, that is, with a 2-digit year (00–99). (TZ) (Calculated from tm_year, tm_yday, and tm_wday.)
	{"%h", "Jan"},
	{"%H", "15"},
	{"%I", "03"},
	// %j     The day of the year as a decimal number (range 001 to 366).  (Calculated from tm_yday.)
	// %k     The hour (24-hour clock) as a decimal number (range 0 to 23); single digits are preceded by a blank.  (See also %H.)  (Calculated from tm_hour.)  (TZ)
	// %l     The hour (12-hour clock) as a decimal number (range 1 to 12); single digits are preceded by a blank.  (See also %I.)  (Calculated from tm_hour.)  (TZ)
	{"%m", "01"},
	{"%M", "04"},
	{"%n", "\n"},
	// %O     Modifier: use alternative format, see below. (SU)
	{"%p", "PM"},
	{"%P", "pm"},
	{"%r", "03:04:05 p.m."},
	{"%R", "15:03"},
	// %s     The number of seconds since the Epoch, 1970-01-01 00:00:00 +0000 (UTC). (TZ) (Calculated from mktime(tm).)
	{"%S", "05"},
	{"%t", "\t"},
	{"%T", "15:04:05"},
	{"%u", "1"},
	// %U     The  week number of the current year as a decimal number, range 00 to 53, starting with the first Sunday as the first day of week 01.  See also %V and %W.  (Calculated from
	//        tm_yday and tm_wday.)
	// %V     The ISO 8601 week number (see NOTES) of the current year as a decimal number, range 01 to 53, where week 1 is the first week that has at least 4 days in the new year.   See
	//        also %U and %W.  (Calculated from tm_year, tm_yday, and tm_wday.)  (SU)
	// %w     The day of the week as a decimal, range 0 to 6, Sunday being 0.  See also %u.  (Calculated from tm_wday.)
	// %W     The week number of the current year as a decimal number, range 00 to 53, starting with the first Monday as the first day of week 01.  (Calculated from tm_yday and tm_wday.)
	// %x     The preferred date representation for the current locale without the time.
	// %X     The preferred time representation for the current locale without the date.
	{"%y", "06"},
	{"%Y", "2006"},
	{"%z", "-7000"},
	{"%Z", "MST"},
	// %+     The date and time in date(1) format. (TZ) (Not supported in glibc2.)
	// %%     A literal '%' character.
}

// TranslateFormat replaces many of the most common strftime(3) formats to go
// formats.  See source for supported and unsupported formats.
func TranslateFormat(in string) string {
	ret := in
	for _, r := range replacements {
		ret = strings.ReplaceAll(ret, r.from, r.to)
	}

	return ret
}
