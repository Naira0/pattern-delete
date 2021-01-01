# pattern-delete
A simple utility tool for deleting files with a pettern (regexp and extensions) also supports recursive deleting.

## how to use
it works on a flag system you must provide either `-p` or `-e` flags or both.


## flags
  `-e` The extension to look for. If not set it will ignore extensions
        
  `-p` The  regexp pattern to delete for. If not set it will either default to the given extension flag or it wont do anything at all.
        
  `-r` If set it will recursively delete files.
