/*
	Cart API in Go
	Uses Riak KV
*/

package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"io/ioutil"
	"time"
	"errors"
	"encoding/json"
	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"github.com/unrolled/render"
)



/* Riak REST Client */
var debug = true
var server1 = "http://35.173.1.123:8098" // set in environment
//var server2 = "http://54.67.116.220:8098" // set in environment
//var server3 = "http://13.52.101.249:8098" // set in environment

type Client struct {
	Endpoint string
	*http.Client
}


var tr = &http.Transport{
	MaxIdleConns:       10,
	IdleConnTimeout:    30 * time.Second,
	DisableCompression: true,
}

func NewClient(server string) *Client {
	return &Client{
		Endpoint:  	server,
		Client: 	&http.Client{Transport: tr},
	}
}

func (c *Client) Ping() (string, error) {
	resp, err := c.Get(c.Endpoint + "/ping" )
	if err != nil {
		fmt.Println("[RIAK DEBUG] " + err.Error())
		return "Ping Error!", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if debug { fmt.Println("[RIAK DEBUG] GET: " + c.Endpoint + "/ping => " + string(body)) }
	return string(body), nil
}

func (c *Client) addItems(key string, value Order) (Order, error) {
	var ord_nil = Order {}
	reqbody, err := json.Marshal(value)
	
	req, _  := http.NewRequest("PUT", c.Endpoint + "/buckets/cart/keys/"+key+"?returnbody=true", strings.NewReader(string(reqbody)) )
	req.Header.Add("Content-Type", "application/json")
	//fmt.Println( req )
	resp, err := c.Do(req)	
	if err != nil {
		fmt.Println("[RIAK DEBUG] " + err.Error())
		return ord_nil, err
	}	
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if debug { fmt.Println("[RIAK DEBUG] PUT: " + c.Endpoint + "/buckets/cart/keys/"+key+"?returnbody=true => " + string(body)) }
	var ord = Order { }
	if err := json.Unmarshal(body, &ord); err != nil {
		fmt.Println("RIAK DEBUG] JSON unmarshaling failed: %s", err)
		return ord_nil, err
	}
	return ord, nil
}

func (c *Client) deleteCart(key string) (ErrorMessage) {
	var ord_nil = ErrorMessage {}
	req, _  := http.NewRequest("DELETE", c.Endpoint + "/buckets/cart/keys/"+key+"?returnbody=true", nil )
	resp, err := c.Do(req)
	fmt.Println(resp.StatusCode)
	
	defer resp.Body.Close()	
	body, err := ioutil.ReadAll(resp.Body)
	fmt.Println("KEY VALUE",key)
	
	if err != nil {
		fmt.Println("[RIAK DEBUG] ===> " + err.Error())
		ord_nil.message = "internal Error occurred"
		return ord_nil
	}
	if resp.StatusCode != 204  {
		ord_nil.message = "Key not found!"
		return ord_nil
	}

	if debug { fmt.Println("[RIAK DEBUG] GET: " + c.Endpoint + "/buckets/cart/keys/"+key +" => " + string(body)) }
	
	return ord_nil

}

func (c *Client) getCart(key string) (Order, error) {
	var ord_nil = Order {}
	resp, err := c.Get(c.Endpoint + "/buckets/cart/keys/"+key )
	fmt.Println(resp.StatusCode)
	if err != nil {
		fmt.Println("[RIAK DEBUG] ===> " + err.Error())
		return ord_nil, err
	}
	if resp.StatusCode != 200 {
		return ord_nil, errors.New("Key not found..")
		
		//return ord_nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if debug { fmt.Println("[RIAK DEBUG] GET: " + c.Endpoint + "/buckets/cart/keys/"+key +" => " + string(body)) }
	var ord = Order { }
	if err := json.Unmarshal(body, &ord); err != nil {
		fmt.Println("RIAK DEBUG] JSON unmarshaling failed: %s", err)
		return ord_nil, err
	}
	return ord, nil
}

func (c *Client) deleteCart(key string) (ErrorMessage) {
	var ord_nil = ErrorMessage {}
	req, _  := http.NewRequest("DELETE", c.Endpoint + "/buckets/cart/keys/"+key+"?returnbody=true", nil )
	resp, err := c.Do(req)
	fmt.Println(resp.StatusCode)
	
	defer resp.Body.Close()	
	body, err := ioutil.ReadAll(resp.Body)
	fmt.Println("KEY VALUE",key)
	
	if err != nil {
		fmt.Println("[RIAK DEBUG] ===> " + err.Error())
		ord_nil.message = "internal Error occurred"
		return ord_nil
	}
	if resp.StatusCode != 204  {
		ord_nil.message = "Key not found!"
		return ord_nil
	}

	if debug { fmt.Println("[RIAK DEBUG] GET: " + c.Endpoint + "/buckets/cart/keys/"+key +" => " + string(body)) }
	
	return ord_nil

}


// NewServer configures and returns a Server.
func NewServer() *negroni.Negroni {
	formatter := render.New(render.Options{
		IndentJSON: true,
	})
	n := negroni.Classic()
	mx := mux.NewRouter()
	initRoutes(mx, formatter)
	n.UseHandler(mx)
	return n
}

// Init Database Connections

func init() {

	// Get Environment Config
	//server1 = os.Getenv("RIAK1")
	//server2 = os.Getenv("RIAK2")
	//server3 = os.Getenv("RIAK3")
	fmt.Println("Riak Server1:", server1 )	
	//fmt.Println("Riak Server2:", server2 )	
	//fmt.Println("Riak Server3:", server3 )	

	// Riak KV Setup	
	c1 := NewClient(server1)
	msg, err := c1.Ping( )
	if err != nil {
		log.Fatal(err)
	} else {
		log.Println("Riak Ping Server1: ", msg)		
	}

	// c2 := NewClient(server2)
	// msg, err = c2.Ping( )
	// if err != nil {
	// 	log.Fatal(err)
	// } else {
	// 	log.Println("Riak Ping Server2: ", msg)		
	// }

	// c3 := NewClient(server3)
	// msg, err = c3.Ping( )
	// if err != nil {
	// 	log.Fatal(err)
	// } else {
	// 	log.Println("Riak Ping Server3: ", msg)		
	// }

}

func setupResponse(w *http.ResponseWriter, req *http.Request) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
    (*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, PUT, DELETE, OPTIONS")
    (*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}

// API Routes
func initRoutes(mx *mux.Router, formatter *render.Render) {
	mx.HandleFunc("/cart/add", addItemsToCartHandler(formatter)).Methods("PUT")
	mx.HandleFunc("/cart/{key}", getCartHandler(formatter)).Methods("GET")
	mx.HandleFunc("/cart/{key}", deleteCartHandler(formatter)).Methods("DELETE")
	//mx.HandleFunc("/cart", optionsHandler(formatter)).Methods("OPTIONS")
	//mx.HandleFunc("/checkout/{userid}", CheckoutCartHandler(formatter)).Methods("GET")
	mx.HandleFunc("/ping", pingHandler(formatter)).Methods("GET")

}

// API Ping Handler
func pingHandler(formatter *render.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		formatter.JSON(w, http.StatusOK, struct{ Test string }{"Cart APIs version 1.0 alive!"})
	}
}

// Helper Functions
func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}

//API Ping Handler
func optionsHandler(formatter *render.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		setupResponse(&w, req)
		fmt.Println("options handler PREFLIGHT Request")
			return
	}
}


