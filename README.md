# sp-limiter

## 設定檔
- 路徑： sp-limiter/configs/config.yaml
- 重要參數
    - host.mode 控制限流實作類型：counter(計數器)、tokenBucket(令牌桶)、redisCounter(計數器 Redis)
    - host.interval 以秒為單位，e.g. 控制每 60s 重新解除限流條件
    - host.maxCount 控制上限次數

## 思路
in-memory 版本：counter(計數器)、tokenBucket(令牌桶)
<br/>
資料庫版本：redisCounter(計數器 Redis)
<br/>
選擇方案的思路
- 情境：
  <br/>
  服務尚未進入到集群狀態，有 HA 機制但同時執行時只會有一個服務，在初期可以用 in-memory 資料結構的程式處理限流的邏輯。
  - 優點：適合初期擴建中的服務，未來容易替換為不同實作限流。
  - 缺點：因為服務是有狀態的，進入同時多個服務集群，限流邏輯就無法控制正確總數。
  
- 情境：
  <br/>
  運行多個 API 服務的集群，集群要做限流就無法各個服務自行處理 (總數會超過上限)，需要一個資料庫做狀態的儲存，舉例用 Redis 做資料的狀態儲存。
  - 優點：服務無狀態，服務從單個擴容為多個服務，限流邏輯也可正常運行。
  - 缺點：每個服務各自連 Redis，Redis 成了服務的耦合點，Redis 單線程可能會成為效能的瓶頸。

- 總結：
  <br/>
  依照服務的架構、擁有的 DB 建置、時程等因素確定優先級，在評估適合當下的方案，
  <br/>
  如果單以技術角度選擇建議讓服務屬於無狀態，不要耦合資料狀態讓服務可彈性的做調整，讓資料庫做有狀態的服務。

## 啟動方式
- Heroku 平台
    - https://sp-limiter.herokuapp.com/
      <br/>
      網頁版本執行是 redisCounter (計數器 Redis)
      
- 本地
    - 在目錄位置 github.com/ql31j45k3/sp-limiter
      <br/>
      執行語法 ： go run cmd/sp-limiter/main.go
      <br/>
      本地如果要執行 redisCounter 版本，請確認 config.yaml 的 redis 帳密與位置參數
      
## 測試方式
- 控制 IP 數量
  檔案路徑：github.com/ql31j45k3/sp-limiter/api_test.go
  <br/>
  調整宣告最上方的 ip 字串陣列參數資料
  
- 在目錄位置 github.com/ql31j45k3/sp-limiter
  <br/>
  執行語法 (測試覆蓋率) ： go test -cover ./internal/modules/limiter
  <br/>
  執行語法 (單元測試輸出結果) ： go test ./internal/modules/limiter
  <br/>
  執行語法 (單元測試包含細節) ： go test -v ./internal/modules/limiter
  