package cmd

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/egichandrap/opsdoctor/utils"
	"github.com/spf13/cobra"
)

type HostDef struct {
	Name string `json:"name"`
	Addr string `json:"addr"`
}

var healthCmd = &cobra.Command{
	Use:   "health",
	Short: "Health checks",
}

var healthRunCmd = &cobra.Command{
	Use:   "run [hosts.json]",
	Short: "Run health checks for hosts defined in JSON file",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		fn := args[0]
		b, err := os.ReadFile(fn)
		if err != nil {
			return err
		}
		var hosts []HostDef
		if err := json.Unmarshal(b, &hosts); err != nil {
			return err
		}
		utils.VPrintf("Running health check for %d hosts\n", len(hosts))

		type Resp struct {
			Host    HostDef       `json:"host"`
			OK      bool          `json:"ok"`
			Detail  string        `json:"detail"`
			Latency time.Duration `json:"latency"`
		}

		var mu sync.Mutex
		results := make([]Resp, 0, len(hosts))
		wg := sync.WaitGroup{}
		for _, h := range hosts {
			h := h
			wg.Add(1)
			go func() {
				defer wg.Done()
				start := time.Now()
				ok, detail := doBasicCheck(h.Addr)
				lat := time.Since(start)
				mu.Lock()
				results = append(results, Resp{Host: h, OK: ok, Detail: detail, Latency: lat})
				mu.Unlock()
			}()
		}
		wg.Wait()

		if utils.IsJSONOutput() {
			b, _ := json.MarshalIndent(results, "", "  ")
			fmt.Println(string(b))
			return nil
		}
		for _, r := range results {
			if r.OK {
				utils.PrintColor("green", fmt.Sprintf("%s (%s) OK in %v", r.Host.Name, r.Host.Addr, r.Latency))
			} else {
				utils.PrintColor("red", fmt.Sprintf("%s (%s) FAIL: %s", r.Host.Name, r.Host.Addr, r.Detail))
			}
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(healthCmd)
	healthCmd.AddCommand(healthRunCmd)
}

func doBasicCheck(addr string) (bool, string) {
	timeout := 3 * time.Second
	conn, err := netDialTimeout(addr, timeout)
	if err != nil {
		return false, err.Error()
	}
	_ = conn.Close()
	return true, "connected"
}

func netDialTimeout(addr string, timeout time.Duration) (net.Conn, error) {
	conn, err := net.DialTimeout("tcp", addr, timeout)
	if err == nil {
		return conn, nil
	}
	if !strings.Contains(addr, ":") {
		if conn2, err2 := net.DialTimeout("tcp", addr+":443", timeout); err2 == nil {
			return conn2, nil
		}
		if conn3, err3 := net.DialTimeout("tcp", addr+":80", timeout); err3 == nil {
			return conn3, nil
		}
	}
	return nil, err
}
