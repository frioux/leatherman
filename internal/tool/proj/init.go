package proj

import (
	"crypto/sha1"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path"
	"time"
)

var workdir string

func init() {
	var err error
	workdir, err = os.Getwd() // XXX I bet more needs to be done here
	if err != nil {
		panic("error getting workdir: " + err.Error())
	}
}

type multiError []error

func (e *multiError) Error() string {
	ret := ""
	for _, err := range *e {
		ret += err.Error() + "; "
	}
	return ret
}

func initialize(args []string) error {
	flags := flag.NewFlagSet(args[0], flag.ContinueOnError)

	vimMP := managedPath{
		name:    "vim",
		path:    func(m managedPath) string { return vimSessions + "/" + m.proj },
		content: initVim,
	}
	noteMP := managedPath{
		name:    "note",
		path:    func(m managedPath) string { return notes + "/" + m.proj + ".md" },
		content: initNote,
	}
	smartcdMP := managedPath{
		name:    "smartcd",
		path:    func(m managedPath) string { return smartcd + "/" + workdir + "/bash_enter" },
		content: initSmartCD,
	}

	mps := []managedPath{vimMP, noteMP, smartcdMP}

	flags.BoolVar(&vimMP.skip, "skip-vim", false, "skips creation of vim session")
	flags.BoolVar(&noteMP.skip, "skip-note", false, "skips creation of note")
	flags.BoolVar(&smartcdMP.skip, "skip-smartcd", false, "skips creation of smartcd")

	flags.BoolVar(&vimMP.force, "force-vim", false, "forces creation of vim session")
	flags.BoolVar(&noteMP.force, "force-note", false, "forces creation of note")
	flags.BoolVar(&smartcdMP.force, "force-smartcd", false, "forces creation of smartcd")

	if err := flags.Parse(args[1:]); err != nil {
		return err
	}

	if len(flags.Args()) != 1 {
		return errors.New(args[0] + " requires at least one argument")
	}

	name := flags.Args()[0]
	for i := range mps {
		mps[i].proj = name
	}

	var errs multiError
	for _, m := range mps {
		if m.skip {
			continue // remove from list?
		}
		exists, err := m.fileExists()
		if err != nil {
			errs = append(errs, fmt.Errorf(m.name+" exist check: %w", err))
			continue
		}

		if !m.force && exists {
			errs = append(errs, errors.New("file already exists: "+m.path(m)))
			continue
		}
	}

	if len(errs) != 0 {
		return &errs
	}

	for _, m := range mps {
		if m.skip {
			continue
		}
		if err := m.manage(); err != nil {
			return err
		}
	}

	return nil
}

func (m managedPath) fileExists() (bool, error) {
	_, err := os.Stat(m.path(m))
	if err == nil {
		return true, nil
	}

	if os.IsNotExist(err) {
		return false, nil
	}

	return false, err
}

type managedPath struct {
	path        func(managedPath) string
	content     func(managedPath, io.Writer) error
	force, skip bool
	proj, name  string
}

func (m managedPath) manage() error {
	if err := os.MkdirAll(path.Dir(m.path(m)), os.FileMode(0755)); err != nil {
		return err
	}
	f, err := os.Create(m.path(m))
	if err != nil {
		return err
	}
	defer f.Close()

	return m.content(m, f)
}

func initSmartCD(m managedPath, w io.Writer) error {
	if _, err := fmt.Fprintln(w, "autostash PROJ="+m.proj); err != nil {
		return err
	}

	return nil
}

func initVim(m managedPath, w io.Writer) error {
	// blindly copy pasted from a fresh session
	session := `let SessionLoad = 1
if &cp | set nocp | endif
let s:so_save = &so | let s:siso_save = &siso | set so=0 siso=0
let v:this_session=expand("<sfile>:p")
silent only
silent tabonly
cd ` + workdir + `
if expand('%') == '' && !&modified && line('$') <= 1 && getline(1) == ''
  let s:wipebuf = bufnr('%')
endif
set shortmess=aoO
argglobal
%argdel
set splitbelow splitright
set nosplitbelow
set nosplitright
wincmd t
set winminheight=0
set winheight=1
set winminwidth=0
set winwidth=1
tabnext 1
if exists('s:wipebuf') && len(win_findbuf(s:wipebuf)) == 0
  silent exe 'bwipe ' . s:wipebuf
endif
unlet! s:wipebuf
set winheight=1 winwidth=20 shortmess=filnxtToOS
set winminheight=1 winminwidth=1
let s:sx = expand("<sfile>:p:r")."x.vim"
if file_readable(s:sx)
  exe "source " . fnameescape(s:sx)
endif
let &so = s:so_save | let &siso = s:siso_save
nohlsearch
let g:this_session = v:this_session
let g:this_obsession = v:this_session
let g:this_obsession_status = 2
doautoall SessionLoadPost
unlet SessionLoad
" vim: set ft=vim : 
`
	if _, err := fmt.Fprint(w, session); err != nil {
		return err
	}

	return nil
}

func initNote(m managedPath, w io.Writer) error {
	sum := sha1.Sum([]byte(m.proj))

	// XXX could just use the post from .projections.json
	session := `---
title: ` + m.proj + `
date: ` + time.Now().Format("2006-01-02T15:04:05") + `
tags: [ project ]
guid: ` + fmt.Sprintf("%x", sum[:]) + `
---
`
	if _, err := fmt.Fprint(w, session); err != nil {
		return err
	}

	return nil
}
