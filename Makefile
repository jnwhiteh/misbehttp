include $(GOROOT)/src/Make.inc

TARG=misbehttp
GOFILES=\
		main.go \
		counter.go \
		clients.go \
#		trickle_req.go \
		trickle_resp.go \
		norequest.go \
		disconnect.go \
		badrequest.go \

include $(GOROOT)/src/Make.cmd
