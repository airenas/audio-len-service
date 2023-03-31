# audio-len-service

![](https://github.com/airenas/audio-len-service/workflows/Go/badge.svg) [![Coverage Status](https://coveralls.io/repos/github/airenas/audio-len-service/badge.svg)](https://coveralls.io/github/airenas/audio-len-service)

Simple service to return an audio duration of *wav*, *mp3*, *m4a* files. The service is written in *Go*. It wraps *sox* and *ffprobe* to get the real audio duration.

To test the service look into [examples/docker-compose](examples/docker-compose).  

```bash
    cd examples/docker-compose
    docker-compose up -d
    curl -X POST http://localhost:8003/duration -H 'content-type: multipart/form-data' -F file=@1.mp3
```

The result is in seconds:

```json
{"duration":5.20127}
```

---

## License

Copyright © 2023, [Airenas Vaičiūnas](https://github.com/airenas).

Released under the [The 3-Clause BSD License](LICENSE).

---
