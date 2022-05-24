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
const URL = "http://162.55.235.87/?id="
const letterBytes = "0123456789"
const IDLENGTH = 5
const BLOCK_VERIFIER_VALID_AMOUNT = 9
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
 
func turbo_tx_verify(class TurboTxSave) (int,int,string) {
  // variables
  results := make(chan string)
  delegate_count := 0
  delegate_selection_count := 0
  var network_data_nodes_results [5]string
  var tx TXResults;
  block_status := "false"
 
  // get a list of each current delegate
  string := get_http_data(GET_DELEGATES_URL)
 fmt.Printf("current delegates: %s\n", string)
  // parse the delegates data
  var delegates = []delegatesArray{}
    if err := json.Unmarshal([]byte(string), &delegates); err != nil {
        return 0,0,block_status
    }

fmt.Printf("TX HASH: %s\n", class.TX_Hash)
 
  // get the tx for each delegate on a separate thread
  for count,val := range delegates {
       count++
        go func(val delegatesArray) {
fmt.Printf("IP = %s\n", "http://" + val.IPAddress + ":18281/get_transaction_pool")
            results <- send_http_data("http://" + val.IPAddress + ":18281/get_transaction_pool","")
        }(val)
    }
 
  // Receive results from the channel and use them.
    for count,_ := range delegates {
        if strings.Contains(<-results, class.TX_Hash) {
          delegate_count++
          if delegate_selection_count == 0 {
            delegate_selection_count = count
          }
        }
    }

  if delegate_count < BLOCK_VERIFIER_VALID_AMOUNT {
  // the delegates did not have a majority but it could already be in a block
fmt.Printf("str1: %s\n", "TX invalid checking if tx is in a block")

    // check if this tx is already in a block  
    network_data_nodes_results[0] = send_http_data("http://127.0.0.1:18285/json_rpc",`{"jsonrpc":"2.0","id":"0","method":"check_tx_key","params":{"txid":"` + class.TX_Hash + `","tx_key":"` + class.TX_Key + `","address":"` + class.Receiver + `"}}`)
    network_data_nodes_results[1] = send_http_data("http://127.0.0.1:18286/json_rpc",`{"jsonrpc":"2.0","id":"0","method":"check_tx_key","params":{"txid":"` + class.TX_Hash + `","tx_key":"` + class.TX_Key + `","address":"` + class.Receiver + `"}}`)
    network_data_nodes_results[2] = send_http_data("http://127.0.0.1:18287/json_rpc",`{"jsonrpc":"2.0","id":"0","method":"check_tx_key","params":{"txid":"` + class.TX_Hash + `","tx_key":"` + class.TX_Key + `","address":"` + class.Receiver + `"}}`)
    network_data_nodes_results[3] = send_http_data("http://127.0.0.1:18288/json_rpc",`{"jsonrpc":"2.0","id":"0","method":"check_tx_key","params":{"txid":"` + class.TX_Hash + `","tx_key":"` + class.TX_Key + `","address":"` + class.Receiver + `"}}`)
    network_data_nodes_results[4] = send_http_data("http://127.0.0.1:18289/json_rpc",`{"jsonrpc":"2.0","id":"0","method":"check_tx_key","params":{"txid":"` + class.TX_Hash + `","tx_key":"` + class.TX_Key + `","address":"` + class.Receiver + `"}}`)
    if strings.Contains(network_data_nodes_results[0], "\"in_pool\": true") {
      json.Unmarshal([]byte(network_data_nodes_results[0]), &tx)
      block_status = "true"
      goto TXVALID
    } else if strings.Contains(network_data_nodes_results[1], "\"in_pool\": true") {
      json.Unmarshal([]byte(network_data_nodes_results[1]), &tx)
      block_status = "true"
      goto TXVALID
    } else if strings.Contains(network_data_nodes_results[2], "\"in_pool\": true") {
      json.Unmarshal([]byte(network_data_nodes_results[2]), &tx)
      block_status = "true"
      goto TXVALID
    } else if strings.Contains(network_data_nodes_results[3], "\"in_pool\": true") {
      json.Unmarshal([]byte(network_data_nodes_results[3]), &tx)
      block_status = "true"
      goto TXVALID
    } else if strings.Contains(network_data_nodes_results[4], "\"in_pool\": true") {
      json.Unmarshal([]byte(network_data_nodes_results[4]), &tx)
      block_status = "true"
      goto TXVALID
    } else {
      return 0,delegate_count,block_status
    }  
  }
  



 /*// get the delegate selection
 for count, val := range delegates {
  if count == delegate_selection_count {
    delegate_selection = val.IPAddress
   }
 }*/

fmt.Printf("str1: %s\n", "TX valid")
 
  // the majority of delegates verified the tx, now check if both the sender and receiver are in the tx and the amount is correct
  network_data_nodes_results[0] = send_http_data("http://127.0.0.1:18285/json_rpc",`{"jsonrpc":"2.0","id":"0","method":"check_tx_key","params":{"txid":"` + class.TX_Hash + `","tx_key":"` + class.TX_Key + `","address":"` + class.Receiver + `"}}`)
    network_data_nodes_results[1] = send_http_data("http://127.0.0.1:18286/json_rpc",`{"jsonrpc":"2.0","id":"0","method":"check_tx_key","params":{"txid":"` + class.TX_Hash + `","tx_key":"` + class.TX_Key + `","address":"` + class.Receiver + `"}}`)
    network_data_nodes_results[2] = send_http_data("http://127.0.0.1:18287/json_rpc",`{"jsonrpc":"2.0","id":"0","method":"check_tx_key","params":{"txid":"` + class.TX_Hash + `","tx_key":"` + class.TX_Key + `","address":"` + class.Receiver + `"}}`)
    network_data_nodes_results[3] = send_http_data("http://127.0.0.1:18288/json_rpc",`{"jsonrpc":"2.0","id":"0","method":"check_tx_key","params":{"txid":"` + class.TX_Hash + `","tx_key":"` + class.TX_Key + `","address":"` + class.Receiver + `"}}`)
    network_data_nodes_results[4] = send_http_data("http://127.0.0.1:18289/json_rpc",`{"jsonrpc":"2.0","id":"0","method":"check_tx_key","params":{"txid":"` + class.TX_Hash + `","tx_key":"` + class.TX_Key + `","address":"` + class.Receiver + `"}}`)
    if strings.Contains(network_data_nodes_results[0], "\"received\"") {
      json.Unmarshal([]byte(network_data_nodes_results[0]), &tx)
      block_status = "false"
      goto TXVALID
    } else if strings.Contains(network_data_nodes_results[1], "\"received\"") {
      json.Unmarshal([]byte(network_data_nodes_results[1]), &tx)
      block_status = "false"
      goto TXVALID
    } else if strings.Contains(network_data_nodes_results[2], "\"received\"") {
      json.Unmarshal([]byte(network_data_nodes_results[2]), &tx)
      block_status = "false"
      goto TXVALID
    } else if strings.Contains(network_data_nodes_results[3], "\"received\"") {
      json.Unmarshal([]byte(network_data_nodes_results[3]), &tx)
      block_status = "false"
      goto TXVALID
    } else if strings.Contains(network_data_nodes_results[4], "\"received\"") {
      json.Unmarshal([]byte(network_data_nodes_results[4]), &tx)
      block_status = "false"
      goto TXVALID
    } else {
      return 0,delegate_count,block_status
    } 

  TXVALID:
 
fmt.Printf("Amount received from %s and %s %d\n",class.ID,class.TX_Hash,tx.Result.Received) 
  return tx.Result.Received,delegate_count,block_status
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
    CREATEID:

    // create the id
    id := RandStringBytes(IDLENGTH);
 
    // check if the id is already in the database
    val, err := rdb.Get(ctx, id).Result()
    if val != "" {
      goto CREATEID
    }

    // get the timestamp
    timestamp := strconv.FormatInt(time.Now().UTC().UnixNano(), 10)
 
    // save the data in the database
    data := string(`{"id": "` + id + `", "tx_hash": "` + c.Query("tx_hash") + `", "tx_key": "` + c.Query("tx_key") + `", "timestamp": "` + timestamp + `", "sender": "` + c.Query("sender") + `", "receiver": "` + c.Query("receiver") + `", "amount": "` + c.Query("amount") + `"}`)
fmt.Printf("str1: %s\n", data)

    err = rdb.Set(ctx, id, data, 1*time.Hour).Err()
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
   amount,delegate_count,block_status := turbo_tx_verify(data)

   if amount < datamount || amount <= 0 {
      error := ErrorResults{"error"}
      return c.JSON(error)
  } 
 
  result := TurboTxOut{id, data.TX_Hash, data.Timestamp, data.Sender, data.Receiver, strconv.FormatInt(int64(amount), 10),strconv.FormatInt(int64(delegate_count), 10),block_status}
  return c.JSON(result)
})

app.Static("/", "/var/www/html/")
 
app.Get("/*", func(c *fiber.Ctx) error {
  return c.SendString("Invalid URL")
})
 
  app.Listen(":8000")
}
