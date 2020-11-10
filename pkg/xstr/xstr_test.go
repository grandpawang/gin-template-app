package xstr

import (
	"testing"
)

func TestJoinInts(t *testing.T) {
	// test empty slice
	is := []uint{}
	s := JoinUints(is)
	if s != "" {
		t.Errorf("input:%v,output:%s,result is incorrect", is, s)
	} else {
		t.Logf("input:%v,output:%s", is, s)
	}
	// test len(slice)==1
	is = []uint{1}
	s = JoinUints(is)
	if s != "1" {
		t.Errorf("input:%v,output:%s,result is incorrect", is, s)
	} else {
		t.Logf("input:%v,output:%s", is, s)
	}
	// test len(slice)>1
	is = []uint{1, 2, 3}
	s = JoinUints(is)
	if s != "1,2,3" {
		t.Errorf("input:%v,output:%s,result is incorrect", is, s)
	} else {
		t.Logf("input:%v,output:%s", is, s)
	}
}

func TestSplitInts(t *testing.T) {
	// test empty slice
	s := ""
	is, err := SplitInts(s)
	if err != nil || len(is) != 0 {
		t.Error(err)
	}
	// test split uint
	s = "1,2,3"
	is, err = SplitInts(s)
	if err != nil || len(is) != 3 {
		t.Error(err)
	}
}

func BenchmarkJoinInts(b *testing.B) {
	is := make([]uint, 10000, 10000)
	for i := uint(0); i < 10000; i++ {
		is[i] = i
	}
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			JoinUints(is)
		}
	})
}
