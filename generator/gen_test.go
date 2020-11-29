package generator

import (
	"generator/lib"
	"generator/testhelper"
	"testing"
	"time"
)

func TestStart(t *testing.T)  {
	server := testhelper.NewTcpServer()
	defer server.Close()
	serverAddr := "127.0.0.1:8000"
	t.Logf("Startup TCP server(%s)...\n", serverAddr)
	err := server.Listen(serverAddr)
	if err != nil {
		t.Fatalf("TCP Server startup failing! (addr=%s)!\n", serverAddr)
	}

	pset := Param{
		Caller:     testhelper.NewTcpClient(serverAddr),
		TimeoutNS:  50 * time.Millisecond,
		LPS:        uint32(1000),
		DurationNS: 10 * time.Second,
		ResultCh:   make(chan * lib.CallResult, 50),
	}
	t.Logf("Initialize load generator (timeoutNS=%v, lps=%d, durationNS=%v)...",
		pset.TimeoutNS, pset.LPS, pset.DurationNS)

	gen, err := NewGenerator(pset)
	if err != nil {
		t.Fatalf("Load generator initialization failing: %s\n",
			err)
	}

	t.Log("Start load generator...")
	gen.Start()

	retMap := make(map[lib.RetCode] int)
	for r := range pset.ResultCh {
		retMap[r.Code] += 1
	}

	var total int
	for k,v := range  retMap {
		codePlain := lib.GetRetCodePlain(k)
		t.Logf("  Code plain: %s (%d), Count: %d.\n",
			codePlain, k, v)
		total += v
	}
	t.Logf("Total: %d.\n", total)

	successCount := retMap[lib.RET_CODE_SUCCESS]
	tps := float64(successCount) / float64(pset.DurationNS/1e9)
	t.Logf("Loads per second: %d; Treatments per second: %f.\n", pset.LPS, tps)
}