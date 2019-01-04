package main

import "net/http"

const (
	healthOK   = `[{"epoch":"1546598830","timestamp":"10:47:10","cluster":"docker-cluster","status":"green","node.total":"1","node.data":"1","shards":"0","pri":"0","relo":"0","init":"0","unassign":"0","pending_tasks":"0","max_task_wait_time":"-","active_shards_percent":"100.0%"}]`
	healthFail = `[{"epoch":"1546603983","timestamp":"12:13:03","cluster":"docker-cluster","status":"red","node.total":"1","node.data":"1","shards":"0","pri":"0","relo":"0","init":"0","unassign":"9","pending_tasks":"0","max_task_wait_time":"-","active_shards_percent":"0.0%"}]`

	settingsOK   = `{"acknowledged":true,"persistent":{"cluster":{"routing":{"allocation":{"enable":"none"}}}},"transient":{}}`
	settingsFAIL = `{"acknowledged":false,"persistent":{"cluster":{"routing":{"allocation":{"enable":"none"}}}},"transient":{}}`

	flushOK   = `{"_shards":{"total":0,"successful":0,"failed":0}}`
	flushFAIL = `{"_shards":{"total":9,"successful":0,"failed":9},"test":{"total":9,"successful":0,"failed":9,"failures":[{"shard":0,"reason":"no active shards"},{"shard":1,"reason":"no active shards"},{"shard":2,"reason":"no active shards"}]}}`
)

func esHealthHandler(w http.ResponseWriter, r *http.Request) {
	if r.Context().Value(isFailing).(bool) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(healthFail))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(healthOK))
}

func esSettingsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Context().Value(isFailing).(bool) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(settingsFAIL))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(settingsOK))
}

func esFlushHandler(w http.ResponseWriter, r *http.Request) {
	if r.Context().Value(isFailing).(bool) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(flushFAIL))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(flushOK))
}
