# Epify

Epify categorizes shows using the
[Jellyfin naming scheme](https://jellyfin.org/docs/general/server/media/shows/).

Usage:

    epify show name year tvdbid dir
    epify season [-m index] seasonnum showdir episode...
    epify add [-m index] seasondir episode...


`epify show` creates a show directory like "Series Name (2018) [tvdbid-65567]".

`epify season` populates a season directory with episodes. Episodes are labeled
like "Series Name S01E01.mkv".

`epify add` adds episodes to a season directory, continuing at the previous
episode increment.

The `-m` flag specifies the index of the episode number in filenames for the
`epify season` and `epify add` commands.

## Examples

Create show directory `/media/shows/The Office (2005) [tvdbid-73244]`:

```sh
$ epify show 'The Office' 2005 73244 '/media/shows'
```

Populate season directory
`/media/shows/The Office (2005) [tvdbid-73244]/Season 03`:

```sh
$ epify season 3 '/media/shows/The Office (2005) [tvdbid-73244]' /downloads/the_office_s3_p1/ep*.mkv
```

Populate season directory
`/media/shows/Breaking Bad (2008) [tvdbid-81189]/Season 04`:

```sh
$ epify season -m 1 4 '/media/shows/Breaking Bad (2008) [tvdbid-81189]' /downloads/breaking_bad_s4_p1/s4ep*.mkv
```

Add episodes to `/media/shows/The Office (2005) [tvdbid-73244]/Season 03`:

```sh
$ epify add '/media/shows/The Office (2005) [tvdbid-73244]/Season 03' /downloads/the_office_s3_p2/ep*.mkv
```

Add episodes to `/media/shows/Breaking Bad (2008) [tvdbid-81189]/Season 04`:

```sh
$ epify add -m 1 '/media/shows/Breaking Bad (2008) [tvdbid-81189]/Season 04' /downloads/breaking_bad_s4_p2/s4ep*.mkv
```
