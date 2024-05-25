# vertigo
track if it is a hotdog or is not a hot dog

Where to go? -> vertigo!

`go run ./cmd/vertigo`


Build commandline tool 

`go build -o vertigo ./cmd/vertigo`

List all shoes

`./vertigo -list shoes`

Add shoe

`go run ./cmd/vertigo/ -discord -add https://stockx.com/air-jordan-1-retro-high-travis-scott`

Set up discord bot tokens etc in a `.env` file in the root directory, just as the `.env.template`.
After doing that, with `--discord` as a flag while executing `.vertigo` with an `image_url`, you will send a notification to the channel set up in the `.env` file. 

e.g: `./vertigo --discord --add shoes name="Mars Yard" brand=Nike silhouette="Mars Yard" image_url=https://content.deadstock.de/media/pages/uploads/2017/07/136e53244e-1706280229/nikelab-tom-sachs-mars-yard-2-global-release-info-1-750x450-crop.webp"`
