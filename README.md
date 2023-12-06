# Arb Updater

A simple tool to update locale arb file (flutter localizations) by template arb file. This tool will keep translated entries and add new entries to arb file.

## Build

Install go and make first, and run the following command to compile binary which placed in bin folder:

```bash
make build
```

And run the binary to check version:

```bash
./bin/arb-updater -v
```

## Usage

Run the arb-updater with -h parameter for details:

```bash
arb-updater -h
```

Update the locale arb file by arb template file (default is intl_en.arb) and print result in the console (dry run):

```bash
arb-updater -t path/to/template.arb -l path/to/locale.arb -sort -print-only
```

Update the locale arb file by arb template file with max 3000 corutines, and backup origin locale arb file as *.bak file before overwrite it:

```bash
arb-updater -t path/to/template.arb -l path/to/locale.arb -sort -c 3000
```

Update the locale arb file to specific file without overwrite it:

```bash
arb-updater -t path/to/template.arb -l path/to/locale.arb -o path/to/updated.arb -sort
```

Update the locale arb file without max corutines limit (or with '-c 0' parameter):

```bash
arb-updater -t path/to/template.arb -l path/to/locale.arb -sort -full-on
```
