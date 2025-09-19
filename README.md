# elfGetArgs

`elfGetArgs` is a CLI tool for analyzing ELF (Executable and Linkable Format) files, specifically designed to extract argument information of specified functions.  
It is useful for static analysis, reverse engineering, and security research.

---

## 📦 Installation

Clone the repository: 
```bash
git clone https://github.com/MercuryNearTheMoon/elfGetArgs.git
``` 

### Run Locally

To compile the project:

```bash
cd elfGetArgs
go build .
```

---

### 🐳 Docker Usage

You can run `elfGetArgs` via Docker without installing Go locally.

#### Build Docker image

```bash
cd elfGetArgs
docker build -t elfgetargs .
```

#### Run Docker container
```bash
docker run --rm \
  -v /path/to/binaries:/target \
  -v /path/to/output:/app/output \
  elfgetargs \
  -p /target -A amd64 -f open -a 0 -o /app/output/output.csv -w 4
```
Replace `/path/to/binaries` with the directory containing ELF files,
and `/path/to/output` with the directory where you want to save the output CSV.

---

## ⚙️ Usage
Main command-line flags:

```text
Usage of ./elfGetArgs:

Required flags:
  -p, --path string       Target directory or file to scan
  -A, --arch string       Target architecture (amd64, arm64)
  -f, --func string       Function name to search (can repeat, at least one required)
  -a, --arg int           Argument index for corresponding function (can repeat, start from 0, at least one required)

Optional flags:
  -o, --out string        Output CSV file (optional)
  -w, --worker int        Number of workers (optional, default 4)
```

### Example

```bash
  ./elfGetArgs -p ./binaries -A amd64 -f open -a 0 -o output.csv -w 8
```

This scans the `./binaries` directory for ELF files, looks for the first argument of the `open` function, outputs results to `output.csv`, and uses 8 worker threads.