#!/bin/sh

go list -f '{{$dir := .Dir}}{{range .GoFiles}}{{$dir}}/{{.}}{{"\n"}}{{end}}' ./internal/tool/... |
   xargs -n1 -I{} goblin -file {} |
   jq '.declarations[] | select(.type == "function") | select(.comments[] | match("Command: ")) | .comments' -c
