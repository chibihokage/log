# log
log package for golang

# USAGE 

``` 
func InitLog(urlPath string, c echo.Context) (log Logs.LogLevel, trnslog, trnsEndpointlog Logs.Log, err error) {
	serviceName := viper.GetString("serviceName")
	filename, err := getfileName(urlPath, serviceName)
	if err != nil {
		err = errors.Wrap(err, "cannot getfileName")
		return
	}
	funcName := urlPath + ".POST"
	sourceSysID := viper.GetString("sourceSysID")
	requestIP := c.RealIP()
	trnsID := GetTrnsID()
	sesID := GetSessionID(c)
	log = Logs.InitDebuglog(filename, sesID, trnsID, sourceSysID)
	trnslog = Logs.InitTrnslog(filename, sourceSysID, sesID, trnsID, requestIP, serviceName, funcName)
	trnsEndpointlog = Logs.InitEndpointTrnslog(filename, sourceSysID, sesID, trnsID, requestIP, serviceName, funcName)
	return
}


log, trnslog, trnsEPlog, err := common.InitLog(urlPath, c)
log.Infolog.Printf("*#*#*#*#*#*#*#*#* BEGIN :: %s.CALL *#*#*#*#*#*#*#*#*", urlPath)
log.Infolog.Printf("INPUT HEADER| %v", c.Request().Header)
