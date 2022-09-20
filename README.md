# cmk-teamspeak
Check_MK agent check for Teamspeak3 virtual server instances


## Setup
Install the mkp on your checkmk with:
```bash
omd su site
mkp install teamspeak3-0.2.mkp
```

Then copy `agents/plugins/Teamspeak3` to `/usr/lib/check_mk_agent/plugins/Teamspeak3` on your checkmk agent.


Create your configuration (destination: `/etc/check_mk/teamspeak3.cfg`)


After that restart checkmk and the agent and rescan for services for the specific host.
