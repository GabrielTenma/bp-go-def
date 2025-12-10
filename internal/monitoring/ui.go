package monitoring

const IndexHTML = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Go App Monitor</title>
    <style>
        :root {
            --bg-color: #121212;
            --text-color: #e0e0e0;
            --accent-color: #bb86fc;
            --secondary-bg: #1e1e1e;
            --border-color: #333;
        }
        
        [data-theme="light"] {
             --bg-color: #f5f5f5;
             --text-color: #121212;
             --accent-color: #6200ee;
             --secondary-bg: #ffffff;
             --border-color: #e0e0e0;
        }

        body { font-family: 'Inter', sans-serif; background: var(--bg-color); color: var(--text-color); margin: 0; padding: 20px; transition: background 0.3s, color 0.3s; }
        .container { max-width: 1200px; margin: 0 auto; }
        header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 20px; border-bottom: 1px solid var(--border-color); padding-bottom: 10px; }
        h1 { margin: 0; font-size: 1.5rem; }
        
        .grid { display: grid; grid-template-columns: repeat(auto-fit, minmax(300px, 1fr)); gap: 20px; margin-bottom: 20px; }
        .card { background: var(--secondary-bg); padding: 20px; border-radius: 8px; border: 1px solid var(--border-color); }
        .card h2 { margin-top: 0; font-size: 1.2rem; color: var(--accent-color); }
        
        pre { background: #000; color: #0f0; padding: 10px; border-radius: 4px; overflow-x: auto; white-space: pre-wrap; font-family: 'Consolas', monospace; font-size: 0.9rem; }
        #logs-container { height: 400px; overflow-y: auto; background: #0a0a0a; border: 1px solid var(--border-color); padding: 10px; font-family: 'Consolas', monospace; font-size: 0.85rem; }
        .log-entry { margin-bottom: 2px; }
        
        .badge { padding: 4px 8px; border-radius: 4px; font-size: 0.8rem; font-weight: bold; }
        .badge.on { background: #2e7d32; color: #fff; }
        .badge.off { background: #c62828; color: #fff; }
        
        button { background: var(--secondary-bg); color: var(--text-color); border: 1px solid var(--border-color); padding: 8px 16px; cursor: pointer; border-radius: 4px; }
        button:hover { border-color: var(--accent-color); }
        
        /* Minimalist Scrollbar */
        ::-webkit-scrollbar { width: 8px; }
        ::-webkit-scrollbar-track { background: var(--bg-color); }
        ::-webkit-scrollbar-thumb { background: #555; border-radius: 4px; }
    </style>
</head>
<body>
    <div class="container">
        <header>
            <h1>📊 System Monitor</h1>
            <button onclick="toggleTheme()">Toggle Theme</button>
        </header>
        
        <div class="grid">
            <div class="card">
                <h2>Services Status</h2>
                <div id="services-status">Loading...</div>
            </div>
            <div class="card">
                <h2>Infrastructure</h2>
                <div id="infra-status">Loading...</div>
            </div>
            <div class="card">
                <h2>Real-time CPU</h2>
                <h3 id="cpu-val" style="font-size: 2rem; margin: 0;">0%</h3>
            </div>
        </div>

        <div class="card">
            <div style="display:flex; justify-content:space-between; align-items:center; margin-bottom:10px;">
                <h2>Live Logs</h2>
                <button onclick="clearLogs()">Clear</button>
            </div>
            <div id="logs-container"></div>
        </div>
    </div>

    <script>
        // Theme
        function toggleTheme() {
            const current = document.body.getAttribute('data-theme');
            document.body.setAttribute('data-theme', current === 'light' ? 'dark' : 'light');
        }

        // Status
        async function fetchStatus() {
            try {
                const res = await fetch('/api/status');
                const data = await res.json();
                
                // Services
                let sHtml = '<ul>';
                for (const [k, v] of Object.entries(data.services || {})) {
                    sHtml += '<li>' + k + ': <span class="badge ' + (v ? 'on' : 'off') + '">' + (v ? 'ONLINE' : 'OFFLINE') + '</span></li>';
                }
                sHtml += '</ul>';
                document.getElementById('services-status').innerHTML = sHtml;

                // Infra
                let iHtml = '<ul>';
                for (const [k, v] of Object.entries(data.infrastructure || {})) {
                    iHtml += '<li>' + k + ': <span class="badge ' + (v ? 'on' : 'off') + '">' + (v ? 'ENABLED' : 'DISABLED') + '</span></li>';
                }
                iHtml += '</ul>';
                document.getElementById('infra-status').innerHTML = iHtml;
            } catch (e) {
                console.error(e);
            }
        }
        fetchStatus();
        setInterval(fetchStatus, 5000);

        // Logs SSE
        const logContainer = document.getElementById('logs-container');
        const logSource = new EventSource('/api/logs');
        logSource.onmessage = function(event) {
            const div = document.createElement('div');
            div.className = 'log-entry';
            // Simple ANSI stripping for web view if raw bytes come in with colors
            // Use a regex/library for full ansi support, minimal strip here:
            const cleanMsg = event.data.replace(/\u001b\[\d+m/g, ""); 
            div.textContent = cleanMsg; 
            logContainer.appendChild(div);
            logContainer.scrollTop = logContainer.scrollHeight;
            
            // Limit logs
            if (logContainer.children.length > 500) {
                logContainer.removeChild(logContainer.firstChild);
            }
        };

        // CPU SSE
        const cpuSource = new EventSource('/api/cpu');
        cpuSource.onmessage = function(event) {
            document.getElementById('cpu-val').textContent = event.data + '%';
        };

        function clearLogs() {
            logContainer.innerHTML = '';
        }
    </script>
</body>
</html>
`
