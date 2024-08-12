# Towel Tornado - The Towel Waving Game

Towel Tornado is a game for Sauna Masters (Aufguss) where you use a towel to blow the sauna hats off three test dummies in a sauna. It's a fun spin on the art of waving a towel where technique, power, stamina, and a bit of luck all contribute to success. 

Check it out in action in this short video:
https://www.youtube.com/watch?v=WVP7yY2LHeM

## Requirements
- A computer (Mac, Windows, Linux) to run the game software
- USB Hub
- Screen
- Three Anemometers (recommended: Habotest Digital Anemometer HT625B USB, other types may work as well)
- Structure to hold the anemometers

## Setup Instructions
- Build a structure to hold the three anemometers, spaced 60 cm apart, and positioned approximately 130 cm above the ground.
- Connect the anemometers with USB cables to a USB hub.
- Connect the USB hub to a computer running the Towel Tornado application.
- Connect the computer to a screen positioned behind the anemometers.
- Open a web browser and navigate to the URL provided by the Towel Tornado application when you started it (likely http://localhost:8080/). Display the game view in full-screen kiosk mode (Mac: fn+F and Shift+Cmd+F).
- Use another browser tab or device to control the game at http://localhost:8080/admin.html. You can control the game from the computer or a phone/tablet.
- Turn on the anemometers one at a time, starting with the left-most one. This will ensure they register in the correct order, and you will see each anemometer reported in the console of Towel Tornado (which runs from the command line).
- On the Habotest Digital Anemometer HT625B, enable the USB function by holding the "H" button for a few seconds. After about 20 minutes, the anemometers will turn off automatically.

## Notes
- The game hi-score is saved into a json file in the current directory. If you want to run a new session, simply delete or rename this file.

# Software
In this repo you will find the software for the Towel Tornado cross platform application. It's written in the awesome language, Go-lang (Google language). You are welcome to improve the software. It's licensed under MIT License. 

In the binaries/ folder you will find downloadable executables in the zip-files. Run them from the Command window / Terminal on your particular system. The UI is found in a web browser when the application is started.

# Background

The game is developed voluntarily by Kasper Skårhøj for Gusmesterforeningen (Sauna Master Association) in Denmark. Therefore the logo of Gusmesterforeningen is integrated in the graphics.

# Development Notes

## Todo
- Issue: Registering a time of 0.00 as well as 60.0 seconds has been seen - felt like a bug that had to be fixed.
- Would be nice to have alternative ways to register the anemometers
- Additional game ideas:
    - Tetris with Kerkes sauna stones - using the wind data to rotate left / right
    - The classic Hammer Amusement Machine

## Usage on a Raspberry Pi notes (slightly internal)
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