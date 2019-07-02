package preset_function

//导入json

import (
  "net/http"
)

func RegisterMonitorBusiness() {
   RegisterAuth2Method("/server/manager/stop", stopServer)        //关闭 服务
   RegisterAuth2Method("/server/query/status", queryStatusServer) //查询 服务器基础信息
   RegisterAuth2Method("/server/query/system_log", querySystemLogServer) //查询服务器上的系统日志
   RegisterAuth2Method("/server/query/acotr_table", queryAcotrTableServer) //查询服务器上的Actor列表带翻页
   RegisterAuth2Method("/server/query/actor_info", queryActorInfoServer)  //查询服务器上某个Actor详细信息
}

func stopServer(w http.ResponseWriter, r *http.Request) {

}

func queryStatusServer(w http.ResponseWriter, r *http.Request) {

}

func querySystemLogServer(w http.ResponseWriter, r *http.Request) {

}

func queryAcotrTableServer(w http.ResponseWriter, r *http.Request) {

}

func queryActorInfoServer(w http.ResponseWriter, r *http.Request) {

}
