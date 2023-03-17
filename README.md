# gpuinfo

A simple CLI for getting information about GPUs from you servers listed
in SSH config (`.ssh/config`).

## Installation

```bash
git clone https://github.com/fkdosilovic/gpuinfo.git
cd gpuinfo
go build -o gpuinfo main.go
```

## ToDo

- [ ] Use different for free GPUs and GPUs in use
- [ ] Add an option to show only free GPUs
<!-- - [ ]: Add an option to show only GPUs in use -->
- [ ] Add an option to sort by different criteria
<!-- - [ ]: Add an option to show only GPUs with a certain amount of memory -->