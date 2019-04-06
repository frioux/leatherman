/*
mozcookiejar gives access to cookies stored in SQLite by Mozilla products.

The general idea is that you can log into some website via Firefox (or whatever)
and then have your CLI tools effectively be logged in as well.

It is up to you to create a cookie jar and to connect to the SQLite database.
The example below should make it clear how to safely and correctly do that.  I
have tested the github.com/mattn/go-sqlite3 SQLite driver, but I am pretty sure
any of the various SQLite drivers will work.

This library does not support any of the following Cookie fields:

 * MaxAge
 * HttpOnly
 * Raw
 * Unparsed

See documentation for the underlying format here: http://kb.mozillazine.org/Cookies.sqlite

*/
package mozcookiejar

/*
Copyright 2018 Arthur Axel fREW Schmidt

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

import (
	"database/sql"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"time"
)

// LoadIntoJar populates the cookie jar from values in database.
func LoadIntoJar(db *sql.DB, jar *cookiejar.Jar) error {
	rows, err := db.Query(
		"SELECT host, path, name, value, isSecure, expiry FROM moz_cookies")
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var host, path, value, name string
		var isSecure bool
		var expiry int64

		err = rows.Scan(&host, &path, &name, &value, &isSecure, &expiry)
		if err != nil {
			return err
		}

		jar.SetCookies(&url.URL{Scheme: "http", Host: host}, []*http.Cookie{{
			Name:   name,
			Value:  value,
			Secure: isSecure,
			Path:   path,
			Domain: host,

			Expires:    time.Unix(expiry, 0),
			RawExpires: fmt.Sprintf("%d", expiry),

			// Intentionally Left Blank
			// * MaxAge
			// * HttpOnly
			// * Raw
			// * Unparsed
		}})
	}
	err = rows.Err()
	if err != nil {
		return err
	}

	return nil
}

// func Store(jar *cookiejar.Jar) error {
// 	return nil
// }
