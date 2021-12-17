Idea from https://github.com/modderfawker/ApexLegendsIPLogger/blob/master/iplogger.py and https://gamebanana.com/tools/5684

This tool responds to a hotkey (Ctrl+Shift+P) and then searches for UDP packets in the port range 37005-38515 (Apex servers).
This IP will be pinged, if the latency is too high (100ms) this IP will be blocked in the windows firewall after the game so you can't connect to this server anymore.
If the packet loss is too high you can also block the last IP manually with (Ctrl+Shift+O).

![screenshot](https://i.imgur.com/fHYvVJH.png)

#### For color in the Console:
`REG ADD HKEY_CURRENT_USER\Console /v VirtualTerminalLevel /t REG_DWORD /d 0x00000001 /f`
