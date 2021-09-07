# prstat -- report on all processes, running or done

This is based on process accounting, so you can capture the results of processes
which ran for quite a short time, as well as running processes looked at by
polling via ps.

It's something of a Linux flavor of solaris prstat(1), but not closely comparable.
See also dump-acct(1) and ps(1)

Written so I could fairly compare long-running multi-threaded processes and short-lived
parent-and-child processes, both of which I have.
