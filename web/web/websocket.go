package web

import "net/http"

func cors(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")                                                            // 允许访问所有域，可以换成具体url，注意仅具体url才能带cookie信息
		w.Header().Add("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, Token") //header的类型
		w.Header().Add("Access-Control-Allow-Credentials", "true")                                                    //设置为true，允许ajax异步请求带cookie信息
		w.Header().Add("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")                             //允许请求方法
		w.Header().Set("content-type", "application/json;charset=UTF-8")                                              //返回数据格式是json
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		f(w, r)
	}
}

//type connection struct {
//	ws   *websocket.Conn
//	sc   chan []byte
//	data *data
//}
//
//var wu = &websocket.Upgrader{ReadBufferSize: 512,
//	WriteBufferSize: 512, CheckOrigin: func(r *http.Request) bool { return true }}
//
//func myws(w http.ResponseWriter, r *http.Request) {
//	ws, err := wu.Upgrade(w, r, nil)
//	if err != nil {
//		return
//	}
//	c := &connection{sc: make(chan []byte, 256), ws: ws, data: &data{}}
//	h.r <- c
//	go c.writer()
//	c.reader()
//	defer func() {
//		c.data.Type = "logout"
//		user_list = del(user_list, c.data.User)
//		c.data.UserList = user_list
//		c.data.Content = c.data.User
//		data_b, _ := json.Marshal(c.data)
//		h.b <- data_b
//		h.r <- c
//	}()
//}
//
//func (c *connection) writer() {
//	for message := range c.sc {
//		c.ws.WriteMessage(websocket.TextMessage, message)
//	}
//	c.ws.Close()
//}
//
//var user_list = []string{}
//
//func (c *connection) reader() {
//	for {
//		_, message, err := c.ws.ReadMessage()
//		if err != nil {
//			h.r <- c
//			break
//		}
//		json.Unmarshal(message, &c.data)
//		switch c.data.Type {
//		case "login":
//			c.data.User = c.data.Content
//			c.data.From = c.data.User
//			user_list = append(user_list, c.data.User)
//			c.data.UserList = user_list
//			data_b, _ := json.Marshal(c.data)
//			h.b <- data_b
//		case "user":
//			c.data.Type = "user"
//			data_b, _ := json.Marshal(c.data)
//			h.b <- data_b
//		case "logout":
//			c.data.Type = "logout"
//			user_list = del(user_list, c.data.User)
//			data_b, _ := json.Marshal(c.data)
//			h.b <- data_b
//			h.r <- c
//		default:
//			fmt.Print("========default================")
//		}
//	}
//}
//
//func del(slice []string, user string) []string {
//	count := len(slice)
//	if count == 0 {
//		return slice
//	}
//	if count == 1 && slice[0] == user {
//		return []string{}
//	}
//	var n_slice = []string{}
//	for i := range slice {
//		if slice[i] == user && i == count {
//			return slice[:count]
//		} else if slice[i] == user {
//			n_slice = append(slice[:i], slice[i+1:]...)
//			break
//		}
//	}
//	fmt.Println(n_slice)
//	return n_slice
//}
