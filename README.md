# gpuinfo

A simple CLI for getting information about GPUs from you servers listed
in SSH config (`.ssh/config`).

**Example output:**

![Example output](gpuinfo-example-output.png)

## Installation

```bash
git clone https://github.com/fkdosilovic/gpuinfo.git
cd gpuinfo
go build -o gpuinfo main.go
```

You can add the compiled binary to your `$HOME/.local/bin` or just run it from
the directory where you cloned the repo.

## ToDo

- [ ] Use different for free GPUs and GPUs in use
- [ ] Add an option to show only free GPUs
<!-- - [ ]: Add an option to show only GPUs in use -->
- [ ] Add an option to sort by different criteria
<!-- - [ ]: Add an option to show only GPUs with a certain amount of memory -->