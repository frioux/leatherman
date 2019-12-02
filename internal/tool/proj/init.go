package proj

import (
	"errors"
	"flag"
	"fmt"
	"time"
	"crypto/sha1"
	"os"
)

func initialize(args []string) error {
	flags := flag.NewFlagSet(args[0], flag.ContinueOnError)
	var skipVim, skipNote, skipSmartCD, forceVim, forceNote, forceSmartCD bool

	flags.BoolVar(&skipVim, "skip-vim", false, "skips creation of vim session")
	flags.BoolVar(&skipNote, "skip-note", false, "skips creation of note")
	flags.BoolVar(&skipSmartCD, "skip-smartcd", false, "skips creation of smartcd")

	flags.BoolVar(&forceVim, "force-vim", false, "forces creation of vim session")
	flags.BoolVar(&forceNote, "force-note", false, "forces creation of note")
	flags.BoolVar(&forceSmartCD, "force-smartcd", false, "forces creation of smartcd")

	if err := flags.Parse(args[1:]); err != nil {
		return err
	}

	if len(flags.Args()) != 1 {
		return errors.New(args[0] + " requires at least one argument")
	}

	name := flags.Args()[0]

	if err := initSmartCD(name, skipSmartCD, forceSmartCD); err != nil {
		return err
	}

	if err := initVim(name, skipVim, forceVim); err != nil {
		return err
	}

	if err := initNote(name, skipNote, forceNote); err != nil {
		return err
	}

	return nil
}

func fileExists(name string) (bool, error) {
	_, err := os.Stat(name)
	if err == nil {
		return true, nil
	}

	if os.IsNotExist(err) {
		return false, nil
	}

	return false, err
}

func initSmartCD(name string, skip, force bool) error {
	if skip {
		return nil
	}

	workdir, err := os.Getwd() // XXX I bet more needs to be done here
	if err != nil {
		return err
	}

	dir := smartcd + "/" + workdir
	if err := os.MkdirAll(dir, os.FileMode(0755)); err != nil {
		return err
	}

	p := dir + "/bash_enter"
	e, err := fileExists(p)
	if err != nil {
		return err
	}

	if !force && e {
		return errors.New("file already exists: " + p)
	}

	f, err := os.Create(p)
	if err != nil {
		return err
	}

	if _, err := fmt.Fprintln(f, "autostash PROJ="+name); err != nil {
		return err
	}

	return nil
}

func initVim(name string, skip, force bool) error {
	if skip {
		return nil
	}

	workdir, err := os.Getwd() // XXX I bet more needs to be done here
	if err != nil {
		return err
	}

	p := vimSessions + "/" + name
	e, err := fileExists(p)
	if err != nil {
		return err
	}

	if !force && e {
		return errors.New("file already exists: " + p)
	}

	f, err := os.Create(p)
	if err != nil {
		return err
	}
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
	if _, err := fmt.Fprint(f, session); err != nil {
		return err
	}

	return nil
}

func initNote(name string, skip, force bool) error {
	if skip {
		return nil
	}

	p := notes + "/" + name + ".md"
	e, err := fileExists(p)
	if err != nil {
		return err
	}

	if !force && e {
		return errors.New("file already exists: " + p)
	}

	f, err := os.Create(p)
	if err != nil {
		return err
	}

	sum := sha1.Sum([]byte(name))

	// XXX could just use the post from .projections.json
	session := `---
title: ` + name + `
date: ` + time.Now().Format("2006-01-02T15:04:05") + `
tags: [ project ]
guid: ` + fmt.Sprintf("%x", sum[:]) + `
---
`
	if _, err := fmt.Fprint(f, session); err != nil {
		return err
	}

	return nil
}
