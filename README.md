# dns-agent
Agent for updating IP record into DnsPod.cn

Golang写的用于定时更新DnsPod中的IP记录的定时更新器
现在设定为每10分钟更新一次。

参数配置请修改config.json文件

config.json
```json
{
"LoginEmail": "***@163.com", //登陆DNSPOD用到的邮箱名
"LoginPassword": "****",     //登陆DNSPOD的密码
"Domains": [
	{"DomainId": 17039700,   //域名ID，可以在域名管理界面，域名前面的checkbox中的value
	"Records": [{"RecordId": 70808122, "SubDomain": "@"}, //域名记录, RecordId同样也可以在配置IP的界面，前面复选框中获得
				{"RecordId": 70808125, "SubDomain": "www"}]
	}
  ]
}
```

如有任何疑问，请联系coraldane@163.com