# Ts3telegrambot
Teamspeak 3 2 Telegram Bot written in go 

####Whats this? What does it do?
Bot sends messages via Telegram ( https://telegram.org/) when a Client joins your Teamspeak Server directly to your Telegram account. So your phone may ring if someone joins your Teamspeak server. 
It also stores times of how long a client has been on your server. (Identified by username currently)


#####How does it work?
Bot connects to your Teamspeak Server using the Query Port and listens to the update channel. 

####How to use?
Create a TOML config file with your TS3 Credentials and your Telegram Chat id / Bot Api key. Compile it with go and start it. Join your TS Server
Done

Questions?

