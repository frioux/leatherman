#!/bin/sh

go list -f '{{$dir := .Dir}}{{range .GoFiles}}{{$dir}}/{{.}}{{"\n"}}{{end}}' ./internal/tool/... |
   xargs -n1 -I{} goblin -file {} |
   jq '.path as $path | .declarations[] | select(.type == "function") | select(.name.value | test("^[A-Z]")) | {"path": $path, "comments": .comments[]}' -c
