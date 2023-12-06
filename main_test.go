package main

import (
	"encoding/json"
	"testing"
	arb "arb-updater/pkg/arb"
)

const (
	templateFile string = "test-data/intl_en.arb"
	localeFile   string = "test-data/intl_zh_Hant.arb"
	maxCorutines int    = 3000
	isSort       bool   = true
)

func TestUpdateArb(t *testing.T) {
	t.Run("Test update arb file with 3000 corutines", func(t *testing.T) {
		update := arb.NewArbUpdate(templateFile, localeFile, "test/unit_tested_with_limit.arb")

		ents, err := update.ReadAllEntries()
		if err != nil {
			t.Errorf("Unit Test (Load Arb Files) Fail: %v\n", err)
		}

		merged, err := ents.Merge(maxCorutines, isSort)
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
	})

	t.Run("Test update arb file without max corutines limit", func(t *testing.T) {
		update := arb.NewArbUpdate(templateFile, localeFile, "test/unit_tested_without_limit.arb")

		ents, err := update.ReadAllEntries()
		if err != nil {
			t.Errorf("Unit Test (Load Arb Files) Fail: %v\n", err)
		}

		merged, err := ents.MergeFullOn(isSort)
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
	})
}

