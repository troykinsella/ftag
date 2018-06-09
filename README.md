# `ftag`

[![GitHub (pre-)release](https://img.shields.io/github/release/troykinsella/ftaG/all.svg)](https://github.com/troykinsella/ftag/releases)
[![License](https://img.shields.io/github/license/troykinsella/ftag.svg)](https://github.com/troykinsella/ftag/blob/master/LICENSE)
[![Build Status](https://travis-ci.org/troykinsella/ftag.svg?branch=master)](https://travis-ci.org/troykinsella/ftag)
[![Go Report](https://goreportcard.com/badge/github.com/troykinsella/ftag)](https://goreportcard.com/report/github.com/troykinsella/ftag)

> "file tag"

A command line utility for managing the assignment of tags to files.

## Installation

Head over to [releases](https://github.com/troykinsella/ftag/releases) and download the appropriate binary for your system.
Put the binary in a convenient place, such as `/usr/local/bin/ftag`.

## Usage

`ftag` maintains a bi-directional mapping of file names to tags in a `.ftag` file as JSON.

### Tag a File

```bash
$ ftag add my_file.txt awesome
```

### Remove a Tag on a File

```bash
$ ftag rm my_file.txt cool
```

### List Tags on a File

```bash
$ ftag list my_file.txt
awesome
```

### Clear Tags on a File

```bash
$ ftag clear my_file.txt
```

### Lookup Files by Tag

```bash
$ ftag find awesome
my_file.txt
my_other_file.txt
```

### Lookup Files Having Multiple Tags (AND)

```bash
$ ftag find awesome wicked
my_awesome_and_wicked_file.txt
```

### Move a File

When you move a file, `ftag` needs to be notified so it can update its mapping file.

```bash
$ mv oldfile.txt newfile.txt
$ ftag mv oldfile.txt newfile.txt
```

## License

MIT © Troy Kinsella