// API to delete cart
func deleteCartHandler(formatter *render.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		setupResponse(&w, req)
		if (*req).Method == "OPTIONS" {
			fmt.Println("PREFLIGHT Request")
			return
		}
		params := mux.Vars(req)
		var ord_err = ErrorMessage{}
		var key string = params["key"]
		fmt.Println("Key : ", key)
		if key == ""  {
			formatter.JSON(w, http.StatusBadRequest, "Invalid Request. Cart Key Missing.")
		} else {
			c1 := NewClient(server1)
			ord := c1.deleteCart(key+"_cart")
			if ord != ord_err {
				
				fmt.Println("NOT SUCCESSFULL")
				formatter.JSON(w, http.StatusBadRequest, struct{ Message string }{ ord.message })

			}else{
				ord.message = "Deleted Cart Successfully!"
				fmt.Println("SUCCESSFULL")
				formatter.JSON(w, http.StatusOK, struct{ Message string }{ ord.message })
			}
		}
	}
}

//API to get cart details
func getCartHandler(formatter *render.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		setupResponse(&w, req)
		if (*req).Method == "OPTIONS" {
			fmt.Println("PREFLIGHT Request")
			return
		}
		params := mux.Vars(req)
		var key string = params["key"]
		fmt.Println("ID : ", key)
		if key == ""  {
			formatter.JSON(w, http.StatusBadRequest, "Invalid Request. Cart Key Missing.")
		} else {
			c1 := NewClient(server1)
			ord, err := c1.getCart(key+"_cart")
			fmt.Println(ord)
			//fmt.Println("Key ---> : ", err != nil)
			if err != nil {
				fmt.Println("err : ", err)
				formatter.JSON(w, http.StatusBadRequest, struct{ Test string }{"Cart not found!"})
			} else {
				formatter.JSON(w, http.StatusOK, ord)
			}
		}

	}
}

