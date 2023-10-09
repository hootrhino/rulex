package cron_task

// CRON_TASK_TYPE
const (
	CRON_TASK_TYPE_LINUX_SHELL = 1
	CRON_TASK_TYPE_WINDOWS_CMD = CRON_TASK_TYPE_LINUX_SHELL + 1
)

const PERM_0777 = 0777

const CRON_ASSETS = "cron_assets"

const CRON_LOGS = "cron_logs"
