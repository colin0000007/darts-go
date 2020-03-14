package dat

import (
	"fmt"
	"log"
	"math/rand"
	"sort"
	"strconv"
	"testing"
	"time"
)

func TestFetch(t *testing.T) {
	dat := NewDoubleArrayTrie()
	cases := map[string][]string{
		"ex1_pass":  []string{"一举", "一举一动", "一举成名", "万能", "万能胶"},
		"ex2_fail":  []string{"b", "a"},
		"ex3_pass":  []string{"a", "ac", "ab"},
		"ex4_pass:": []string{"ab", "abc", "be", "bfg", "c"},
		"ex5_pass":  []string{"a", "abcd", "d"},
	}
	for name, case_ := range cases {
		t.Run(name, func(t *testing.T) {
			err := dat.Build1(case_)
			if err != nil {
				t.Error(err)
			}
		})
	}
}

// 对于len和cap，初始时给定cap的容量就会在内存中分配cap这么大的数组
// 直到扩容时才会生成新的数组
func TestMake(t *testing.T) {
	list := make([]int, 0, 2)
	list = append(list, 1)
	list = append(list, 1)
	fmt.Println(list[1])
	fmt.Printf("%p %d %d\n", list, len(list), cap(list))
	list = append(list, 1)
	fmt.Printf("%p %d %d\n", list, len(list), cap(list))
	list = append(list, 2)
	fmt.Printf("%p %d %d\n", list, len(list), cap(list))
	list = append(list, 3)
	fmt.Printf("%p %d %d\n", list, len(list), cap(list))
}

func TestBuild(t *testing.T) {
	samples := makeSample(1000000, 3, 8)
	dat := NewDoubleArrayTrie()
	err := dat.Build1(samples)
	if err != nil {
		log.Fatal(err)
	}
	//log.Printf("%v", strings.Join(samples, "\n"))
	count := 100
	errCount := 0
	for i := 0; i < count; i++ {
		fmt.Print(samples[i], ":")
		id, ok := dat.IndexOf(samples[i])
		if !ok {
			errCount++
			fmt.Printf("-1,may common prefix or not exists, index = %d", id)
		} else {
			fmt.Println(samples[id] == samples[i])
		}
	}

	for i := count; i < len(samples); i++ {
		_, ok := dat.IndexOf(samples[i])
		if !ok {
			errCount++
			//fmt.Println(samples[i])
		}
	}
	log.Printf("build done, indexOf error num %d", errCount)
}

// 构建词库的字符集，会随机使用
var dict = [...]rune{
	'a', 'b', 'c', 'd', 'e', 'f', 'g',
	'h', 'i', 'j', 'l', 'm', 'n', 'o', 'p', 'q',
	'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z',
}

/*
	产生词库样本
	@param keyMinLen 最小key长度
	@param keyMaxLen 最大key长度
	@return []string 返回样本结果（排好序）
*/
func makeSample(keySize int, keyMinLen, keyMaxLen int) []string {
	rand.Seed(time.Now().Unix())
	keyMap := make(map[string]struct{}, keySize)
	keys := make([]string, keySize, keySize)
	dictLen := len(dict)
	kRang := keyMaxLen - keyMinLen + 1
	keyCount := 0
	for keyCount != keySize {
		kLen := rand.Intn(kRang) + keyMinLen
		rs := make([]rune, 0, kLen)
		for j := 0; j < kLen; j++ {
			rs = append(rs, dict[rand.Intn(dictLen)])
		}
		key := string(rs)
		if _, ok := keyMap[key]; !ok {
			keyMap[key] = struct{}{}
			keys[keyCount] = key
			keyCount++
		}
	}
	sort.Strings(keys)
	return keys
}

func TestMakeSample(t *testing.T) {
	samples := makeSample(1000, 3, 8)
	for _, sam := range samples {
		fmt.Println(len(sam), sam)
	}
}

func TestRand(t *testing.T) {
	for i := 0; i < 100; i++ {
		fmt.Print(rand.Intn(5)+1, ",")
	}
}

func BenchmarkLen(b *testing.B) {
	var sli = make([]byte, 1024*1024, 1024*1024)
	lenCount := 10
	for k := 0; k < lenCount; k++ {
		b.Run("len_"+fmt.Sprint(k), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = len(sli)
			}
		})
	}

}

// 测试保存和加载DAT
func TestStoreLoad(t *testing.T) {
	samples := makeSample(1000000, 3, 8)
	vals := make([]int, len(samples), len(samples))
	for i := 0; i < len(samples); i++ {
		vals[i] = i
	}
	dat := NewDoubleArrayTrie()
	err := dat.Build2(samples, vals)
	if err != nil {
		t.Fatal(err)
	}
	err = dat.Store("/Users/didi/Documents/go/test.dic")
	if err != nil {
		t.Fatal(err)
	}
	dat2 := NewDoubleArrayTrie()
	err = dat2.Load("/Users/didi/Documents/go/test.dic")
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 100; i++ {
		index, ok := dat2.IndexOf(samples[i])
		fmt.Print(samples[i] + ":")
		if ok {
			fmt.Print(i == index)
		} else {
			fmt.Print("not ok")
		}
		value := dat2.GetValue(samples[i])
		fmt.Println(",val:", value)
	}
	fmt.Println("allocSize:", dat.allocSize)
}

func TestExactMatch(t *testing.T) {
	dat := new(DoubleArrayTrie)
	dat.Build1([]string{"1", "2", "3"})
	index, _ := dat.IndexOf("2")
	fmt.Println(index)
}

// 测试直接声明的切片的len和cap
// 测试nil 切片是否调用len是否长度为0
func TestSlice(t *testing.T) {
	//var sli = []int{}
	var sli = make([]int, 0)
	t.Log("len:", len(sli), "cap:", cap(sli))
	var sli2 []*int = nil
	t.Log("nil len:", len(sli2))
}

type NilWriter struct {
}

func (w *NilWriter) Write(p []byte) (n int, err error) {
	return 0, nil
}

func BenchmarkBuild(b *testing.B) {
	sizes := []int64{
		1000000,
		2000000,
		3000000,
		4000000,
		5000000,
		6000000,
		7000000,
		8000000,
		9000000,
		10000000,
	}
	// 不打印日志
	log.SetOutput(&NilWriter{})
	for _, size := range sizes {
		b.Run("keySize_"+strconv.FormatInt(size, 10), func(b *testing.B) {
			samples := makeSample(int(size), 3, 8)
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				dat := NewDoubleArrayTrie()
				dat.Build1(samples)
			}
		})
	}
}
