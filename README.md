
# Setup

- ***NOTE: This guide assumes you already have some knowledge on how to navigate an Ubuntu linux server***

## Installation

- Install Golang
- Install nginx for **reverse proxy**
- Install supervisor to **load your site**



## Golang site

1. `cd ~/go/src`
2. `git clone http://github.com/user/repo`
3. `cd your_git_folder`
4. Build site executable
    - `go build -o site_executable_name cmd/web/*.go`

## NGINX

1. Create config
    1. `cd /etc/nginx/sites-available`
    2. `sudo nano *config_name*`
    3. Fill the file with this base config:
    ```
    server {
        listen 80;
        listen [::]:80;
        listen 443;
        server_name example.example.com;

        location / {
            proxy_pass http://localhost:8080;
            proxy_http_version 1.1;
            proxy_set_header Upgrade $http_upgrade;
            proxy_set_header Connection keep-alive;
            proxy_set_header Host $host;
            proxy_cache_bypass $http_upgrade;
            proxy_set_header X-Real-IP $remote_addr;
            proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
            proxy_set_header X-Forwarded-Proto $scheme;
        }
    }
    ```
    4. Change `example.example.com` to your domain
    5. Change `8080` to your port
    6. Save and exit 
        - `ctrl+x`, `y`, `enter`
2. Link the file to sited-enabled
    - `sudo ln -s /etc/nginx/sites-available/config_name /etc/nginx/sites-enabled/config_name`
3. Test the config
    - `sudo service nginx configtest`

3. Restart nginx
    - `sudo service nginx restart`


## Supervisor

1. Create config
    1. `cd /etc/supervisor/conf.d`
    2. `sudo nano name.conf`
    3. Fill the file with this base config:
    ```
    [program:example_name]
    directory=/home/youruser/go/src/site_folder
    command=/home/youruser/go/src/site_folder/site_executable_name
    autorestart=true
    user=youruser
    redirect_stderr=true
    stdout_logfile=/home/youruser/logs/site.log
    stdout_logfile_maxbytes=50MB
    startretries=3
    startsecs=0
    ```
    4. Replace the following parameters
        - `example_name`
        - `youruser`
        - `site_folder`
        - `site_executable_name`
    6. Save and exit 
        - `ctrl+x`, `y`, `enter`

2. `sudo supervisorctl reread`
3. `sudo supervisorctl update`
4. `sudo supervisorctl restart program_name`

<hr>

Hopefully everything will work with this

Good luck!
