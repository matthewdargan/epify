# Epify

Epify categorizes shows using the
[Jellyfin naming scheme](https://jellyfin.org/docs/general/server/media/shows/).

Usage:

    epify show name year tvdbid dir
    epify season season-num episode-start show-dir {episode-dir | episode...}
    epify add season-dir {episode-dir | episode...}


`epify show` creates a show directory like "Series Name (2018) [tvdbid-65567]".

`epify season` creates a season directory for a show. The directory will
contain episodes from `dir` labeled like "Series Name S01E01.mkv".

`epify add` adds episodes to a season directory, continuing at the last
episode increment.

## Examples

Create show directory `/media/shows/The Office (2005) [tvdbid-73244]`:

```sh
$ epify show 'The Office' 2005 73244 '/media/shows'
```

Create season directory
`/media/shows/The Office (2005) [tvdbid-73244]/Season 03` from episodes in
`/downloads/the_office_s3`:

```sh
$ epify season 3 1 '/media/shows/The Office (2005) [tvdbid-73244]' '/downloads/the_office_s3'
```

Create season directory
`/media/shows/The Office (2005) [tvdbid-73244]/Season 03` from individual
episodes:

```sh
$ epify season 3 1 '/media/shows/The Office (2005) [tvdbid-73244]' '/downloads/the_office_s3/ep1' '/downloads/the_office_s3/ep2'
```

Add episodes in `/downloads/the_office_s3_p2` to
`/media/shows/The Office (2005) [tvdbid-73244]/Season 03`:

```sh
$ epify add '/media/shows/The Office (2005) [tvdbid-73244]/Season 03' '/downloads/the_office_s3_p2'
```

Add individual episodes to
`/media/shows/The Office (2005) [tvdbid-73244]/Season 03`:

```sh
$ epify add '/media/shows/The Office (2005) [tvdbid-73244]/Season 03' '/downloads/the_office_s3/ep3' '/downloads/the_office_s3/ep4'
```
