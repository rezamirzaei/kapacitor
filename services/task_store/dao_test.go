package task_store_test

import (
	"testing"

	"github.com/influxdata/kapacitor/services/task_store"
)

func TestDBRPEqualAsSet(t *testing.T) {
	tt := []struct {
		name string
		ds   []task_store.DBRP
		bs   []task_store.DBRP
		eq   bool
	}{
		{
			name: "two sets different order", // TODO: better name
			ds: []task_store.DBRP{
				{"telegraf", "autogen"},
				{"telegraf", "not_autogen"},
			},
			bs: []task_store.DBRP{
				{"telegraf", "not_autogen"},
				{"telegraf", "autogen"},
			},
			eq: true,
		},
		{
			name: "unequal size", // TODO: better name
			ds: []task_store.DBRP{
				{"telegraf", "not_autogen"},
			},
			bs: []task_store.DBRP{
				{"telegraf", "not_autogen"},
				{"telegraf", "autogen"},
			},
			eq: false,
		},
		{
			name: "one element different rp", // TODO: better name
			ds: []task_store.DBRP{
				{"telegraf", "not_autogen"},
			},
			bs: []task_store.DBRP{
				{"telegraf", "autogen"},
			},
			eq: false,
		},
		{
			name: "one element different db", // TODO: better name
			ds: []task_store.DBRP{
				{"not_telegraf", "autogen"},
			},
			bs: []task_store.DBRP{
				{"telegraf", "autogen"},
			},
			eq: false,
		},
	}

	for _, tst := range tt {
		t.Run(tst.name, func(t *testing.T) {
			if exp, got := tst.eq, task_store.EqualAsSets(tst.ds, tst.bs); exp != got {
				t.Fatalf("Expected sets to be equal") // TODO: better message
			}
			if exp, got := tst.eq, task_store.EqualAsSets(tst.bs, tst.ds); exp != got {
				t.Fatalf("Expected sets to be equal") // TODO: better message
			}
		})
	}
}
