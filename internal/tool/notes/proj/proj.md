Proj integrates my various project management tools.

Usage is:

```bash
$ cd my/cool-project
$ proj init cool-project
```

The following flags are supported before the project name:

 * `-skip-vim`: does not create vim session
 * `-force-vim`: overwrites any existing vim session
 * `-skip-note`: does not create note
 * `-force-note`: overwrites any existing note
 * `-skip-smartcd`: does not create smartcd enter script
 * `-force-smartcd`: overwrites any existing smartcd enter script

Once a project has been initialized, you can run:

```bash
$ proj vim
```

To open the vim session for that project.

I use [vim sessions][vim], [a notes system][notes], and of course checkouts of
code all over the place.  Proj is meant to make creation of a vim session and a
note easy and eventually allow jumping back and forth between the three.  As of
2019-12-02 it is almost painfully specific to my personal setup, but as I
discover the actual patterns I'll probably generalize.

Proj uses uses [smartcd][smartcd] both as a mechanism and as the means to
add functionality to projects within shell sessions.

[vim]: https://blog.afoolishmanifesto.com/posts/vim-session-workflow/
[notes]: https://blog.afoolishmanifesto.com/posts/a-love-letter-to-plain-text/#notes
[smartcd]: https://github.com/cxreg/smartcd
