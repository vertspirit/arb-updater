package arb_update

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"sync"
)

type ArbUpdate struct {
	TemplateFile string
	LocaleFile   string
	OutputFile   string
}

type ArbEntries struct {
	TemplateEntries map[string]interface{}
	LocaleEntries   map[string]interface{}
}

func NewArbUpdate(tempFile, localeFile, outFile string) *ArbUpdate {
	return &ArbUpdate{
		TemplateFile: tempFile,
		LocaleFile:   localeFile,
		OutputFile:   outFile,
	}
}

func (au *ArbUpdate) ReadAllEntries() (*ArbEntries, error) {
	ent := &ArbEntries{}

	err := ReadArbEntriesFromFile(au.TemplateFile, &ent.TemplateEntries)
	if err != nil {
		return nil, err
	}

	err = ReadArbEntriesFromFile(au.LocaleFile, &ent.LocaleEntries)
	if err != nil {
		return nil, err
	}

	return ent, nil
}

func (ae *ArbEntries) Merge(ctx context.Context, maxWorkers int, sortByKey bool) ([]byte, error) {
	var merged *sync.Map

	c, cancel := context.WithCancel(ctx)

	wp := NewWorkerPool(maxWorkers)
	go wp.Allocate(ae.TemplateEntries)
	go func() {
		merged = wp.ReadData(c)
		defer cancel()
	}()
	wp.Run(ae.LocaleEntries)
	<- wp.Done

	for {
		select {
		case <-c.Done():
			return MarshalJSON(merged, sortByKey)
		}
	}
}

func (ae *ArbEntries) MergeFullOn(sortByKey bool) ([]byte, error) {
	merged := MergeEntriesFullOn(ae.TemplateEntries, ae.LocaleEntries)
	return MarshalJSON(merged, sortByKey)
}

func MarshalJSON(m *sync.Map, sortByKey bool) ([]byte, error) {
	data := make(map[string]interface{})
	m.Range(func(k, v interface{}) bool {
		data[k.(string)] = v
		return true
	})
	if sortByKey {
		keys := make([]string, 0, len(data))
		for k := range data {
			keys = append(keys, k)
		}

		sort.Strings(keys)

		sortedData := make(map[string]interface{})
		for _, k := range keys {
			sortedData[k] = data[k]
		}

		return json.MarshalIndent(sortedData, "", "    ")
	}
	return json.MarshalIndent(data, "", "    ")
}

func ReadArbEntriesFromFile(file string, data *map[string]interface{}) error {
	f, err := os.Open(file)
	if err != nil {
		return err
	}
	defer f.Close()

	decTemp := json.NewDecoder(f)
	if err = decTemp.Decode(data); err != nil {
		return err
	}
	return nil
}

func PrintEntries(data []byte) {
	fmt.Println(string(data))
}

func SaveEntriesToFile(au *ArbUpdate, data []byte) error {
	if au.LocaleFile == au.OutputFile {
		err := os.Rename(au.LocaleFile, fmt.Sprintf("%s.bak", au.LocaleFile))
		if err != nil {
			return err
		}
	}

	dir := filepath.Dir(au.OutputFile)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err = os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	}

	if err := os.WriteFile(au.OutputFile, data, 0644); err != nil {
		return err
	}
	return nil
}

func CompareEntries(localeEntries map[string]interface{}, key string) bool {
	for ent, _ := range localeEntries {
		if ent == key {
			return true
		}
	}
	return false
}

type entryItem struct {
	entry   string
	content interface{}
}

// run merge job without worker pool
func MergeEntriesFullOn(templateEntries, localeEntries map[string]interface{}) *sync.Map {
	merged := new(sync.Map)
	wg := sync.WaitGroup{}
	ch := make(chan entryItem, 64)
	done := make(chan bool)

	if len(templateEntries) > 0 && len(localeEntries) > 0 {
		for k, _ := range templateEntries {
			wg.Add(1)
			go func(ch chan <- entryItem, key string) {
				defer wg.Done()
				check := CompareEntries(localeEntries, key)
				ent := entryItem { entry:   key }
				if check {
					ent.content = localeEntries[key]
				} else {
					ent.content = templateEntries[key]
				}
				ch <- ent
			}(ch, k)
		}
		go func(mgd *sync.Map) {
			for {
				select {
				case m, ok := <- ch:
					if !ok {
						done <-true
						break
					} else {
						mgd.Store(m.entry, m.content)
					}
				default:
				}
			}
		}(merged)
		wg.Wait()
		close(ch)
	}

	<-done
	return merged
}
