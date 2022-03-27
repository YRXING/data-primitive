package util

import (
	"fmt"
	"github.com/hashicorp/consul/api"
	log "github.com/sirupsen/logrus"
	"time"
)

type ConsulService struct {
	IP   string
	Port int
	Tag  []string
	Name string
}

func RegisterService(consulAddress string, service *ConsulService) {
	consulConfig := api.DefaultConfig()
	consulConfig.Address = consulAddress
	client, err := api.NewClient(consulConfig)
	if err != nil {
		log.Errorf("New consul client err \n:%v", err)
		return
	}
	agent := client.Agent()

	reg := &api.AgentServiceRegistration{
		ID:      fmt.Sprintf("%v-%v-%v", service.Name, service.IP, service.Port), // service id
		Name:    service.Name,                                                    // service name
		Tags:    service.Tag,
		Port:    service.Port,
		Address: service.IP,
		Check: &api.AgentServiceCheck{ // health check
			//HTTP: fmt.Sprintf("http://%s:%d",service.IP,service.Port),
			Interval:                       (time.Duration(10) * time.Second).String(),                      // check interval
			GRPC:                           fmt.Sprintf("%v:%v/%v", service.IP, service.Port, service.Name), // support grpc
			DeregisterCriticalServiceAfter: (time.Duration(1) * time.Minute).String(),                       // deregister time
		},
	}
	log.Printf("register to %v\n", consulAddress)
	if err := agent.ServiceRegister(reg); err != nil {
		log.Printf("service register error %v", err)
		return
	}
}

func DeregisterService(consulAddress, serviceName string) {
	cfg := api.DefaultConfig()
	cfg.Address = consulAddress
	client, err := api.NewClient(cfg)
	if err != nil {
		log.Errorf("New consul client err :%v", err)
		return
	}
	log.Printf("deregister from %v\n", consulAddress)
	if err := client.Agent().ServiceDeregister(serviceName); err != nil {
		log.Printf("service deregister error: %v", err)
		return
	}
}

func FilterServiceByID(consulAddress, filterId string) map[string]*api.AgentService {
	cfg := api.DefaultConfig()
	cfg.Address = consulAddress
	client, err := api.NewClient(cfg)
	if err != nil {
		log.Errorf("New consul client err :%v", err)
		return nil
	}
	data, err := client.Agent().ServicesWithFilter(fmt.Sprintf("Service ==%s", filterId))
	if err != nil {
		panic(err)
	}
	return data
}

func FindServiceByID(consulAddress, serviceID string) *api.AgentService {
	cfg := api.DefaultConfig()
	cfg.Address = consulAddress
	client, err := api.NewClient(cfg)
	if err != nil {
		log.Errorf("New consul client err :%v", err)
		return nil
	}

	service, _, err := client.Agent().Service(serviceID, nil)

	//var lastIndex uint64
	//client.Health().Service("serviceName","serviceTag",true,&api.QueryOptions{
	//	WaitIndex: lastIndex,
	//})

	if err != nil {
		return nil
	}
	return service
}

