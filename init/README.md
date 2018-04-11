Examples for integrating gjfy with various startup systems
==========================================================

Examples provided for: upstart, systemd

All example configurations assume that you have:

  - a system user account with the name "gjfy" and primary group "nogroup"
    (although you could of course choose totally different names
    for those)
  - a directory /usr/local/gjfy that is readable by the gjfy user
  - all required gjfy files installed into that directory

Installation for systemd
------------------------

Copy `init/systemd/gjfy.service` to `/etc/systemd/system`. Change as necessary.
Run `systemctl enable gjfy` and `systemctl start gjfy`. To watch the logs run `journalctl -u gjfy`
