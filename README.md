# Epify

Epify categorizes shows using the
[Jellyfin naming scheme](https://jellyfin.org/docs/general/server/media/shows/).

Usage:

    epify show name year tvdbid dir
    epify season seasonnum showdir episode...
    epify add seasondir episode...


`epify show` creates a show directory like "Series Name (2018) [tvdbid-65567]".

`epify season` populates a season directory with episodes. Episodes are labeled
like "Series Name S01E01.mkv".

`epify add` adds episodes to a season directory, continuing at the previous
episode increment.

## Examples

Create show directory `/media/shows/The Office (2005) [tvdbid-73244]`:

```sh
$ epify show 'The Office' 2005 73244 '/media/shows'
```

Populate season directory
`/media/shows/The Office (2005) [tvdbid-73244]/Season 03`:

```sh
$ epify season 3 '/media/shows/The Office (2005) [tvdbid-73244]' '/downloads/the_office_s3_part_1' '/downloads/the_office_s3_ep13.mkv'
```

Add episodes to `/media/shows/The Office (2005) [tvdbid-73244]/Season 03`:

```sh
$ epify add '/media/shows/The Office (2005) [tvdbid-73244]/Season 03' '/downloads/the_office_s3_ep23.mkv' '/downloads/the_office_s3_p2'
```
