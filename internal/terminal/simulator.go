package terminal

import "strings"

func SimulateCommandOutput(command string) string {
	normalized := strings.TrimSpace(strings.ToLower(command))

	switch {
	case normalized == "":
		return ""

	case strings.HasPrefix(normalized, "curl -i"):
		return "HTTP/1.1 403 Forbidden\nServer: nginx\nContent-Type: text/html\n\n"

	case strings.HasPrefix(normalized, "curl"):
		return "HTTP/1.1 403 Forbidden\nServer: nginx\nContent-Type: text/html\n"

	case strings.Contains(normalized, "tail") && strings.Contains(normalized, "/var/log/nginx/error.log"):
		return `2026/07/10 22:05:12 [error] 1234#1234: *567 open() "/var/www/html/index.html" failed (13: Permission denied), client: 172.20.0.1
				2026/07/10 22:05:15 [error] 1234#1234: *568 directory index of "/var/www/html/" is forbidden, client: 172.20.0.1
				2026/07/10 22:05:20 [crit] 1234#1234: *569 SSL: error:02001002: system library:fopen:No such file or directory
				`

	case strings.Contains(normalized, "grep") && strings.Contains(normalized, "permission"):
		return `2026/07/10 22:05:12 [error] 1234#1234: *567 open() "/var/www/html/index.html" failed (13: Permission denied), client: 172.20.0.1
				`

	case strings.HasPrefix(normalized, "ls -la /var/www/html"):
		return `total 16
				drwxr-xr-x 2 root root     4096 Jul 10 22:03 .
				drwxr-xr-x 3 root root     4096 Jul 10 22:03 ..
				-rw-r--r-- 1 root root      612 Jul 10 22:03 index.html
				`

	case strings.HasPrefix(normalized, "chown -r www-data:www-data /var/www/html"):
		return ""

	case normalized == "nginx -t":
		return `nginx: the configuration file /etc/nginx/nginx.conf syntax is ok
				nginx: configuration file /etc/nginx/nginx.conf test is successful
				`

	case normalized == "systemctl reload nginx":
		return ""

	case normalized == "systemctl restart nginx":
		return ""

	case normalized == "systemctl status nginx":
		return `● nginx.service - A high performance web server
				Loaded: loaded (/lib/systemd/system/nginx.service; enabled)
				Active: active (running)
				`

	case normalized == "free -h":
		return `              total        used        free      shared  buff/cache   available
				Mem:           1.9Gi       1.8Gi        74Mi        12Mi        60Mi        80Mi
				Swap:             0B          0B          0B
				`

	case normalized == "htop":
		return "Interactive process viewer opened. High memory usage detected earlier.\n"

	case strings.Contains(normalized, "dmesg") && strings.Contains(normalized, "oom"):
		return `[2026-07-10 22:12:45] Out of memory: Killed process 4567 (postgres) total-vm:2845678kB, anon-rss:1245678kB
				[2026-07-10 22:12:46] Out of memory: Killed process 3789 (gunicorn)
				`

	case normalized == "systemctl status postgresql":
		return `● postgresql.service - PostgreSQL database server
				Loaded: loaded
				Active: failed
				`

	case normalized == "systemctl restart postgresql":
		return ""

	case normalized == "docker compose restart backend":
		return "Container backend restarted\n"

	case normalized == "docker compose ps":
		return `NAME        STATUS
				backend     running
				postgres    running
				`

	case normalized == "ping google.com":
		return "ping: google.com: Temporary failure in name resolution\n"

	case normalized == "ping 8.8.8.8":
		return `PING 8.8.8.8 (8.8.8.8): 56 data bytes
				64 bytes from 8.8.8.8: icmp_seq=0 ttl=117 time=12.4 ms
				`

	case normalized == "nslookup google.com":
		return ";; connection timed out; no servers could be reached\n"

	case normalized == "cat /etc/resolv.conf":
		return "nameserver 8.8.3.8\n"

	case strings.Contains(normalized, "nameserver 8.8.8.8") && strings.Contains(normalized, "/etc/resolv.conf"):
		return ""

	case normalized == "systemctl restart systemd-resolved":
		return ""

	case normalized == "docker ps":
		return `CONTAINER ID   IMAGE          COMMAND                  STATUS
				a1b2c3d4e5f6   auth-service   "npm start"              Restarting (1) 42 seconds ago
				b2c3d4e5f6a1   postgres       "docker-entrypoint..."   Up 10 minutes
				`

	case strings.HasPrefix(normalized, "docker logs auth-service"):
		return `ERROR: connection to database "auth" failed: could not connect to server: Connection refused
				Is the server running on host "db" and accepting TCP/IP connections on port 5432?
				`

	case normalized == "cat .env | grep db_":
		return `DB_HOST=db
				DB_PORT=5433
				DB_NAME=auth
				DB_USER=auth
				`

	case strings.Contains(normalized, "grep -r") && strings.Contains(normalized, "5433"):
		return "./.env:DB_PORT=5433\n"

	case strings.Contains(normalized, "sed -i") && strings.Contains(normalized, "5433") && strings.Contains(normalized, "5432"):
		return ""

	case strings.HasPrefix(normalized, "docker compose up -d --force-recreate auth-service"):
		return "Container auth-service recreated\n"

	case strings.Contains(normalized, "sed -n") && strings.Contains(normalized, "/etc/nginx/sites-available/default"):
		return `40: server {
				41:     listen 80;
				42:     server_name example.com
				43:
				44:     location / {
				45:         root /var/www/html;
				46:     }
				47: }
				`

	case strings.HasPrefix(normalized, "vim"):
		return "File opened and saved.\n"

	default:
		return "command executed\n"
	}
}
