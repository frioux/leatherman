`update` checks to see if there's an update from github and installs it if there
is.  If LM_GH_TOKEN is set to a personal access token this can be called more
frequently without exhausting github api limits.
