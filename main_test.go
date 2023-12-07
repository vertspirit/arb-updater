package main

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	arb "arb-updater/pkg/arb"
)

const (
	templateFile string = "test-data/intl_en.arb"
	localeFile   string = "test-data/intl_zh_Hant.arb"
)

func updateArbTestingSortedLimited(t *testing.T) {
	update := arb.NewArbUpdate(templateFile, localeFile, "test/unit_tested_with_limit.arb")

	ents, err := update.ReadAllEntries()
	if err != nil {
		t.Errorf("Unit Test (Load Arb Files) Fail: %v\n", err)
	}

	ctx := context.TODO()
	merged, err := ents.Merge(ctx, 3000, true)
	if err != nil {
		t.Errorf("Unit Test (Merge Arb Entries) Fail: %v\n", err)
	}

	var data map[string]interface{}
	if err = json.Unmarshal(merged, &data); err != nil {
		t.Errorf("Unit Test (Parse Updated Entries) Fail: %v\n", err)
	}

	if data["appName"].(string) != "This as translated entry" {
		t.Errorf(
			"Unit Test (Check Updated Entries) Fail: %s\n",
			"The content of existing entry from origin arb file is not preserved.",
		)
	}

	val, ok := data["email"]
	if !ok {
		t.Errorf(
			"Unit Test (Check Updated Entries) Fail: %s\n",
			"The new entry from template arb file is not merged.",
		)
	}
	if val.(string) != "Email" {
		t.Errorf(
			"Unit Test (Check Updated Entries) Fail: %s\n",
			"The content of new entry from template arb file is not correct.",
		)
	}

	if err = arb.SaveEntriesToFile(update, merged); err != nil {
		t.Errorf("Unit Test (Write Updated Entries to File) Fail: %v\n", err)
	}
}

func updateArbTestingSortedUnLimited(t *testing.T) {
	update := arb.NewArbUpdate(templateFile, localeFile, "test/unit_tested_without_limit.arb")

	ents, err := update.ReadAllEntries()
	if err != nil {
		t.Errorf("Unit Test (Load Arb Files) Fail: %v\n", err)
	}

	merged, err := ents.MergeFullOn(true)
	if err != nil {
		t.Errorf("Unit Test (Merge Arb Entries - FullOn) Fail: %v\n", err)
	}

	var data map[string]interface{}
	if err = json.Unmarshal(merged, &data); err != nil {
		t.Errorf("Unit Test (Parse Updated Entries - FullOn) Fail: %v\n", err)
	}

	if data["appName"].(string) != "This as translated entry" {
		t.Errorf(
			"Unit Test (Check Updated Entries - FullOn) Fail: %s\n",
			"The content of existing entry from origin arb file is not preserved.",
		)
	}

	val, ok := data["email"]
	if !ok {
		t.Errorf(
			"Unit Test (Check Updated Entries - FullOn) Fail: %s\n",
			"The new entry from template arb file is not merged.",
		)
	}
	if val.(string) != "Email" {
		t.Errorf(
			"Unit Test (Check Updated Entries - FullOn) Fail: %s\n",
			"The content of new entry from template arb file is not correct.",
		)
	}

	if err = arb.SaveEntriesToFile(update, merged); err != nil {
		t.Errorf("Unit Test (Write Updated Entries to File) Fail: %v\n", err)
	}
}

func TestUpdateArb(t *testing.T) {
	t.Run("Test update arb file with 3000 corutines", updateArbTestingSortedLimited)
	t.Run("Test update arb file without max corutines limit", updateArbTestingSortedUnLimited)
}

func updateArbLimited(maxWorkers int, sorted bool) {
	update := arb.NewArbUpdate(templateFile, localeFile, "test/bechmark_sorted_limited.arb")

	ents, err := update.ReadAllEntries()
	if err != nil {
		fmt.Printf("Benchmark Fail (sorted and limited): %v\n", err)
	}

	ctx := context.TODO()
	_, err = ents.Merge(ctx, maxWorkers, sorted)
	if err != nil {
		fmt.Printf("Benchmark Fail (sorted and limited): %v\n", err)
	}
}

func updateArbUnLimited(sorted bool) {
	update := arb.NewArbUpdate(templateFile, localeFile, "test/bechmark_sorted_limited.arb")

	ents, err := update.ReadAllEntries()
	if err != nil {
		fmt.Printf("Benchmark Fail (sorted and limited): %v\n", err)
	}

	_, err = ents.MergeFullOn(sorted)
	if err != nil {
		fmt.Printf("Benchmark Fail (sorted and limited): %v\n", err)
	}
}

func BenchmarkUpdateArbLimitedSorted(b *testing.B) {
	maxWorkers := 3000
	for n := 0; n < b.N; n++ {
		updateArbLimited(maxWorkers, true)
	}
}

func BenchmarkUpdateArbLimitedUnSorted(b *testing.B) {
	maxWorkers := 3000
	for n := 0; n < b.N; n++ {
		updateArbLimited(maxWorkers, false)
	}
}

func BenchmarkUpdateArbUnLimitedSorted(b *testing.B) {
	for n := 0; n < b.N; n++ {
		updateArbUnLimited(true)
	}
}

func BenchmarkUpdateArbUnLimitedUnSorted(b *testing.B) {
	for n := 0; n < b.N; n++ {
		updateArbUnLimited(false)
	}
}
