import { writable } from "svelte/store";
import {
  type Station,
  PlayerService,
} from "$bindings/github.com/SladkyCitron/vlna/service";
import { Events } from "@wailsio/runtime";

export interface SongMetadata {
  title: string;
  artist: string;
}

export const currentStation = writable<Station | null>(null);
export const isPlaying = writable<boolean>(false);
export const currentSongMetadata = writable<SongMetadata | null>(null);
export const volume = writable<number>(75);

Events.On("player:icy-metadata", (event) => {
  const data = event.data;
  if (data?.StreamTitle) {
    const parts = data.StreamTitle.split(" - ");
    if (parts.length >= 2) {
      currentSongMetadata.set({
        artist: parts[0].trim(),
        title: parts.slice(1).join(" - ").trim(),
      });
    } else {
      currentSongMetadata.set({ artist: "", title: data.StreamTitle });
    }
  }
});

export async function playStation(station: Station) {
  currentStation.set(station);
  currentSongMetadata.set(null);
  isPlaying.set(true);
  await PlayerService.Play(station.url_resolved || station.url);
}

export async function togglePlay() {
  let playing = false;
  isPlaying.subscribe((v) => (playing = v))();

  if (playing) {
    isPlaying.set(false);
    await PlayerService.Pause();
  } else {
    let station: Station | null = null;
    currentStation.subscribe((v) => (station = v))();
    if (station) {
      isPlaying.set(true);
      await PlayerService.Resume();
    }
  }
}

export async function setPlayerVolume(val: number) {
  volume.set(val);
  await PlayerService.SetVolume(val / 100);
}
