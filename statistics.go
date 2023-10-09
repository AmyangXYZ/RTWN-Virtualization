package main

type Statistics struct {
	Network              `json:"-"`
	PartitionSuccessApps []int              `json:"partition_success_apps"`
	PartitionSuccess     bool               `json:"partition_success"`
	TaskSuccessApps      []int              `json:"task_success_apps"`
	TaskSuccess          bool               `json:"task_success"`
	SupplyFunc           map[int]SupplyStat `json:"supply_functions"`
}

type SupplyStat struct {
	ActualSupply [][2]int `json:"actual"` // [slot id, allocated num_tx]
	IdealSupply  [][2]int `json:"ideal"`
}

type TaskStat struct {
	AppID           int   `json:"app_id"`
	TaskID          int   `json:"task_id"`
	Period          int   `json:"period"`
	FinishTime      []int `json:"finish_time"`
	IdealFinishTime []int `json:"ideal_finish_time"`
}

func (nw *Network) NewStat() *Statistics {
	return &Statistics{
		Network:    *nw,
		SupplyFunc: make(map[int]SupplyStat),
	}
}

// collect result
func (nw *Network) CollectStat() {
	if nw.Settings.Verbose {
		nw.logger.Println("Collecting results...")
	}
	nw.Manager.Report()
	nw.Stat.TaskSuccess = true
	for _, app := range nw.Apps {
		if !app.Report() {
			nw.Stat.TaskSuccess = false
		}
	}
}
