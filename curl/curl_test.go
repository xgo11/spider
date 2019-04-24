package curl

import "encoding/json"
import "testing"

func TestAnalyzeCurl(t *testing.T) {
	url, kwArgs, err := AnalyzeCurl(`curl 'https://www.csdn.net/js/tingyun-rum-feed.js?v=20181030' -H 'pragma: no-cache' -H 'cookie: uuid_tt_dd=10_19601600960-1552016852948-858934; _ga=GA1.2.1189684765.1552019215; dc_session_id=10_1552206016662.236562; ADHOC_MEMBERSHIP_CLIENT_ID1.0=c2c53bff-dcfe-62b2-e9de-50b1393a00be; UserName=whucyl; UserInfo=7e7e04aa187d40c49fb41d4713c740d1; UserToken=7e7e04aa187d40c49fb41d4713c740d1; UserNick=PlusPlus1; AU=009; UN=whucyl; BT=1555599510081; Hm_ct_6bcd52f51e9b3dce32bec4a3997715ac=6525*1*10_19601600960-1552016852948-858934!5744*1*whucyl; Hm_lvt_6bcd52f51e9b3dce32bec4a3997715ac=1555936908,1555937468,1555937475,1556009314; TY_SESSION_ID=7546db2a-40d3-4ae8-a52b-78eb3f3f51e3; Hm_lpvt_6bcd52f51e9b3dce32bec4a3997715ac=1556119539; dc_tos=pqh1ir' -H 'dnt: 1' -H 'accept-encoding: gzip, deflate, br' -H 'accept-language: zh-CN,zh;q=0.9,en;q=0.8' -H 'user-agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/72.0.3626.121 Safari/537.36' -H 'accept: */*' -H 'cache-control: no-cache' -H 'authority: www.csdn.net' -H 'referer: https://www.csdn.net/' --compressed`)
	if err != nil {
		t.Log(err)
	} else {
		t.Logf("url = %v", url)
		kwBytes, _ := json.MarshalIndent(kwArgs, "", "    ")
		t.Logf("kwargs = %v", string(kwBytes))
	}
}
