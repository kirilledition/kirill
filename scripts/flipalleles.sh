#!/bin/bash

build/kirill flipalleles \
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
