package cmd

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/kevwan/mapreduce"
	"github.com/spf13/cast"
	"github.com/spf13/pflag"
)

var (
	ProtocolPorts []string
	LifeCycleTime int // second
	IsListen      bool
	IP            string
	V             bool
)

type ProtocolPort struct {
	Protocol string
	Port     int
}

type Listener struct {
	Listener interface{}

	Protocol string
}

func main() {
	bindFlags()

	checkPorts, err := preprocessing()
	if err != nil {
		fmt.Println(err)
		return
	}

	sig := make(chan os.Signal, 1)
	if LifeCycleTime != 0 && IsListen {
		time.AfterFunc(time.Duration(LifeCycleTime)*time.Second, func() {
			sig <- syscall.SIGQUIT
		})
	}

	if !IsListen {
		// client to request
		reduce, _ := mapreduce.MapReduce(func(source chan<- interface{}) {
			for _, v := range checkPorts {
				source <- v
			}
		}, func(item interface{}, writer mapreduce.Writer, cancel func(error)) {
			v := item.(ProtocolPort)
			switch v.Protocol {
			case "tcp":
				resp, err := http.Get(fmt.Sprintf("http://%s:%d", IP, v.Port))
				if err != nil {
					writer.Write(fmt.Sprintf("%s:%d", IP, v.Port))
					return
				}
				all, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					writer.Write(fmt.Sprintf("%s:%d", IP, v.Port))
					return
				}
				_ = resp.Body.Close()
				if V {
					fmt.Println(string(all))
				}
			case "udp":
				// todo
			}
		}, func(pipe <-chan interface{}, writer mapreduce.Writer, cancel func(error)) {
			var errAddr []string
			for v := range pipe {
				errAddr = append(errAddr, v.(string))
			}
			writer.Write(errAddr)
		})
		if len(reduce.([]string)) != 0 {
			fmt.Printf("err addr: %s", strings.Join(reduce.([]string), ","))
		}
		return
	}

	var conflictPorts []string
	var listeners []Listener

	_ = mapreduce.MapReduceVoid(func(source chan<- interface{}) {
		for _, v := range checkPorts {
			source <- v
		}
	}, func(item interface{}, writer mapreduce.Writer, cancel func(error)) {
		v := item.(ProtocolPort)

		var listener interface{}
		var err error
		switch v.Protocol {
		case "tcp":
			listener, err = net.Listen("tcp", fmt.Sprintf("%s:%d", IP, v.Port))
		case "udp":
			addr, _ := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", IP, v.Port))
			listener, err = net.ListenUDP("udp", addr)
			if err != nil {
				fmt.Println(err)
			}
		}
		if err != nil {
			writer.Write(v.Port)
		} else {
			writer.Write(Listener{
				Listener: listener,
				Protocol: v.Protocol,
			})
		}
	}, func(pipe <-chan interface{}, cancel func(error)) {
		for v := range pipe {
			if value, ok := v.(int); ok && value != 0 {
				conflictPorts = append(conflictPorts, cast.ToString(v))
			} else {
				if value, ok := v.(Listener); ok {
					listeners = append(listeners, value)
				}
			}
		}
	})

	if len(conflictPorts) != 0 {
		for _, v := range listeners {
			switch v.Protocol {
			case "tcp":
				_ = v.Listener.(net.Listener).Close()
			case "udp":
				_ = v.Listener.(*net.UDPConn).Close()
			}
		}
		fmt.Printf("conflict ports: %s", strings.Join(conflictPorts, ","))
		os.Exit(1)
		return
	}

	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		_, _ = writer.Write([]byte("success"))
	})

	for _, v := range listeners {
		v := v
		go func() {
			switch v.Protocol {
			case "tcp":
				_ = http.Serve(v.Listener.(net.Listener), nil)
			case "udp":
				for {
					data := make([]byte, 1024)
					_, rAddr, err := v.Listener.(*net.UDPConn).ReadFromUDP(data)
					if err != nil {
						sig <- syscall.SIGQUIT
						return
					}
					strData := string(data)
					if V {
						fmt.Print("Received:", strData)
					}

					_, err = v.Listener.(*net.UDPConn).WriteToUDP([]byte(strData), rAddr)
					if err != nil {
						sig <- syscall.SIGQUIT
						return
					}
					if V {
						fmt.Print("Send:", strData)
					}
				}
			}
		}()
	}

	signalHandler(listeners, sig)
}

func preprocessing() ([]ProtocolPort, error) {
	var protocolPorts []ProtocolPort
	for _, v := range ProtocolPorts {
		var protocol, port string
		split := strings.Split(v, "/")
		if len(split) == 2 {
			port = split[0]
			protocol = strings.ToLower(split[1])
		} else {
			port = split[0]
			protocol = "tcp"
		}

		portInt, err := cast.ToIntE(port)
		if err != nil {
			return nil, err
		}

		protocolPorts = append(protocolPorts, ProtocolPort{
			Protocol: protocol,
			Port:     portInt,
		})
	}
	return protocolPorts, nil
}

func signalHandler(listeners []Listener, sig chan os.Signal) {
	// signal handler
	signal.Notify(sig, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		s := <-sig
		fmt.Printf("get a signal %s\n", s.String())
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			for _, v := range listeners {
				switch v.Protocol {
				case "tcp":
					_ = v.Listener.(net.Listener).Close()
				case "udp":
					_ = v.Listener.(*net.UDPConn).Close()
				}
			}
			fmt.Printf("close network successfully")
			return
		default:
			return
		}
	}
}

func bindFlags() {
	pflag.StringSliceVarP(&ProtocolPorts, "ports", "p", nil, "set ports to pre check, default protocol is tcp")
	pflag.IntVarP(&LifeCycleTime, "lifecycle-time", "t", 0, "set listen lifecycle time when listen")
	pflag.StringVarP(&IP, "ip", "", "0.0.0.0", "set ip")
	pflag.BoolVarP(&IsListen, "listen", "l", false, "is listen")
	pflag.BoolVarP(&V, "", "v", false, "show info")

	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.Parse()
}
