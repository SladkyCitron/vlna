import { writable } from "svelte/store";
import { type Station } from "$bindings/github.com/SladkyCitron/vlna/service";

export const selectedStation = writable<Station | null>(null);
export const isStationDetailsOpen = writable<boolean>(false);

export function openStationDetails(station: Station) {
  console.log("openStationDetails", station);
  selectedStation.set(station);
  isStationDetailsOpen.set(true);
}

export function closeStationDetails() {
  console.log("closeStationDetails");
  selectedStation.set(null);
  isStationDetailsOpen.set(false);
}
