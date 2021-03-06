package dao

import (
	client "github.com/influxdata/influxdb/client/v2"
	"kr/paasta/monitoring/openstack/models"
	"kr/paasta/monitoring/openstack/utils"
	"fmt"
)


type InstanceDao struct {
	influxClient 	client.Client
}

func GetInstanceDao(influxClient client.Client) *InstanceDao {
	return &InstanceDao{
		influxClient: 	influxClient,
	}
}


//Instance의 현재 CPU사용률을 조회한다.
func (d InstanceDao) GetInstanceCpuUsageList(request models.DetailReq)(_ client.Response, errMsg models.ErrMessage){
	var errLogMsg string
	defer func() {
		if r := recover(); r != nil {
			errMsg = models.ErrMessage{
				"Message": errLogMsg ,
			}
		}
	}()

	cpuUsageSql := "select mean(value) as usage from \"cpu.utilization_norm_perc\"  where resource_id = '%s' ";

	var q client.Query
	if request.DefaultTimeRange != "" {

		cpuUsageSql += " and time > now() - %s  group by time(%s);"

		q = client.Query{
			Command:  fmt.Sprintf( cpuUsageSql,
				request.InstanceId, request.DefaultTimeRange, request.GroupBy),
			Database: models.MetricDBName,
		}
	}else{

		cpuUsageSql += " and time < now() - %s and time > now() - %s  group by time(%s);"

		q = client.Query{
			Command:  fmt.Sprintf( cpuUsageSql,
				request.InstanceId, request.TimeRangeFrom, request.TimeRangeTo, request.GroupBy),
			Database: models.MetricDBName,
		}
	}

	models.MonitLogger.Debug("GetInstanceCpuUsageList Sql==>", q)
	resp, err := d.influxClient.Query(q)

	return utils.GetError().CheckError(*resp, err)
}


//Instance의 현재 CPU사용률을 조회한다.
func (d InstanceDao) GetInstanceMemoryUsageList(request models.DetailReq)(_ client.Response, errMsg models.ErrMessage){
	var errLogMsg string
	defer func() {
		if r := recover(); r != nil {

			errMsg = models.ErrMessage{
				"Message": errLogMsg ,
			}
		}
	}()

	memoryTotalSql := "select mean(value) as usage from \"mem.total_mb\"  where resource_id = '%s' ";
	memoryFreeSql := "select mean(value) as usage from \"mem.free_mb\"  where resource_id = '%s' ";

	models.MonitLogger.Debug("defaultTimeRange: %s, timeRangeFrom: %s, timeRangeTo:%s", request.DefaultTimeRange, request.TimeRangeFrom, request.TimeRangeTo)

	var q client.Query
	if request.DefaultTimeRange != "" {

		memoryTotalSql += " and time > now() - %s  group by time(%s);"
		memoryFreeSql += " and time > now() - %s  group by time(%s);"

		q = client.Query{
			Command:  fmt.Sprintf( memoryTotalSql + memoryFreeSql,
				request.InstanceId, request.DefaultTimeRange, request.GroupBy,
				request.InstanceId, request.DefaultTimeRange, request.GroupBy),
			Database: models.MetricDBName,
		}
	}else{

		memoryTotalSql += " and time < now() - %s and time > now() - %s  group by time(%s);"
		memoryFreeSql  += " and time < now() - %s and time > now() - %s  group by time(%s);"

		q = client.Query{
			Command:  fmt.Sprintf( memoryTotalSql + memoryFreeSql,
				request.InstanceId, request.TimeRangeFrom, request.TimeRangeTo, request.GroupBy,
				request.InstanceId, request.TimeRangeFrom, request.TimeRangeTo, request.GroupBy),
			Database: models.MetricDBName,
		}
	}

	models.MonitLogger.Debug("GetInstanceMemoryUsageList Sql==>", q)
	resp, err := d.influxClient.Query(q)

	return utils.GetError().CheckError(*resp, err)
}

