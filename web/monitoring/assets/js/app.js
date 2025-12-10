function app() {
    return {
        activeTab: 'dashboard',
        sidebarOpen: true,
        isDark: false, // synced in init

        tabs: [
            { id: 'dashboard', label: 'Dashboard', icon: '<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect x="3" y="3" width="7" height="7"></rect><rect x="14" y="3" width="7" height="7"></rect><rect x="14" y="14" width="7" height="7"></rect><rect x="3" y="14" width="7" height="7"></rect></svg>' },
            { id: 'endpoints', label: 'Endpoints', icon: '<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><polyline points="22 12 18 12 15 21 9 3 6 12 2 12"></polyline></svg>' },
            { id: 'redis', label: 'Redis', icon: '<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="10"></circle><line x1="12" y1="8" x2="12" y2="12"></line><line x1="12" y1="16" x2="12.01" y2="16"></line></svg>' },
            { id: 'postgres', label: 'Postgres', icon: '<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M2 12h20"></path><path d="M12 2v20"></path><path d="M20 20a1 1 0 1 0 2 0a1 1 0 1 0-2 0"></path><path d="M4 20a1 1 0 1 0 2 0a1 1 0 1 0-2 0"></path><path d="M20 4a1 1 0 1 0 2 0a1 1 0 1 0-2 0"></path><path d="M4 4a1 1 0 1 0 2 0a1 1 0 1 0-2 0"></path></svg>' },
            { id: 'kafka', label: 'Kafka', icon: '<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"></path><polyline points="7 10 12 15 17 10"></polyline><line x1="12" y1="15" x2="12" y2="3"></line></svg>' },
            { id: 'cron', label: 'Cron Jobs', icon: '<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="10"></circle><polyline points="12 6 12 12 16 14"></polyline></svg>' },
            { id: 'config', label: 'Config', icon: '<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M12.22 2h-.44a2 2 0 0 0-2 2v.18a2 2 0 0 1-1 1.73l-.43.25a2 2 0 0 1-2 0l-.15-.08a2 2 0 0 0-2.73.73l-.22.38a2 2 0 0 0 .73 2.73l.15.1a2 2 0 0 1 1 1.72v.51a2 2 0 0 1-1 1.74l-.15.09a2 2 0 0 0-.73 2.73l.22.38a2 2 0 0 0 2.73.73l.15-.08a2 2 0 0 1 2 0l.43.25a2 2 0 0 1 1 1.73V20a2 2 0 0 0 2 2h.44a2 2 0 0 0 2-2v-.18a2 2 0 0 1 1-1.73l.43-.25a2 2 0 0 1 2 0l.15.08a2 2 0 0 0 2.73-.73l.22-.39a2 2 0 0 0-.73-2.73l-.15-.1a2 2 0 0 1-1-1.74v-.47a2 2 0 0 1 1-1.74l.15-.1a2 2 0 0 0 .73-2.73l-.22-.38a2 2 0 0 0-2.73-.73l-.15.08a2 2 0 0 1-2 0l-.43-.25a2 2 0 0 1-1-1.73V4a2 2 0 0 0-2-2z"></path><circle cx="12" cy="12" r="3"></circle></svg>' },
            { id: 'banner', label: 'Banner', icon: '<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M4 19.5v-15A2.5 2.5 0 0 1 6.5 2H20v20H6.5a2.5 2.5 0 0 1 0-5H20"></path></svg>' }
        ],

        // Dashboard Data
        // Dashboard Data
        serviceCount: 0,
        cpuUsage: 0,
        logs: [],
        cpuChart: null,
        endpoints: [],
        dummyLogActive: false,
        cronJobs: [],
        appConfig: {},
        configContent: '', // New
        bannerContent: '',

        // System Data
        infraStats: { total: 0, active: 0, items: [] },
        infraStatus: {}, // New
        sysInfo: { hostname: '', ip: '', disk: {} },

        // Redis Data // New

        // Redis Data
        redisPattern: '*',
        redisKeys: [],
        redisModalOpen: false,
        selectedRedisKey: '',
        selectedRedisValue: '',

        // Postgres Data
        pgInfo: {},
        pgQueries: [],

        // Kafka
        kafkaMsg: '',

        init() {
            // Theme init
            this.isDark = localStorage.getItem('theme') === 'dark' || (!('theme' in localStorage) && window.matchMedia('(prefers-color-scheme: dark)').matches);
            if (this.isDark) {
                document.documentElement.classList.add('dark');
            } else {
                document.documentElement.classList.remove('dark');
            }

            // SSE
            this.connectLogs();

            // Periodic
            this.fetchStatus();
            this.fetchDummyStatus();
            setInterval(() => this.fetchStatus(), 5000);

            // Watch tab changes to load data
            this.$watch('activeTab', (val) => {
                if (val === 'endpoints') this.fetchEndpoints();
                if (val === 'redis') this.fetchRedisKeys();
                if (val === 'postgres') { this.fetchPgQueries(); this.fetchPgInfo(); }
                if (val === 'kafka') this.fetchKafka();
                if (val === 'cron') this.fetchCronJobs();
                if (val === 'config') this.fetchConfig();
                if (val === 'banner') this.fetchBanner();
            });
        },

        get activeTabLabel() {
            return this.tabs.find(t => t.id === this.activeTab).label;
        },

        toggleTheme() {
            this.isDark = !this.isDark;
            if (this.isDark) {
                document.documentElement.classList.add('dark');
                localStorage.setItem('theme', 'dark');
            } else {
                document.documentElement.classList.remove('dark');
                localStorage.setItem('theme', 'light');
            }
        },

        logout() {
            // Simple approach for basic auth: 
            // 1. Redirect to URL with invalid credentials (often works to clear cache)
            // 2. Or just reload/home.
            // A common trick is to fetch with bad creds then reload.
            // For now, naive reload to home.
            window.location.reload();
        },

        connectLogs() {
            const es = new EventSource('/api/logs');
            let buffer = [];
            const MAX_LOGS = 100;

            es.onmessage = (e) => {
                const msg = e.data.replace(/\u001b\[\d+m/g, "");
                buffer.push(msg);
            };

            // Flush buffer periodically to avoid UI freezing (Throttled 1s)
            setInterval(() => {
                if (buffer.length > 0) {
                    // Append new logs
                    this.logs.push(...buffer);
                    buffer = [];

                    // Trim roughly
                    if (this.logs.length > MAX_LOGS) {
                        this.logs = this.logs.slice(-MAX_LOGS);
                    }

                    // Auto scroll
                    this.$nextTick(() => {
                        const box = document.getElementById('logs-box');
                        if (box) box.scrollTop = box.scrollHeight;
                    });
                }
            }, 1000);
        },

        formatLog(logLine) {
            try {
                // Try parsing JSON first (Zerolog format)
                const data = JSON.parse(logLine);
                const time = data.time ? new Date(data.time).toLocaleTimeString() : '';
                const level = (data.level || 'UNKNOWN').toUpperCase();
                const msg = data.message || JSON.stringify(data);

                let levelColor = 'text-gray-400';
                if (level === 'INFO') levelColor = 'text-blue-400';
                if (level === 'WARN' || level === 'WARNING') levelColor = 'text-yellow-400';
                if (level === 'ERROR' || level === 'FATAL') levelColor = 'text-red-400';
                if (level === 'DEBUG') levelColor = 'text-purple-400';

                return `<span class="text-gray-600 dark:text-gray-500 w-[80px] inline-block shrink-0">${time}</span>
                        <span class="${levelColor} font-bold w-[60px] inline-block shrink-0">${level}</span>
                        <span class="text-gray-300">${msg}</span>`;
            } catch (e) {
                // Fallback for raw text
                let color = 'text-gray-300';
                if (logLine.includes('INFO')) color = 'text-blue-300';
                if (logLine.includes('WARN')) color = 'text-yellow-300';
                if (logLine.includes('ERROR') || logLine.includes('FAIL')) color = 'text-red-300';
                if (logLine.includes('DEBUG')) color = 'text-purple-300';

                return `<span class="${color}">${logLine}</span>`;
            }
        },

        async fetchStatus() {
            try {
                const res = await fetch('/api/status');
                const data = await res.json();

                // Services
                this.serviceCount = Object.values(data.services || {}).filter(Boolean).length;

                // Infrastructure
                const infra = data.infrastructure || {};
                this.infraStatus = infra;
                const infraKeys = Object.keys(infra);
                const activeInfra = Object.values(infra).filter(Boolean).length;
                this.infraStats = {
                    total: infraKeys.length,
                    active: activeInfra,
                    items: infraKeys.map(k => ({ name: k, active: infra[k] }))
                };

                // System
                if (data.system) {
                    this.sysInfo = {
                        hostname: data.system.network?.hostname || 'Unknown',
                        ip: data.system.network?.ip || 'Unknown',
                        disk: data.system.disk || { used_percent: 0, total_gb: 0, used_gb: 0 }
                    };
                }

            } catch (e) { }
        },

        async fetchDummyStatus() {
            try {
                const res = await fetch('/api/logs/dummy/status');
                const data = await res.json();
                this.dummyLogActive = data.active;
            } catch (e) { }
        },

        async toggleDummyLog() {
            try {
                const res = await fetch('/api/logs/dummy', {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({ enable: !this.dummyLogActive })
                });
                const data = await res.json();
                this.dummyLogActive = !this.dummyLogActive;
                // Force sync
                this.fetchDummyStatus();
            } catch (e) { }
        },

        async fetchEndpoints() {
            try {
                const res = await fetch('/api/endpoints');
                this.endpoints = await res.json();
            } catch (e) { this.endpoints = []; }
        },

        async fetchCronJobs() {
            try {
                const res = await fetch('/api/cron');
                this.cronJobs = await res.json();
            } catch (e) { this.cronJobs = []; }
        },

        async fetchConfig() {
            try {
                // Fetch raw for editor
                const res = await fetch('/api/config/raw');
                const data = await res.json();
                this.configContent = data.content || '';

                // Keep appConfig for viewing if needed, but we focus on editor now
                // this.appConfig = ...
            } catch (e) { this.configContent = '# Error loading config.yaml'; }
        },

        async saveConfigText() {
            try {
                const res = await fetch('/api/config', {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({ content: this.configContent })
                });
                const data = await res.json();
                if (res.ok) {
                    alert(data.message);
                } else {
                    alert('Error: ' + data.error);
                }
            } catch (e) { alert('Failed to save config'); }
        },

        async backupConfig() {
            try {
                const res = await fetch('/api/config/backup', { method: 'POST' });
                const data = await res.json();
                if (res.ok) {
                    alert(data.message);
                } else {
                    alert('Error: ' + data.error);
                }
            } catch (e) { alert('Failed to backup config'); }
        },

        async fetchBanner() {
            try {
                const res = await fetch('/api/banner');
                const data = await res.json();
                this.bannerContent = data.content || '';
            } catch (e) { this.bannerContent = 'Error loading banner'; }
        },

        async saveBanner() {
            try {
                const res = await fetch('/api/banner', {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({ content: this.bannerContent })
                });
                if (res.ok) {
                    alert('Banner saved successfully!');
                } else {
                    alert('Failed to save banner.');
                }
            } catch (e) { alert('Error saving banner'); }
        },

        async fetchRedisKeys() {
            try {
                const res = await fetch(`/api/redis/keys?pattern=${encodeURIComponent(this.redisPattern)}`);
                const data = await res.json();
                this.redisKeys = Array.isArray(data) ? data : [];
            } catch (e) { this.redisKeys = []; }
        },

        async viewRedisValue(key) {
            this.selectedRedisKey = key;
            this.selectedRedisValue = 'Loading...';
            this.redisModalOpen = true;
            try {
                const res = await fetch(`/api/redis/key/${encodeURIComponent(key)}`);
                const data = await res.json();
                this.selectedRedisValue = data.value;
            } catch (e) { this.selectedRedisValue = 'Error fetching value'; }
        },

        async fetchPgQueries() {
            try {
                const res = await fetch('/api/postgres/queries');
                this.pgQueries = await res.json();
            } catch (e) { this.pgQueries = []; }
        },

        async fetchPgInfo() {
            try {
                const res = await fetch('/api/postgres/info');
                this.pgInfo = await res.json();
            } catch (e) { }
        },

        async fetchKafka() {
            try {
                const res = await fetch('/api/kafka/topics');
                const data = await res.json();
                this.kafkaMsg = JSON.stringify(data, null, 2);
            } catch (e) { }
        }
    }
}
