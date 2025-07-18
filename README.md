# PhantomFlow

<img width="1920" height="1080" alt="Black and White Minimalist Gothic Music Playlist Instagram Post" src="https://github.com/user-attachments/assets/621605cf-6cd4-4ad3-9130-811f0a4a0c91" />


**phantomflow** is a high-performance HTTP request simulator designed for load testing specific API endpoints by generating large volumes of web analytics event data. It is written in Go and leverages parallel processing to achieve high requests per second (RPS).

---

## Key Features

* Parallel HTTP requests using multiple workers
* Randomized web event data generation
* Performance optimization using memory pools (strings.Builder)
* High-performance HTTP client using Go
* Real-time RPS monitoring
* Integrated performance profiling (pprof)

---

## Usage

### Installation and Build

```bash
git clone https://github.com/YOUR_USERNAME/phantomflow.git
cd phantomflow
go build phantomflow.go
```

### Running the Simulator

```bash
./phantomflow [-d duration_seconds] [-n total_requests] [-c concurrent_workers] [ENDPOINT]
```

#### Examples

* Continuously send requests for 5 minutes (300 seconds) with 256 concurrent workers:

```bash
./phantomflow -d 300 -c 256 http://klicklab-nlb-0f6efee8fd967688.elb.ap-northeast-2.amazonaws.com/api/analytics/collect
```

* Send 100,000 requests using 64 concurrent workers:

```bash
./phantomflow -n 100000 -c 64 http://localhost:8080/api
```

---

## Real-time Performance Monitoring

You can monitor RPS (requests per second) in real-time and perform profiling by accessing the URL below:

* **Local Environment:**

```
http://localhost:6060/debug/pprof/
```

* **Remote Environment (e.g., EC2):**

Direct access on EC2 instance:

```bash
curl http://localhost:6060/debug/pprof/
```

Access via SSH tunneling from local PC:

```bash
ssh -L 6060:localhost:6060 -i YOUR_KEY.pem ubuntu@YOUR_EC2_IP
```

Then, in your local browser:

```
http://localhost:6060/debug/pprof/
```

---

## Important Notes

* Enhance error handling to diagnose issues clearly in production environments.
* Ensure HTTP protocol settings match your environment (default is HTTP/1.1).

---

## License

MIT License