//Instance의 현재 CPU사용률을 조회한다.
func (d InstanceDao) GetInstanceCpuUsage(request models.InstanceReq)(_ client.Response, errMsg models.ErrMessage){

	var errLogMsg string
	defer func() {
		if r := recover(); r != nil {

			errMsg = models.ErrMessage{
				"Message": errLogMsg ,
			}
		}
	}()

	cpuUsageSql := "select value from \"cpu.utilization_norm_perc\"  where time > now() - 2m and resource_id = '%s' order by time desc limit 1";

	var q client.Query

	q = client.Query{
		Command:  fmt.Sprintf( cpuUsageSql,
			request.InstanceId),
		Database: models.MetricDBName,
	}

	models.MonitLogger.Debug("GetInstanceCpuUsage Sql======>", q)
	resp, err := d.influxClient.Query(q)
	if err != nil{
		errLogMsg = err.Error()
	}

	return utils.GetError().CheckError(*resp, err)
}


//Node의 현재 Total Memory을 조회한다.
func (d InstanceDao) GetInstanceTotalMemoryUsage(request models.InstanceReq)(_ client.Response, errMsg models.ErrMessage){

	var errLogMsg string
	defer func() {
		if r := recover(); r != nil {

			errMsg = models.ErrMessage{
				"Message": errLogMsg ,
			}
		}
	}()

	totalMemSql := "select value from \"mem.total_mb\" where time > now() - 2m and resource_id = '%s' order by time desc limit 1;"

	var q client.Query

	q = client.Query{
		Command:  fmt.Sprintf( totalMemSql ,
			request.InstanceId),
		Database: models.MetricDBName,
	}

	models.MonitLogger.Debugf("GetInstanceTotalMemoryUsage Sql======>", q)
	resp, err := d.influxClient.Query(q)
	if err != nil{
		errLogMsg = err.Error()
	}

	return utils.GetError().CheckError(*resp, err)
}

//Node의 현재 Total Memory을 조회한다.
func (d InstanceDao) GetInstanceFreeMemoryUsage(request models.InstanceReq)(_ client.Response, errMsg models.ErrMessage){

	var errLogMsg string
	defer func() {
		if r := recover(); r != nil {

			errMsg = models.ErrMessage{
				"Message": errLogMsg ,
			}
		}
	}()

	freeMemSql  := "select value from \"mem.free_mb\"  where time > now() - 2m and resource_id = '%s' order by time desc limit 1;"

	var q client.Query

	q = client.Query{
		Command:  fmt.Sprintf( freeMemSql ,
			request.InstanceId),
		Database: models.MetricDBName,
	}

	models.MonitLogger.Debugf("GetInstanceFreeMemoryUsage Sql======>", q)
	resp, err := d.influxClient.Query(q)
	if err != nil{
		errLogMsg = err.Error()
	}

	return utils.GetError().CheckError(*resp, err)
}


//Node의 disk io read Kbyte를 조회한다.
func (d InstanceDao) GetInstanceDiskIoKbyte(request models.DetailReq, gubun string)(_ client.Response, errMsg models.ErrMessage){

	var errLogMsg string
	defer func() {
		if r := recover(); r != nil {

			errMsg = models.ErrMessage{
				"Message": errLogMsg ,
			}
		}
	}()

	var diskUsageSql string

	if gubun == "read"{
		diskUsageSql  = "select sum(value)/1024 as usage from \"io.read_bytes_sec\"  where resource_id = '%s' "
	}else if gubun == "write"{
		diskUsageSql  = "select sum(value)/1024 as usage from \"io.write_bytes_sec\"  where resource_id = '%s' "
	}


	models.MonitLogger.Debugf("defaultTimeRange: %s, timeRangeFrom: %s, timeRangeTo:%s", request.DefaultTimeRange, request.TimeRangeFrom, request.TimeRangeTo)

	var q client.Query
	if request.DefaultTimeRange != "" {

		diskUsageSql += " and time > now() - %s  group by time(%s);"

		q = client.Query{
			Command:  fmt.Sprintf( diskUsageSql,
				request.InstanceId,  request.DefaultTimeRange, request.GroupBy),
			Database: models.MetricDBName,
		}
	}else{

		diskUsageSql += " and time < now() - %s and time > now() - %s  group by time(%s);"


		q = client.Query{
			Command:  fmt.Sprintf( diskUsageSql,
				request.InstanceId,  request.TimeRangeFrom, request.TimeRangeTo, request.GroupBy),
			Database: models.MetricDBName,
		}
	}
	models.MonitLogger.Debug("GetInstanceDiskIoReadKbyte Sql==>", q)
	resp, err := d.influxClient.Query(q)

	return utils.GetError().CheckError(*resp, err)

}


