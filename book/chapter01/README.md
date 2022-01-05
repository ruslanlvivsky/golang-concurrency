# 동시성 소개

## 동시성이 어려운 이유

### 레이스 컨디션

- 둘 이상의 작업이 올바른 순서로 실행돼야 하지만 프로그램이 그렇게 작성되지 않아 순서 유지가 보장되지 않을 때 발생한다
- 대부분의 경우 동일 변수에 값을 쓰려고 하는 데이터 레이스 상황이다

```go
var data int
go func () {
    data++
}()
if data == 0 {
    fmt.Println("the value is %v.\n".data)
}
```

sleep 사용은 단지 데이터 레이스가 발생할 가능성을 낮출 뿐이다

### 원자성

동작하는 컨텍스트 안에서 나누어지거나 중단되지 않는다는 것을 의미한다. Context 컨텍스트라는 용어가 중요하다. 어떤 컨텍스트에서는 원자적인 것이 다른 컨텍스트에서는 아닐 수 있다. 연산의 원자성은 현재 정의된
범위에 따라 달라질 수 있다

불가분 indivisible, 중단 불가 uninterruptible 이라는 용어도 중요하다 사용자가 정의한 컨텍스트 내부에서 원자적인 요소가 통째로 발생하며, 해당 컨텍스트 내에서는 해당 요소 외에 어떤 것도
동시에 이루어지지 않는다는 것을 의미한다.

i++

- i의 값을 가져온다
- i의 값을 증가시킨다
- i의 값을 저장한다

원자적 연산을 조합한다고 해서 반드시 더 큰 원자적 연산이 생성되는 것은 아니다

연산을 원자적으로 만드는 것은 사용자가 어떤 컨텍스트에서 원자성을 얻고자 하는지에 달려 있다.

수행되는 프로세스들이 없는 프로그램의 컨텍스트라면, 이 코드는 해당 컨텍스트 내에서 원자적이다

i 값을 다른 고루틴들에게 노출하지 않는 고루틴의 컨텍스트인 경우에도 이 코드는 원자적이다

#### 원자성이 중요한 이유는 무언가가 원자적이라면 암묵적으로 동시에 실행되는 컨텍스트들 내에서는 안전하다는 것을 의미하기 때문이다

### 메모리 접근 동기화

```go
var data int
go func () {
    data++
}()
if data == 0 {
    fmt.Println("the value is 0.")
} else {
    fmt.Printf("the value is %v\n", data)
}
```

임계 영역은 프로그램에서 공유 리소스에 독점적으로 접근해야 하는 영역을 가리킨다

##### 세가지 임계 영역

- data 변수를 증가시키는 goroutine
- data 값이 0인지 확인하는 if 구문
- 출력할 data의 값을 가져오는 fmt.Printf 구문 



### 데드락, 라이브락, 기아 상태

#### 데드락

프로그램에서 동시에 실행 중인 모든 프로세스는 자신이 아닌 다른 프로세스가 끝나기만을 기다린다 

데드락 상태는 외부 개입 없이 프로그램을 복구할 수 없다 

```go
type value struct {
    mu sync.Mutext
    value int
}

var wg sync.WaitGroup
printSum := func(v1, v2 *value){
    defer wg.Done()
    v1.mu.Lock()
    defer v2.mu.Unlock()
    
    v2.mu.Lock()
    defer v2.mu.Unlock()
    
    fmt.Printf("sum=%v\n", v1.value + v2.value)
}

var a, b value
wg.Add(2)
go printSum(&a, &b)
go printSum(&b, &a)
wg.Wait()
```

두 고루틴은 서로를 무한히 기다린다 

#### 데드락 발생 조건

1. 상호 배제 

   동시에 실행되는 프로세스가 어떤 임의의 시점에 하나의 리소스에 대한 배타적 권리를 보유한다 

2. 대기 조건

   동시에 실행되는 프로세스는 하나의 리소스를 보유하고 있는 동시에 또 다른 추가 리소스를 기다리고 있다 

3. 비선점

   동시에 실행되는 프로세스 중 하나를 보유하고 있는 리소스는 해당 프로세스에 의해서만 사용 해제될 수 있으므로 이 조건도 만족한다.

4. 순환 대기 

   동시에 실행되는 프로세스 중 하나가 다른 동시 프로세스로 이어지는 체인에서 기다려야 하며, P2는 최종적으로 P1을 기다려야 하는데 마지막 조건 역시 충족하다 

이 조건들 중 하나라도 참이 아니라면 데드락 발생을 예방할 수 있다 



### 라이브락

프로그램들이 활동적으로 동시에 연산을 수행하고 있지만, 이 연산들이 실제로 프로그램의 상태를 진행시키는데 아무런 영향을 주지 못하는 의미없는 연산 상태를 의미한다

프로그램이 마치 동작하는 것처럼 보이기 때문에 라이브락은 데드락보다 더 알아보기 힘들다 



### 기아 상태

어떤 동시 프로세스가 작업을 수행하는 데 필요한 모든 리소스를 얻을 수 없는 모든 상황을 의미한다 

> 다른 동시 프로세스 혹은 프로세스들이 가능한 효율적으로 작업을 수행하는 것을 부당하게 방해하거나, 작업을 전혀 수행하지 못하게 만드는 욕심 많은 동시 프로세스가 하나 이상 존재한다는 것을 암시한다 



### 동시실행 안전성 판단

- 누가 동시성을 책임지는가?
- 문제 공간은 동시성 기본 요소에 어떻게 매핑되는가?
- 동기화는 누가 담당하는가?


