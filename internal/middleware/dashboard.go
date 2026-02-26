package middleware

import (
	"bufio"
	"encoding/json"
	"html/template"
	"net/http"
	"os"

	"github.com/princetheprogrammer/apisentinel/internal/logger"
)

// DashboardData represents all the information for the HTML template.
type DashboardData struct {
	Stats      *Metrics
	RecentLogs []logger.AuditEvent
	BlockedIPs []string
}

const dashboardTemplate = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>🛡️ API Sentinel Dashboard</title>
    <style>
        body { font-family: 'Courier New', Courier, monospace; background: #f0f0f0; color: #000; padding: 2rem; }
        h1 { border-bottom: 5px solid #000; padding-bottom: 0.5rem; display: inline-block; }
        .grid { display: grid; grid-template-columns: 1fr 1fr; gap: 2rem; margin-top: 2rem; }
        .card { border: 4px solid #000; background: #fff; padding: 1.5rem; box-shadow: 10px 10px 0px #000; }
        .card.danger { border-color: #ff0000; box-shadow: 10px 10px 0px #ff0000; }
        h2 { border-bottom: 2px solid #000; margin-top: 0; }
        table { width: 100%; border-collapse: collapse; margin-top: 1rem; }
        th, td { border: 2px solid #000; padding: 0.5rem; text-align: left; }
        th { background: #000; color: #fff; }
        .stat-value { font-size: 3rem; font-weight: bold; }
        .violation { color: #ff0000; font-weight: bold; }
    </style>
</head>
<body>
    <h1>🛡️ API SENTINEL: SYSTEM STATUS</h1>
    
    <div class="grid">
        <div class="card">
            <h2>📊 METRICS</h2>
            <div>TOTAL REQUESTS: <span class="stat-value">{{.Stats.TotalRequests}}</span></div>
            <div style="color: #ff0000;">BLOCKED ATTACKS: <span class="stat-value">{{.Stats.BlockedRequests}}</span></div>
        </div>
        
        <div class="card">
            <h2>🚫 CURRENT BLOCKLIST</h2>
            <ul>
                {{range .BlockedIPs}}
                <li><strong>{{.}}</strong></li>
                {{else}}
                <li>No IPs currently blocked.</li>
                {{end}}
            </ul>
        </div>
    </div>

    <div class="card danger" style="margin-top: 3rem;">
        <h2>📜 RECENT SECURITY AUDIT LOGS</h2>
        <table>
            <thead>
                <tr>
                    <th>TIMESTAMP</th>
                    <th>REQUEST ID</th>
                    <th>LOCATION</th>
                    <th>SOURCE IP</th>
                    <th>VIOLATION</th>
                    <th>PATH</th>
                </tr>
            </thead>
            <tbody>
                {{range .RecentLogs}}
                <tr>
                    <td>{{.Timestamp.Format "2006-01-02 15:04:05"}}</td>
                    <td><code>{{.RequestID}}</code></td>
                    <td>{{.Location}}</td>
                    <td>{{.SourceIP}}</td>
                    <td class="violation">{{.ViolationType}}</td>
                    <td>{{.Path}}</td>
                </tr>
                {{else}}
                <tr>
                    <td colspan="4" style="text-align: center;">No security events recorded.</td>
                </tr>
                {{end}}
            </tbody>
        </table>
    </div>

    <div style="margin-top: 2rem; font-size: 0.8rem;">
        API SENTINEL v1.0 | PrinceTheProgrammer | STATUS: RUNNING
    </div>
</body>
</html>
`

// DashboardHandler renders the HTML dashboard.
func DashboardHandler(bl *IPBlocklist, auditPath string) http.HandlerFunc {
	tmpl := template.Must(template.New("dashboard").Parse(dashboardTemplate))

	return func(w http.ResponseWriter, r *http.Request) {
		data := DashboardData{
			Stats:      GlobalMetrics,
			BlockedIPs: make([]string, 0),
			RecentLogs: make([]logger.AuditEvent, 0),
		}

		// 1. Get Blocked IPs
		bl.mu.RLock()
		for ip := range bl.blocked {
			data.BlockedIPs = append(data.BlockedIPs, ip)
		}
		bl.mu.RUnlock()

		// 2. Read Recent Logs
		data.RecentLogs = readLastLogs(auditPath, 10)

		w.Header().Set("Content-Type", "text/html")
		tmpl.Execute(w, data)
	}
}

func readLastLogs(path string, count int) []logger.AuditEvent {
	file, err := os.Open(path)
	if err != nil {
		return nil
	}
	defer file.Close()

	var logs []logger.AuditEvent
	scanner := bufio.NewScanner(file)
	
	// Simple approach: read everything and take the last N.
	// In pro tier, we would read the file backwards.
	var allLogs []logger.AuditEvent
	for scanner.Scan() {
		var event logger.AuditEvent
		if err := json.Unmarshal(scanner.Bytes(), &event); err == nil {
			allLogs = append(allLogs, event)
		}
	}

	start := len(allLogs) - count
	if start < 0 {
		start = 0
	}
	
	// Reverse the logs so latest is first
	for i := len(allLogs) - 1; i >= start; i-- {
		logs = append(logs, allLogs[i])
	}

	return logs
}
