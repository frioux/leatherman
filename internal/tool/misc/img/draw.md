Draws images with lua.

```bash
$ draw 'rect(10, 10, 118, 118, red, yellow)' > x.png
```

Inspired by pico-8.  This tool takes lua scripts as strings and writes a png to
standard out.

Consider this tool unstable, I'll likely make it read scripts from either
standard in, or files, or both, and make arguments no longer the default.

## Lua API

### `set(x, y, c)`

Takes an x, y coordinate and sets it to a color.

### `rgb(r, g, b)`

Takes a red, green, and blue value (from 0 to 255 or as floating points from 0
to 1), returns a color value.

The following colors are defined as globals for you:

 * black
 * white
 * red
 * blue
 * yellow
 * green
 * orange
 * purple
 * cyan
 * magenta

### `sin(t)`

Returns sine of t, in terms of pi, not degrees.

### `cos(t)`

Returns cosine of t, in terms of pi, not degrees.

### `tan(t)`

Returns tangent of t, in terms of pi, not degrees.

### `PI`

Constant for pi.

### `rect(x1, y1, x2, y2, bordercolor, fillcolor)`

Draws a rectangle from (x1, y1) to (x2, y2) with a border of bordercolor and
filled with fillcolor.

### `circ(x, y, r, bordercolor, fillcolor)`

Draws a circle around (x, y) with radius r with a border of bordercolor and
filled with fillcolor.

### `line(x1, y1, x2, y2, color)`

Draws a line from (x1, y1) to (x2, y2) in color.

## BUGS

Something is wrong with `line` in certain situations (I'm assuming infinity or
NaN or something is causing the issue.)  `line` is used when drawing `circle`s,
so you can see the bug by drawing a circle and there will be weird gaps in the
top and bottom of them.
