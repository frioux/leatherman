package notes

// 1. load, parse, and render now.md
//   * load(db) (io.Reader, error)
//   * parse(io.Reader) ([]byte, error)
//   * render(rw http.ResponseWriter, r *http.Request)
// 2. allow ~~ toggling in now.go for `## $TODAY ##`
//   * modify parse to add /toggle post link
//   * add /toggle?item=md5sum(item)
// 3. list and sort all files with `todo-` prefix
//   * list(db) ([]file, error)
//   * render
