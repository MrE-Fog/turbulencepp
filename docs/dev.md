# Development

Source code is located in `src/github.com/cppforlife/turbulence`. Use `manifests/example.yml` and `manifests/dummy.yml` to deploy API server and dummy deployment.

Run `./src/github.com/cppforlife/turbulence/bin/test` for unit tests.

Run `cd tests && ./run.sh` for an integration test.

## Dependencies

Run `./update-deps` to update `github.com/cppforlife/turbulence` package dependencies. `deps.txt` will be updated with Git SHAs for each dependency.

## Planned tasks

- lock up whole machine
- remount disk as readonly
- corrupt disks

https://www.kernel.org/doc/Documentation/sysrq.txt might be useful...
http://blog.hut8labs.com/gorillas-before-monkeys.html
http://techblog.netflix.com/2011/04/lessons-netflix-learned-from-aws-outage.html
