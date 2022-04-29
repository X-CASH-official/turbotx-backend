package main
 
import (
"fmt"
"context"
"strconv"
"bytes"
"io/ioutil"
"net/http"
"time"
"encoding/json"
"github.com/gofiber/fiber/v2"
"github.com/go-redis/redis/v8"
)
 
// global structures
type TurboTxIn struct {
    Sender string `form:"sender"`
    Receiver string `form:"receiver"`
    Amount string `form:"amount"`
    TX_Hash string `form:"tx_hash"`
    TX_Key string `form:"tx_key"`    
}
 
type TurboTxSave struct {
    ID string `json:"id"`
    TX_Hash string `json:"tx_hash"`
    TX_Key string `form:"tx_key"`
    Timestamp string `json:"timestamp"`
    Sender string `json:"sender"`
    Receiver string `json:"receiver"`
    Amount string `json:"amount"`
}
 
type TurboTxOut struct {
    ID string `json:"id"`
    TX_Hash string `json:"tx_hash"`
    Timestamp string `json:"timestamp"`
    Sender string `json:"sender"`
    Receiver string `json:"receiver"`
    Amount string `json:"amount"`
}
 
type delegatesArray []struct {
    TotalVoteCount string `json:"total_vote_count"`
    IPAddress string `json:"IP_address"`
    DelegateName string `json:"delegate_name"`
    SharedDelegateStatus string `json:"shared_delegate_status"`
    DelegateFee string `json:"delegate_fee"`
    OnlineStatus string `json:"online_status"`
    BlockVerifierTotalRounds string `json:"block_verifier_total_rounds"`
    BlockVerifierOnlinePercentage string `json:"block_verifier_online_percentage"`
    BlockProducerTotalRounds string `json:"block_producer_total_rounds"`
}
 
type TXResults struct {
	ID      string `json:"id"`
	Jsonrpc string `json:"jsonrpc"`
	Result  struct {
		Confirmations int  `json:"confirmations"`
		InPool        bool `json:"in_pool"`
		Received      int  `json:"received"`
	} `json:"result"`
}
 
// global constants
const URL string = "https://xcash.foundation/"
const letterBytes = "0123456789"
const IDLENGTH = 5
const BLOCK_VERIFIER_VALID_AMOUNT = 27
 
 
// Functions
func send_http_data(url string,data string) string
{
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/json")
 
	client := &http.Client{}
	client.Timeout = time.Second * 2
	resp, err := client.Do(req)
	if err != nil {
		return "error"
	}
	defer resp.Body.Close()
        return string(ioutil.ReadAll(resp.Body))
}
 
 
func RandStringBytes(n int) string 
{
    b := make([]byte, n)
    for i := range b {
        b[i] = letterBytes[rand.Intn(len(letterBytes))]
    }
    return string(b)
}
 
func turbo_tx_verify(class TurboTxSave) int
{
  // variables
  results := make(chan string)
  delegate_count := 0
  delegate_selection_count := 0
  delegate_selection := ""
 
  // get a list of each current delegate
  string := send_http_data("http://delegates.xcash.foundation/getdelegates","")
 
  // parse the delegates data
  var delegates []delegatesArray
    if err := json.Unmarshal([]byte(string), &delegates); err != nil {
        return 0
    }
 
  // get the tx for each delegate on a separate thread
  for count, val := range delegates {
        go func() {
            results <- send_http_data(val.IP_address + ":18281/get_transaction_pool","")
        }(count)
    }
 
  // Receive results from the channel and use them.
    for i := 0; i < rep; i++ {
        if strings.Contains(<-results, class.tx_hash) {
          delegate_count++
          if delegate_selection_count == 0 {
            delegate_selection_count = i
          }
        }
    }
 
  if delegate_count < BLOCK_VERIFIER_VALID_AMOUNT {
    return 0
  }
 
 // get the delegate selection
 for count, val := range delegates {
  if count == delegate_selection_count {
    delegate_selection = val.IP_address
   }
 }
 
  // set your wallet to use the delegate selection remote node
  string = send_http_data("127.0.0.1:18285/json_rpc",[]byte(`{"jsonrpc":"2.0","id":"0","method":"set_daemon","params": {"address":delegate_selection + ":18281,"trusted":true}}`))
 
  if strings.Contains(string, "error") {
   return 0
  }
 
  // the majority of delegates verified the tx, now check if both the sender and receiver are in the tx and the amount is correct
  string = send_http_data("127.0.0.1:18285/json_rpc",[]byte(`{"jsonrpc":"2.0","id":"0","method":"check_tx_key","params":{"txid":class.tx_hash,"tx_key":class.tx_key,"address":class.reciever}}`))
 
  if !strings.Contains(string, "\"received\"") {
   return 0
  }
 
   tx := delegatesArray{}
   json.Unmarshal([]byte(string), &tx)
 
  return tx.Result.received
}
 
func main() {
 
// set redis connection
var ctx = context.Background()
rdb := redis.NewClient(&redis.Options{
        Addr:     "localhost:6379",
        Password: "", // no password set
        DB:       0,  // use default DB
    })
 
// setup fiber
app := fiber.New(fiber.Config{
Prefork: true,
DisableStartupMessage: true,
})
 
app.Post("/processturbotx/", func(c *fiber.Ctx) error
{
    // create a new turbo tx in object
    turbotxin := new(TurboTxIn)
 
    // get the post request data
    if err := c.BodyParser(turbotxin); err != nil
    {
      return c.SendString("error")
      return err
    }
 
    CREATEID:
 
    // create the id
    id := RandStringBytes(IDLENGTH);
 
    // check if the id is already in the database
    val, err := rdb.Get(ctx, id).Result()
    if err != nil {
        return c.SendString("error")
        return err
    }
    if err != redis.Nil {
      goto CREATEID
    }
 
    // get the timestamp
    timestamp := strconv.FormatInt(time.Now().UTC().UnixNano(), 10)
 
    // save the data in the database
    data := string(`{"id": id, "tx_hash": turbotxin.tx_hash, "tx_key": turbotxin.tx_key, "timestamp": timestamp, "sender": turbotxin.sender, "receiver": turbotxin.receiver, "amount": turbotxin.amount}`)
 
    err := rdb.Set(ctx, id, data, 15*time.Minute).Err()
    if err != nil {
        return c.SendString("error")
        return err
    }
 
    // return the id
    return c.SendString(URL + id)
})
 
app.Get("/getturbotx/", func(c *fiber.Ctx) error
{
  id := c.Query("id")
 
  // get the id from the database
  val, err := rdb.Get(ctx, id).Result()
    if err != nil || err == redis.Nil {
        return c.SendString("error")
        return err
    }
 
   // convert the string to a json object
   data := TurboTxSave{}
   json.Unmarshal([]byte(val), &data)
 
   // check if the amount is correct and the sender and receiver are in the output
   amount := turbo_tx_verify(data)
   if amount <= data.amount || amount <= 0
  {
      return c.SendString("error")
      return err
  } 
 
  result := TurboTxOut{data.id, data.tx_hash, data.timestamp, data.sender, data.receiver, strconv.FormatInt(int64(amount), 10)}
  return c.JSON(result)
})
 
app.Get("/*", func(c *fiber.Ctx) error
{
  return c.SendString("Invalid URL")
})
 
  app.Listen(":8000")
}