// API Add Items to cart
func addItemsToCartHandler(formatter *render.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		setupResponse(&w, req)
		if (*req).Method == "OPTIONS" {
			fmt.Println("PREFLIGHT Request")
			return
		}
		var cartPayload Order
		if err := json.NewDecoder(req.Body).Decode(&cartPayload); err != nil {
			fmt.Println(" Error: ", err)
			formatter.JSON(w, http.StatusBadRequest, "Invalid request payload")
			return
		}
		fmt.Println("Adding Items to Cart: ", cartPayload)

		var uuid string = cartPayload.Username+"_cart";
		fmt.Println( "Order Params ID: ", uuid )

		if cartPayload.Username == ""  {
			formatter.JSON(w, http.StatusBadRequest, "Invalid Request. Order ID Missing.")
		} else {
			c1 := NewClient(server1)
			ord, err := c1.addItems(uuid, cartPayload)
			if err != nil {
				log.Fatal(err)
				formatter.JSON(w, http.StatusBadRequest, err)
			} else {
				formatter.JSON(w, http.StatusOK, ord)
			}
		}
	}
}


// API Add Items to cart
func addItemsToCartHandler(formatter *render.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		setupResponse(&w, req)
		if (*req).Method == "OPTIONS" {
			fmt.Println("PREFLIGHT Request")
			return
		}
		var cartPayload Order
		if err := json.NewDecoder(req.Body).Decode(&cartPayload); err != nil {
			fmt.Println(" Error: ", err)
			formatter.JSON(w, http.StatusBadRequest, "Invalid request payload")
			return
		}
		fmt.Println("Adding Items to Cart: ", cartPayload)

		var uuid string = cartPayload.Username+"_cart";
		fmt.Println( "Order Params ID: ", uuid )
		//value := "Order Processed"

		if cartPayload.Username == ""  {
			formatter.JSON(w, http.StatusBadRequest, "Invalid Request. Order ID Missing.")
		} else {
			c1 := NewClient(server1)
			ord, err := c1.addItems(uuid, cartPayload)
			if err != nil {
				log.Fatal(err)
				formatter.JSON(w, http.StatusBadRequest, err)
			} else {
				formatter.JSON(w, http.StatusOK, ord)
			}
		}
	}
}

  

