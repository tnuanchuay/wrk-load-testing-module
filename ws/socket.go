package ws

import (
	"github.com/kataras/iris"
	"encoding/json"
	"time"
	"github.com/kataras/iris/adaptors/websocket"
)

type GroupSocket struct{
	Sockets	[]*websocket.Connection
}

func (r *GroupSocket) BroadcastAllExcept(channel string, data map[string]interface{}, so websocket.Connection){
	b, _ := json.Marshal(data)
	for _, s := range r.Sockets{
		if (*s).ID() != so.ID(){
			(*s).Emit(channel, string(b))
			time.Sleep(10 * time.Millisecond)
		}
	}
}


func (r *GroupSocket) DisconnectAll(){
	for _, so := range r.Sockets{
		(*so).Disconnect()
	}
	r = new(GroupSocket)
}

func (r *GroupSocket) DisconnectAllExcept(so websocket.Connection){
	for i, s := range r.Sockets{
		if (*s).ID() != so.ID(){
			r.Sockets = append(r.Sockets[:i], r.Sockets[i+1:]...)
			(*s).Disconnect()
			break;
		}

	}
}

func (r *GroupSocket) Disconnect(so websocket.Connection){

	for i, s := range r.Sockets{
		if (*s).ID() == so.ID(){
			r.Sockets = append(r.Sockets[:i], r.Sockets[i+1:]...)
			so.Disconnect()
			break;
		}

	}
}

func (r *GroupSocket) BroadCast(channel string, data map[string]interface{}){
	b, _ := json.Marshal(data)
	for _, s := range r.Sockets{
		(*s).Emit(channel, string(b))
		time.Sleep(10 * time.Millisecond)
	}
}