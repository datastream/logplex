A Go package for reading syslog streams

current it just support RSYSLOG_SyslogProtocol23Format.

    "<%PRI%>1 %TIMESTAMP:::date-rfc3339% %HOSTNAME% %APP-NAME% %PROCID% %MSGID% %STRUCTURED-DATA% %msg%\n\"
