package fake

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/rand/v2"
	"net"
	"os"
	"slices"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/traPtitech/piscon-portal-v2/server/domain"
	"github.com/traPtitech/piscon-portal-v2/server/services/instance"
)

const (
	privateIPCidr = "192.168.0.1/24"
	publicIPCidr  = "192.0.2.1/24"

	ipInfoFileName   = "ipinfo.json"
	instanceFileName = "instances.json"
)

type IPInfo struct {
	PrivateIPs []string `json:"private_ips"`
	PublicIPs  []string `json:"public_ips"`
}

// Manager is a fake implementation of the instance.Manager interface.
// It simulates the behavior of an instance manager for development environments.
// It stores instance data in a local file system.
type Manager struct {
	mu   sync.Mutex
	root *os.Root
}

var _ instance.Manager = &Manager{}

func NewManager(root *os.Root) (*Manager, error) {
	m := &Manager{
		root: root,
	}

	// Ensure the instances file exists
	if _, err := m.root.Stat(instanceFileName); errors.Is(err, os.ErrNotExist) {
		f, err := m.root.Create(instanceFileName)
		if err != nil {
			return nil, fmt.Errorf("create instances file: %w", err)
		}
		defer f.Close()

		if err := json.NewEncoder(f).Encode([]domain.InfraInstance{}); err != nil {
			return nil, fmt.Errorf("initialize instances file: %w", err)
		}
	}

	// Ensure the IP info file exists
	if _, err := m.root.Stat(ipInfoFileName); errors.Is(err, os.ErrNotExist) {
		ipInfo := IPInfo{
			PrivateIPs: []string{},
			PublicIPs:  []string{},
		}
		if err := m.storeIPInfo(ipInfo); err != nil {
			return nil, fmt.Errorf("initialize IP info file: %w", err)
		}
	}

	return m, nil
}

func (m *Manager) Create(_ context.Context, _ string, _ []string) (string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	id := uuid.NewString()
	privateIP, err := m.generatePrivateIP()
	if err != nil {
		return "", fmt.Errorf("generate private IP: %w", err)
	}

	publicIP, err := m.generatePublicIP()
	if err != nil {
		return "", fmt.Errorf("generate public IP: %w", err)
	}

	instance := domain.InfraInstance{
		ProviderInstanceID: id,
		Status:             domain.InstanceStatusBuilding,
		PrivateIP:          &privateIP,
	}

	instances, err := m.readInstances()
	if err != nil {
		return "", fmt.Errorf("read instances: %w", err)
	}
	instances = append(instances, instance)

	if err := m.storeInstances(instances); err != nil {
		return "", fmt.Errorf("store instances: %w", err)
	}

	// Simulate a delay for building the instance
	time.AfterFunc(5*time.Second, func() {
		log.Println("try to lock instance for update")
		m.mu.Lock()
		defer m.mu.Unlock()
		log.Println("instance locked, updating status to running")
		instance.Status = domain.InstanceStatusRunning
		instance.PublicIP = &publicIP
		if err := m.updateInstance(instance); err != nil {
			log.Printf("update instance status: %v", err)
		}
	})

	return instance.ProviderInstanceID, nil
}

func (m *Manager) Get(_ context.Context, id string) (domain.InfraInstance, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	instances, err := m.readInstances()
	if err != nil {
		return domain.InfraInstance{}, fmt.Errorf("read instances: %w", err)
	}

	instance, ok := lo.Find(instances, func(i domain.InfraInstance) bool {
		return i.ProviderInstanceID == id
	})
	if !ok {
		return domain.InfraInstance{}, fmt.Errorf("instance not found: %s", id)
	}
	return instance, nil
}

func (m *Manager) GetAll(_ context.Context) ([]domain.InfraInstance, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	instances, err := m.readInstances()
	if err != nil {
		return nil, fmt.Errorf("read instances: %w", err)
	}
	return instances, nil
}

func (m *Manager) GetByIDs(_ context.Context, ids []string) ([]domain.InfraInstance, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	instances, err := m.readInstances()
	if err != nil {
		return nil, fmt.Errorf("read instances: %w", err)
	}

	result := make([]domain.InfraInstance, 0, len(instances))
	for _, instance := range instances {
		if slices.Contains(ids, instance.ProviderInstanceID) {
			result = append(result, instance)
		}
	}

	return result, nil
}

func (m *Manager) Delete(_ context.Context, instance domain.InfraInstance) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	instance.Status = domain.InstanceStatusDeleted

	instances, err := m.readInstances()
	if err != nil {
		return fmt.Errorf("read instances: %w", err)
	}

	instances = slices.DeleteFunc(instances, func(i domain.InfraInstance) bool {
		return i.ProviderInstanceID == instance.ProviderInstanceID
	})

	// Remove the IPs from the IP info
	ipInfo, err := m.readIPInfo()
	if err != nil {
		return fmt.Errorf("read IP info: %w", err)
	}
	ipInfo.PrivateIPs = lo.Without(ipInfo.PrivateIPs, *instance.PrivateIP)
	ipInfo.PublicIPs = lo.Without(ipInfo.PublicIPs, *instance.PublicIP)

	if err := m.storeIPInfo(ipInfo); err != nil {
		return fmt.Errorf("store IP info: %w", err)
	}

	if err := m.storeInstances(instances); err != nil {
		return fmt.Errorf("store instances: %w", err)
	}

	return nil
}

