package main
 
import (
"fmt"
"math/rand"
"strings"
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
    Delegate_Count string `json:"delegate_count"`
    Block_Status string `json:"block_status"`
}
 
type delegatesArray struct {
    TotalVoteCount string `json:"total_vote_count"`
    IPAddress string `json:"ip_address"`
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

type ErrorResults struct {
    Error string `json:"Error"`
}
 
// global constants
const URL = "http://turbotx.xcash.foundation/?id="
const letterBytes = "0123456789"
const IDLENGTH = 5
const BLOCK_VERIFIER_TOTAL_AMOUNT = 14
const BLOCK_VERIFIER_VALID_AMOUNT = 9
const TX_HASH_LENGTH = 64
const PUBLIC_ADDRESS_LENGTH = 98
const GET_DELEGATES_URL = "http://dpops-test-internal-2.xcash.foundation/getdelegates"
 
// Functions
func send_http_data(url string,data string) string {
	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(data)))
if err != nil {
		return "error1"
	}
	req.Header.Set("Content-Type", "application/json")
 
	client := &http.Client{}
	client.Timeout = time.Second * 2
	resp, err := client.Do(req)
	if err != nil {
		return "error2"
	}
	defer resp.Body.Close()
        body, _ := ioutil.ReadAll(resp.Body)
fmt.Printf("for %s sending %s received %s\n", url,data,body)

        return string(body)
}
 

func get_http_data(url string) string {
	req, err := http.NewRequest("GET", url, bytes.NewBuffer([]byte("")))
if err != nil {
		return "error"
	}	
req.Header.Set("Content-Type", "application/json")
 
	client := &http.Client{}
	client.Timeout = time.Second * 2
	resp, err := client.Do(req)
	if err != nil {
		return "error"
	}
	defer resp.Body.Close()
        body, _ := ioutil.ReadAll(resp.Body)
        return string(body)
}
func RandStringBytes(n int) string {
    b := make([]byte, n)
    for i := range b {
        b[i] = letterBytes[rand.Intn(len(letterBytes))]
    }
    return string(b)
}
 
func turbo_tx_verify(class TurboTxSave) (int,int,string,string) {
  // variables
  results := make(chan string)
  delegate_count := 0
  var tx TXResults;
  block_status := "false"
  delegates_data_sender := ""
  delegates_data_receiver := ""
  var delegates_results [BLOCK_VERIFIER_TOTAL_AMOUNT]int
 
  // get a list of each current delegate
  string := get_http_data(GET_DELEGATES_URL)
 fmt.Printf("current delegates: %s\n", string)
  // parse the delegates data
  var delegates = []delegatesArray{}
    if err := json.Unmarshal([]byte(string), &delegates); err != nil {
        return 0,0,block_status,"0"
    }

  if (len(delegates) < BLOCK_VERIFIER_TOTAL_AMOUNT) {
    return 0,0,block_status,"0"
  }

fmt.Printf("TX HASH: %s\n", class.TX_Hash)
 
  // get the tx for each delegate on a separate thread
  for count,val := range delegates {
       //count++
       if (count == BLOCK_VERIFIER_TOTAL_AMOUNT) {
        break;
       }
        go func(val delegatesArray) {
fmt.Printf("IP = %s\n", "http://" + val.IPAddress + ":18281/get_transaction_pool")
            results <- send_http_data("http://" + val.IPAddress + ":18281/get_transaction_pool","")
        }(val)
    }
 
  // Receive results from the channel and use them.
    for count,_ := range delegates {
        if (count == BLOCK_VERIFIER_TOTAL_AMOUNT) {
        break;
       }
        if strings.Contains(<-results, class.TX_Hash) {
          delegate_count++
          delegates_results[count] = 1
        } else {
        delegates_results[count] = 0
       }
    }

  if delegate_count < BLOCK_VERIFIER_VALID_AMOUNT {
  // the delegates did not have a majority but it could already be in a block
fmt.Printf("str1: %s\n", "TX invalid checking if tx is in a block")

    // check if this tx is already in a block 
    for count,val := range delegates {
      if (count == BLOCK_VERIFIER_TOTAL_AMOUNT) {
        break;
       }

      // since this tx could already be in a block, the delegates_results might be 0 on all of them, so just try each one in order since the first 5 are the seed nodes
       delegates_data_receiver = send_http_data("http://" + val.IPAddress + ":18286/json_rpc",`{"jsonrpc":"2.0","id":"0","method":"check_tx_key","params":{"txid":"` + class.TX_Hash + `","tx_key":"` + class.TX_Key + `","address":"` + class.Receiver + `"}}`)
       delegates_data_sender = send_http_data("http://" + val.IPAddress + ":18286/json_rpc",`{"jsonrpc":"2.0","id":"0","method":"check_tx_key","params":{"txid":"` + class.TX_Hash + `","tx_key":"` + class.TX_Key + `","address":"` + class.Sender + `"}}`)
       if strings.Contains(delegates_data_receiver, "\"in_pool\": false") && strings.Contains(delegates_data_sender, "\"in_pool\": false") {
      json.Unmarshal([]byte(delegates_data_receiver), &tx)
      block_status = "true"

      // get the timestamp
      timestamp_data := send_http_data("http://" + val.IPAddress + ":18281/get_transactions",`{"txs_hashes":["` + class.TX_Hash + `"]}`)
      timestamp_blockchain_data := timestamp_data[strings.Index(timestamp_data, "\"block_timestamp\""):len(timestamp_data)]
      timestamp := timestamp_blockchain_data[19:strings.Index(timestamp_blockchain_data, ",")]
      return tx.Result.Received,delegate_count,block_status,timestamp
      }
    } 
    return 0,delegate_count,block_status,"0"
  }

fmt.Printf("str1: %s\n", "TX valid")
 
  // the majority of delegates verified the tx, now check if both the sender and receiver are in the tx and the amount is correct
  for count,val := range delegates {
      if (count == BLOCK_VERIFIER_TOTAL_AMOUNT) {
        break;
       }
      if delegates_results[count] == 1 {
       delegates_data_receiver = send_http_data("http://" + val.IPAddress + ":18286/json_rpc",`{"jsonrpc":"2.0","id":"0","method":"check_tx_key","params":{"txid":"` + class.TX_Hash + `","tx_key":"` + class.TX_Key + `","address":"` + class.Receiver + `"}}`)
      delegates_data_sender = send_http_data("http://" + val.IPAddress + ":18286/json_rpc",`{"jsonrpc":"2.0","id":"0","method":"check_tx_key","params":{"txid":"` + class.TX_Hash + `","tx_key":"` + class.TX_Key + `","address":"` + class.Sender + `"}}`)
       if strings.Contains(delegates_data_receiver, "\"in_pool\": true") && strings.Contains(delegates_data_sender, "\"in_pool\": true") {
      json.Unmarshal([]byte(delegates_data_receiver), &tx)
      block_status = "false"
      return tx.Result.Received,delegate_count,block_status,"0"
      }
      }
    }
 return 0,delegate_count,block_status,"0"
}
 
