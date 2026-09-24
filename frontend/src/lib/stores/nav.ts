import { writable } from "svelte/store";

export type View = "explore" | "search" | "favorites" | "settings";

export const activeView = writable<View>("explore");
