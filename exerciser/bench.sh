go test -bench=. -cpuprofile profile.out
go tool pprof profile.out