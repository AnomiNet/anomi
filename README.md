# Anomi

Backend API for [Anomi](http://anomi.net), a discussion and content sharing website.

## Usage

Install go, then

```bash
go get github.com/anominet/anomi
go install github.com/anominet/anomi
```

Run `anomi -h` for help.


### Building docker container

Make sure you have [docker installed](https://docs.docker.com/installation/) and working.

```bash
cd ~/$GOPATH/src/github.com/anominet/anomi
make
```

### Running docker container

First start a redis container

```bash
docker run --name anomi-redis -v $DATA_DIR:/data -d redis redis-server --appendonly yes
```

## Deploying Anomi

Generate a rails secret key and create the file `.env.web`:

```bash
RAILS_SECRET_KEY=YOUR_SECRET_KEY_HERE
```

Then run:

```bash
docker-compose up
```

## License

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this project except in compliance with the License.
You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
