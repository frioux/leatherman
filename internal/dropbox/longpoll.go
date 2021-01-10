package dropbox

import (
	"context"
	"fmt"
	"os"
	"time"
)

// Longpoll signals changes in dir by sending on ch.
func (db Client) Longpoll(ctx context.Context, dir string, ch chan<- struct{}) {
OUTER:
	for {
		res, err := db.ListFolder(ListFolderParams{
			Path: dir,
		})
		if err != nil {
			fmt.Fprintln(os.Stderr, "ListFolder:", err)
			continue
		}

		for res.HasMore {
			res, err = db.ListFolderContinue(res.Cursor)
			if err != nil {
				fmt.Fprintln(os.Stderr, "ListFolderContinue:", err)
				continue
			}
		}

		cu := res.Cursor

		changed, backoff, err := db.ListFolderLongPoll(ctx, cu, 480)
		if err != nil {
			fmt.Fprintln(os.Stderr, "ListFolderLongPoll:", err)
			continue
		}

		if backoff != 0 {
			time.Sleep(time.Second * time.Duration(backoff))
		}

		if !changed {
			continue
		}

		res = ListFolderResult{HasMore: true, Cursor: cu}
		res, err = db.ListFolderContinue(res.Cursor)
		if err != nil {
			fmt.Fprintln(os.Stderr, "ListFolderContinue:", err)
			continue
		}

		if len(res.Entries) > 0 {
			ch <- struct{}{}
		}
		continue OUTER
	}
}
