# ai-proxy-poc
* 과제 전형을 위해 구현된 로드 밸런서입니다.
* pkg/load-balancer의 Algorithm 인터페이스를 구현하여 밸런싱 알고리즘을 결정할 수 있습니다.
* Round Robin 방식만 구현되었습니다.

## Usage
1. make [all|windows|linux|macos-arm] 명령으로 OS 환경에 맞춰 빌드할 수 있습니다.(바이너리는 bin/)
2. config.yaml 파일은 바이너리와 같은 경로에 존재하도록 구현되었습니다.
3. configs/config.yaml의 내용을 기반으로 노드와 알고리즘을 설정합니다.
4. 바이너리를 실행하면 서버가 동작합니다.

## TODO
1. 동시성 처리를 위한 mutex를 Atomic 연산으로 변경, 테스트
2. 각 노드의 helth check와 노드 확장
3. 노드로 요청이 노드의 문제로 인한 에러인 경우 Rate Limiter의 처리