func (m *Manager) Stop(_ context.Context, instance domain.InfraInstance) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	instance.Status = domain.InstanceStatusStopping
	if err := m.updateInstance(instance); err != nil {
		return fmt.Errorf("update instance status to stopping: %w", err)
	}

	// Simulate a delay for stopping the instance
	time.AfterFunc(5*time.Second, func() {
		m.mu.Lock()
		defer m.mu.Unlock()
		instance.Status = domain.InstanceStatusStopped
		if err := m.updateInstance(instance); err != nil {
			log.Printf("update instance status to stopped: %v", err)
		}
	})

	return nil
}

func (m *Manager) Start(_ context.Context, instance domain.InfraInstance) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	instance.Status = domain.InstanceStatusStarting
	if err := m.updateInstance(instance); err != nil {
		return fmt.Errorf("update instance status to starting: %w", err)
	}

	// Simulate a delay for starting the instance
	time.AfterFunc(5*time.Second, func() {
		m.mu.Lock()
		defer m.mu.Unlock()
		instance.Status = domain.InstanceStatusRunning
		if err := m.updateInstance(instance); err != nil {
			log.Printf("update instance status to running: %v", err)
		}
	})

	return nil
}

func (m *Manager) readInstances() ([]domain.InfraInstance, error) {
	f, err := m.root.Open(instanceFileName)
	if err != nil {
		return nil, fmt.Errorf("open instances file: %w", err)
	}
	defer f.Close()

	var instances []domain.InfraInstance
	if err := json.NewDecoder(f).Decode(&instances); err != nil {
		return nil, fmt.Errorf("decode instances: %w", err)
	}
	return instances, nil
}

func (m *Manager) storeInstances(instances []domain.InfraInstance) error {
	f, err := m.root.Create(instanceFileName)
	if err != nil {
		return fmt.Errorf("create instances file: %w", err)
	}
	defer f.Close()

	if err := json.NewEncoder(f).Encode(instances); err != nil {
		return fmt.Errorf("encode instances: %w", err)
	}
	return nil
}

func (m *Manager) updateInstance(instance domain.InfraInstance) error {
	f, err := m.root.OpenFile(instanceFileName, os.O_RDWR, 0644)
	if err != nil {
		return fmt.Errorf("open instances file: %w", err)
	}
	defer f.Close()

	var instances []domain.InfraInstance
	if err := json.NewDecoder(f).Decode(&instances); err != nil {
		return fmt.Errorf("decode instances: %w", err)
	}

	for i, inst := range instances {
		if inst.ProviderInstanceID == instance.ProviderInstanceID {
			instances[i] = instance
			break
		}
	}

	// Truncate the file and write the updated instances
	if err := f.Truncate(0); err != nil {
		return fmt.Errorf("truncate instances file: %w", err)
	}
	if _, err := f.Seek(0, 0); err != nil {
		return fmt.Errorf("seek to beginning of instances file: %w", err)
	}
	if err := json.NewEncoder(f).Encode(instances); err != nil {
		return fmt.Errorf("encode updated instances: %w", err)
	}
	return nil
}

func (m *Manager) readIPInfo() (IPInfo, error) {
	f, err := m.root.Open(ipInfoFileName)
	if err != nil {
		return IPInfo{}, fmt.Errorf("open IP info file: %w", err)
	}
	defer f.Close()

	var ipInfo IPInfo
	if err := json.NewDecoder(f).Decode(&ipInfo); err != nil {
		return IPInfo{}, fmt.Errorf("decode IP info: %w", err)
	}
	return ipInfo, nil
}

func (m *Manager) storeIPInfo(ipInfo IPInfo) error {
	f, err := m.root.Create(ipInfoFileName)
	if err != nil {
		return fmt.Errorf("create IP info file: %w", err)
	}
	defer f.Close()

	if err := json.NewEncoder(f).Encode(ipInfo); err != nil {
		return fmt.Errorf("encode IP info: %w", err)
	}
	return nil
}

func (m *Manager) generatePrivateIP() (string, error) {
	ipInfo, err := m.readIPInfo()
	if err != nil {
		return "", fmt.Errorf("read IP info: %w", err)
	}

	for {
		ip := randomIPv4(privateIPCidr)
		if slices.Contains(ipInfo.PrivateIPs, ip) {
			continue // IP already exists, generate a new one
		}
		ipInfo.PrivateIPs = append(ipInfo.PrivateIPs, ip)

		if err := m.storeIPInfo(ipInfo); err != nil {
			return "", fmt.Errorf("store IP info: %w", err)
		}
		return ip, nil
	}
}

func (m *Manager) generatePublicIP() (string, error) {
	ipInfo, err := m.readIPInfo()
	if err != nil {
		return "", fmt.Errorf("read IP info: %w", err)
	}

	for {
		ip := randomIPv4(publicIPCidr)
		if slices.Contains(ipInfo.PublicIPs, ip) {
			continue // IP already exists, generate a new one
		}
		ipInfo.PublicIPs = append(ipInfo.PublicIPs, ip)

		if err := m.storeIPInfo(ipInfo); err != nil {
			return "", fmt.Errorf("store IP info: %w", err)
		}
		return ip, nil
	}
}

func randomIPv4(cidr string) string {
	_, ipnet, _ := net.ParseCIDR(cidr)
	ones, _ := ipnet.Mask.Size()
	numHosts := 1 << (32 - ones)
	randomOffset := rand.IntN(numHosts-2) + 1 // Exclude network and broadcast addresses
	randomIP := make(net.IP, 4)

	copy(randomIP, ipnet.IP)
	for i := 3; i >= 0; i-- {
		randomIP[i] += byte(randomOffset & 0xFF)
		randomOffset >>= 8
		if randomOffset == 0 {
			break
		}
	}
	return randomIP.String()
}
