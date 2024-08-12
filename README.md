# toweltornado
Source code for the Towel Tornado game


***************
DEVELOPMENT:
***************

Compile and upload (no USB support, 32 bit and 64 bit):
env GOOS=linux GOARCH=arm64 go build . && scp toweltornado scp://pi@192.168.11.23




sudo nano /etc/systemd/system/toweltornado.service

[Unit]
Description=Towel Tornado

[Service]
ExecStart=/home/pi/toweltornado
Restart=always
User=pi
Group=pi
WorkingDirectory=/home/pi

[Install]
WantedBy=multi-user.target


sudo systemctl daemon-reload
sudo systemctl enable toweltornado.service
sudo systemctl restart toweltornado.service
sudo systemctl status toweltornado.service
journalctl -u toweltornado.service -f



*******************
http://192.168.11.23:8031/

Changes to Media Player .23:
- Has toweltornado on it
- Has toweltornado.service running
- Has 'sudo systemctl disable piplayer.service'
- Has browser kiosk on start-up:
    sudo nano /etc/xdg/autostart/display.desktop
    
    [Desktop Entry]
    Name=Chrome
    Exec=chromium-browser --app=http://192.168.11.23:8031 --kiosk --app-shell-hosted


- sudo nano /etc/lightdm/lightdm.conf
    Added xserver-command=X -s 0 -dpms