func main() {
// set the random number generator
rand.Seed(time.Now().UTC().UnixNano())

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
 
app.Post("/processturbotx/", func(c *fiber.Ctx) error {
    // Variables
    var id string
    tx_hash := c.Query("tx_hash")
    tx_key := c.Query("tx_key")
    sender := c.Query("sender")
    receiver := c.Query("receiver")
    amount := c.Query("amount")

    // error check
    if (len(tx_hash) != TX_HASH_LENGTH || len(tx_key) != TX_HASH_LENGTH || len(sender) != PUBLIC_ADDRESS_LENGTH || len(receiver) != PUBLIC_ADDRESS_LENGTH) {
      error := ErrorResults{"error"}
      return c.JSON(error)
    }

    if _, err := strconv.Atoi(amount); err != nil {
      error := ErrorResults{"error"}
      return c.JSON(error)
    }

    // get the id
    id = tx_hash[:IDLENGTH]

    // get the timestamp
    timestamp := strconv.FormatInt(time.Now().UTC().UnixNano(), 10)
 
    // save the data in the database
    data := string(`{"id": "` + id + `", "tx_hash": "` + tx_hash + `", "tx_key": "` + tx_key + `", "timestamp": "` + timestamp + `", "sender": "` + sender + `", "receiver": "` + receiver + `", "amount": "` + c.Query("amount") + `"}`)
fmt.Printf("str1: %s\n", data)

    err := rdb.Set(ctx, id, data, 1*time.Hour).Err()
    if err != nil {
        error := ErrorResults{"error"}
        return c.JSON(error)
    }
 
    // return the id
    return c.SendString(URL + id + "}")
})
 
app.Get("/getturbotx/", func(c *fiber.Ctx) error {
  id := c.Query("id")
val, _ := rdb.Get(ctx, id).Result()
    if val == "" {
      error := ErrorResults{"error"}
      return c.JSON(error)
    }
fmt.Printf("%s\n", val)
   // convert the string to a json object
   var data TurboTxSave
   json.Unmarshal([]byte(val), &data)
 fmt.Printf("str1: %s\n", "checking data")
fmt.Println("Struct is:", data)

   // check if the amount is correct and the sender and receiver are in the output
   datamount, _ := strconv.Atoi(data.Amount)
   amount,delegate_count,block_status,timestamp := turbo_tx_verify(data)

   if amount < datamount || amount <= 0 {
      error := ErrorResults{"error"}
      return c.JSON(error)
  } 

  if timestamp == "0" {
    timestamp = data.Timestamp
  }
 
  result := TurboTxOut{id, data.TX_Hash, timestamp, data.Sender, data.Receiver, strconv.FormatInt(int64(amount), 10),strconv.FormatInt(int64(delegate_count), 10),block_status}
  return c.JSON(result)
})

app.Static("/", "/var/www/html/turbotx/")
 
app.Get("/*", func(c *fiber.Ctx) error {
  return c.SendString("Invalid URL")
})
 
  app.Listen(":3000")
}
