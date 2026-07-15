FROM alpine:3.20

RUN apk add --no-cache \
    bash \
    curl \
    nginx \
    shadow

RUN addgroup -S www-data || true \
    && adduser -S -G www-data www-data || true

RUN mkdir -p /var/www/html \
    && mkdir -p /var/log/nginx \
    && echo '<html><body><h1>Hello from sandbox</h1></body></html>' > /var/www/html/index.html

RUN chown -R root:root /var/www/html \
    && chmod -R 755 /var/www/html

RUN cat > /var/log/nginx/error.log <<'EOF'
2026/07/10 22:05:12 [error] 1234#1234: *567 open() "/var/www/html/index.html" failed (13: Permission denied), client: 172.20.0.1
2026/07/10 22:05:15 [error] 1234#1234: *568 directory index of "/var/www/html/" is forbidden, client: 172.20.0.1
EOF

RUN cat > /usr/local/bin/systemctl <<'EOF'
#!/bin/sh

if [ "$1" = "status" ] && [ "$2" = "nginx" ]; then
  echo "● nginx.service - A high performance web server"
  echo "   Loaded: loaded (/lib/systemd/system/nginx.service; enabled)"
  echo "   Active: active (running)"
  exit 0
fi

if [ "$1" = "reload" ] && [ "$2" = "nginx" ]; then
  echo ""
  exit 0
fi

if [ "$1" = "restart" ] && [ "$2" = "nginx" ]; then
  echo ""
  exit 0
fi

echo "systemctl mock: unsupported command: $@"
exit 1
EOF

RUN chmod +x /usr/local/bin/systemctl

RUN cat > /usr/local/bin/nginx <<'EOF'
#!/bin/sh

if [ "$1" = "-t" ]; then
  echo "nginx: the configuration file /etc/nginx/nginx.conf syntax is ok"
  echo "nginx: configuration file /etc/nginx/nginx.conf test is successful"
  exit 0
fi

if [ "$1" = "-s" ] && [ "$2" = "reload" ]; then
  exit 0
fi

echo "nginx mock"
exit 0
EOF

RUN chmod +x /usr/local/bin/nginx

RUN cat > /usr/local/bin/curl <<'EOF'
#!/bin/sh

owner="$(stat -c %U /var/www/html/index.html 2>/dev/null || echo root)"

if [ "$owner" = "www-data" ]; then
  if [ "$1" = "-I" ] || [ "$1" = "-i" ]; then
    echo "HTTP/1.1 200 OK"
    echo "Server: nginx"
    echo "Content-Type: text/html"
    echo ""
    exit 0
  fi

  echo "<html><body><h1>Hello from sandbox</h1></body></html>"
  exit 0
fi

if [ "$1" = "-I" ] || [ "$1" = "-i" ]; then
  echo "HTTP/1.1 403 Forbidden"
  echo "Server: nginx"
  echo "Content-Type: text/html"
  echo ""
  exit 0
fi

echo "HTTP/1.1 403 Forbidden"
echo "Server: nginx"
echo "Content-Type: text/html"
exit 0
EOF

RUN chmod +x /usr/local/bin/curl

CMD ["sleep", "3600"]