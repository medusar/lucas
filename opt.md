# Steps
1. redis-benchmark -p 6379 -t set -n 10000000 -q
2. go tool pprof --seconds=5 localhost:8080/debug/pprof/profile

```bash
➜  lucas git:(basic) ✗ go tool pprof -seconds=5 localhost:8080/debug/pprof/profile
Fetching profile over HTTP from http://localhost:8080/debug/pprof/profile?seconds=5
Please wait... (5s)
Saved profile in /Users/jasonartka/pprof/pprof.samples.cpu.006.pb.gz
Type: cpu
Time: Dec 1, 2019 at 6:03pm (CST)
Duration: 5.10s, Total samples = 12.14s (237.87%)
Entering interactive mode (type "help" for commands, "o" for options)
(pprof) web
(pprof)
(pprof) top10
Showing nodes accounting for 11.49s, 94.65% of 12.14s total
Dropped 71 nodes (cum <= 0.06s)
Showing top 10 nodes out of 63
      flat  flat%   sum%        cum   cum%
     4.57s 37.64% 37.64%      4.58s 37.73%  syscall.Syscall
     1.62s 13.34% 50.99%      1.63s 13.43%  runtime.pthread_cond_wait
     1.37s 11.29% 62.27%      1.37s 11.29%  runtime.kevent
     1.13s  9.31% 71.58%      1.13s  9.31%  runtime.nanotime
     0.98s  8.07% 79.65%      0.98s  8.07%  runtime.pthread_cond_timedwait_relative_np
     0.79s  6.51% 86.16%      0.79s  6.51%  runtime.pthread_cond_signal
     0.54s  4.45% 90.61%      1.92s 15.82%  runtime.netpoll
     0.22s  1.81% 92.42%      0.22s  1.81%  runtime.freedefer
     0.18s  1.48% 93.90%      0.19s  1.57%  runtime.usleep
     0.09s  0.74% 94.65%      0.12s  0.99%  runtime.scanobject
```
