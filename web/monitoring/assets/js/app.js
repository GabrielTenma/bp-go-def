// Initialize Notyf Globally
window.notyf = new Notyf({
    duration: 2500,
    dismissible: true,
    position: { x: 'center', y: 'bottom' },
    ripple: false,
    types: [
        {
            type: 'success',
            background: '#16a34a',
            icon: {
                className: 'notyf__icon--success',
                tagName: 'i',
                color: 'white'
            }
        },
        {
            type: 'error',
            background: '#dc2626',
            duration: 5000
        },
        {
            type: 'info',
            background: '#3b82f6',
            icon: false
        }
    ]
});

function app() {
    return {
        activeTab: 'dashboard',
        sidebarOpen: false, // Mobile sidebar
        sidebarCollapsed: false, // Desktop sidebar (new)
        isDark: false, // synced in init

        menuCategories: [
            {
                name: 'General',
                items: [
                    { id: 'dashboard', label: 'Dashboard', icon: '<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect x="3" y="3" width="7" height="7"></rect><rect x="14" y="3" width="7" height="7"></rect><rect x="14" y="14" width="7" height="7"></rect><rect x="3" y="14" width="7" height="7"></rect></svg>' },
                    { id: 'endpoints', label: 'Endpoints', icon: '<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><polyline points="22 12 18 12 15 21 9 3 6 12 2 12"></polyline></svg>' }
                ]
            },
            {
                name: 'Infrastructure',
                items: [
                    { id: 'redis', label: 'Redis', icon: '<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="10"></circle><line x1="12" y1="8" x2="12" y2="12"></line><line x1="12" y1="16" x2="12.01" y2="16"></line></svg>' },
                    { id: 'postgres', label: 'Postgres', icon: '<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M2 12h20"></path><path d="M12 2v20"></path><path d="M20 20a1 1 0 1 0 2 0a1 1 0 1 0-2 0"></path><path d="M4 20a1 1 0 1 0 2 0a1 1 0 1 0-2 0"></path><path d="M20 4a1 1 0 1 0 2 0a1 1 0 1 0-2 0"></path><path d="M4 4a1 1 0 1 0 2 0a1 1 0 1 0-2 0"></path></svg>' },
                    { id: 'kafka', label: 'Kafka', icon: '<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"></path><polyline points="7 10 12 15 17 10"></polyline><line x1="12" y1="15" x2="12" y2="3"></line></svg>' },
                    { id: 'storage', label: 'Storage', icon: '<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"></path><polyline points="21 12 21 12"></polyline><rect width="20" height="8" x="2" y="4" rx="2" ry="2"></rect><line x1="10" y1="8" x2="14" y2="8"></line></svg>' },
                    { id: 'system', label: 'System', icon: '<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect x="2" y="2" width="20" height="8" rx="2" ry="2"></rect><rect x="2" y="14" width="20" height="8" rx="2" ry="2"></rect><line x1="6" y1="6" x2="6.01" y2="6"></line><line x1="6" y1="18" x2="6.01" y2="18"></line></svg>' },
                    { id: 'external', label: 'External', icon: '<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="10"></circle><line x1="2" y1="12" x2="22" y2="12"></line><path d="M12 2a15.3 15.3 0 0 1 4 10 15.3 15.3 0 0 1-4 10 15.3 15.3 0 0 1-4-10 15.3 15.3 0 0 1 4-10z"></path></svg>' }
                ]
            },
            {
                name: 'Periodic',
                items: [
                    { id: 'cron', label: 'Cron Jobs', icon: '<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="10"></circle><polyline points="12 6 12 12 16 14"></polyline></svg>' }
                ]
            },
            {
                name: 'Other',
                items: [
                    { id: 'config', label: 'Config', icon: '<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M12.22 2h-.44a2 2 0 0 0-2 2v.18a2 2 0 0 1-1 1.73l-.43.25a2 2 0 0 1-2 0l-.15-.08a2 2 0 0 0-2.73.73l-.22.38a2 2 0 0 0 .73 2.73l.15.1a2 2 0 0 1 1 1.72v.51a2 2 0 0 1-1 1.74l-.15.09a2 2 0 0 0-.73 2.73l.22.38a2 2 0 0 0 2.73.73l.15-.08a2 2 0 0 1 2 0l.43.25a2 2 0 0 1 1 1.73V20a2 2 0 0 0 2 2h.44a2 2 0 0 0 2-2v-.18a2 2 0 0 1 1-1.73l.43-.25a2 2 0 0 1 2 0l.15.08a2 2 0 0 0 2.73-.73l.22-.39a2 2 0 0 0-.73-2.73l-.15-.1a2 2 0 0 1-1-1.74v-.47a2 2 0 0 1 1-1.74l.15-.1a2 2 0 0 0 .73-2.73l-.22-.38a2 2 0 0 0-2.73-.73l-.15.08a2 2 0 0 1-2 0l-.43-.25a2 2 0 0 1-1-1.73V4a2 2 0 0 0-2-2z"></path><circle cx="12" cy="12" r="3"></circle></svg>' },
                    { id: 'banner', label: 'Banner', icon: '<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M4 19.5v-15A2.5 2.5 0 0 1 6.5 2H20v20H6.5a2.5 2.5 0 0 1 0-5H20"></path></svg>' },
                    { id: 'settings', label: 'User Settings', icon: '<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M19 21v-2a4 4 0 0 0-4-4H9a4 4 0 0 0-4 4v2"></path><circle cx="12" cy="7" r="4"></circle></svg>' },
                    { id: 'maintenance', label: 'Maintenance', icon: '<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M14.7 6.3a1 1 0 0 0 0 1.4l1.6 1.6a1 1 0 0 0 1.4 0l3.77-3.77a6 6 0 0 1-7.94 7.94l-6.91 6.91a2.12 2.12 0 0 1-3-3l6.91-6.91a6 6 0 0 1 7.94-7.94l-3.76 3.76z"></path></svg>' }
                ]
            }
        ],

        // Flat tabs array for backward compatibility
        tabs: [
            { id: 'dashboard', label: 'Dashboard', icon: '<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect x="3" y="3" width="7" height="7"></rect><rect x="14" y="3" width="7" height="7"></rect><rect x="14" y="14" width="7" height="7"></rect><rect x="3" y="14" width="7" height="7"></rect></svg>' },
            { id: 'endpoints', label: 'Endpoints', icon: '<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><polyline points="22 12 18 12 15 21 9 3 6 12 2 12"></polyline></svg>' },
            { id: 'redis', label: 'Redis', icon: '<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="10"></circle><line x1="12" y1="8" x2="12" y2="12"></line><line x1="12" y1="16" x2="12.01" y2="16"></line></svg>' },
            { id: 'postgres', label: 'Postgres', icon: '<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M2 12h20"></path><path d="M12 2v20"></path><path d="M20 20a1 1 0 1 0 2 0a1 1 0 1 0-2 0"></path><path d="M4 20a1 1 0 1 0 2 0a1 1 0 1 0-2 0"></path><path d="M20 4a1 1 0 1 0 2 0a1 1 0 1 0-2 0"></path><path d="M4 4a1 1 0 1 0 2 0a1 1 0 1 0-2 0"></path></svg>' },
            { id: 'kafka', label: 'Kafka', icon: '<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"></path><polyline points="7 10 12 15 17 10"></polyline><line x1="12" y1="15" x2="12" y2="3"></line></svg>' },
            { id: 'storage', label: 'Storage', icon: '<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"></path><polyline points="21 12 21 12"></polyline><rect width="20" height="8" x="2" y="4" rx="2" ry="2"></rect><line x1="10" y1="8" x2="14" y2="8"></line></svg>' },
            { id: 'system', label: 'System', icon: '<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect x="2" y="2" width="20" height="8" rx="2" ry="2"></rect><rect x="2" y="14" width="20" height="8" rx="2" ry="2"></rect><line x1="6" y1="6" x2="6.01" y2="6"></line><line x1="6" y1="18" x2="6.01" y2="18"></line></svg>' },
            { id: 'external', label: 'External', icon: '<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="10"></circle><line x1="2" y1="12" x2="22" y2="12"></line><path d="M12 2a15.3 15.3 0 0 1 4 10 15.3 15.3 0 0 1-4 10 15.3 15.3 0 0 1-4-10 15.3 15.3 0 0 1 4-10z"></path></svg>' },
            { id: 'cron', label: 'Cron Jobs', icon: '<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="10"></circle><polyline points="12 6 12 12 16 14"></polyline></svg>' },
            { id: 'config', label: 'Config', icon: '<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M12.22 2h-.44a2 2 0 0 0-2 2v.18a2 2 0 0 1-1 1.73l-.43.25a2 2 0 0 1-2 0l-.15-.08a2 2 0 0 0-2.73.73l-.22.38a2 2 0 0 0 .73 2.73l.15.1a2 2 0 0 1 1 1.72v.51a2 2 0 0 1-1 1.74l-.15.09a2 2 0 0 0-.73 2.73l.22.38a2 2 0 0 0 2.73.73l.15-.08a2 2 0 0 1 2 0l.43.25a2 2 0 0 1 1 1.73V20a2 2 0 0 0 2 2h.44a2 2 0 0 0 2-2v-.18a2 2 0 0 1 1-1.73l.43-.25a2 2 0 0 1 2 0l.15.08a2 2 0 0 0 2.73-.73l.22-.39a2 2 0 0 0-.73-2.73l-.15-.1a2 2 0 0 1-1-1.74v-.47a2 2 0 0 1 1-1.74l.15-.1a2 2 0 0 0 .73-2.73l-.22-.38a2 2 0 0 0-2.73-.73l-.15.08a2 2 0 0 1-2 0l-.43-.25a2 2 0 0 1-1-1.73V4a2 2 0 0 0-2-2z"></path><circle cx="12" cy="12" r="3"></circle></svg>' },
            { id: 'banner', label: 'Banner', icon: '<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M4 19.5v-15A2.5 2.5 0 0 1 6.5 2H20v20H6.5a2.5 2.5 0 0 1 0-5H20"></path></svg>' },
            { id: 'settings', label: 'User Settings', icon: '<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M19 21v-2a4 4 0 0 0-4-4H9a4 4 0 0 0-4 4v2"></path><circle cx="12" cy="7" r="4"></circle></svg>' }
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
        monitoringConfig: { title: 'GoBP Admin', subtitle: 'Go Echo Boilerplate' }, // New

        // User Settings
        userSettings: { username: '', photoPath: '' },
        passwordForm: { current: '', new: '', confirm: '' },

        // System Data
        infraStats: { total: 0, active: 0, items: [] },
        infraStatus: {}, // New
        pgInfo: {},
        pgQueries: [],
        kafkaMsg: '',
        sysInfo: { hostname: '', ip: '', disk: {} },

        // Redis Data // New

        // Redis Data
        redisPattern: '*',
        redisKeys: [],
        redisModalOpen: false,
        selectedRedisKey: '',
        selectedRedisValue: '',

        // Infrastructure Data
        redis: {},
        postgres: {},
        pgQueries: [],
        sqlQuery: "SELECT tablename FROM pg_catalog.pg_tables WHERE schemaname = 'public';",
        queryResults: null,
        queryError: null,
        isQueryRunning: false,
        kafka: {},
        cronJobs: [],

        // New Infrastructure
        storage: {},
        system: { cpu: {}, memory: {}, disk: {} },
        external: [],

        // System Graphs History
        sysHistory: {
            cpu: Array(20).fill(0),
            ram: Array(20).fill(0)
        },
        graphThrottle: 1000, // Default 1s
        lastGraphUpdate: 0,

        // Logs
        logThrottle: 1000,
        logInterval: null,

        // Notyf instance (Moved to window.notyf)

        // CodeMirror Instances
        cmInstances: {},

        async init() {
            // Theme init
            this.isDark = localStorage.getItem('theme') === 'dark' || (!('theme' in localStorage) && window.matchMedia('(prefers-color-scheme: dark)').matches);
            if (this.isDark) {
                document.documentElement.classList.add('dark');
            } else {
                document.documentElement.classList.remove('dark');
            }

            console.log("Dashboard App Initialized");

            // Initialize CodeMirror after Alpine mounts
            this.$nextTick(async () => {
                this.initCodeMirror('configEditor', 'yaml', 'configContent');
                this.initCodeMirror('bannerEditor', 'shell', 'bannerContent');
                this.initCodeMirror('sqlEditor', 'sql', 'sqlQuery');

                // Watchers to sync API data -> CodeMirror (Async Fetch)
                this.$watch('configContent', (val) => {
                    const cm = this.cmInstances['configEditor'];
                    if (cm && cm.getValue() !== val) cm.setValue(val);
                });
                this.$watch('bannerContent', (val) => {
                    const cm = this.cmInstances['bannerEditor'];
                    if (cm && cm.getValue() !== val) cm.setValue(val);
                });
                this.$watch('sqlQuery', (val) => {
                    const cm = this.cmInstances['sqlEditor'];
                    if (cm && cm.getValue() !== val) cm.setValue(val);
                });

                // SSE
                this.connectLogs();

                // Periodic
                await this.fetchStatus();
                this.fetchDummyStatus();
                this.fetchMonitoringConfig();
                this.fetchUserSettings(); // Load user settings for header
                setInterval(() => this.fetchStatus(), 5000);

                // Load data for badges
                this.fetchEndpoints();
                this.fetchCronJobs();

                // Watch tab changes to load data & refresh CodeMirror
                this.$watch('activeTab', (val) => {
                    // CodeMirror Refresh
                    this.$nextTick(() => {
                        if (val === 'config' && this.cmInstances['configEditor']) {
                            this.cmInstances['configEditor'].refresh();
                        }
                        if (val === 'banner' && this.cmInstances['bannerEditor']) {
                            this.cmInstances['bannerEditor'].refresh();
                        }
                        if (val === 'postgres' && this.cmInstances['sqlEditor']) {
                            this.cmInstances['sqlEditor'].refresh();
                        }
                    });

                    // Data Load
                    if (val === 'endpoints') this.fetchEndpoints();
                    if (val === 'redis') this.fetchRedisKeys();
                    if (val === 'postgres') { this.fetchPgQueries(); this.fetchPgInfo(); }
                    if (val === 'kafka') this.fetchKafka();
                    if (val === 'cron') this.fetchCronJobs();
                    if (val === 'config') this.fetchConfig();
                    if (val === 'banner') this.fetchBanner();
                    if (val === 'settings') this.fetchUserSettings();
                });
            });
        },


        get activeTabLabel() {
            return this.tabs.find(t => t.id === this.activeTab).label;
        },

        get activeEndpointsCount() {
            return this.endpoints.filter(e => e.active).length;
        },

        get activeCronCount() {
            return this.cronJobs.length;
        },

        async logout() {
            try {
                // POST to logout endpoint to clear session
                await fetch('/logout', {
                    method: 'POST'
                });
            } catch (error) {
                console.error('Logout error:', error);
            } finally {
                // Always redirect to login page (replace history to prevent back button)
                window.location.replace('/');
            }
        },

        async runQuery() {
            this.isQueryRunning = true;
            this.queryError = null;
            this.queryResults = null;

            try {
                const res = await fetch('/api/postgres/query', {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({ query: this.sqlQuery })
                });

                const data = await res.json();

                if (!res.ok) {
                    this.queryError = data.error || 'Query failed';
                } else {
                    this.queryResults = data;
                    if (Array.isArray(data) && data.length === 0) {
                        this.queryError = "No results found.";
                    }
                }
            } catch (err) {
                this.queryError = "Network error: " + err.message;
            } finally {
                this.isQueryRunning = false;
            }
        },

        async restartService() {
            if (!confirm('Are you sure you want to restart the service? This will briefly interrupt availability.')) return;

            try {
                const res = await fetch('/api/restart', { method: 'POST' });
                if (res.ok) {
                    this.showToast('Service is restarting...', 'success');
                    setTimeout(() => {
                        window.location.reload();
                    }, 3000);
                } else {
                    this.showToast('Failed to restart service', 'error');
                }
            } catch (err) {
                console.error("Restart error:", err);
                this.showToast('Failed to connect to server', 'error');
            }
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

        // Toast Methods (Wrapper for Notyf)
        showToast(message, type = 'info', title = '') {
            // Title is ignored in standard Notyf simple calls, appending if present
            const msg = title ? `<b>${title}</b><br>${message}` : message;

            if (!window.notyf) return; // Safety check

            if (type === 'success') {
                window.notyf.success(msg);
            } else if (type === 'error') {
                window.notyf.error(msg);
            } else {
                window.notyf.open({
                    type: 'info',
                    message: msg
                });
            }
        },

        // Logs
        updateThrottle() {
            if (this.logInterval) clearInterval(this.logInterval);
            this.setupLogFlush();
        },

        setupLogFlush() {
            const MAX_LOGS = 100;
            this.logInterval = setInterval(() => {
                if (this.logBuffer && this.logBuffer.length > 0) {
                    this.logs.push(...this.logBuffer);
                    this.logBuffer = [];

                    if (this.logs.length > MAX_LOGS) {
                        this.logs = this.logs.slice(-MAX_LOGS);
                    }

                    this.$nextTick(() => {
                        const box = document.getElementById('logs-box');
                        if (box) box.scrollTop = box.scrollHeight;
                    });
                }
            }, parseInt(this.logThrottle));
        },

        connectLogs() {
            const es = new EventSource('/api/logs');
            this.logBuffer = [];

            es.onmessage = (e) => {
                const msg = e.data.replace(/\u001b\[\d+m/g, "");
                this.logBuffer.push(msg);
            };

            this.setupLogFlush();
        },

        formatLog(logLine) {
            try {
                // Try parsing JSON first (Zerolog format)
                const data = JSON.parse(logLine);
                const time = data.time ? new Date(data.time).toLocaleTimeString() : '';
                const level = (data.level || 'UNKNOWN').toUpperCase();
                const msg = data.message || JSON.stringify(data);

                let badgeClass = 'bg-gray-100 text-gray-700 dark:bg-gray-800 dark:text-gray-300';
                if (level === 'INFO') badgeClass = 'bg-blue-100 text-blue-700 dark:bg-blue-900/30 dark:text-blue-300';
                if (level === 'WARN' || level === 'WARNING') badgeClass = 'bg-yellow-100 text-yellow-700 dark:bg-yellow-900/30 dark:text-yellow-300';
                if (level === 'ERROR' || level === 'FATAL') badgeClass = 'bg-red-100 text-red-700 dark:bg-red-900/30 dark:text-red-300';
                if (level === 'DEBUG') badgeClass = 'bg-purple-100 text-purple-700 dark:bg-purple-900/30 dark:text-purple-300';

                return `<div class="flex items-start gap-4 text-xs font-mono leading-relaxed">
                            <span class="text-muted-foreground w-[85px] shrink-0 pt-0.5">${time}</span>
                            <span class="px-1.5 py-0.5 rounded-[4px] font-semibold text-[10px] shrink-0 w-[50px] text-center ${badgeClass}">${level}</span>
                            <span class="text-foreground break-all pt-0.5">${msg}</span>
                        </div>`;
            } catch (e) {
                // Fallback for raw text
                let badgeClass = 'bg-gray-100 text-gray-700 dark:bg-gray-800 dark:text-gray-300';
                if (logLine.includes('INFO')) badgeClass = 'bg-blue-100 text-blue-700 dark:bg-blue-900/30 dark:text-blue-300';
                if (logLine.includes('WARN')) badgeClass = 'bg-yellow-100 text-yellow-700 dark:bg-yellow-900/30 dark:text-yellow-300';
                if (logLine.includes('ERROR') || logLine.includes('FAIL')) badgeClass = 'bg-red-100 text-red-700 dark:bg-red-900/30 dark:text-red-300';

                return `<div class="flex items-start gap-4 text-xs font-mono leading-relaxed">
                            <span class="text-muted-foreground w-[85px] shrink-0 pt-0.5">${new Date().toLocaleTimeString()}</span>
                             <span class="px-1.5 py-0.5 rounded-[4px] font-semibold text-[10px] shrink-0 w-[50px] text-center ${badgeClass}">RAW</span>
                            <span class="text-foreground break-all pt-0.5">${logLine}</span>
                        </div>`;
            }
        },

        async fetchStatus() {
            try {
                const res = await fetch('/api/status');
                const data = await res.json();

                // Services
                this.services = data.services || [];
                this.serviceCount = this.services.filter(s => s.active).length;

                // Infrastructure
                // Backend returns keys: redis, postgres, kafka, minio(storage), external
                // We map them to infraStatus for simple TRUE/FALSE checks or specific logic
                const infra = {
                    redis: data.redis && data.redis.connected,
                    postgres: data.postgres && data.postgres.connected,
                    kafka: data.kafka && data.kafka.connected,
                    minio: data.storage && data.storage.connected,
                    external: data.external && data.external.length > 0
                };

                // data.external is likely a map or list? 
                // HttpManager.GetStatus() returns map[string]interface{ "services": []... }?
                // No, HttpManager GetStatus returns map with "services" list.
                // Let's assume connected if at least one check ran? 
                // Actually External tab shows individual status.

                this.infraStatus = infra;

                // Count active infrastructure
                // defined as connected=true
                const infraKeys = Object.keys(infra);
                const activeInfra = Object.values(infra).filter(Boolean).length;
                this.infraStats = {
                    total: infraKeys.length,
                    active: activeInfra,
                    items: infraKeys.map(k => ({ name: k, active: infra[k] }))
                };

                // Update specific data sections
                this.redis = data.redis || {};
                this.postgres = data.postgres || {};
                this.kafka = data.kafka || {};
                this.storage = data.storage || {};
                this.external = data.external || [];

                // System Graphs Data
                if (data.system) {
                    this.system = data.system; // Fix: Assign entire system object for graphs

                    // Update History for Graphs
                    const now = Date.now();
                    if (now - this.lastGraphUpdate > (this.graphThrottle || 1000)) {
                        this.sysHistory.cpu.push(data.system.cpu?.usage_percent || 0);
                        this.sysHistory.cpu.shift();

                        this.sysHistory.ram.push(data.system.memory?.used_percent || 0);
                        this.sysHistory.ram.shift();

                        this.lastGraphUpdate = now;
                    }
                }

                // System Info (Host/IP) - Mapped from system_info (added in backend)
                if (data.system_info) {
                    this.sysInfo = {
                        hostname: data.system_info.hostname || 'Unknown',
                        ip: data.system_info.ip || 'Unknown',
                        disk: data.system?.disk || {}
                    };
                }

            } catch (e) { console.error("Fetch status error", e); }
        },

        async fetchDummyStatus() {
            try {
                const res = await fetch('/api/logs/dummy/status');
                const data = await res.json();
                this.dummyLogActive = data.active;
            } catch (e) { }
        },

        async fetchMonitoringConfig() {
            try {
                const res = await fetch('/api/monitoring/config');
                const data = await res.json();
                if (data.title) this.monitoringConfig.title = data.title;
                if (data.subtitle) this.monitoringConfig.subtitle = data.subtitle;
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
                    this.showToast(data.message, 'success', 'Saved');
                } else {
                    this.showToast(data.error || 'Failed to save', 'error', 'Error');
                }
            } catch (e) { this.showToast('Failed to save config', 'error', 'Error'); }
        },

        async backupConfig() {
            try {
                const res = await fetch('/api/config/backup', { method: 'POST' });
                const data = await res.json();
                if (res.ok) {
                    this.showToast(data.message, 'success', 'Backup Created');
                } else {
                    this.showToast(data.error || 'Backup failed', 'error', 'Error');
                }
            } catch (e) { this.showToast('Failed to backup config', 'error', 'Error'); }
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
                    this.showToast('Banner saved successfully!', 'success', 'Saved');
                } else {
                    this.showToast('Failed to save banner.', 'error', 'Error');
                }
            } catch (e) { this.showToast('Error saving banner', 'error', 'Error'); }
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
        },

        // User Settings Methods
        async fetchUserSettings() {
            try {
                const res = await fetch('/api/user/settings');
                const data = await res.json();
                this.userSettings.username = data.username || 'Admin';
                this.userSettings.photoPath = data.photo_path || '';
            } catch (e) { }
        },

        async updateUsername() {
            try {
                const res = await fetch('/api/user/settings', {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({ username: this.userSettings.username })
                });
                const data = await res.json();
                this.showToast(data.message || 'Username updated', 'success', 'Updated');
            } catch (e) {
                this.showToast('Failed to update username', 'error', 'Error');
            }
        },

        async changePassword() {
            if (this.passwordForm.new !== this.passwordForm.confirm) {
                this.showToast('New passwords do not match', 'error', 'Validation');
                return;
            }
            try {
                const res = await fetch('/api/user/password', {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({
                        current_password: this.passwordForm.current,
                        new_password: this.passwordForm.new
                    })
                });
                const data = await res.json();
                if (res.ok) {
                    this.showToast(data.message || 'Password changed', 'success', 'Success');
                    this.passwordForm = { current: '', new: '', confirm: '' };
                } else {
                    this.showToast(data.error || 'Failed to change password', 'error', 'Error');
                }
            } catch (e) {
                this.showToast('Failed to change password', 'error', 'Error');
            }
        },

        async uploadPhoto(event) {
            const file = event.target.files[0];
            if (!file) return;

            const formData = new FormData();
            formData.append('photo', file);

            try {
                const res = await fetch('/api/user/photo', {
                    method: 'POST',
                    body: formData
                });
                const data = await res.json();
                if (res.ok) {
                    this.userSettings.photoPath = data.photo_path;
                    this.showToast(data.message || 'Photo uploaded', 'success', 'Success');
                } else {
                    this.showToast(data.error || 'Upload failed', 'error', 'Error');
                }
            } catch (e) {
                this.showToast('Upload failed', 'error', 'Error');
            }
        },

        async deletePhoto() {
            if (!confirm('Delete profile photo?')) return;
            try {
                const res = await fetch('/api/user/photo', { method: 'DELETE' });
                const data = await res.json();
                if (res.ok) {
                    this.userSettings.photoPath = '';
                    this.showToast(data.message || 'Photo deleted', 'success', 'Deleted');
                } else {
                    this.showToast(data.error || 'Delete failed', 'error', 'Error');
                }
            } catch (e) {
                this.showToast('Delete failed', 'error', 'Error');
            }
        },

        initCodeMirror(id, mode, model) {
            const el = document.getElementById(id);
            if (!el) return;

            // Prevent double init
            if (this.cmInstances[id]) return;

            const cm = CodeMirror.fromTextArea(el, {
                mode: mode,
                theme: 'dracula',
                lineNumbers: true,
                lineWrapping: true
            });

            // Two-way binding: Update Alpine data on change
            cm.on('change', () => {
                this[model] = cm.getValue();
            });

            // Set initial value from Alpine data
            if (this[model]) {
                cm.setValue(this[model]);
            }

            this.cmInstances[id] = cm;
        }
    }
}
