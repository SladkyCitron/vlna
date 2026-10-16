# 🎵 Vlna

[![GitHub license](https://img.shields.io/github/license/SladkyCitron/vlna)](LICENSE) ![Made in Slovakia](https://raw.githubusercontent.com/pedromxavier/flag-badges/refs/heads/main/badges/SK.svg)

**Vlna** (pronounced _[ˈvl̩na]_, meaning "wave" in Slovak) is a lightweight internet radio player built with Go + Wails 3 + Svelte 5.

**Video Demo:** <URL HERE>

## ✨ Features

- Play internet radio stations from all over the world
- Search and favorite stations
- Real-time display of currently playing song title and artist (ICY metadata)
- Light and dark theme and language with persistent storage
- Full MPRIS support and media keys support
- Discord Rich Presence integration
- Cross-platform support
- Super lightweight and resource-friendly

## 📦 Installation (from source)

1. Install [Go](https://go.dev/doc/install)
2. Install a C compiler
3. Install Node.js and npm
4. Install [Wails 3](https://v3.wails.io/quick-start/installation)
5. Clone the repository:

```sh
git clone https://github.com/SladkyCitron/vlna.git
cd vlna
```

5. Build the app:

```sh
wails3 build
```

The built app will be located in the `bin` directory.

## ⚖️ License

MIT License (see [LICENSE](LICENSE))
