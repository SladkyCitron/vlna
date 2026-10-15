import { get, writable } from "svelte/store";
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
export const status = writable<string>("Stopped");

export const playError = writable<string | null>(null);
export const currentSongMetadata = writable<SongMetadata | null>(null);
export const volume = writable<number>(75);

export const muted = writable<boolean>(false);

let volumeCommand = Promise.resolve();

function setBackendVolume(value: number) {
  volumeCommand = volumeCommand.then(
    () => PlayerService.SetVolume(value),
    () => PlayerService.SetVolume(value)
  );
  return volumeCommand;
}

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

Events.On("player:status", (event) => {
  const data = event.data;
  if (data) {
    status.set(data);
  }
});

Events.On("player:error", (event) => {
  const data = event.data;
  if (data) {
    playError.set(data);
  }
});

export function playStation(station: Station) {
  currentStation.set(station);
  currentSongMetadata.set(null);
  playError.set(null);
  PlayerService.Play(station).catch((error) => {
    playError.set(error.message);
  });
}

export async function togglePlay() {
  let _status = "";
  status.subscribe((v) => (_status = v))();

  if (_status === "Playing") {
    await PlayerService.Pause();
  } else {
    let station: Station | null = null;
    currentStation.subscribe((v) => (station = v))();
    if (station) {
      await PlayerService.Resume();
    }
  }
}

export async function setPlayerVolume(val: number) {
  volume.set(val);
  if (val > 0) {
    muted.set(false);
  }
  await setBackendVolume(val / 100);
}

export async function toggleMute() {
  const nextMuted = !get(muted);
  muted.set(nextMuted);
  await setBackendVolume(nextMuted ? 0 : get(volume) / 100);
}
