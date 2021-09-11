package steamsrv

import (
	"fmt"
	"io/fs"
	"os"
	"regexp"
	"sort"
	"time"

	"github.com/frioux/leatherman/internal/steam"
)

// shotPattern decomposes steam screenshots filenames
//                                     1             2                   3          4     5    6     7     8      9
//                                     ?            appid               year      month  day   hour  min  sec     i
//                                    760           319630              2021       04    15    21    05    01     1
var shotPattern = regexp.MustCompile(`([^/]+)/remote/([^/]+)/screenshots/(\d\d\d\d)(\d\d)(\d\d)(\d\d)(\d\d)(\d\d)_(\d+).jpg`)

type screenshot struct {
	Name, Thumbnail string

	AppID int
	Date  time.Time
}

func (s screenshots) Screenshots(wantAppID int) ([]screenshot, error) {
	filenames, err := fs.Glob(s.fss, "*/remote/*/screenshots/*.jpg")
	if err != nil {
		return nil, err
	}

	ret := make([]screenshot, 0, 1000)
	for _, filename := range filenames {
		m := shotPattern.FindStringSubmatch(filename)
		if len(m) == 0 {
			fmt.Fprintf(os.Stderr, "path didn't match pattern: %s\n", filename)
			continue
		}

		appID := mustAtoi(m[2])
		if wantAppID != 0 && wantAppID != appID {
			continue
		}

		date := time.Date(
			// year         month                       day
			mustAtoi(m[3]), time.Month(mustAtoi(m[4])), mustAtoi(m[5]),
			// hour         minute          second
			mustAtoi(m[6]), mustAtoi(m[7]), mustAtoi(m[8]),
			// ns timezone
			0, time.Local)

		thumbnail := fmt.Sprintf("%s/remote/%s/screenshots/thumbnails/%s%s%s%s%s%s_%s.jpg", m[1], m[2], m[3], m[4], m[5], m[6], m[7], m[8], m[9])

		if f, err := s.fss.Open(thumbnail); err == nil {
			f.Close()
		} else {
			thumbnail = ""
		}
		ret = append(ret, screenshot{
			AppID:     appID,
			Name:      filename,
			Thumbnail: thumbnail,
			Date:      date,
		})
	}

	sort.Slice(ret, func(i, j int) bool { return ret[i].Date.Before(ret[j].Date) })
	return ret, nil
}

type screenshots struct {
	a   *steam.AppIDs
	fss fs.FS
}
