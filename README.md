# phantomflow

**phantomflow**는 웹 애널리틱스 이벤트 데이터를 대량으로 생성하여 특정 API 엔드포인트에 부하 테스트를 수행하는 고성능 HTTP 요청 시뮬레이터입니다. Go 언어로 작성되었으며, 병렬 처리를 통해 높은 초당 요청수(RPS)를 제공합니다.

---

## 특징

* 다중 워커를 이용한 병렬 HTTP 요청
* 랜덤화된 웹 이벤트 데이터 생성
* 메모리 풀(strings.Builder)을 활용한 성능 최적화
* Go 언어의 고성능 HTTP 클라이언트 사용
* 실시간 RPS 모니터링
* 성능 프로파일링(pprof) 내장

---

## 사용법

### 설치 및 빌드

```bash
git clone https://github.com/YOUR_USERNAME/phantomflow.git
cd phantomflow
go build phantomflow.go
```

### 실행 명령어

```bash
./phantomflow [-d 지속시간(초)] [-n 총요청수] [-c 동시워커수] [ENDPOINT]
```

#### 예시

* 5분(300초) 동안 256개의 동시 워커로 지속적으로 요청 보내기:

```bash
./phantomflow -d 300 -c 256 http://klicklab-nlb-0f6efee8fd967688.elb.ap-northeast-2.amazonaws.com/api/analytics/collect
```

* 10만 건의 요청을 64개의 동시 워커로 보내기:

```bash
./phantomflow -n 100000 -c 64 http://localhost:8080/api
```

---

## 실시간 성능 모니터링

실행 중 매초 RPS(초당 요청 수)를 확인할 수 있으며, 성능 프로파일링을 원할 경우 아래 URL을 브라우저에서 접속하세요.

```
http://localhost:6060/debug/pprof/
```

---


## 주의 사항

* 에러 처리를 강화하여 운영 환경에서의 문제 원인을 파악할 수 있도록 개선이 필요합니다.
* HTTP 프로토콜 설정이 환경에 맞게 구성되었는지 확인하세요(기본 HTTP/1.1).
