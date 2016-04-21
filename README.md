# dns-agent
Agent for updating IP record into DnsPod.cn

Golang写的用于定时更新DnsPod中的IP记录的定时更新器
现在设定为每60秒更新一次。

参数配置请修改config.json文件

config.json
```json
{
  "debug": true,
  "interval": 60,
  "login_email": "***@163.com",//登陆DNSPOD用到的邮箱名
  "login_pwd": "******",  //登陆DNSPOD的密码
  "domains": [
    {
      "domain_name": "12345.com",  //域名名称,该域名的解析DNS要首先转入到DnsPod
      "record_names": ["@", "www", "jenkins", "kubeapi", "etcd", "ops"] //子域名的前缀
    },
    {
      "domain_name": "yuming.org",
      "record_names": ["@", "www"]
    }
  ],
  "redis": {     //在IP变更的时候用作消息队列
      "enabled": true,
      "dsn": "******",
      "passwd": "******",
      "maxIdle": 5,
      "connTimeout": 5000,
      "readTimeout": 5000,
      "writeTimeout": 5000
  }
}
```

如有任何疑问，请联系coraldane@163.com