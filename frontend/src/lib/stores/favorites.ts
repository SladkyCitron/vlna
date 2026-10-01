import { writable } from "svelte/store";
import {
  type Station,
  type Stations,
  ConfigService,
} from "$bindings/github.com/SladkyCitron/vlna/service";

export const favorites = writable<Stations>([]);
export const favoritesLoading = writable(true);

let configLoaded = false;
let loadPromise: Promise<void> | undefined;

export function loadFavorites(): Promise<void> {
  if (loadPromise) {
    return loadPromise;
  }

  loadPromise = ConfigService.GetFavorites()
    .then((stations) => {
      favorites.set(stations ?? []);
    })
    .catch((error) => {
      favorites.set([]);
      console.error("Failed to load favorites:", error);
    })
    .finally(() => {
      configLoaded = true;
      favoritesLoading.set(false);
    });

  return loadPromise;
}

export function toggleFavorite(station: Station) {
  if (!configLoaded) {
    void loadFavorites().then(() => toggleFavorite(station));
    return;
  }

  favorites.update((currentFavorites) => {
    const storedFavorites = currentFavorites ?? [];
    const stationIndex = storedFavorites.findIndex(
      (favorite) => favorite.stationuuid === station.stationuuid
    );

    if (stationIndex === -1) {
      return [...storedFavorites, station];
    }

    return storedFavorites.filter((_, index) => index !== stationIndex);
  });
}

favorites.subscribe((value) => {
  if (!configLoaded || value === undefined) {
    return;
  }

  ConfigService.SetFavorites(value).catch((error) => {
    console.error("Failed to save favorites:", error);
  });

  ConfigService.SaveConfig().catch((error) => {
    console.error("Failed to save config:", error);
  });
});
