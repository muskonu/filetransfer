package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"log"
	"signaling/proto"
)

func Connection(c *gin.Context) {
	client, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		logrus.Errorln("upgrade:", err)
		return
	}
	defer client.Close()
	id, err := register(client)
	if err != nil {
		returnError(client, err)
		return
	}
	defer hub.Delete(id)
	err = EventLoop(client)
	log.Println(err)
	returnError(client, err)
	return
}

func register(c *websocket.Conn) (string, error) {
	_, msg, err := c.ReadMessage()
	if err != nil {
		return "", errors.Wrap(err, "register readMessage error")
	}
	req, err := proto.PayloadRegisterRequest(msg)
	if err != nil {
		return "", errors.Wrap(err, "register message parse error")
	}
	switch req.Cmd {
	case proto.CmdRegister:
		if !hub.Register(req.SourceID, c) {
			return "", errors.New("the id is already registered")
		}
		return req.SourceID, nil
	}
	return "", errors.New("not a valid cmd")
}

func EventLoop(c *websocket.Conn) error {
	for {
		_, msg, err := c.ReadMessage()
		if err != nil {
			return errors.Wrap(err, "readMessage error")
		}

		req, err := proto.PayloadRequest(msg)
		if err != nil {
			err := c.WriteMessage(websocket.TextMessage, proto.NewResponse(proto.BadPayloadResponse, err.Error()))
			if err != nil {
				return errors.Wrap(err, "payload parse error")
			}
			continue
		}
		switch req.Cmd {
		case proto.CmdOffer, proto.CmdAnswer, proto.CmdCandidate:
			err := hub.SafeTransport(req, req.Cmd)
			if err != nil {
				return errors.Wrap(err, "transport error")
			}
		}
	}
}

func returnError(c *websocket.Conn, e error) {
	logrus.Errorln(e)
	c.WriteMessage(websocket.CloseMessage, proto.NewResponse(proto.CloseResponse, e.Error()))
	return
}
