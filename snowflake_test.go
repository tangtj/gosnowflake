package gosnowflake

import (
	"testing"
)

func TestSnowflake_NextId(t *testing.T) {
	sn := NewSnowFlake(10)
	for i := 0; i < 100000; i++ {
		sn.NextId()
	}

}

func BenchmarkName(b *testing.B) {
	sn := NewSnowFlake(10)
	for i := 0; i < b.N; i++ {

		go func() {
			for i := 0; i < 100000; i++ {
				sn.NextId()
			}
		}()

	}
}

func TestSnowflake_getNextMill(t *testing.T) {
	sn := NewSnowFlake(10)
	n := now()
	sn.lastTime = n
	if sn.getNextMill() == n {
		t.Error("获取下一秒异常")
	}
}
