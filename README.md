# 🦍 kirill: Yet another bioinformatics toolbox 

Kirill is a command-line interface (CLI) application that provides a collection of tools for bioinformatics. This repository contains the source code and documentation for the application. Kirill currently consists of two commands: `fetchpdb` and `flipalleles`.

## Installation

To install Kirill, you can download the source code from this repository and build the binary using the Go compiler. Make sure you have Go installed on your system.

```sh
git clone https://github.com/kirilledition/kirill.git
cd kirill
make build
```

This will create an executable binary named `kirill` in the `kirill/build` directory.

## Usage

### fetchpdb

`fetchpdb` is a command-line tool to download protein structures from the Protein Data Bank (PDB). It accepts a list of PDB IDs or an input file containing PDB IDs, one per line.

**Example usage:**

1. Download structures for a list of PDB IDs:

```sh
kirill fetchpdb 1abc 2def 3ghi
```

2. Download structures from an input file (each PDB ID on a separate line):

```sh
kirill fetchpdb pdb_ids.txt
```

3. Download structures from an input file and save them to a specific output directory:

```sh
kirill fetchpdb pdb_ids.txt -o /path/to/output
```

### 2. flipalleles

`flipalleles` is a command-line tool designed to process and modify genetic summary statistics data by flipping alleles and their corresponding effects according to a reference summary statistics file. The primary use case for this program is to harmonize the data from two separate summary statistics files, ensuring consistency in allele representation and effects direction.

**Example usage:**

```sh
kirill flipalleles \
	--sumstats test_data/wrong_alleles.tsv \
	--sumstats-effect-allele Allele1 \
	--sumstats-other-allele Allele2 \
	--sumstats-snp MarkerName \
	--sumstats-effect Zscore \
	--effect-type BETA \
	--reference test_data/reference.tsv \
	--reference-effect-allele A1 \
	--reference-other-allele A2 \
	--reference-snp SNP \
	--output test_data/flipped.tsv
```

## Contributing

Contributions to Kirill are welcome! If you would like to add new features or improve existing ones, please create a fork of this repository and submit a pull request.
