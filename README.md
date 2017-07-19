# Kitsu

Kitsu is a real-time files sync tool which handle versioning.

Features:

- Real-time files synchronization
- Cleaner that delete files when a limit is reach
- Easily exclude patterns
- Very lightweight (thanks to go!)

## Install

`go get -u github.com/toxinu/kitsu`

## Usage

`kitsu folder-to-save output`

## Help

```
$ kitsu --help
Usage: kitsu [OPTIONS] <source> <destination>
  -cleaner-interval int
    	Cleaner interval in seconds. (minimum 5) (default 60)
  -exclude value
    	Exclude pattern as regexp (repeatable).
  -max-versions int
    	Maximum file versions to keep (will delete olders). (default 20)
```

## License

License is MIT.