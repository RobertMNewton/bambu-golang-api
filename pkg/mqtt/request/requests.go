package request

import (
	"fmt"
	"time"
)

func CreateGetVersionRequest(sequenceID string) Request {
	return CreateRequest("info", "get_version", sequenceID, nil)
}

func CreatePushAllRequest(sequenceID string) Request {
	params := map[string]interface{}{
		"version":     1,
		"push_target": 1,
	}
	return CreateRequest("pushing", "pushall", sequenceID, params)
}

func UpgradeConfirmRequest(sequenceID string) Request {
	params := map[string]interface{}{
		"src_id": 1,
	}
	return CreateRequest("upgrade", "upgrade_confirm", sequenceID, params)
}

func UpgradeConsistencyConfirmRequest(sequenceID string) Request {
	params := map[string]interface{}{
		"src_id": 1,
	}
	return CreateRequest("upgrade", "consistency_confirm", sequenceID, params)
}

func UpgradeStartRequest(sequenceID, url, module, version string) Request {
	params := map[string]interface{}{
		"src_id":  1,
		"url":     url,
		"module":  module,
		"version": version,
	}
	return CreateRequest("upgrade", "start", sequenceID, params)
}

func UpgradeGetHistoryRequest(sequenceID string) Request {
	return CreateRequest("upgrade", "get_history", sequenceID, nil)
}

func CreateStopPrintRequest(sequenceID string) Request {
	params := map[string]interface{}{
		"param": "",
	}
	return CreateRequest("print", "stop", sequenceID, params)
}

func CreatePausePrintRequest(sequenceID string) Request {
	params := map[string]interface{}{
		"param": "",
	}
	return CreateRequest("print", "pause", sequenceID, params)
}

func CreateResumePrintRequest(sequenceID string) Request {
	params := map[string]interface{}{
		"param": "",
	}
	return CreateRequest("print", "resume", sequenceID, params)
}

func CreateAMSChangeFilamentRequest(sequenceID string, target int, currTemp, tarTemp float64) Request {
	params := map[string]interface{}{
		"target":    target,
		"curr_temp": currTemp,
		"tar_temp":  tarTemp,
	}
	return CreateRequest("print", "ams_change_filament", sequenceID, params)
}

func AmsUserSettingRequest(sequenceID string, amsID int, startupReadOption, trayReadOption bool) Request {
	params := map[string]interface{}{
		"ams_id":              amsID,
		"startup_read_option": startupReadOption,
		"tray_read_option":    trayReadOption,
	}
	return CreateRequest("print", "ams_user_setting", sequenceID, params)
}

func PrintAmsFilamentSettingRequest(sequenceID string, amsID, trayID int, trayInfoIdx, trayColor, trayType string, nozzleTempMin, nozzleTempMax int) Request {

	params := map[string]interface{}{
		"ams_id":          amsID,
		"tray_id":         trayID,
		"tray_info_idx":   trayInfoIdx,
		"tray_color":      trayColor,
		"nozzle_temp_min": nozzleTempMin,
		"nozzle_temp_max": nozzleTempMax,
		"tray_type":       trayType,
	}
	return CreateRequest("print", "ams_filament_setting", sequenceID, params)
}

func PrintAmsControlRequest(sequenceID string, param string) Request {
	params := map[string]interface{}{
		"param": param,
	}
	return CreateRequest("print", "ams_control", sequenceID, params)
}

func CreatePrintSpeedRequest(sequenceID string, speedLevel int) Request {
	params := map[string]interface{}{
		"param": fmt.Sprintf("%d", speedLevel),
	}
	return CreateRequest("print", "print_speed", sequenceID, params)
}

func CreateGCodeFileRequest(sequenceID, filename string) Request {
	params := map[string]interface{}{
		"param": filename,
	}
	return CreateRequest("print", "gcode_file", sequenceID, params)
}

func CreateGCodeLineRequest(sequenceID, gcode string) Request {
	params := map[string]interface{}{
		"param": gcode,
	}
	return CreateRequest("print", "gcode_line", sequenceID, params)
}

func CreateCalibrationRequest(sequenceID string) Request {
	return CreateRequest("print", "calibration", sequenceID, nil)
}

func CreateUnloadFilamentRequest(sequenceID string) Request {
	return CreateRequest("print", "unload_filament", sequenceID, nil)
}

func CreateProjectFileRequest(sequenceID, filename string) Request {
	params := map[string]interface{}{
		"param":          filename,
		"project_id":     "0",
		"profile_id":     "0",
		"task_id":        "0",
		"subtask_id":     "0",
		"subtask_name":   "",
		"file":           "",
		"url":            "file:///mnt/sdcard",
		"md5":            "",
		"timelapse":      true,
		"bed_type":       "auto",
		"bed_levelling":  true,
		"flow_cali":      true,
		"vibration_cali": true,
		"layer_inspect":  true,
		"ams_mapping":    "",
		"use_ams":        false,
	}
	return CreateRequest("print", "project_file", sequenceID, params)
}

func CreateSkipObjectsRequest(sequenceID string, objList []int) Request {
	params := map[string]interface{}{
		"timestamp": time.Now().UnixMilli(),
		"obj_list":  objList,
	}
	return CreateRequest("print", "skip_objects", sequenceID, params)
}

// Example function to create a "ledctrl" request
func CreateLEDControlRequest(sequenceID, ledNode, ledMode string, ledOnTime, ledOffTime, loopTimes, intervalTime int) Request {
	params := map[string]interface{}{
		"led_node":      ledNode,
		"led_mode":      ledMode,
		"led_on_time":   ledOnTime,
		"led_off_time":  ledOffTime,
		"loop_times":    loopTimes,
		"interval_time": intervalTime,
	}
	return CreateRequest("system", "ledctrl", sequenceID, params)
}

// Primarily Used for Cloud Printers
func CreateGetAccessCodeRequest(sequenceID string) Request {
	return CreateRequest("system", "get_access_code", sequenceID, nil)
}

// control can be either 'enable' or 'disable'
func CreateIPCamRecordSetRequest(sequenceID, control string) Request {
	params := map[string]interface{}{
		"control": control,
	}
	return CreateRequest("camera", "ipcam_record_set", sequenceID, params)
}

// control can be either 'enable' or 'disable'
func CreateIPCamTimelapseRequest(sequenceID, control string) Request {
	params := map[string]interface{}{
		"control": control,
	}
	return CreateRequest("camera", "ipcam_timelapse", sequenceID, params)
}

func CreateXCamControlSetRequest(sequenceID, moduleName string, control, printHalt bool) Request {
	params := map[string]interface{}{
		"module_name": moduleName,
		"control":     control,
		"print_halt":  printHalt,
	}
	return CreateRequest("xcam", "xcam_control_set", sequenceID, params)
}