//Instance의 network io read Kbyte를 조회한다.
func (d InstanceDao) GetInstanceNetworkKbyte(request models.DetailReq , inOut string)(_ client.Response, errMsg models.ErrMessage){

	var errLogMsg string
	defer func() {
		if r := recover(); r != nil {

			errMsg = models.ErrMessage{
				"Message": errLogMsg ,
			}
		}
	}()

	var networkSql string
	if inOut == "in"{
		networkSql  = "select sum(value)/1024 as usage from \"net.in_bytes_sec\"  where resource_id = '%s' "
	}else if inOut == "out"{
		networkSql  = "select sum(value)/1024 as usage from \"net.out_bytes_sec\"  where resource_id = '%s'"
	}

	models.MonitLogger.Debugf("defaultTimeRange: %s, timeRangeFrom: %s, timeRangeTo:%s", request.DefaultTimeRange, request.TimeRangeFrom, request.TimeRangeTo)

	var q client.Query
	if request.DefaultTimeRange != "" {

		networkSql += " and time > now() - %s  group by time(%s);"

		q = client.Query{
			Command:  fmt.Sprintf( networkSql,
				request.InstanceId, request.DefaultTimeRange, request.GroupBy),
			Database: models.MetricDBName,
		}
	}else{

		networkSql += " and time < now() - %s and time > now() - %s  group by time(%s);"

		q = client.Query{
			Command:  fmt.Sprintf( networkSql,
				request.InstanceId, request.TimeRangeFrom, request.TimeRangeTo, request.GroupBy),
			Database: models.MetricDBName,
		}
	}
	models.MonitLogger.Debug("GetInstanceNetworkKbyte Sql==>", q)
	resp, err := d.influxClient.Query(q)

	return utils.GetError().CheckError(*resp, err)

}

//Instance의 network io read Kbyte를 조회한다.
func (d InstanceDao) GetInstanceNetworkPackets(request models.DetailReq , inOut string)(_ client.Response, errMsg models.ErrMessage){

	var errLogMsg string
	defer func() {
		if r := recover(); r != nil {

			errMsg = models.ErrMessage{
				"Message": errLogMsg ,
			}
		}
	}()

	var networkSql string
	if inOut == "in"{
		networkSql  = "select sum(value) as usage from \"net.in_packets_sec\"  where resource_id = '%s' "
	}else if inOut == "out"{
		networkSql  = "select sum(value) as usage from \"net.out_packets_sec\"  where resource_id = '%s'"
	}

	models.MonitLogger.Debugf("defaultTimeRange: %s, timeRangeFrom: %s, timeRangeTo:%s", request.DefaultTimeRange, request.TimeRangeFrom, request.TimeRangeTo)

	var q client.Query
	if request.DefaultTimeRange != "" {

		networkSql += " and time > now() - %s  group by time(%s);"

		q = client.Query{
			Command:  fmt.Sprintf( networkSql,
				request.InstanceId, request.DefaultTimeRange, request.GroupBy),
			Database: models.MetricDBName,
		}
	}else{

		networkSql += " and time < now() - %s and time > now() - %s  group by time(%s);"

		q = client.Query{
			Command:  fmt.Sprintf( networkSql,
				request.InstanceId, request.TimeRangeFrom, request.TimeRangeTo, request.GroupBy),
			Database: models.MetricDBName,
		}
	}
	models.MonitLogger.Debug("GetInstanceNetworkPackets Sql==>", q)
	resp, err := d.influxClient.Query(q)

	return utils.GetError().CheckError(*resp, err)

}
