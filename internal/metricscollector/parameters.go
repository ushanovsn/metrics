package metricscollector

// used gauge parameters names
var mGaugeNames map[string]bool = map[string]bool {
	"Alloc": true,
	"BuckHashSys": true,
	"Frees": true,
	"GCCPUFraction": true,
	"GCSys": true,
	"HeapAlloc": true,
	"HeapIdle": true,
	"HeapInuse": true,
	"HeapObjects": true,
	"HeapReleased": true,
	"HeapSys": true,
	"LastGC": true,
	"Lookups": true,
	"MCacheInuse": true,
	"MCacheSys": true,
	"MSpanInuse": true,
	"MSpanSys": true,
	"Mallocs": true,
	"NextGC": true,
	"NumForcedGC": true,
	"NumGC": true, 
	"OtherSys": true,
	"PauseTotalNs": true,
	"StackInuse": true,
	"StackSys": true,
	"Sys": true,
	"TotalAlloc": true,
}



// used counters parameters names
var mCounterNames map[string]bool = map[string]bool {
	"PollCount": true,
}



// used gauge randoms parameters names
var mRandomNames map[string]bool = map[string]bool {
	"RandomValue": true,
}

