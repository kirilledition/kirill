#!/bin/bash

build/kirill fetchpdb 7cel 2mwk --output test_data

build/kirill fetchpdb test_data/pdb_list.txt --output test_data