package cmd

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/spf13/cobra"
)

func indexOf[T comparable](collection []T, el T) int {
	for i, x := range collection {
		if x == el {
			return i
		}
	}
	return -1
}

type EffectType string

const (
	BETA EffectType = "BETA"
	OR   EffectType = "OR"
)

func flipBeta(beta float64) float64 {
	return -1 * beta
}

func flipOR(or float64) float64 {
	return 1 / or
}

type Alleles struct {
	Effect string
	Other  string
}

func parseSumstatsFileToMap(
	filename,
	SNPFieldName,
	effectAlleleFieldName,
	otherAlleleFieldName string) (map[string]Alleles, error) {

	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = '\t'

	referenceSNPMapping := make(map[string]Alleles)

	header, err := reader.Read()
	if err != nil {
		return nil, err
	}

	SNPIndex := indexOf(header, SNPFieldName)
	effectAlleleIndex := indexOf(header, effectAlleleFieldName)
	otherAlleleIndex := indexOf(header, otherAlleleFieldName)

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}

		snp := record[SNPIndex]
		referenceSNPMapping[snp] = Alleles{
			Effect: record[effectAlleleIndex],
			Other:  record[otherAlleleIndex],
		}
	}

	return referenceSNPMapping, nil
}

func processAndWriteFlippedStats(
	inputFilename,
	outputFilename,
	SNPFieldName,
	effectAlleleFieldName,
	otherAlleleFieldName,
	effectFieldName,
	effectType string,
	referenceSNPMapping map[string]Alleles,
) error {

	var flippingFunction func(effect float64) float64
	switch effectType {
	case string(BETA):
		flippingFunction = flipBeta
	case string(OR):
		flippingFunction = flipOR
	default:
		return fmt.Errorf("unknown effect type: %s", effectType)
	}

	inFile, err := os.Open(inputFilename)
	if err != nil {
		return err
	}
	defer inFile.Close()

	outFile, err := os.Create(outputFilename)
	if err != nil {
		return err
	}
	defer outFile.Close()

	reader := csv.NewReader(inFile)
	reader.Comma = '\t'

	writer := csv.NewWriter(outFile)
	writer.Comma = '\t'

	header, err := reader.Read()
	if err != nil {
		return err
	}

	SNPIndex := indexOf(header, SNPFieldName)
	effectAlleleIndex := indexOf(header, effectAlleleFieldName)
	otherAlleleIndex := indexOf(header, otherAlleleFieldName)
	effectIndex := indexOf(header, effectFieldName)

	err = writer.Write(header)
	if err != nil {
		return err
	}

	var snp, effectAllele, otherAllele string
	var effect float64

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		snp = record[SNPIndex]
		effectAllele = record[effectAlleleIndex]
		otherAllele = record[otherAlleleIndex]
		effect, err = strconv.ParseFloat(record[effectIndex], 64)
		if err != nil {
			return err
		}

		if referenceAlleles, ok := referenceSNPMapping[snp]; ok {
			if effectAllele == referenceAlleles.Other && otherAllele == referenceAlleles.Effect {

				fmt.Println("Flipping alleles for SNP", snp)
				record[effectAlleleIndex] = otherAllele
				record[otherAlleleIndex] = effectAllele
				record[effectIndex] = fmt.Sprintf("%.7f", flippingFunction(effect))

			}
		}

		err = writer.Write(record)
		if err != nil {
			return err
		}
	}

	writer.Flush()
	return nil
}

func flipAlleles(
	referenceFilename,
	referenceSNPFieldName,
	referenceEffectAlleleFieldName,
	referenceOtherAlleleFieldName,
	sumstatsFilename,
	outputFilename,
	sumstatsSNPFieldName,
	sumstatsEffectAlleleFieldName,
	sumstatsOtherAlleleFieldName,
	sumstatsEffectFieldName,
	effectType string,
) error {

	referenceSNPMapping, err := parseSumstatsFileToMap(
		referenceFilename,
		referenceSNPFieldName,
		referenceEffectAlleleFieldName,
		referenceOtherAlleleFieldName,
	)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	err = processAndWriteFlippedStats(
		sumstatsFilename,
		outputFilename,
		sumstatsSNPFieldName,
		sumstatsEffectAlleleFieldName,
		sumstatsOtherAlleleFieldName,
		sumstatsEffectFieldName,
		effectType,
		referenceSNPMapping,
	)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return nil
}

var flipallelesCmd = &cobra.Command{
	Use:   "flipalleles",
	Short: "Flip alleles in a summary statistics file",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		sumstatsFilename, _ := cmd.Flags().GetString("sumstats")
		sumstatsEffectAlleleFieldName, _ := cmd.Flags().GetString("sumstats-effect-allele")
		sumstatsOtherAlleleFieldName, _ := cmd.Flags().GetString("sumstats-other-allele")
		sumstatsSNPFieldName, _ := cmd.Flags().GetString("sumstats-snp")
		sumstatsEffectFieldName, _ := cmd.Flags().GetString("sumstats-effect")
		effectType, _ := cmd.Flags().GetString("effect-type")

		referenceFilename, _ := cmd.Flags().GetString("reference")
		referenceEffectAlleleFieldName, _ := cmd.Flags().GetString("reference-effect-allele")
		referenceOtherAlleleFieldName, _ := cmd.Flags().GetString("reference-other-allele")
		referenceSNPFieldName, _ := cmd.Flags().GetString("reference-snp")

		outputFilename, _ := cmd.Flags().GetString("output")

		flipAlleles(
			referenceFilename,
			referenceSNPFieldName,
			referenceEffectAlleleFieldName,
			referenceOtherAlleleFieldName,
			sumstatsFilename,
			outputFilename,
			sumstatsSNPFieldName,
			sumstatsEffectAlleleFieldName,
			sumstatsOtherAlleleFieldName,
			sumstatsEffectFieldName,
			effectType,
		)
	},
}

func init() {
	rootCmd.AddCommand(flipallelesCmd)

	flipallelesCmd.Flags().StringP("sumstats", "", "", "Summary statistics file")
	flipallelesCmd.Flags().StringP("sumstats-effect-allele", "", "", "Effect allele field name in summary statistics file")
	flipallelesCmd.Flags().StringP("sumstats-other-allele", "", "", "Other allele field name in summary statistics file")
	flipallelesCmd.Flags().StringP("sumstats-snp", "", "", "SNP field name in summary statistics file")
	flipallelesCmd.Flags().StringP("sumstats-effect", "", "", "Effect field name in summary statistics file")
	flipallelesCmd.Flags().StringP("effect-type", "", "", "Effect type (beta or odds ratio)")

	flipallelesCmd.Flags().StringP("reference", "", "", "Reference file")
	flipallelesCmd.Flags().StringP("reference-effect-allele", "", "", "Effect allele field name in reference file")
	flipallelesCmd.Flags().StringP("reference-other-allele", "", "", "Other allele field name in reference file")
	flipallelesCmd.Flags().StringP("reference-snp", "", "", "SNP field name in reference file")

	flipallelesCmd.Flags().StringP("output", "", "", "Output file")
}
