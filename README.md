# gpustatus

A simple CLI for getting information about GPUs from servers listed in SSH config (`.ssh/config`).

**Note**: Currently, only NVIDIA GPUs are supported.

**Example output:**

![Example output](gpustatus-example-output.png)

## Installation

You should have Go installed on your system. If you don't have it, you can install it from [here](https://golang.org/doc/install).

Recommended way of installing `gpustatus` is to use the Makefile:

```bash
git clone https://github.com/fkdosilovic/gpustatus.git
cd gpustatus
make install
```

Run `gpustatus --help` to see all available options.

## To Do

- [ ] Add an option to sort by different criteria
- [ ] Add different output options (json, csv)
